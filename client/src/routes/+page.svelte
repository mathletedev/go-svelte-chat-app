<script lang="ts">
    import { SERVER_URL, WS_URL } from "$lib/config";
    import { onMount } from "svelte";

    let socket: WebSocket;

    let count = 0;

    onMount(async () => {
        const res = await fetch(`${SERVER_URL}/api/count`);
        const data = await res.json();

        count = data;

        socket = new WebSocket(WS_URL);

        socket.onmessage = event => {
            const i = parseInt(event.data);
            if (!isNaN(i)) {
                count = i;
            }
        };
    });

    const increment = () => {
        socket.send((count + 1).toString());
    };

    const sendMessage = () => {
        socket.send("Hello, world!");
    };
</script>

<h1>Welcome to SvelteKit</h1>
<p>
    Visit <a href="https://kit.svelte.dev">kit.svelte.dev</a> to read the documentation
</p>
<p>Count: {count}</p>
<button on:click={increment}>Increment</button>
<button on:click={sendMessage}>Send message</button>
