package global

type RetData struct {
	Name      string `json:"Name"`
	FuncName  string `json:"FuncName"`
	ClassName string `json:"ClassName"`
	ToString  string `json:"2String"`
}
type Data struct {
	Name      string `json:"Name"`
	FuncName  string `json:"FuncName"`
	ClassName string `json:"ClassName"`
	ToString  string `json:"2String"`
	Time      string `json:"Time"`
}

type UserInfo struct {
	USERNAME string  `json:"Username"`
	PASSWORD string  `json:"Password"`
	HWID     string  `json:"HWID"`
	TIME     string  `json:"Time"`
	RANK     int     `json:"Rank"`
	VERSION  string  `json:"Version"`
	DATA     RetData `json:"Data"`
}
