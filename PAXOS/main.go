package PAXOS

import (
	"../Members"
	"fmt"
	"math/rand"
	"time"
)
import "../Omega"

// Status ENUM
type Status int
const (
	Idle Status = 0
	Trying Status = 1
	Polling Status = 2
)

// Structs
type PAXOS_Module struct {

	// SubModules
	messages *PAXOS_Messages_Module
	leaderElection Omega.Omega_Module

	// Variables
	leader *Members.Member

	hasPropose bool
	propose int

	// Persistent Variables
	LastTried int
	Outcome   int
	PrevBal   int
	PrevDec   int
	NextBal   int

	// Non-Persistent Variables
	CurrentStatus Status
	PrevVotes []int
	Quorum []Members.Member
	Voters []Members.Member
	Decree int

}

// Initializers
func Init() *PAXOS_Module {

	var mod = PAXOS_Module{}

	mod.leader = nil
	mod.hasPropose = false
	mod.messages = Init_Messages()
	mod.leaderElection = Omega.Init()

	mod.CurrentStatus = Idle

	mod.LastTried = -1
	mod.Outcome   = -1
	mod.PrevBal   = -1
	mod.PrevDec   = -1
	mod.NextBal   = -1

	go Start(&mod)
	return &mod;

}
func Start(self *PAXOS_Module) {

	go sendNextBallotMessage(self)

	for {
		select {
			case y := <- self.messages.Ind:
				if y.messageType == BeginBallot { self.recvBeginBallot(y.from, y.ballot, y.outcome) }
				if y.messageType == LastVote { self.recvLastVote(y.from, y.ballot, y.outcome) }
				if y.messageType == NextBallot { self.recvNextBallot(y.from, y.ballot) }
				if y.messageType == Success { self.recvSuccess(y.from, y.outcome) }
				if y.messageType == Voted { self.recvVoted(y.from, y.ballot) }
				break;
			case y := <- self.leaderElection.Ind:
				changeLeader(self, y.Member)
				break;
		}
	}
}

// Always Enabled
func sendNextBallotMessage(self *PAXOS_Module) {
	for {
		if self.CurrentStatus == Trying {
			Members.ForEach(func(member Members.Member) {
				fmt.Println("Sim! " + member.Name)
				self.messages.sendMessage_NextBallot(member, self.LastTried)
			})

			time.Sleep(5 * time.Second)
		}

		time.Sleep(500 * time.Millisecond)
	}
}

// Auxiliary Methods
func changeLeader(self *PAXOS_Module, member *Members.Member) {

	self.leader = member;
	fmt.Println("The leader is " + self.leader.Name)
	if member.Name == Members.GetSelf().Name {
		fmt.Println("I'm the leader!");
		tryNewBallot(self);
	}

}

func tryNewBallot(self *PAXOS_Module) {

	self.hasPropose = true
	self.propose = rand.Intn(10) // For testing purposes only

	self.LastTried++
	self.CurrentStatus = Trying
	self.PrevVotes = self.PrevVotes[:0]

}

// Events
func (self PAXOS_Module) recvBeginBallot(from Members.Member, ballot int, outcome int) {

}

func (self PAXOS_Module) recvLastVote(from Members.Member, ballot int, outcome int) {

}

func (self PAXOS_Module) recvNextBallot(from Members.Member, ballot int) {
	if (ballot > self.NextBal) {
		self.NextBal = ballot
		fmt.Println("I'm " + Members.GetSelf().Name + " and my ballot is " + string(ballot))
	}
}

func (self PAXOS_Module) recvSuccess(from Members.Member, outcome int) {

}

func (self PAXOS_Module) recvVoted(from Members.Member, ballot int) {

}