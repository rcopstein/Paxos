package PAXOS

import (
	"../Members"
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
	outcome int
	ballot int

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
	result.link.Init(Members.GetSelf().Address)

	go func() {
		for {
			var y = <- result.link.Ind
			result.Ind <- recvMessage(y.Message)
		}
	}()

	return &result;

}

func buildMessage(messageType MessageType, ballot int, outcome int) string {

	var message string;

	message += Members.GetSelf().Name // Copy Current Member's Name

	message += "/"
	message += string(messageType); // Copy the Message Type

	message += "/";
	message += string(ballot); // Copy the Ballot Number

	message += "/";
	message += string(outcome) // Copy the Outcome

	return message;

}
func recvMessage(message string) PAXOS_Message {

	var elements = strings.Split(message, "/")
	var result = PAXOS_Message{}

	var from = Members.Find(elements[0]);
	result.from = *from

	var i, _ = strconv.Atoi(elements[1])
	result.messageType = MessageType(i);

	var j, _ = strconv.Atoi(elements[2])
	result.ballot = j

	var k, _ = strconv.Atoi(elements[3])
	result.ballot = k

	return result

}

func (self PAXOS_Messages_Module) sendMessage_BeginBallot(member Members.Member, ballot int, outcome int) {

	var message = P2PLink.P2PLink_Req_Message{
		To: member.Address,
		Message: buildMessage(BeginBallot, ballot, outcome),
	}

	self.link.Req <- message

}

func (self PAXOS_Messages_Module) sendMessage_LastVote(member Members.Member, ballot int, outcome int) {

	var message = P2PLink.P2PLink_Req_Message{
		To: member.Address,
		Message: buildMessage(LastVote, ballot, outcome),
	}

	self.link.Req <- message

}

func (self PAXOS_Messages_Module) sendMessage_NextBallot(member Members.Member, ballot int) {

	var message = P2PLink.P2PLink_Req_Message{
		To: member.Address,
		Message: buildMessage(NextBallot, ballot, 0),
	}

	self.link.Req <- message

}

func (self PAXOS_Messages_Module) sendMessage_Success(member Members.Member, outcome int) {

	var message = P2PLink.P2PLink_Req_Message{
		To: member.Address,
		Message: buildMessage(BeginBallot, 0, outcome),
	}

	self.link.Req <- message

}

func (self PAXOS_Messages_Module) sendMessage_Voted(member Members.Member, ballot int) {

	var message = P2PLink.P2PLink_Req_Message{
		To: member.Address,
		Message: buildMessage(Voted, ballot, 0),
	}

	self.link.Req <- message

}