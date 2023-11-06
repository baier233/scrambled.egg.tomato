package proxy

import (
	"ScrambledEggwithTomato/global"
	"ScrambledEggwithTomato/mylogger"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/fatih/color"

	mcnet "github.com/Tnze/go-mc/net"
	"github.com/Tnze/go-mc/net/packet"
)

const (
	HANDSHAKE_PROTO_STATUS = 1
	HANDSHAKE_PROTO_LOGIN  = 2
	debug                  = false
)

type MinecraftProxyServer struct {
	Listen string
	Remote string
	MOTD   string

	running bool
	server  *mcnet.Listener

	HandleEncryption func(serverId string) error
	HandleLogin      func(packet *PacketLoginStart)

	// Middleware : 用于修改/解析包的中间件
	Middleware []func(packet *packet.Packet, clientConn *mcnet.Conn, serverConn *mcnet.Conn) bool
}

func Test(pk *packet.Packet, clientConn *mcnet.Conn, serverConn *mcnet.Conn) bool {

	var id packet.VarInt
	pk.Scan(&id)
	if /*id == 0x18 || id == 0x46 || id == 0x22 || */ id == 0x9 /* || id == 0x17 || pk.ID == 0x04*/ {
		var (
			channel packet.Identifier
			data    packet.ByteArray
		)
		_ = pk.Scan(
			&channel,
			&data,
		)
		yellow := color.New(color.FgYellow).SprintFunc()
		red := color.New(color.FgRed).SprintFunc()
		blue := color.New(color.FgBlue).SprintFunc()
		info := color.New(color.FgMagenta, color.Bold).SprintFunc()

		if len(string(channel)) != 0 && debug {
			fmt.Printf("%s : %s %s %s: %s %s %s:0x%x\n",
				yellow("Message"), data,
				info("->"),
				red("Channel"), channel,
				info("->"),
				blue("ID"), id)
			/*return true*/
		}

		if strings.Contains(string(channel), "MC|Brand") && strings.Contains(string(data), "vanilla") {
			fmt.Printf("%s", red("Detected"))
			fmt.Println(red(pk.Data))
			/*p := packet.Marshal(
				id,
				channel,
				packet.ByteArray("fml,forge"),
			)
			*pk = p*/
			numbers := []int{9, 8, 77, 67, 124, 66, 114, 97, 110, 100, 9, 102, 109, 108, 44, 102, 111, 114, 103, 101}
			bytes := make([]byte, len(numbers))
			for i, num := range numbers {
				bytes[i] = byte(num)
			}
			(*pk).Data = bytes
			return true
		}
		/*ori
		if strings.Contains(string(channel), "MC|Brand") && strings.Contains(string(data), "vanilla") {
			fmt.Printf("%s", red("Detected"))
			p := packet.Marshal(
				id,
				channel,
				packet.ByteArray("fml,forge"),
			)
			*pk = p
			return true
		}
		*/
	}

	return true
}
func (s *MinecraftProxyServer) StartServer() error {
	var err error
	s.server, err = mcnet.ListenMC(s.Listen)
	if err != nil {
		return err
	}
	s.running = true
	defer func() {
		s.CloseServer()
		mylogger.Log("Proxy已结束")
	}()
	for s.running && EnabledProxy {
		conn, err := s.server.Accept()
		if err != nil {
			continue
		}
		go func() {
			s.handleConnection(&conn)
		}()
	}

	return nil
}

func (s *MinecraftProxyServer) CloseServer() {
	if s.running {
		s.server.Close()
	}
}

func (s *MinecraftProxyServer) handleConnection(conn *mcnet.Conn) error {
	defer conn.Close()

	handshake, err := ReadHandshake(conn)
	if err != nil {
		return err
	}
	if handshake.NextState == HANDSHAKE_PROTO_LOGIN {
		err = s.forwardConnection(conn, *handshake)
		return err
	} else if handshake.NextState == HANDSHAKE_PROTO_STATUS {
		err = s.handlePing(conn, *handshake)
		return err
	}

	return nil
}

