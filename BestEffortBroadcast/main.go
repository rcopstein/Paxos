package main

import "os"
import "fmt"
import PP2PLink "../PP2PLink"

type BestEffortBroadcast_Message struct {
	From string
	Data string
}

type BestEffortBroadcast_Module struct {
	Req chan BestEffortBroadcast_Message
	Ind chan BestEffortBroadcast_Message
	Addresses []string
	Ports []string

	Pp2plink PP2PLink.PP2PLink
}

func (module BestEffortBroadcast_Module) Init(addresses []string, ports []string) {

	fmt.Println("Init BEB!")
	module.Ports = ports
	module.Addresses = addresses
	module.Pp2plink = PP2PLink.PP2PLink{
		Req: make(chan PP2PLink.PP2PLink_Message),
		Ind: make(chan PP2PLink.PP2PLink_Message) }

	module.Pp2plink.Init(addresses[0], ports[0]);
	module.Start()

}

func (module BestEffortBroadcast_Module) Start() {

	go func () {

		for {
			select{
			case y:= <- module.Req:
				module.Broadcast(y);
			case y:= <- module.Pp2plink.Ind:
				module.Deliver(PP2PLink2BEB(y)); // Deserializar Isso
			}
		}

	}()

}

func (module BestEffortBroadcast_Module) Broadcast(message BestEffortBroadcast_Message) {

	fmt.Println("Broadcast: " + message.Data)
	for i := 0; i < len(module.Addresses); i++ {
		fmt.Println("Here!")
		message.From = module.Addresses[i] + ":" + module.Ports[i]
		module.Pp2plink.Req <- BEB2PP2PLink(message);
	}

}

func (module BestEffortBroadcast_Module) Deliver(message BestEffortBroadcast_Message) {

	fmt.Println("Received: " + message.Data + " from " + message.From)
	module.Ind <- message

}

func BEB2PP2PLink(message BestEffortBroadcast_Message) PP2PLink.PP2PLink_Message {

	return PP2PLink.PP2PLink_Message{
		Address: message.From,
		Message: message.Data }

}

func PP2PLink2BEB(message PP2PLink.PP2PLink_Message) BestEffortBroadcast_Message {

	return BestEffortBroadcast_Message{
		From: message.Address,
		Data: message.Message }

}


func main() {

	size := len(os.Args) / 2
	addresses := make([]string, 0, size)
	ports     := make([]string, 0, size)

	for i := 1; i < len(os.Args); i += 2 {
		addresses = append(addresses, os.Args[i])
		ports     = append(ports, os.Args[i + 1])
	}

	mod := BestEffortBroadcast_Module{
		Req: make(chan BestEffortBroadcast_Message),
		Ind: make(chan BestEffortBroadcast_Message) }
	mod.Init(addresses, ports)

	msg := BestEffortBroadcast_Message{
		From: addresses[0] + ":" + ports[0],
		Data: "BATATA!" }

	yy := make(chan string)
	mod.Req <- msg
	<- yy
}