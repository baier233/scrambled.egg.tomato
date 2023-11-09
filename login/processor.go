package login

import (
	"ScrambledEggwithTomato/VMProtect"
	"ScrambledEggwithTomato/utils"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

var (
	ErrorInternalImpossibleInputData = errors.New("Impossible input data internally")
	ErrorIllegalInputData            = errors.New("Illegal input data")
	ErrorNonexistUser                = errors.New("User does not exist")
	ErrorIllegalReturnData           = errors.New("Illegle return data")
	ErrorIncorrectActivationCode     = errors.New("Incorrect & Expired Activation Code")
	ErrorInvalidTime                 = errors.New("Invalid Time")
	ErrorUserRegitered               = errors.New("This username is registered")
)

type UserInfo struct {
	USERNAME string `json:"Username"`
	PASSWORD string `json:"Password"`
	HWID     string `json:"HWID"`
	TIME     string `json:"Time"`
	RANK     int    `json:"Rank"`
}

type User struct {
	Data    []string
	RetData []string
	Mark    bool
}

const (
	VerifyType     = "Verify"
	VerifyLength   = len(VerifyType)
	RegisterType   = "Register"
	RegisterLength = len(RegisterType)
)
const (
	MsgHeader = "Baier#1337"

	HeaderLength = len(MsgHeader)
)
const (
	TypeRegister = iota
	TypeLogin
)

// NewUser creates a User.
func NewUser(data []string, _type int) (*User, error) {

	if _type == TypeLogin {
		if len(data) != 3 {
			return nil, ErrorInternalImpossibleInputData
		}
		user := &User{Data: data}
		user.RetData = make([]string, 3)

		return user, nil
	}

	if _type == TypeRegister {
		if len(data) != 4 {
			return nil, ErrorInternalImpossibleInputData
		}
		user := &User{Data: data}
		user.RetData = make([]string, 3)

		return user, nil
	}

	return nil, ErrorInternalImpossibleInputData
}
func packLogin(username, password, hwid string) []byte {
	VMProtect.BeginUltra("packLogin\x00")
	timestr := strconv.FormatInt(time.Now().Unix(), 10)

	preData := username + "|" + password + "|" + hwid
	preDataHex := base64.StdEncoding.EncodeToString([]byte(preData))
	fmt.Println(preDataHex)
	data2send := MsgHeader + VerifyType + preDataHex + "|" + timestr
	VMProtect.End()
	return []byte(data2send)
}
func packRegister(username, password, hwid, activationCode string) []byte {
	VMProtect.BeginUltra("packRegister\x00")
	timestr := strconv.FormatInt(time.Now().Unix(), 10)

	preData := username + "|" + password + "|" + hwid + "|" + activationCode
	preDataHex := base64.StdEncoding.EncodeToString([]byte(preData))
	fmt.Println(preDataHex)
	data2send := MsgHeader + RegisterType + preDataHex + "|" + timestr

	VMProtect.End()
	return []byte(data2send)
}
func (user *User) _processLogin() error {
	conn, err := net.Dial("tcp", "localhost:9999")
	if err != nil {
		return err
	}

	VMProtect.BeginUltra("processLogin\x00")
	username := user.Data[0]
	password := user.Data[1]
	hwid := user.Data[2]
	VMProtect.End()
	err = utils.WriteN(conn, packLogin(username, password, hwid), utils.PacketHerderLen_32)
	if err != nil {

		return err
	}
	bytesRead, err := utils.ReadN(conn, utils.PacketHerderLen_32)
	if err != nil {
		return err
	}

	if string(bytesRead) == "VerifyCode:-1" {
		return ErrorNonexistUser
	}

	var userinfo UserInfo

	err = json.Unmarshal(bytesRead, &userinfo)

	if err != nil {
		return ErrorIllegalReturnData
	}

	user.Mark = true
	user.RetData[0] = userinfo.USERNAME

	conn.Close()
	VMProtect.End()
	return nil
}

func (user *User) ProcessLogin() error {

	user.Mark = false
	if len(user.Data) != 3 {
		return ErrorInternalImpossibleInputData
	}

	return user._processLogin()
}

func (user *User) _processRegister() error {
	conn, err := net.Dial("tcp", "localhost:9999")
	if err != nil {
		return err
	}
	VMProtect.BeginUltra("processRegister\x00")
	username := user.Data[0]
	password := user.Data[1]
	hwid := user.Data[2]
	activationCode := user.Data[3]
	err = utils.WriteN(conn, packRegister(username, password, hwid, activationCode), utils.PacketHerderLen_32)
	if err != nil {

		return err
	}
	bytesRead, err := utils.ReadN(conn, utils.PacketHerderLen_32)
	if err != nil {
		return err
	}
	if strings.Contains(string(bytesRead), "VerifyCode") {
		if string(bytesRead) == VMProtect.GoString(VMProtect.DecryptStringA("VerifyCode:-1\x00")) {
			fmt.Println("1")
			return ErrorUserRegitered
		}

		if string(bytesRead) == "VerifyCode:-3" {

			return ErrorIncorrectActivationCode
		}
		return ErrorInternalImpossibleInputData
	}
	var userinfo UserInfo

	err = json.Unmarshal(bytesRead, &userinfo)

	if err != nil {
		return ErrorIllegalReturnData
	}
	user.Mark = true
	user.RetData[0] = userinfo.USERNAME
	VMProtect.End()
	return nil
}

func (user *User) ProcessRegister() error {
	if len(user.Data) != 4 {
		user.Mark = false
		return ErrorInternalImpossibleInputData
	}

	return user._processRegister()
}
