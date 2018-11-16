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
	largestProposedOrDecided int
	smallestNonUsed          int
	valueToPropose           int
	omega                    Omega.Omega_Module
	instances                map[int]*SinglePaxos

	mutex *sync.Mutex
}

func NewMulti() *MultiPaxos {

	result := MultiPaxos{}

	result.largestProposedOrDecided = -1
	result.smallestNonUsed = 0

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
				fmt.Println("Decided", y.Values[0], "for instance", snumber)

				if y.Values[0] > mp.largestProposedOrDecided {
					mp.largestProposedOrDecided = number
				}

				if number < mp.largestProposedOrDecided && number > mp.smallestNonUsed {
					for i := mp.smallestNonUsed; i < mp.largestProposedOrDecided; i++ {
						inst, ok := mp.instances[i]
						if !ok {
							inst = mp.CreateInstance(i)
						}
						if inst.Outcome == -1 {
							inst.propose = -2
							inst.hasPropose = true
						}
					}
				}

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
	if !ok {
		instance = mp.CreateInstance(message.Instance)
	}
	instance.MsgReq <- message

}
func (mp *MultiPaxos) SendMessage(message PaxosMessages.Message, to *Members.Member) {

	PaxosMessages.Send(message, to)

}

func (mp *MultiPaxos) CreateInstance(number int) *SinglePaxos {

	mp.mutex.Lock()

	instance := NewSingle(number)
	mp.instances[number] = instance

	mp.mutex.Unlock()

	return instance

}
func (mp *MultiPaxos) ProposeValue(value int) {

	instance, ok := mp.instances[mp.smallestNonUsed]
	if !ok {
		instance = mp.CreateInstance(mp.smallestNonUsed)
	}

	instance.propose = value
	instance.hasPropose = true

	mp.smallestNonUsed++
	if mp.smallestNonUsed > mp.largestProposedOrDecided {
		mp.largestProposedOrDecided = mp.smallestNonUsed
	}

}
