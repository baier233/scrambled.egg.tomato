package localserver

import (
	"ScrambledEggwithTomato/clientlauncher"
	"ScrambledEggwithTomato/modloader"
	"ScrambledEggwithTomato/mylogger"
	"ScrambledEggwithTomato/utils"
	"encoding/json"
	"log"
	"net"
	"strconv"
	"strings"
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
func process(conn *net.Conn) {
	defer (*conn).Close()
	for {
		bytesRead, err := utils.ReadN(*conn, utils.PacketHerderLen_32)
		if err != nil {
			continue
		}
		if strings.Compare(string(bytesRead), "exit") == 0 {
			log.Println("Connection sent exit")
			return
		}

		if string(bytesRead) == "GetData|CL" {
			if clientlauncher.EnabledCL {
				err := utils.WriteN(*conn, []byte(GetData()), utils.PacketHerderLen_32)
				if err != nil {
					continue
				}
				mylogger.Log("CL加载成功[2]...")
				continue
			}
			utils.WriteN(*conn, []byte("Baier#1337"), utils.PacketHerderLen_32)
		}
		if string(bytesRead) == "GetData|Mod" {
			if modloader.EnablleMod {
				err := utils.WriteN(*conn, []byte("E8 ?? ?? ?? ?? 90 48 8B 4D ?? FF 15 ?? ?? ?? ?? BA 01 00 00 00 48 8B 4D ?? E8 ?? ?? ?? ?? 90 48 8B 4D ?? FF 15 ?? ?? ?? ?? BA 01 00 00 00 48 8B 4D ?? E8 ?? ?? ?? ??"), utils.PacketHerderLen_32)
				if err != nil {
					continue
				}
				continue
			}
			utils.WriteN(*conn, []byte("Baier#1337"), utils.PacketHerderLen_32)
		}
	}

}

func BeginListen() {
	mylogger.Log("LocalServer监听中...")
	listen, err := net.Listen("tcp", ":14889")
	if err != nil {
		return
	}
	for {
		conn, err := listen.Accept()
		if err != nil {
			return
		}

		go process(&conn)

	}

}
