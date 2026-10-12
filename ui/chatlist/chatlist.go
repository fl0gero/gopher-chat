package chatlist

import (
	"errors"
	"fmt"
	"gopher-chat/internal/handler"
	"gopher-chat/internal/server"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// CreateChatDialogue builds the left-hand panel of the application: an input field at the top,
// a list of created chats in the centre, and an "Add chat" button at the bottom, which
// opens the room creation form.
func CreateChatDialogue(myWindow fyne.Window,
	pool *handler.ConnectionPool,
	onRoomSelected func(string),
	userName string,
	onMessageReceived func(string, server.Message)) *fyne.Container {

	var chatNames []string

	input := widget.NewEntry()
	input.SetPlaceHolder("Enter...")

	listOfChats := widget.NewList(func() int {
		return len(chatNames)
	}, func() fyne.CanvasObject {
		return widget.NewLabel("template")
	}, func(id widget.ListItemID, object fyne.CanvasObject) {
		object.(*widget.Label).SetText(chatNames[id])
	})

	// OnSelected is triggered when a user clicks on a chat in the list:
	listOfChats.OnSelected = func(id widget.ListItemID) {
		name := chatNames[id]
		_, err := pool.GetOrConnect(name, userName, func(message server.Message) {
			onMessageReceived(name, message)
		})
		if err != nil {
			fmt.Println("Error connecting:", err)
			return
		}

		onRoomSelected(name)
	}

	//Popped up window with all logic
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Name")
	// The validator is called by Fyne when the field is checked: an empty name is treated as an error,
	// and the form prevents the ‘Create’ button from being clicked.
	nameEntry.Validator = func(s string) error {
		if s == "" {
			return errors.New("empty name")
		}
		return nil
	}
	passwordEntry := widget.NewEntry()
	passwordEntry.SetPlaceHolder("Password")

	items := []*widget.FormItem{
		widget.NewFormItem("Name", nameEntry),
		widget.NewFormItem("Password (optional)", passwordEntry),
	}

	// Chat creation form
	form := dialog.NewForm("Creating New Chat", "Create", "Cancel", items, func(b bool) {
		if !b {
			return
		}

		roomName := strings.TrimSpace(nameEntry.Text)

		nameEntry.SetText("")
		passwordEntry.SetText("")

		if roomName == "" {
			dialog.ShowError(errors.New("the name of the room can't be empty"), myWindow)
			return
		}

		for _, existingName := range chatNames {
			if strings.EqualFold(existingName, roomName) {
				dialog.ShowError(errors.New("you have already connected to this room"), myWindow)
				return
			}
		}

		_, err := pool.GetOrConnect(roomName, userName, func(message server.Message) {
			onMessageReceived(roomName, message)
		})
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}

		chatNames = append(chatNames, roomName)
		listOfChats.Refresh()
	}, myWindow)
	form.Resize(fyne.NewSize(400, 200))

	button := widget.NewButton("Add chat", func() { form.Show() })
	listOfChats.Resize(fyne.NewSize(200, 600))

	chatContainer := container.NewBorder(
		input,       // top
		button,      // bottom
		nil,         // left
		nil,         // right
		listOfChats, // center
	)
	return chatContainer

}
