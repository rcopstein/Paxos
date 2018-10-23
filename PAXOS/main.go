package PAXOS

import (
	"../Members"
	"fmt"
)
import "../Omega"

type Status int
const (
	Idle Status = 0
	Trying Status = 1
	Polling Status = 2
)

type PAXOS_Module struct {

	// SubModules
	messages PAXOS_Messages_Module
	leaderElection Omega.Omega_Module

	// Variables
	leader *Members.Member

	// Persistent Variables
	LastTried *int
	Outcome *int
	PrevBal *int
	PrevDec *int
	NextBal *int

	// Non-Persistent Variables
	CurrentStatus Status
	PrevVotes []int
	Quorum []Members.Member
	Voters []Members.Member
	Decree int

}

func Init() PAXOS_Module {

	var mod = PAXOS_Module{}

	mod.leader = nil
	mod.messages = Init_Messages()
	mod.leaderElection = Omega.Init()

	go mod.Start()
	return mod;

}
func (self PAXOS_Module) Start() {
	for {
		select {
			case y := <- self.messages.Ind:
				if y.messageType == BeginBallot { self.recvBeginBallot(y.from, *y.ballot, *y.outcome) }
				if y.messageType == LastVote { self.recvLastVote(y.from, *y.ballot, *y.outcome) }
				if y.messageType == NextBallot { self.recvNextBallot(y.from, *y.ballot) }
				if y.messageType == Success { self.recvSuccess(y.from, *y.outcome) }
				if y.messageType == Voted { self.recvVoted(y.from, *y.ballot) }
				break;
			case y := <- self.leaderElection.Ind:
				self.changeLeader(y.Member)
				break;
		}
	}
}

//
func (self PAXOS_Module) changeLeader(member *Members.Member) {

	self.leader = member;
	fmt.Println("The leader is " + self.leader.Name)

	if (member.Name == Members.GetSelf().Name) { fmt.Println("I'm the leader!") }

}

// Events
func (self PAXOS_Module) recvBeginBallot(from Members.Member, ballot int, outcome int) {

}

func (self PAXOS_Module) recvLastVote(from Members.Member, ballot int, outcome int) {

}

func (self PAXOS_Module) recvNextBallot(from Members.Member, ballot int) {

}


func (self PAXOS_Module) recvSuccess(from Members.Member, outcome int) {

}


func (self PAXOS_Module) recvVoted(from Members.Member, ballot int) {

}