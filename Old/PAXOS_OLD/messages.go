package PAXOS_OLD

import (
	"../../Members"
	"strconv"
	"strings"
)
import "../../P2PLink"

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
type Message struct {

	instance int
	messageType MessageType
	from Members.Member
	vals []int

}
type MessagesModule struct {

	Ind      chan Message
	link     *P2PLink.P2PLink

}

// Functions
func InitMessages() *MessagesModule {

	var result = MessagesModule{
		Ind: make(chan Message),
	}
	result.link = P2PLink.Init(Members.GetSelf().Address)

	go func() {
		for {
			var y = <- result.link.Ind
			result.Ind <- recvMessage(y.Message)
		}
	}()

	return &result

}

func buildMessage(instance int, messageType MessageType, vals ...int) string {

	var message string

	var inst = strconv.Itoa(instance)
	message += string(inst)
	message += "/"

	var name = strconv.Itoa(Members.GetSelf().Name)
	message += string(name) // Copy Current Member's Name

	message += "/"
	message += strconv.Itoa(int(messageType)) // Copy the Message Type

	for _, val := range vals {
		message += "/"
		message += strconv.Itoa(val) // Copy Every Value
	}

	return message

}
func recvMessage(message string) Message {

	var elements = strings.Split(message, "/")
	var result = Message{}

	l, _ := strconv.Atoi(elements[0])
	result.instance = l

	k, _ := strconv.Atoi(elements[1])
	var from = Members.Find(k)
	result.from = *from

	var i, _ = strconv.Atoi(elements[2])
	result.messageType = MessageType(i)

	for j := 3; j < len(elements); j++ {
		var k, _ = strconv.Atoi(elements[j])
		result.vals = append(result.vals, k)
	}

	return result

}

func (mod *MessagesModule) SendMessageBeginBallot(instance int, member Members.Member, ballot int, outcome int) {

	var message = P2PLink.P2PLink_Req_Message{
		To: member.Address,
		Message: buildMessage(instance, BeginBallot, ballot, outcome),
	}

	mod.link.Req <- message
	// fmt.Println(Members.GetSelf().Name, " sent BeginBallot to:", member.Name)

}

func (mod *MessagesModule) SendMessageLastVote(instance int, member Members.Member, lastVote int, prevVote int, prevDec int) {

	var message = P2PLink.P2PLink_Req_Message{
		To: member.Address,
		Message: buildMessage(instance, LastVote, lastVote, prevVote, prevDec),
	}

	mod.link.Req <- message
	// fmt.Println(Members.GetSelf().Name, " sent LastVote to:", member.Name)

}

func (mod *MessagesModule) SendMessageNextBallot(instance int, member Members.Member, ballot int) {

	var message = P2PLink.P2PLink_Req_Message{
		To: member.Address,
		Message: buildMessage(instance, NextBallot, ballot),
	}

	mod.link.Req <- message
	// fmt.Println(Members.GetSelf().Name, " sent NextBallot to:", member.Name)

}

func (mod *MessagesModule) SendMessageSuccess(instance int, member Members.Member, outcome int) {

	var message = P2PLink.P2PLink_Req_Message{
		To: member.Address,
		Message: buildMessage(instance, Success, outcome),
	}

	mod.link.Req <- message
	// fmt.Println(Members.GetSelf().Name, " sent Success to:", member.Name)

}

func (mod *MessagesModule) SendMessageVoted(instance int, member Members.Member, ballot int) {

	var message = P2PLink.P2PLink_Req_Message{
		To: member.Address,
		Message: buildMessage(instance, Voted, ballot),
	}

	mod.link.Req <- message
	// fmt.Println(Members.GetSelf().Name, " sent Voted to:", member.Name)

}