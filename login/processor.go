package login

import (
	"ScrambledEggwithTomato/VMProtect"
	"ScrambledEggwithTomato/global"
	"ScrambledEggwithTomato/mylogger"
	"ScrambledEggwithTomato/utils"
	"encoding/base64"
	"encoding/json"
	"net"
	"strconv"
	"strings"
	"time"
)

func packLogin(username, password, hwid string) []byte {
	VMProtect.BeginUltra("PackLoginMark\x00")
	timestr := strconv.FormatInt(time.Now().Unix(), 10)

	preData := username + "|" + password + "|" + hwid
	preDataHex := base64.StdEncoding.EncodeToString([]byte(preData))
	data2send := MsgHeader + VerifyType + preDataHex + "|" + timestr
	VMProtect.End()
	return []byte(data2send)
}
func packRegister(username, password, hwid, activationCode string) []byte {
	VMProtect.BeginUltra("PackRegisterMark\x00")
	timestr := strconv.FormatInt(time.Now().Unix(), 10)

	preData := username + "|" + password + "|" + hwid + "|" + activationCode
	preDataHex := base64.StdEncoding.EncodeToString([]byte(preData))
	data2send := MsgHeader + RegisterType + preDataHex + "|" + timestr

	VMProtect.End()
	return []byte(data2send)
}
func (user *User) _processLogin() error {
	conn, err := net.Dial("tcp", "111.180.205.168:9999")
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
	mylogger.Log("Stage 0")
	bytesRead, err := utils.ReadN(conn, utils.PacketHerderLen_32)
	if err != nil {
		return err
	}
	mylogger.Log("Stage 1")
	if string(bytesRead) == VMProtect.GoString(VMProtect.DecryptStringA("VerifyCode:-2\x00")) {
		return ErrorNonexistUser
	}

	if string(bytesRead) == VMProtect.GoString(VMProtect.DecryptStringA("VerifyCode:-3\x00")) {
		return ErrorIncorrectHWID
	}

	var userinfo global.UserInfo

	err = json.Unmarshal(bytesRead, &userinfo)

	if err != nil {
		return ErrorIllegalReturnData
	}
	global.CurrentUserInfo = &userinfo
	user.Mark = true
	user.RetData[0] = userinfo.USERNAME
	user.RetData[1] = userinfo.VERSION

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
	conn, err := net.Dial("tcp", "111.180.205.168:9999")
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
			return ErrorUserRegistered
		}

		if string(bytesRead) == "VerifyCode:-3" {
			return ErrorIncorrectActivationCode
		}
		return ErrorInternalImpossibleInputData
	}
	var userinfo global.UserInfo

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
