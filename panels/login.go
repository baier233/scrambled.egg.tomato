package panels

import (
	"ScrambledEggwithTomato/global"
	"ScrambledEggwithTomato/localserver"
	"ScrambledEggwithTomato/login"
	"ScrambledEggwithTomato/mylogger"
	"ScrambledEggwithTomato/utils"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fyne.io/fyne/v2/dialog"
	"io"
	"io/ioutil"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const credentialsFile = "credentials.json"
const key = "BaierOops#133777" // Change this to a secure key in a real application

// Credentials struct to store username and password
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

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
		err := saveCredentials(username, password)
		if err != nil {
			return err
		}
		return nil
	}

	return login.ErrorNonexistUser

}

func LoginScreen(_ fyne.Window) fyne.CanvasObject {
	usernameInput := widget.NewEntry()
	passwordInput := widget.NewPasswordEntry()
	form := widget.NewForm(widget.NewFormItem("用户名", usernameInput), widget.NewFormItem("密码", passwordInput))
	savedUsername, savedPassword, err := loadCredentials()
	if err == nil {
		usernameInput.SetText(savedUsername)
		passwordInput.SetText(savedPassword)
	}
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

func saveCredentials(username, password string) error {
	creds := Credentials{
		Username: username,
		Password: password,
	}

	data, err := json.MarshalIndent(creds, "", "    ")
	if err != nil {
		return err
	}

	// Encrypt the data before saving
	encryptedData, err := encrypt(data, []byte(key))
	if err != nil {
		return err
	}

	return ioutil.WriteFile(credentialsFile, encryptedData, 0600)
}

func loadCredentials() (string, string, error) {
	data, err := ioutil.ReadFile(credentialsFile)
	if err != nil {
		return "", "", err
	}

	// Decrypt the data before loading
	decryptedData, err := decrypt(data, []byte(key))
	if err != nil {
		return "", "", err
	}

	var creds Credentials
	err = json.Unmarshal(decryptedData, &creds)
	if err != nil {
		return "", "", err
	}

	return creds.Username, creds.Password, nil
}

func encrypt(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)

	return ciphertext, nil
}

func decrypt(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(data) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(data, data)

	return data, nil
}