// forward connection to real server
func (s *MinecraftProxyServer) forwardConnection(conn *mcnet.Conn, handshake PacketHandshake) error {
	remoteConn, err := mcnet.DialMC(s.Remote)
	if err != nil {
		return err
	}
	defer remoteConn.Close()

	// modify & send handshake packet
	if strings.Contains(handshake.ServerAddress, "\x00FML\x00") {
		handshake.ServerAddress = strings.SplitN(s.Remote, ":", 2)[0] + "\u0000FML\u0000"
	} else {
		handshake.ServerAddress = strings.SplitN(s.Remote, ":", 2)[0]
	}
	port := 25565
	{
		slice := strings.SplitN(s.Remote, ":", 2)
		if len(slice) > 1 {
			port, err = strconv.Atoi(slice[1])
			if err != nil {
				port = 25565
			}
		}
	}
	handshake.ServerPort = uint16(port)
	WriteHandshake(remoteConn, handshake)

	// read username
	loginStart, err := ReadLoginStart(conn)
	if err != nil {
		return err
	}
	if s.HandleLogin != nil {
		s.HandleLogin(loginStart)
	}
	mylogger.Log("处理进服:" + loginStart.Name)

	WriteLoginStart(remoteConn, *loginStart)

	if s.HandleEncryption != nil {
		err = s.handleEncryption(conn, remoteConn)
		if err != nil {
			return err
		}
	}

	mylogger.Log("Forwarding packets")

	// forward connection
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		forwardPacket(conn, remoteConn, s.Middleware)
		wg.Done()
	}()
	go func() {
		forwardPacket(remoteConn, conn, s.Middleware)
		wg.Done()
	}()
	wg.Wait()

	return nil
}

func forwardPacket(srcConn *mcnet.Conn, dstConn *mcnet.Conn, middleware []func(packet *packet.Packet, clientConn *mcnet.Conn, serverConn *mcnet.Conn) bool) {
	for global.CurrentServer != nil {
		var p packet.Packet
		err := srcConn.ReadPacket(&p)
		if err != nil {
			break
		}

		for _, mw := range middleware {
			if !mw(&p, srcConn, dstConn) {
				return
			}
		}

		err = dstConn.WritePacket(p)
		if err != nil {
			break
		}
	}
}

func (s *MinecraftProxyServer) handleEncryption(conn *mcnet.Conn, remoteConn *mcnet.Conn) error {
	var p packet.Packet

	err := remoteConn.ReadPacket(&p)
	if err != nil {
		return err
	}

	if p.ID != 0x01 { // not an encryption request packet
		conn.WritePacket(p)
		return nil
	}

	pk, err := ReadEncryptionRequest(p)
	if err != nil {
		return err
	}

	key, encoStream, decoStream := newSymmetricEncryption()
	realServerId := authDigest(pk.ServerID, key, pk.PublicKey)
	err = s.HandleEncryption(realServerId)
	if err != nil {
		return err
	}

	p, err = genEncryptionKeyResponse(key, pk.PublicKey, pk.VerifyToken)
	if err != nil {
		return fmt.Errorf("gen encryption key response fail: %v", err)
	}

	err = remoteConn.WritePacket(p)
	if err != nil {
		return err
	}

	remoteConn.SetCipher(encoStream, decoStream)

	return nil
}

// handle ping request
func (s *MinecraftProxyServer) handlePing(conn *mcnet.Conn, handshake PacketHandshake) error {
	for {
		var p packet.Packet
		err := conn.ReadPacket(&p)
		if err != nil {
			return err
		}

		switch p.ID {
		case 0x00: // status request
			resp := StatusResponse{}

			resp.Version.Name = "Homo"
			resp.Version.Protocol = int(handshake.ProtocolVersion)
			resp.Players.Max = 666666
			resp.Players.Online = 114514
			resp.Description = s.MOTD
			if resp.Description == "" {
				resp.Description = "HomoEase_ProxyServer"
			}

			bytes, err := json.Marshal(resp)
			if err != nil {
				return nil
			}

			err = WriteStatusResponse(conn, PacketStatusResponse{
				Response: string(bytes),
			})
			if err != nil {
				return err
			}

		case 0x01: // ping
			var payload packet.Long
			err := p.Scan(&payload)
			if err != nil {
				return err
			}

			err = conn.WritePacket(packet.Marshal(
				0x01,
				packet.Long(payload)),
			)
			if err != nil {
				return err
			}
		}
	}

}
