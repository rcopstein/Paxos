package Paxos

import (
	"../PaxosMessages"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)
import "../Members"

type Status int
const (
	Idle Status = 0
	Trying Status = 1
	Polling Status = 2
)

type Vote struct {
	member *Members.Member
	decision int
	ballot int
}
type SinglePaxos struct {

	// Channels
	Ind chan PaxosMessages.Message
	MsgInd chan PaxosMessages.Message

	Req chan PaxosMessages.Message
	MsgReq chan PaxosMessages.Message
	LeadReq chan PaxosMessages.Message

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

	Owners map[int]*Members.Member

	// Non-Persistent Variables
	CurrentStatus Status
	PrevVotes []Vote
	Quorum []*Members.Member
	Voters []*Members.Member
	Decree int

}

func NewSingle(instanceNumber int) *SinglePaxos {

	result := SinglePaxos{}

	result.Ind = make(chan PaxosMessages.Message)
	result.MsgInd = make(chan PaxosMessages.Message)

	result.Req = make(chan PaxosMessages.Message)
	result.LeadReq = make(chan PaxosMessages.Message)

	result.leader         = nil
	result.hasPropose     = false
	result.instanceNumber = instanceNumber

	result.CurrentStatus = Idle
	result.Owners        = make(map[int]*Members.Member)

	result.PrevVotes = []Vote{}
	result.Quorum    = []*Members.Member{}
	result.Voters    = []*Members.Member{}

	result.LastTried = -1
	result.Outcome   = -1
	result.PrevBal   = -1
	result.PrevDec   = -1
	result.NextBal   = -1

	go result.Start()
	return &result

}
func (sp *SinglePaxos) Start() {

	go sp.sendNextBallotMessage()
	go sp.sendLastVoteMessage()
	go sp.pollingMajoritySet()
	go sp.sendBeginBallot()
	go sp.sendVotedMessage()
	go sp.succeed()
	go sp.sendSuccessMessage()

	for {
		select {
		case y := <- sp.MsgReq:
			if y.Type == PaxosMessages.LastVote    { sp.recvLastVote(y.Member, y.Values[0], y.Values[1], y.Values[2]) }
			if y.Type == PaxosMessages.BeginBallot { sp.recvBeginBallot(y.Member, y.Values[0], y.Values[1]) }
			if y.Type == PaxosMessages.NextBallot  { sp.recvNextBallot(y.Member, y.Values[0]) }
			if y.Type == PaxosMessages.Success     { sp.recvSuccess(y.Member, y.Values[0]) }
			if y.Type == PaxosMessages.Voted       { sp.recvVoted(y.Member, y.Values[0]) }
			break;
		case y := <- sp.LeadReq:
			sp.changeLeader(y.Member)
			break;
		}
	}

}

