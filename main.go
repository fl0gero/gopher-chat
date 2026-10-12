package main

import (
	"fmt"
	"gopher-chat/internal/handler"
	"gopher-chat/internal/server"
	"gopher-chat/ui/chatlist"
	"gopher-chat/ui/rightpanel"
	"gopher-chat/ui/username"
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

func main() {
	pool := handler.NewConnectionPool()
	reg := server.NewRegistry()
	mux := http.NewServeMux()

	mux.Handle("/room/{name}", reg)
	go func() {
		err := http.ListenAndServe(":8080", mux)
		if err != nil {
			fmt.Println("Error:", err)
		}
	}()

	myApp := app.New()
	myWindow := myApp.NewWindow("Split Container")

	username.CreateUser(myWindow, func(userName string) {
		right, onRoomSelected, onMessageReceived := rightpanel.CreateRightPanel(pool, userName)
		chat := chatlist.CreateChatDialogue(myWindow, pool, onRoomSelected, userName, onMessageReceived)

		split := container.NewHSplit(chat, right)
		split.SetOffset(0.3)

		myWindow.SetContent(split)
	})

	myWindow.Resize(fyne.NewSize(760, 580))
	myWindow.ShowAndRun()
}
