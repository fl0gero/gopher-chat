package rightpanel

import (
	"fmt"
	"gopher-chat/internal/handler"
	"gopher-chat/internal/server"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// CreateRightPanel builds the right panel of the app: the message history
// for the selected room at the top, and a field for entering a new message at the bottom.
//
// Returns three values:
// - the container itself, which should be inserted into the window
// - onRoomSelected
// - a function for the chat list: called when the user selects a different room,
// and switches the displayed conversation

func CreateRightPanel(pool *handler.ConnectionPool, userName string) (*fyne.Container, func(roomName string), func(string, server.Message)) {
	var currentRoom string
	var mu sync.Mutex

	// Message history for all rooms (not just the current one)
	roomMessages := make(map[string][]server.Message)

	messages := widget.NewLabel("")

	renderCurrentRoom := func() {
		var text string
		for _, msg := range roomMessages[currentRoom] {
			text += msg.From + ": " + msg.Text + "\n"
		}
		messages.SetText(text)
	}

	onRoomSelected := func(roomName string) {
		mu.Lock()
		currentRoom = roomName
		mu.Unlock()
		renderCurrentRoom()
	}

	// onMessageReceived stores incoming messages and schedules a UI update
	// when the message belongs to the currently selected room.
	onMessageReceived := func(roomName string, msg server.Message) {
		mu.Lock()
		defer mu.Unlock()

		roomMessages[roomName] = append(roomMessages[roomName], msg)

		if roomName == currentRoom {
			fyne.Do(renderCurrentRoom)
		}
	}

	input := widget.NewEntry()
	input.SetPlaceHolder("Message...")
	input.OnSubmitted = func(text string) {
		conn, err := pool.GetOrConnect(currentRoom, userName, func(msg server.Message) {
			onMessageReceived(currentRoom, msg)
		})
		if err != nil {
			fmt.Println("Error connecting to room:", currentRoom)
			return
		}
		handler.SendMessage(conn, text)
		input.SetText("")
	}

	panelContainer := container.NewBorder(nil, input, nil, nil, messages)

	return panelContainer, onRoomSelected, onMessageReceived
}
