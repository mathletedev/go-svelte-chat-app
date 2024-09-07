<script lang="ts">
    import { SERVER_URL, WS_URL } from "$lib/config";
    import { onMount } from "svelte";

    let socket: WebSocket;

    let name = "";
    let message = "";
    let messages = [] as string[];

    onMount(async () => {
        const res = await fetch(`${SERVER_URL}/api/messages`);
        const data = await res.json();

        messages = data;

        socket = new WebSocket(WS_URL);

        socket.onmessage = event => {
            messages = [...messages, JSON.parse(event.data)];
        };
    });

    const sendMessage = () => {
        if (name == "" || message == "") {
            alert("Name and message are required");
            return;
        }

        socket.send(`${name}: ${message}`);

        name = message = "";
    };
</script>

<div class="flex h-screen flex-col gap-8">
    <div class="mx-8 mt-8 flex flex-grow flex-col gap-2 overflow-y-scroll">
        {#each messages as message}
            <p class="rounded border-l-2 border-violet-500 pl-4">{message}</p>
        {/each}
    </div>
    <form class="mx-8 mb-8 flex gap-8" on:submit={sendMessage}>
        <input bind:value={name} placeholder="Name" />
        <input class="w-full" bind:value={message} placeholder="Message" />
        <button
            class="bg-violet-500 text-white hover:bg-violet-600"
            type="submit">Send</button
        >
    </form>
</div>
