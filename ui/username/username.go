package username

import (
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// CreateUser displays a username dialog and passes the entered name
// to onUsernameSet when the user logs in.
func CreateUser(myWindow fyne.Window, onUsernameSet func(string)) {
	userEntry := widget.NewEntry()
	userEntry.Validator = func(s string) error {
		if s == "" {
			return errors.New("Username must not be empty")
		} else if len(s) < 4 {
			return errors.New("Username must be at least 4 characters")
		}
		return nil
	}

	content := []*widget.FormItem{widget.NewFormItem("Username", userEntry)}
	usernameForm := dialog.NewForm("Set a Username", "Log in", "Cancel", content, func(b bool) {
		if b {
			onUsernameSet(userEntry.Text)
		} else {
			myWindow.Close()
		}
	}, myWindow)

	usernameForm.Resize(fyne.NewSize(300, 150))
	usernameForm.Show()

}
