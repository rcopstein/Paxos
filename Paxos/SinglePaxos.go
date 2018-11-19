package Paxos

import (
	"../PaxosMessages"
	"fmt"
	"strconv"
	"time"
)
import "../Members"

var printMessages = false

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

func NewSingle(instanceNumber int, leader *Members.Member) *SinglePaxos {

	result := SinglePaxos{}

	result.Ind = make(chan PaxosMessages.Message)
	result.MsgInd = make(chan PaxosMessages.Message)

	result.Req = make(chan PaxosMessages.Message)
	result.MsgReq = make(chan PaxosMessages.Message)
	result.LeadReq = make(chan PaxosMessages.Message)

	result.leader         = leader
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
	go sp.sendSuccessMessage()
	go sp.sendVotedMessage()
	go sp.sendBeginBallot()
	go sp.succeed()

	for {
		select {
		case y := <- sp.MsgReq:
			if y.Type == PaxosMessages.BeginBallot { sp.recvBeginBallot(y.Member, y.Values[0], y.Values[1]) }
			if y.Type == PaxosMessages.NextBallot  { sp.recvNextBallot(y.Member, y.Values[0]) }
			if y.Type == PaxosMessages.LastVote    { sp.recvLastVote(y.Member, y.Values[0], y.Values[1], y.Values[2]) }
			if y.Type == PaxosMessages.Success     { sp.recvSuccess(y.Member, y.Values[0]) }
			if y.Type == PaxosMessages.Voted       { sp.recvVoted(y.Member, y.Values[0]) }
			break
		case y := <- sp.LeadReq:
			sp.changeLeader(y.Member)
			break
		}
	}

}

func (sp *SinglePaxos) sendNextBallotMessage() {
	for {
		if sp.CurrentStatus == Trying {
			Members.ForEach(func(member *Members.Member) {

				message := PaxosMessages.Message{
					Type : PaxosMessages.NextBallot,
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
				Type : PaxosMessages.LastVote,
				Instance : sp.instanceNumber,
				Member : sp.Owners[sp.NextBal],
				Values : []int { sp.NextBal, sp.PrevBal, sp.PrevDec },
			}

			if message.Member != nil { sp.MsgInd <- message }

		}

		time.Sleep(500 * time.Millisecond)
	}
}
func (sp *SinglePaxos) pollingMajoritySet() {
	for {
		if sp.CurrentStatus == Trying && len(sp.PrevVotes) >= (Members.Count() / 2 + 1) {

			if (printMessages) { fmt.Println("I have a majority set!") }

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
					Type : PaxosMessages.Success,
					Instance : sp.instanceNumber,
					Member : member,
					Values : []int { sp.Outcome },
				}

				if message.Member != nil { sp.MsgInd <- message }

			})
		}

		time.Sleep(500 * time.Millisecond)
	}
}
func (sp *SinglePaxos) sendVotedMessage() {
	for {
		if sp.PrevBal != -1 {

			message := PaxosMessages.Message{
				Type : PaxosMessages.Voted,
				Instance : sp.instanceNumber,
				Member : sp.Owners[sp.PrevBal],
				Values : []int { sp.PrevBal },
			}

			if message.Member != nil { sp.MsgInd <- message }

		}

		time.Sleep(500 * time.Millisecond)
	}
}
func (sp *SinglePaxos) sendBeginBallot() {
	for {
		if sp.CurrentStatus == Polling {
			for _, m := range sp.Quorum {

				message := PaxosMessages.Message{
					Type : PaxosMessages.BeginBallot,
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
			if printMessages { fmt.Println("Decided on", sp.Decree, "!") }
			sp.Outcome = sp.Decree

			message := PaxosMessages.Message{
				Instance : sp.instanceNumber,
				Values : []int { sp.Outcome },
			}
			sp.Ind <- message

		}

		time.Sleep(500 * time.Millisecond)
	}
}

// Auxiliary Methods
func (sp *SinglePaxos) setProposal(value int) {

	sp.propose = value
	sp.hasPropose = true

	sp.tryPropose()

}
func (sp *SinglePaxos) tryPropose() {

	if printMessages {
		fmt.Println(sp.leader)
		fmt.Println(sp.Outcome)
		fmt.Println(sp.hasPropose)
	}

	if sp.leader != nil && sp.leader.Name == Members.GetSelf().Name && sp.hasPropose && sp.Outcome == -1 {

		if printMessages { fmt.Println("Proposing: ", sp.propose) }

		sp.LastTried++
		sp.CurrentStatus = Trying
		sp.PrevVotes = sp.PrevVotes[:0]

	}

}

func (sp *SinglePaxos) changeLeader(member *Members.Member) {

	fmt.Println("Leader is", member);

	sp.leader = member
	sp.tryPropose()

}
func checkContains(a []*Members.Member, b []*Members.Member) bool {
	var flag bool

	for _, ma := range a {
		flag = false
		for _, mb := range b { if ma.Name == mb.Name { flag = true } }
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

		for _, vote := range sp.PrevVotes {
			if vote.member.Name == from.Name {
				return
			}
		}
		sp.PrevVotes = append(sp.PrevVotes, v)

		name := strconv.Itoa(from.Name)
		if (printMessages) { fmt.Println("LastVote from", name) }

	}
}
func (sp *SinglePaxos) recvBeginBallot(from *Members.Member, ballot int, outcome int) {
	if ballot == sp.NextBal && sp.NextBal > sp.PrevBal {
		sp.PrevBal = ballot
		sp.PrevDec = outcome

		name := strconv.Itoa(from.Name)
		if (printMessages) { fmt.Println("BeginBallot from", name) }
	}
}
func (sp *SinglePaxos) recvNextBallot(from *Members.Member, ballot int) {
	if ballot > sp.NextBal {
		sp.NextBal = ballot
		sp.Owners[ballot] = from

		name := strconv.Itoa(from.Name)
		if (printMessages) { fmt.Println("NextBallot from", name) }
	}
}
func (sp *SinglePaxos) recvSuccess(from *Members.Member, outcome int) {

	if sp.Outcome != -1 { return }

	sp.Outcome = outcome
	if (printMessages) { fmt.Println("Decided on", outcome, "!") }

	message := PaxosMessages.Message{
		Instance : sp.instanceNumber,
		Values : []int { outcome },
	}
	sp.Ind <- message

}
func (sp *SinglePaxos) recvVoted(from *Members.Member, ballot int) {
	if ballot == sp.LastTried && sp.CurrentStatus == Polling {
		if !checkContains([]*Members.Member{from}, sp.Voters) {
			sp.Voters = append(sp.Voters, from)

			name := strconv.Itoa(from.Name)
			if (printMessages) { fmt.Println("Voted from", name) }
		}
	}
}