package PAXOS

import (
	"../Members"
	"fmt"
	"strconv"
	"strings"
)
import "../P2PLink"

// Message Type Enum
type MessageType int
const (
	BeginBallot MessageType = 0
	NextBallot  MessageType = 1
	LastVote    MessageType = 2
	Success     MessageType = 3
	Voted       MessageType = 4
)

// Structs
type PAXOS_Message struct {

	messageType MessageType
	from Members.Member
	vals []int

}
type PAXOS_Messages_Module struct {

	Ind chan PAXOS_Message
	link *P2PLink.P2PLink

}

// Functions
func Init_Messages() *PAXOS_Messages_Module {

	var p2plink P2PLink.P2PLink

	var result = PAXOS_Messages_Module{
		Ind: make(chan PAXOS_Message),
		link: &p2plink,
	}
	result.link = P2PLink.Init(Members.GetSelf().Address)

	go func() {
		for {
			var y = <- result.link.Ind
			result.Ind <- recvMessage(y.Message)
		}
	}()

	return &result;

}

func buildMessage(messageType MessageType, vals ...int) string {

	var message string;

	message += Members.GetSelf().Name // Copy Current Member's Name

	message += "/"
	message += strconv.Itoa(int(messageType)); // Copy the Message Type

	for _, val := range vals {
		message += "/"
		message += strconv.Itoa(val) // Copy Every Value
	}

	return message;

}
func recvMessage(message string) PAXOS_Message {

	var elements = strings.Split(message, "/")
	var result = PAXOS_Message{}

	var from = Members.Find(elements[0]);
	result.from = *from

	var i, _ = strconv.Atoi(elements[1])
	result.messageType = MessageType(i);

	for j := 2; j < len(elements); j++ {
		var k, _ = strconv.Atoi(elements[j])
		result.vals = append(result.vals, k)
	}

	return result

}

func (self *PAXOS_Messages_Module) SendMessage_BeginBallot(member Members.Member, ballot int, outcome int) {

	var message = P2PLink.P2PLink_Req_Message{
		To: member.Address,
		Message: buildMessage(BeginBallot, ballot, outcome),
	}

	self.link.Req <- message
	fmt.Println(Members.GetSelf().Name, " sent BeginBallot to:", member.Name)

}

func (self *PAXOS_Messages_Module) SendMessage_LastVote(member Members.Member, lastVote int, prevVote int, prevDec int) {

	var message = P2PLink.P2PLink_Req_Message{
		To: member.Address,
		Message: buildMessage(LastVote, lastVote, prevVote, prevDec),
	}

	self.link.Req <- message
	fmt.Println(Members.GetSelf().Name, " sent LastVote to:", member.Name)

}

func (self *PAXOS_Messages_Module) SendMessage_NextBallot(member Members.Member, ballot int) {

	var message = P2PLink.P2PLink_Req_Message{
		To: member.Address,
		Message: buildMessage(NextBallot, ballot),
	}

	self.link.Req <- message
	fmt.Println(Members.GetSelf().Name, " sent NextBallot to:", member.Name)

}

func (self *PAXOS_Messages_Module) SendMessage_Success(member Members.Member, outcome int) {

	var message = P2PLink.P2PLink_Req_Message{
		To: member.Address,
		Message: buildMessage(Success, outcome),
	}

	self.link.Req <- message
	fmt.Println(Members.GetSelf().Name, " sent Success to:", member.Name)

}

func (self *PAXOS_Messages_Module) SendMessage_Voted(member Members.Member, ballot int) {

	var message = P2PLink.P2PLink_Req_Message{
		To: member.Address,
		Message: buildMessage(Voted, ballot),
	}

	self.link.Req <- message
	fmt.Println(Members.GetSelf().Name, " sent Voted to:", member.Name)

}