# Gopher-Chat

If you’re looking for the best alternative to Telegram and have ended up here, then…well, I’m afraid I’ll have to disappoint you – your search must continue. This is a personal project created solely to explore technologies such as HTTP and WebSocket in Go

## Overview
It's a multi-room chat application written in Go, with a desktop client available on [Fyne](https://fyne.io/). The server and client run within a single process: the WebSocket server runs in a background goroutine, whilst the GUI runs in the main thread.

## Features

- Create and select chat rooms from a single window
- Real-time messaging via WebSocket
- Multiple users in a single room, with the author of each message correctly displayed
- Support for multiple open rooms simultaneously - messages in inactive rooms are not lost

## Technologies

- Go (goroutines, channels,`sync.Mutex`, `encoding/json`)
- [gorilla/websocket](https://github.com/gorilla/websocket) - WebSocket communication
- [Fyne v2](https://fyne.io/) - GUI
- `net/http` - HTTP server and routing (`http.ServeMux` with route patterns)

## Architecture

```
internal/
  server/      - rooms, room registry, server-side client
  handler/     - network client layer: connection, sending, connection pool
ui/
  chatlist/    - chat list, creating and selecting a room
  rightpanel/  - chat history for the selected room, input field
  username/    - username input dialogue on launch
main.go        - main function: starts the server in a goroutine, launches the GUI
```

Each WebSocket connection is associated with a specific room (`/room/{name}`). The server determines the author of the message based on the connection itself, rather than on data sent by the client - the client only sends text.

## Running the application

```bash
go run .
```

The application will open a window, prompt for a username and start a server on port `:8080`. To test messaging between multiple users, run another instance in a separate terminal:

```bash
go run .
```

The second process will not be able to use port `:8080` (it is already in use by the first one) — this is to be expected; it will still connect as a client to the server that is already running.

## Known limitations that will be fixed in future (or perhaps not)

- The room password is collected during chat creation but is not used
- There is no search or filtering function for the list of chats
- The message history is not saved between restarts (it is stored in memory only)
- Open connections are not closed until the process terminates
- It is not possible to remove rooms from the list to which the user is connected
