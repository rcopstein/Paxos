package PAXOS

import (
	"../Members"
	"fmt"
	"math/rand"
	"strconv"
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
type Module struct {

	// SubModules
	messages *MessagesModule
	leaderElection Omega.Omega_Module

	// Variables
	leader *Members.Member

	instanceNumber int
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
func Init() *Module {

	var mod = Module{}

	mod.leader         = nil
	mod.hasPropose     = false
	mod.messages       = InitMessages()
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
	return &mod

}
func (mod *Module) Start() {

	go mod.sendNextBallotMessage()
	go mod.sendLastVoteMessage()
	go mod.pollingMajoritySet()
	go mod.sendBeginBallot()
	go mod.sendVotedMessage()
	go mod.succeed()
	go mod.sendSuccessMessage()

	for {
		select {
			case y := <- mod.messages.Ind:
				if y.messageType == LastVote { mod.recvLastVote(y.from, y.vals[0], y.vals[1], y.vals[2]) }
				if y.messageType == BeginBallot { mod.recvBeginBallot(y.from, y.vals[0], y.vals[1]) }
				if y.messageType == NextBallot { mod.recvNextBallot(y.from, y.vals[0]) }
				if y.messageType == Success { mod.recvSuccess(y.from, y.vals[0]) }
				if y.messageType == Voted { mod.recvVoted(y.from, y.vals[0]) }
				break;
			case y := <- mod.leaderElection.Ind:
				mod.changeLeader(y.Member)
				break;
		}
	}
}

// Always Enabled
func (mod *Module) sendNextBallotMessage() {
	for {
		if mod.CurrentStatus == Trying {
			Members.ForEach(func(member Members.Member) {
				mod.messages.SendMessageNextBallot(mod.instanceNumber, member, mod.LastTried)
			})
		}

		time.Sleep(500 * time.Millisecond)
	}
}
func (mod *Module) sendLastVoteMessage() {
	for {
		if mod.NextBal > mod.PrevBal {
			mod.messages.SendMessageLastVote(mod.instanceNumber, mod.Owners[mod.NextBal], mod.NextBal, mod.PrevBal, mod.PrevDec)
		}

		time.Sleep(500 * time.Millisecond)
	}
}
func (mod *Module) pollingMajoritySet() {
	for {
		if mod.CurrentStatus == Trying && len(mod.PrevVotes) >= (Members.Count() / 2 + 1) {

			fmt.Println("I have a majority set!")

			mod.CurrentStatus = Polling
			mod.Quorum = mod.Quorum[:0]
			mod.Decree = mod.propose
			for _, vote := range mod.PrevVotes {
				if vote.ballot != -1 { mod.Decree = vote.decision }
				mod.Quorum = append(mod.Quorum, vote.member)
			}
			mod.Voters = mod.Voters[:0]
		}

		time.Sleep(500 * time.Millisecond)
	}
}
func (mod *Module) sendBeginBallot() {
	for {
		if mod.CurrentStatus == Polling {
			for _, m := range mod.Quorum {
				mod.messages.SendMessageBeginBallot(mod.instanceNumber, m, mod.LastTried, mod.Decree)
			}
		}

		time.Sleep(500 * time.Millisecond)
	}
}
func (mod *Module) sendVotedMessage() {
	for {
		if mod.PrevBal != -1 {
			mod.messages.SendMessageVoted(mod.instanceNumber, mod.Owners[mod.PrevBal], mod.PrevBal)
		}

		time.Sleep(500 * time.Millisecond)
	}
}
func (mod *Module) succeed() {
	for {
		if mod.CurrentStatus == Polling && mod.Outcome == -1 && checkContains(mod.Quorum, mod.Voters) {
			mod.Outcome = mod.Decree
		}

		time.Sleep(500 * time.Millisecond)
	}
}
func (mod *Module) sendSuccessMessage() {
	for {
		if mod.Outcome != -1 {
			Members.ForEach(func(member Members.Member) {
				mod.messages.SendMessageSuccess(mod.instanceNumber, member, mod.Outcome)
			})
		}

		time.Sleep(500 * time.Millisecond)
	}
}

// Auxiliary Methods
func (mod *Module) changeLeader(member *Members.Member) {

	mod.leader = member

	name := strconv.Itoa(mod.leader.Name)
	fmt.Println("The leader is " + name)
	if member.Name == Members.GetSelf().Name {
		fmt.Println("I'm the leader!")
		mod.tryNewBallot();
	}

}
func (mod *Module) tryNewBallot() {

	mod.hasPropose = true
	mod.propose = rand.Intn(10) // For testing purposes only

	mod.LastTried++
	mod.CurrentStatus = Trying
	mod.PrevVotes = mod.PrevVotes[:0]

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
func (mod *Module) recvLastVote(from Members.Member, NextBal int, LastBal int, LastDec int) {
	if NextBal == mod.LastTried && mod.CurrentStatus == Trying {

		var v = Vote{
			decision: LastDec,
			ballot:   LastBal,
			member:   from,
		}

		for _, vote := range mod.PrevVotes { if vote.member == from { return } }
		mod.PrevVotes = append(mod.PrevVotes, v)
	}
}
func (mod *Module) recvBeginBallot(from Members.Member, ballot int, outcome int) {
	if ballot == mod.NextBal && mod.NextBal > mod.PrevBal {
		mod.PrevBal = ballot
		mod.PrevDec = outcome
	}
}
func (mod *Module) recvNextBallot(from Members.Member, ballot int) {
	if ballot > mod.NextBal {
		mod.NextBal = ballot
		mod.Owners[ballot] = from
		fmt.Println("I'm", Members.GetSelf().Name, "and my ballot is", ballot)
	}
}
func (mod *Module) recvSuccess(from Members.Member, outcome int) {
	fmt.Println("Decided on", outcome, "!")
}
func (mod *Module) recvVoted(from Members.Member, ballot int) {
	if ballot == mod.LastTried && mod.CurrentStatus == Polling {
		mod.Voters = append(mod.Voters, from)
	}
}