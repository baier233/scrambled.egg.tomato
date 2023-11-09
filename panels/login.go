package panels

import (
	"ScrambledEggwithTomato/login"
	"ScrambledEggwithTomato/utils"
	"fmt"
	"fyne.io/fyne/v2/dialog"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var button *widget.Button

func handleLogin(username, password string) error {

	if len(username) == 0 || len(password) == 0 {
		return login.ErrorIllegalInputData
	}

	if strings.ToLower(username) == "compass" {
		return login.ErrorNonexistUser
	}

	user, err := login.NewUser([]string{username, password, utils.GetHWID()}, login.TypeLogin)
	if err != nil {
		return err
	}

	err = user.ProcessLogin()
	if err != nil {
		return err
	}

	if user.Mark {
		currentUser = NewCurrentUser(user)
		return nil
	}

	return login.ErrorNonexistUser

}

func LoginScreen(_ fyne.Window) fyne.CanvasObject {
	usernameInput := widget.NewEntry()
	passwordInput := widget.NewPasswordEntry()
	form := widget.NewForm(widget.NewFormItem("Username", usernameInput), widget.NewFormItem("Password", passwordInput))

	Line.StrokeWidth = 5

	button = widget.NewButton("Login!", func() {
		if button.Text != "Login!" {
			return
		}

		button.SetText("Logging in...")

		go func() {

			err := handleLogin(usernameInput.Text, passwordInput.Text)
			if err != nil {
				dialog.ShowError(err, Window)
				button.SetText("Login!")
				return
			}

			if currentUser.IsLoginIn {
				CheckWPFAndInitPanels()
				Init()
			}
			if currentUser.init && currentUser.user.Mark {
				fmt.Println(currentUser.user.RetData[0])
			}
			button.SetText("Login!")
		}()

	})

	mainPanel := container.NewVBox(
		label,
		form,
		widget.NewSeparator(),
		button,
		Line,
	)
	return mainPanel
}