func (sp *SinglePaxos) sendNextBallotMessage() {
	for {
		if sp.CurrentStatus == Trying {
			Members.ForEach(func(member *Members.Member) {

				message := PaxosMessages.Message{
					Instance : sp.instanceNumber,
					Member : member,
					Values : []int { sp.LastTried },
				}

				sp.MsgInd <- message

			})
		}

		time.Sleep(500 * time.Millisecond)
	}
}
func (sp *SinglePaxos) sendLastVoteMessage() {
	for {
		if sp.NextBal > sp.PrevBal {

			message := PaxosMessages.Message{
				Instance : sp.instanceNumber,
				Member : sp.Owners[sp.NextBal],
				Values : []int { sp.NextBal, sp.PrevBal, sp.PrevDec },
			}

			sp.MsgInd <- message

		}

		time.Sleep(500 * time.Millisecond)
	}
}
func (sp *SinglePaxos) pollingMajoritySet() {
	for {
		if sp.CurrentStatus == Trying && len(sp.PrevVotes) >= (Members.Count() / 2 + 1) {

			fmt.Println("I have a majority set!")

			sp.CurrentStatus = Polling
			sp.Quorum = sp.Quorum[:0]
			sp.Decree = sp.propose
			for _, vote := range sp.PrevVotes {
				if vote.ballot != -1 { sp.Decree = vote.decision }
				sp.Quorum = append(sp.Quorum, vote.member)
			}
			sp.Voters = sp.Voters[:0]
		}

		time.Sleep(500 * time.Millisecond)
	}
}
func (sp *SinglePaxos) sendSuccessMessage() {
	for {
		if sp.Outcome != -1 {
			Members.ForEach(func(member *Members.Member) {

				message := PaxosMessages.Message{
					Instance : sp.instanceNumber,
					Member : member,
					Values : []int { sp.Outcome },
				}

				sp.MsgInd <- message

			})
		}

		time.Sleep(500 * time.Millisecond)
	}
}
func (sp *SinglePaxos) sendVotedMessage() {
	for {
		if sp.PrevBal != -1 {

			message := PaxosMessages.Message{
				Instance : sp.instanceNumber,
				Member : sp.Owners[sp.PrevBal],
				Values : []int { sp.PrevBal },
			}

			sp.MsgInd <- message

		}

		time.Sleep(500 * time.Millisecond)
	}
}
func (sp *SinglePaxos) sendBeginBallot() {
	for {
		if sp.CurrentStatus == Polling {
			for _, m := range sp.Quorum {

				message := PaxosMessages.Message{
					Instance : sp.instanceNumber,
					Member : m,
					Values : []int { sp.LastTried, sp.Decree },
				}

				sp.MsgInd <- message

			}
		}

		time.Sleep(500 * time.Millisecond)
	}
}
func (sp *SinglePaxos) succeed() {
	for {
		if sp.CurrentStatus == Polling && sp.Outcome == -1 && checkContains(sp.Quorum, sp.Voters) {
			sp.Outcome = sp.Decree
		}

		time.Sleep(500 * time.Millisecond)
	}
}

// Auxiliary Methods
func (sp *SinglePaxos) tryNewBallot() {

	sp.hasPropose = true
	sp.propose = rand.Intn(10) // For testing purposes only

	sp.LastTried++
	sp.CurrentStatus = Trying
	sp.PrevVotes = sp.PrevVotes[:0]

}
func (sp *SinglePaxos) changeLeader(member *Members.Member) {

	sp.leader = member

	name := strconv.Itoa(sp.leader.Name)
	name = name
	// fmt.Println("The leader is " + name)
	if member.Name == Members.GetSelf().Name {
		// fmt.Println("I'm the leader!")
		sp.tryNewBallot();
	}

}
func checkContains(a []*Members.Member, b []*Members.Member) bool {
	var flag bool

	for _, ma := range a {
		flag = false
		for _, mb := range b { if ma == mb { flag = true } }
		if !flag { return false }
	}

	return true
}

// Events
func (sp *SinglePaxos) recvLastVote(from *Members.Member, NextBal int, LastBal int, LastDec int) {
	if NextBal == sp.LastTried && sp.CurrentStatus == Trying {

		var v = Vote{
			decision: LastDec,
			ballot:   LastBal,
			member:   from,
		}

		for _, vote := range sp.PrevVotes { if vote.member == from { return } }
		sp.PrevVotes = append(sp.PrevVotes, v)
	}
}
func (sp *SinglePaxos) recvBeginBallot(from *Members.Member, ballot int, outcome int) {
	if ballot == sp.NextBal && sp.NextBal > sp.PrevBal {
		sp.PrevBal = ballot
		sp.PrevDec = outcome
	}
}
func (sp *SinglePaxos) recvNextBallot(from *Members.Member, ballot int) {
	if ballot > sp.NextBal {
		sp.NextBal = ballot
		sp.Owners[ballot] = from
		fmt.Println("I'm", Members.GetSelf().Name, "and my ballot is", ballot)
	}
}
func (sp *SinglePaxos) recvSuccess(from *Members.Member, outcome int) {

	if sp.Outcome != -1 { return }

	sp.Outcome = outcome
	fmt.Println("Decided on", outcome, "!")

	message := PaxosMessages.Message{
		Instance : sp.instanceNumber,
		Values : []int { outcome },
	}

	sp.Ind <- message

}
func (sp *SinglePaxos) recvVoted(from *Members.Member, ballot int) {
	if ballot == sp.LastTried && sp.CurrentStatus == Polling {
		sp.Voters = append(sp.Voters, from)
	}
}