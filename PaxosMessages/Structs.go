package PaxosMessages

import "../Members"

type MessageType int
type Message struct {

	Values   []int
	Instance int
	Type     MessageType
	Member   *Members.Member

}

const BeginBallot MessageType = 0
const NextBallot  MessageType = 1
const LastVote    MessageType = 2
const Success     MessageType = 3
const Voted       MessageType = 4
const Leader      MessageType = 5