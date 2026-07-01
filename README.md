# ChatAppGo

A small real-time chat room built with Go, Fiber, WebSockets, HTMX, and server-rendered HTML templates.

The app serves a simple chat page at `/`, opens a WebSocket connection at `/ws`, and broadcasts submitted messages to every connected client.

## Features

- Real-time message broadcasting with WebSockets
- Server-rendered HTML using Fiber's template engine
- HTMX WebSocket extension for sending messages from the browser
- Static CSS served from the `static` folder
- Simple project structure for learning Go web apps

## Tech Stack

- Go
- Fiber
- Fiber WebSocket
- Fiber HTML templates
- HTMX
- Bootstrap

## Getting Started

Clone the project, install dependencies, and run the app:

```bash
go mod tidy
go run .
```

Then open:

```text
http://localhost:3000
```

You can test the chat by opening the page in two browser tabs and sending messages between them.

## Project Structure

```text
.
├── main.go              # App setup, routes, and server start
├── websocket.go         # WebSocket connection and broadcast logic
├── messages.go          # Message data model
├── handlers/
│   └── handlers.go      # Page handlers
├── views/
│   ├── index.html       # Main chat page
│   └── message.html     # Rendered message snippet
├── static/
│   └── style.css        # App styling
├── go.mod
└── go.sum
```

## How It Works

1. `main.go` starts a Fiber server on port `3000`.
2. The `/` route renders `views/index.html`.
3. The browser connects to `/ws` using the HTMX WebSocket extension.
4. When a user submits the chat form, HTMX sends the message over the WebSocket.
5. `websocket.go` receives the message, renders `views/message.html`, and broadcasts it to connected clients.

## Useful Commands

Run the app:

```bash
go run .
```

Check that the project compiles:

```bash
go test ./...
```

Format Go files:

```bash
gofmt -w .
```
