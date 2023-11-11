package clientlauncher

type ServerData struct {
	ServerIP   string
	ServerPort string
	Username   string
}

func NewServerData() *ServerData {
	return &ServerData{
		ServerIP:   "",
		ServerPort: "",
		Username:   "",
	}
}

var ServerDataChan = make(chan *ServerData)
