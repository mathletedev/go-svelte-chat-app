package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/rs/cors"
	"golang.org/x/net/websocket"
)

type Server struct {
	mu       sync.Mutex
	conns    map[*websocket.Conn]bool
	messages []string
}

func NewServer() *Server {
	return &Server{
		conns:    make(map[*websocket.Conn]bool),
		messages: make([]string, 0),
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

		msg := string(buf[:n])
		log.Println("received:", msg)

		s.mu.Lock()
		s.messages = append(s.messages, msg)
		s.mu.Unlock()

		s.BroadcastMessage(msg)
	}
}

func (s *Server) BroadcastMessage(msg string) {
	encoded, err := json.Marshal(msg)
	if err != nil {
		log.Println("json error:", err)
		return
	}

	for ws := range s.conns {
		ws.Write([]byte(encoded))
	}
}

func (s *Server) HandleGetMessages() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		encoded, err := json.Marshal(s.messages)
		if err != nil {
			w.WriteHeader(500)
		}

		s.mu.Lock()
		w.Write([]byte(encoded))
		s.mu.Unlock()
	}
}

func main() {
	s := NewServer()

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
	})
	handler := c.Handler(http.DefaultServeMux)

	http.Handle("/ws", websocket.Handler(s.HandleWs))
	http.HandleFunc("/api/messages", s.HandleGetMessages())
	log.Println("server started!")
	http.ListenAndServe(":8080", handler)
}
