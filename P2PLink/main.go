package P2PLink

import (
	"fmt"
	"strconv"
	"strings"
)
import "net"

type P2PLink_Req_Message struct {
	To string
	Message string
}

type P2PLink_Ind_Message struct {
	From string
	Message string
}

type P2PLink struct {
	Ind chan P2PLink_Ind_Message
	Req chan P2PLink_Req_Message
}

func Init(address string) *P2PLink {

	var module P2PLink

	module.Ind = make(chan P2PLink_Ind_Message)
	module.Req = make(chan P2PLink_Req_Message)
	Start(&module, address)

	return &module
}

func Start(module *P2PLink, address string) {

	go func() {

		var port, _ = strconv.Atoi(strings.Split(address, ":")[1])

		var addr = net.UDPAddr{
			Port: port,
			IP: net.ParseIP(strings.Split(address, ":")[1]),
		}

		var buf = make([]byte, 1024)
		listen, _ := net.ListenUDP("udp", &addr)

		for {

			var length, _, err = listen.ReadFromUDP(buf)
			if err != nil { continue }
	
			content := make([]byte, length)
			copy(content, buf)

			msg := P2PLink_Ind_Message { Message: string(content) }
			module.Ind <- msg

		}
	}()

	go func() {
		for {
			message := <- module.Req
			module.Send(message)
		}
	}()

}

func (module P2PLink) Send(message P2PLink_Req_Message) {

	conn, err := net.Dial("udp", message.To)
	if err != nil { fmt.Println(err); return }
	fmt.Fprintf(conn, message.Message)
	conn.Close()

}
