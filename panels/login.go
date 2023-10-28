package panels

import (
	"errors"
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

var button *widget.Button

func handleLogin(username, password string) error {

	fmt.Println("Username:", username, " Password:", password)
	if len(username) == 0 || len(password) == 0 {
		return errors.New(emptyInputDataError)
	}

	if strings.ToLower(username) == "compass" {
		return errors.New(nonexistentUserError)
	}

	{
		//Galmoxy gal wikioos whopps eoose?
		Panels = map[string]Panel{
			"注入": {"注入", ModInjectPanel},
			"开端": {"开端", ClientLaunchPanel},
		}
		PanelIndex = map[string][]string{
			"": {"注入", "开端"}}

		Init()
	}

	return errors.New(nonexistentUserError)

}

func LoginScreen(_ fyne.Window) fyne.CanvasObject {
	usernameInput := widget.NewEntry()
	passwordInput := widget.NewPasswordEntry()
	form := widget.NewForm(widget.NewFormItem("用户名", usernameInput), widget.NewFormItem("密码", passwordInput))

	Line.StrokeWidth = 5

	button = widget.NewButton("登录!", func() {
		if button.Text != "登录!" {
			return
		}

		button.SetText("登录...")
		go func() {

			//time.Sleep(time.Second)
			err := handleLogin(usernameInput.Text, passwordInput.Text)
			if err != nil {
				dialog.ShowError(err, Window)
				button.SetText("登录!")
				return
			}

			button.SetText("登录!")
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
