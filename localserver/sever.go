package localserver

import (
	"ScrambledEggwithTomato/VMProtect"
	"ScrambledEggwithTomato/clientlauncher"
	"ScrambledEggwithTomato/global"
	"ScrambledEggwithTomato/modloader"
	"ScrambledEggwithTomato/mylogger"
	"ScrambledEggwithTomato/proxy"
	"ScrambledEggwithTomato/utils"
	"encoding/json"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

func GetData() []byte {

	data := &global.Data{
		Name:      global.CurrentUserInfo.DATA.Name,
		FuncName:  global.CurrentUserInfo.DATA.FuncName,
		ToString:  global.CurrentUserInfo.DATA.ToString,
		ClassName: global.CurrentUserInfo.DATA.ClassName,
		Time:      strconv.FormatInt(time.Now().Unix(), 10),
	}
	marshal, err := json.Marshal(data)
	if err != nil {
		return []byte("")
	}
	return marshal
}
func process(conn *net.Conn) {

	defer (*conn).Close()
	readTimeout := 5 * time.Second
	(*conn).SetReadDeadline(time.Now().Add(readTimeout))
	for {
		bytesRead, err := utils.ReadN(*conn, utils.PacketHerderLen_32)
		if err != nil {
			return
		}
		if strings.Compare(string(bytesRead), "exit") == 0 {
			log.Println("Connection sent exit")
			return
		}

		if string(bytesRead) == "GetData|CL" {
			if global.EnabledCL {
				err := utils.WriteN(*conn, GetData(), utils.PacketHerderLen_32)
				if err != nil {
					return
				}
				mylogger.Log("CL加载成功[2]...")
				return
			}
			utils.WriteN(*conn, []byte("Baier#1337"), utils.PacketHerderLen_32)
		}

		if string(bytesRead) == "GetData|Mod" {
			if global.EnabledMod {
				if global.CurrentUserInfo.DATA.Name == "" || global.CurrentUserInfo.USERNAME == "" {
					return
				}
				err := utils.WriteN(*conn, []byte(VMProtect.GoString(VMProtect.DecryptStringA("E8 ?? ?? ?? ?? 90 48 8B 4D ?? FF 15 ?? ?? ?? ?? BA 01 00 00 00 48 8B 4D ?? E8 ?? ?? ?? ?? 90 48 8B 4D ?? FF 15 ?? ?? ?? ?? BA 01 00 00 00 48 8B 4D ?? E8 ?? ?? ?? ??\x00"))), utils.PacketHerderLen_32)
				if err != nil {
					return
				}
				return
			}
			utils.WriteN(*conn, []byte("Baier#1337"), utils.PacketHerderLen_32)
		}

		if string(bytesRead) == "GetData|ModSign" {
			if global.EnabledMod {
				err := utils.WriteN(*conn, []byte("ok"), utils.PacketHerderLen_32)
				if err != nil {
					return
				}
				mylogger.Log("开始注入mod...[1]")
				modloader.InjectModProcessor()
				return
			}
			utils.WriteN(*conn, []byte("Baier#1337"), utils.PacketHerderLen_32)
		}

		if string(bytesRead) == "GetData|CLSign" {
			if global.EnabledCL {
				err := utils.WriteN(*conn, []byte("ok"), utils.PacketHerderLen_32)
				if err != nil {
					return
				}
				mylogger.Log("CL加载成功[3]")
				serverData := <-clientlauncher.ServerDataChan
				data4proxy := make([]string, 4)
				data4proxy[0] = serverData.ServerIP
				data4proxy[1] = serverData.ServerPort
				data4proxy[2] = serverData.Username
				data4proxy[3] = "25565"
				go func() {
					err := proxy.EstablishServer(data4proxy)
					if err != nil {
						mylogger.Log("启动proxy时遇到不可预期的错误:" + err.Error())
					}
				}()
				return
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
