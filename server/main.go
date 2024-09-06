package main

import (
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/rs/cors"
	"golang.org/x/net/websocket"
)

type Server struct {
	mu    sync.Mutex
	conns map[*websocket.Conn]bool
	count int
}

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
		count: 0,
	}
}

func (s *Server) HandleWs(ws *websocket.Conn) {
	log.Println("connected:", ws.RemoteAddr())

	s.mu.Lock()
	s.conns[ws] = true
	s.mu.Unlock()

	s.Listen(ws)
}

func (s *Server) Listen(ws *websocket.Conn) {
	buf := make([]byte, 1024)

	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}

			log.Println("read error:", err)
			continue
		}

		msg := buf[:n]
		log.Println("received:", string(msg))

		i, err := strconv.Atoi(string(msg))
		if err != nil {
			continue
		}

		s.mu.Lock()
		s.count = i
		s.mu.Unlock()

		s.BroadcastCount()
	}
}

func (s *Server) BroadcastCount() {
	for ws := range s.conns {
		ws.Write([]byte(strconv.Itoa(s.count)))
	}
}

func (s *Server) Pong() {
	for range time.Tick(time.Second) {
		s.count++
		s.BroadcastCount()
	}
}

func (s *Server) HandleCount() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		w.Write([]byte(strconv.Itoa(s.count)))
		s.mu.Unlock()
	}
}

func main() {
	s := NewServer()

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
	})
	handler := c.Handler(http.DefaultServeMux)

	go s.Pong()

	http.Handle("/ws", websocket.Handler(s.HandleWs))
	http.HandleFunc("/api/count", s.HandleCount())
	log.Println("server started!")
	http.ListenAndServe(":8080", handler)
}
