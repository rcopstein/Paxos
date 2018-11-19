package Paxos

import (
	"fmt"
	"strconv"
	"sync"

	"../PaxosMessages"
)
import "../Members"
import "../Omega"

type MultiPaxos struct {
	largestSeen              int
	gapSafety                int

	valueToPropose           int
	omega                    *Omega.Omega_Module
	instances                map[int]*SinglePaxos

	mutex *sync.Mutex
}

func NewMulti() *MultiPaxos {

	result := MultiPaxos{}

	result.largestSeen = -1
	result.gapSafety = 0

	result.instances = make(map[int]*SinglePaxos)
	result.mutex = &sync.Mutex{}
	result.omega = Omega.Init()

	go result.CheckMailbox()
	go result.CheckDecision()
	go result.CheckLeadership()

	return &result

}

func (mp *MultiPaxos) CheckMailbox() {

	for {

		mp.mutex.Lock()

		// Send Messages from Instances
		for _, value := range mp.instances {

			select {
			case y := <-value.MsgInd:

				target := y.Member
				y.Member = Members.GetSelf()
				PaxosMessages.Send(y, target)
				break

			default:
				break
			}

		}

		mp.mutex.Unlock()

		// Receive Messages from Messages Module
		select {
		case y := <-PaxosMessages.Channel:
			mp.ReceiveMessage(y)
			break
		default:
			break
		}
	}
}
func (mp *MultiPaxos) CheckDecision() {
	for {

		mp.mutex.Lock()

		for number, instance := range mp.instances {
			select {
			case y := <-instance.Ind:

				snumber := strconv.Itoa(number)
				fmt.Println("Decided", y.Values[0], "for instance", snumber, "👌")

				if y.Values[0] >  mp.largestSeen { mp.largestSeen = number }

				for i := mp.gapSafety; i < mp.largestSeen; i++ {

					inst, ok := mp.instances[i]
					if !ok { inst = mp.CreateInstance(i, false) }
					if inst.Outcome == -1 { inst.setProposal(-2); }

				}

				mp.gapSafety = mp.largestSeen;

				break

			default:
				break
			}
		}

		mp.mutex.Unlock()

	}
}
func (mp *MultiPaxos) CheckLeadership() {

	for {

		y := <-mp.omega.Ind

		mp.mutex.Lock()

		for _, value := range mp.instances {

			message := PaxosMessages.Message{
				Type:   PaxosMessages.Leader,
				Member: y.Member,
			}

			value.LeadReq <- message

		}

		mp.mutex.Unlock()

	}
}

func (mp *MultiPaxos) ReceiveMessage(message PaxosMessages.Message) {

	instance, ok := mp.instances[message.Instance]
	if !ok { instance = mp.CreateInstance(message.Instance, true) }
	instance.MsgReq <- message

}
func (mp *MultiPaxos) SendMessage(message PaxosMessages.Message, to *Members.Member) {

	PaxosMessages.Send(message, to)

}

func (mp *MultiPaxos) CreateInstance(number int, lock bool) *SinglePaxos {

	if lock { mp.mutex.Lock() }

	instance := NewSingle(number, mp.omega.Current)
	mp.instances[number] = instance

	if number > mp.largestSeen { mp.largestSeen = number }

	if lock { mp.mutex.Unlock() }

	return instance

}
func (mp *MultiPaxos) ProposeValue(value int) {

	instance, ok := mp.instances[mp.largestSeen + 1]
	if !ok { instance = mp.CreateInstance(mp.largestSeen + 1, true) }

	instance.setProposal(value)

}
