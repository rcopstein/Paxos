package PAXOS

import "../Members"

type Status int
const (
	Idle Status = 0
	Trying Status = 1
	Polling Status = 2
)

type PAXOS_Module struct {

	// SubModules
	messages PAXOS_Messages_Module

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
	mod.messages = Init_Messages()

	return mod;

}

func (self PAXOS_Module) Init_Ballot(ballot int) {

	// Reinitialize Non-Persistent Variables

}
