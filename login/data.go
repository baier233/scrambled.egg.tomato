package login

import "errors"

var (
	ErrorInternalImpossibleInputData = errors.New("不合法的内部数据")
	ErrorIllegalInputData            = errors.New("非法数据")
	ErrorNonexistUser                = errors.New("不存在的用户")
	ErrorIncorrectHWID               = errors.New("HWID错误")
	ErrorIllegalReturnData           = errors.New("不合法的返回数据")
	ErrorIncorrectActivationCode     = errors.New("不正确 & 超时 的激活码")

	ErrorUserRegistered = errors.New("整个用户名被注册了")
)

const (
	VerifyType     = "Verify"
	VerifyLength   = len(VerifyType)
	RegisterType   = "Register"
	RegisterLength = len(RegisterType)
)
const (
	MsgHeader = "BaierOops"

	HeaderLength = len(MsgHeader)
)
const (
	TypeRegister = iota
	TypeLogin
)

type CurrentUser struct {
	Init      bool
	User      *User
	IsLoginIn bool
}

var (
	MyCurrentUser = &CurrentUser{
		IsLoginIn: false,
		Init:      false,
	}
)
