# gobcast

A WebSocket broadcast server written in Go — project specification from [roadmap.sh](https://roadmap.sh/projects/broadcast-server). Clients connect, authenticate, set a nickname, and broadcast messages to every other connected client.

## Features

- **WebSocket broadcasting** — messages are sent from one client to all connected clients
- **Nicknames** — each client picks a name, prepended to every message they send
- **Authentication** — clients must present a shared token before joining the broadcast pool
- **CLI flags** — `--port`, `--host`, `--token`, `--nickname` on both server (`start`) and client (`connect`)
- **Graceful shutdown** — SIGINT/SIGTERM cleanly tears down connections
- **Connection tracking** — server logs the current number of connected clients on join and leave

## Building and running

Build and run with Go:

```bash
go build -o gobcast .

# Start server
./gobcast start --port 8080 --token mysecret

# Connect a client
./gobcast connect --port 8080 --token mysecret --nickname anakin
```

Or run directly without building:

```bash
go run . start --port 8080 --token mysecret
go run . connect --port 8080 --token mysecret --nickname anakin
```

Type a message into the client's stdin and press enter — it gets broadcast to every connected client as `anakin: your message`.

## Future improvements

- **Message history** — new clients receive the last N messages when they join, so they're not greeted by an empty screen
- **Private messages** — support direct messages to a specific client via some addressing scheme (`/msg bob hello`)