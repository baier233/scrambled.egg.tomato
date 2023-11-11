package panels

import (
	"ScrambledEggwithTomato/global"
	"ScrambledEggwithTomato/localserver"
	"ScrambledEggwithTomato/login"
	"ScrambledEggwithTomato/mylogger"
	"ScrambledEggwithTomato/utils"
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
		login.MyCurrentUser = NewCurrentUser(user)
		return nil
	}

	return login.ErrorNonexistUser

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

		button.SetText("登录中...")

		go func() {

			err := handleLogin(usernameInput.Text, passwordInput.Text)
			if err != nil {
				dialog.ShowError(err, Window)
				button.SetText("登录!")
				return
			}

			if login.MyCurrentUser.IsLoginIn {
				if login.MyCurrentUser.User.RetData[1] == global.ScrambledEggTomatoVersion {
					go localserver.BeginListen()
					CheckWPFAndInitPanels()
					Init()

					if login.MyCurrentUser.Init && login.MyCurrentUser.User.Mark {
						mylogger.Log("已经登录 用户名:", login.MyCurrentUser.User.RetData[0])
					}
				} else {
					mylogger.Log("当前版本 :" + global.ScrambledEggTomatoVersion + " 不符合最新版本 :" + login.MyCurrentUser.User.RetData[1])
				}
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
