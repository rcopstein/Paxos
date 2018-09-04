package PP2PLink

import "fmt"
import "net"


type PP2PLink_Message struct {
	Address string
	Message string
}

type PP2PLink struct {
	Ind chan PP2PLink_Message
	Req chan PP2PLink_Message
	Run bool
}

func (module PP2PLink) Init(address string, port string) {

	fmt.Println("Init PP2PLink!")
	if (module.Run) {
		return
	}

	module.Run = true;
	module.Start(address, port)
}

func (module PP2PLink) Start(address string, port string) {

	go func() {

		var buf = make([]byte, 1024)
		listen, _ := net.Listen("tcp4", address + ":" + port)

		for {

			conn, err := listen.Accept()
			if err != nil { continue }
			len, _ := conn.Read(buf)
			conn.Close()
	
			content := make([]byte, len)
			copy(content, buf)

			msg := PP2PLink_Message {
				Address: conn.RemoteAddr().String(),
				Message: string(content) }

			module.Ind <- msg

		}
	}()

	go func() {
		for {
			message := <- module.Req
			fmt.Println("Sent: " + message.Message)
			module.Send(message)
		}
	}()

}

func (module PP2PLink) Send(message PP2PLink_Message) {

	conn, err := net.Dial("tcp", message.Address)
	if err != nil { fmt.Println(err); return }
	fmt.Fprintf(conn, message.Message)
	conn.Close()

	fmt.Println("Sent: " + message.Message)

}
