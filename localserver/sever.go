package localserver

import (
	"ScrambledEggwithTomato/clientlauncher"
	"ScrambledEggwithTomato/mylogger"
	"ScrambledEggwithTomato/utils"
	"encoding/json"
	"net"
	"strconv"
	"time"
)

type Data struct {
	Name      string `json:"Name"`
	FuncName  string `json:"FuncName"`
	ClassName string `json:"ClassName"`
	ToString  string `json:"2String"`
	Time      string `json:"Time"`
}

func GetData() string {
	data := &Data{
		Name:      "netease",
		FuncName:  "O0O0O0O0OOOO0",
		ToString:  "toString",
		ClassName: "BigInteger",
		Time:      strconv.FormatInt(time.Now().Unix(), 10),
	}
	marshal, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(marshal)
}

func beginListen() {
	mylogger.Log("LocalServer监听中")
	listen, err := net.Listen("tcp", ":12441")
	if err != nil {
		return
	}
	for {
		conn, err := listen.Accept()
		if err != nil {
			return
		}
		bytesRead, err := utils.ReadN(conn, utils.PacketHerderLen_32)
		if err != nil {
			continue
		}
		if string(bytesRead) == "GetData|CL" {
			if clientlauncher.EnabledCL {
				err := utils.WriteN(conn, []byte(GetData()), utils.PacketHerderLen_32)
				if err != nil {
					continue
				}
				continue
			}
			utils.WriteN(conn, []byte("Baier#1337"), utils.PacketHerderLen_32)
		}
		if string(bytesRead) == "GetData|Mod" {
			if clientlauncher.EnabledCL {
				err := utils.WriteN(conn, []byte(GetData()), utils.PacketHerderLen_32)
				if err != nil {
					continue
				}
				continue
			}
			utils.WriteN(conn, []byte("Baier#1337"), utils.PacketHerderLen_32)
		}

	}

}
