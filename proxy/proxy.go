package proxy

import (
	"ScrambledEggwithTomato/global"
	"ScrambledEggwithTomato/mylogger"
	mcnet "github.com/Tnze/go-mc/net"
	"github.com/Tnze/go-mc/net/packet"
	"net"
	"strings"
	"time"
	"unsafe"
)

func EstablishServer(data []string) error {
	if len(data) != 4 {
		return global.ErrorInternalIncorrectData
	}
	serverIp := data[0]
	serverPort := data[1]
	roleName := data[2]
	localPort := data[3]

	if serverIp == "" || serverPort == "" || roleName == "" || localPort == "" {
		return global.ErrorEmptyInputData
	}
	server := MinecraftProxyServer{
		Listen: "0.0.0.0:" + localPort,
		Remote: serverIp + ":" + serverPort,
		MOTD:   "§西红柿炒鸡蛋§w-§6§l代理服务\n§w目标服务器：" + serverIp + "  角色：" + roleName,
		HandleEncryption: func(serverId string) error {

			conn, err := net.Dial("tcp", "127.0.0.1:55996")
			if err != nil {
				mylogger.Log("无法连接到CL服务器 预期之外的错误　: " + err.Error())
				return err
			} else {

				conn.Write([]byte(serverId + "\u0000"))

				defer conn.Close()

				bytes := make([]byte, 1024)
				conn.Read(bytes)
			}
			return nil

		},
		HandleLogin: func(packet *PacketLoginStart) {
			packet.Name = roleName
		},
		Middleware: []func(packet *packet.Packet, clientConn *mcnet.Conn, serverConn *mcnet.Conn) bool{},
	}
	go handleLocalPing(server.MOTD, localPort)
	defer mylogger.Log("代理服务已结束")
	mylogger.Log("代理服务器已启动")
	global.CurrentServer = unsafe.Pointer(&server)
	err := server.StartServer()
	if err != nil {
		return err
	}
	return nil

}
func handleLocalPing(description string, port string) {
	for EnabledProxy {
		time.Sleep(1000)
		connudp, err := net.Dial("udp", "224.0.2.60:4445")
		if err != nil {
			time.Sleep(1)
		} else {
			connudp.Write([]byte("[MOTD]" + strings.Replace(strings.Replace(description, "\n", " ", -1), "目标", "", -1) + "[/MOTD][AD]" + port + "[/AD]"))
			connudp.Close()
		}
	}

}
