package PP2PLink

import "fmt"
import "net"

type PP2PLink_Req_Message struct {
	To string
	Message string
}

type PP2PLink_Ind_Message struct {
	From string
	Message string
}

type PP2PLink struct {
	Ind chan PP2PLink_Ind_Message
	Req chan PP2PLink_Req_Message
	Run bool
}

func (module PP2PLink) Init(address string) {

	fmt.Println("Init PP2PLink!")
	if (module.Run) {
		return
	}

	module.Run = true;
	module.Start(address)
}

func (module PP2PLink) Start(address string) {

	go func() {

		var buf = make([]byte, 1024)
		listen, _ := net.Listen("tcp4", address)

		for {

			conn, err := listen.Accept()
			if err != nil { continue }
			len, _ := conn.Read(buf)
			conn.Close()
	
			content := make([]byte, len)
			copy(content, buf)

			msg := PP2PLink_Ind_Message {
				From: conn.RemoteAddr().String(),
				Message: string(content) }

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

func (module PP2PLink) Send(message PP2PLink_Req_Message) {

	conn, err := net.Dial("tcp", message.To)
	if err != nil { fmt.Println(err); return }
	fmt.Fprintf(conn, message.Message)
	conn.Close()

}
