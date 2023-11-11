package panels

import (
	"ScrambledEggwithTomato/login"
	"ScrambledEggwithTomato/utils"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"strings"
)

func handleRegister(username, password, activationCode string) (*login.User, error) {

	if len(username) == 0 || len(password) == 0 || len(activationCode) == 0 {
		return nil, login.ErrorIllegalInputData
	}

	if strings.ToLower(username) == "compass" {
		return nil, login.ErrorIncorrectActivationCode
	}

	user, err := login.NewUser([]string{username, password, utils.GetHWID(), activationCode}, login.TypeRegister)
	if err != nil {
		return nil, err
	}

	err = user.ProcessRegister()
	if err != nil {
		return nil, err
	}

	return user, nil

}

func RegisterScreen(_ fyne.Window) fyne.CanvasObject {
	usernameInput := widget.NewEntry()
	passwordInput := widget.NewPasswordEntry()
	activationCodeInput := widget.NewEntry()
	form := widget.NewForm(widget.NewFormItem("用户名", usernameInput), widget.NewFormItem("密码", passwordInput), widget.NewFormItem("卡密", activationCodeInput))

	button = widget.NewButton("注册!", func() {
		if button.Text != "注册!" {
			return
		}
		button.SetText("注册中...")

		go func() {
			defer button.SetText("注册!")
			user, err := handleRegister(usernameInput.Text, passwordInput.Text, strings.TrimSpace(activationCodeInput.Text))
			if err != nil {
				dialog.ShowError(err, Window)
				return
			}
			if user != nil && user.Mark {
				dialog.ShowConfirm("注册", "注册成功", func(b bool) {
				}, Window)
				return
			}

			dialog.ShowConfirm("注册", login.ErrorIncorrectActivationCode.Error(), func(b bool) {}, Window)
		}()
	})

	Line.StrokeWidth = 5
	return container.NewVBox(

		label,
		form,
		widget.NewSeparator(),
		button,
		Line,
	)
}
