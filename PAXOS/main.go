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
type Vote struct {
	member Members.Member
	decision int
	ballot int
}
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

	Owners map[int]Members.Member

	// Non-Persistent Variables
	CurrentStatus Status
	PrevVotes []Vote
	Quorum []Members.Member
	Voters []Members.Member
	Decree int

}

// Initializers
func Init() *PAXOS_Module {

	var mod = PAXOS_Module{}

	mod.leader         = nil
	mod.hasPropose     = false
	mod.messages       = Init_Messages()
	mod.leaderElection = Omega.Init()

	mod.CurrentStatus = Idle
	mod.Owners        = make(map[int]Members.Member)

	mod.PrevVotes = []Vote{}
	mod.Quorum    = []Members.Member{}
	mod.Voters    = []Members.Member{}

	mod.LastTried = -1
	mod.Outcome   = -1
	mod.PrevBal   = -1
	mod.PrevDec   = -1
	mod.NextBal   = -1

	go mod.Start()
	return &mod;

}
func (self *PAXOS_Module) Start() {

	go self.sendNextBallotMessage()
	go self.sendLastVoteMessage()
	go self.pollingMajoritySet()
	go self.sendBeginBallot()
	go self.sendVotedMessage()
	go self.succeed()
	go self.sendSuccessMessage()

	for {
		select {
			case y := <- self.messages.Ind:
				if y.messageType == LastVote { self.recvLastVote(y.from, y.vals[0], y.vals[1], y.vals[2]) }
				if y.messageType == BeginBallot { self.recvBeginBallot(y.from, y.vals[0], y.vals[1]) }
				if y.messageType == NextBallot { self.recvNextBallot(y.from, y.vals[0]) }
				if y.messageType == Success { self.recvSuccess(y.from, y.vals[0]) }
				if y.messageType == Voted { self.recvVoted(y.from, y.vals[0]) }
				break;
			case y := <- self.leaderElection.Ind:
				self.changeLeader(y.Member)
				break;
		}
	}
}

// Always Enabled
func (self *PAXOS_Module) sendNextBallotMessage() {
	for {
		if self.CurrentStatus == Trying {
			Members.ForEach(func(member Members.Member) {
				self.messages.SendMessage_NextBallot(member, self.LastTried)
			})
		}

		time.Sleep(500 * time.Millisecond)
	}
}
func (self *PAXOS_Module) sendLastVoteMessage() {
	for {
		if self.NextBal > self.PrevBal {
			self.messages.SendMessage_LastVote(self.Owners[self.NextBal], self.NextBal, self.PrevBal, self.PrevDec)
		}

		time.Sleep(500 * time.Millisecond)
	}
}
func (self *PAXOS_Module) pollingMajoritySet() {
	for {
		if self.CurrentStatus == Trying && len(self.PrevVotes) >= (Members.Count() / 2 + 1) {

			fmt.Println("I have a majority set!")

			self.CurrentStatus = Polling
			self.Quorum = self.Quorum[:0]
			self.Decree = self.propose
			for _, vote := range self.PrevVotes {
				if vote.ballot != -1 { self.Decree = vote.decision }
				self.Quorum = append(self.Quorum, vote.member)
			}
			self.Voters = self.Voters[:0]
		}

		time.Sleep(500 * time.Millisecond)
	}
}
func (self *PAXOS_Module) sendBeginBallot() {
	for {
		if self.CurrentStatus == Polling {
			for _, m := range self.Quorum {
				self.messages.SendMessage_BeginBallot(m, self.LastTried, self.Decree)
			}
		}

		time.Sleep(500 * time.Millisecond)
	}
}
func (self *PAXOS_Module) sendVotedMessage() {
	for {
		if self.PrevBal != -1 {
			self.messages.SendMessage_Voted(self.Owners[self.PrevBal], self.PrevBal)
		}

		time.Sleep(500 * time.Millisecond)
	}
}
func (self *PAXOS_Module) succeed() {
	for {
		if self.CurrentStatus == Polling && self.Outcome == -1 && checkContains(self.Quorum, self.Voters) {
			self.Outcome = self.Decree
		}

		time.Sleep(500 * time.Millisecond)
	}
}
func (self *PAXOS_Module) sendSuccessMessage() {
	for {
		if (self.Outcome != -1) {
			Members.ForEach(func(member Members.Member) {
				self.messages.SendMessage_Success(member, self.Outcome)
			})
		}

		time.Sleep(500 * time.Millisecond)
	}
}

// Auxiliary Methods
func (self *PAXOS_Module) changeLeader(member *Members.Member) {

	self.leader = member;
	fmt.Println("The leader is " + self.leader.Name)
	if member.Name == Members.GetSelf().Name {
		fmt.Println("I'm the leader!");
		self.tryNewBallot();
	}

}
func (self *PAXOS_Module) tryNewBallot() {

	self.hasPropose = true
	self.propose = rand.Intn(10) // For testing purposes only

	self.LastTried++
	self.CurrentStatus = Trying
	self.PrevVotes = self.PrevVotes[:0]

}
func checkContains(a []Members.Member, b []Members.Member) bool {
	var flag bool

	for _, ma := range a {
		flag = false
		for _, mb := range b { if ma == mb { flag = true } }
		if !flag { return false }
	}

	return true
}

// Events
func (self *PAXOS_Module) recvLastVote(from Members.Member, NextBal int, LastBal int, LastDec int) {
	if (NextBal == self.LastTried && self.CurrentStatus == Trying) {

		var v = Vote{
			decision: LastDec,
			ballot:   LastBal,
			member:   from,
		}

		for _, vote := range self.PrevVotes { if vote.member == from { return } }
		self.PrevVotes = append(self.PrevVotes, v)
	}
}
func (self *PAXOS_Module) recvBeginBallot(from Members.Member, ballot int, outcome int) {
	if ballot == self.NextBal && self.NextBal > self.PrevBal {
		self.PrevBal = ballot
		self.PrevDec = outcome
	}
}
func (self *PAXOS_Module) recvNextBallot(from Members.Member, ballot int) {
	if (ballot > self.NextBal) {
		self.NextBal = ballot
		self.Owners[ballot] = from
		fmt.Println("I'm", Members.GetSelf().Name, "and my ballot is", ballot)
	}
}
func (self *PAXOS_Module) recvSuccess(from Members.Member, outcome int) {
	fmt.Println("Decided on", outcome, "!")
}
func (self *PAXOS_Module) recvVoted(from Members.Member, ballot int) {
	if ballot == self.LastTried && self.CurrentStatus == Polling {
		self.Voters = append(self.Voters, from)
	}
}