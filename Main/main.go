package main

import "../Members"
import "../PaxosMessages"
import "../Paxos"
import "strconv"
import "fmt"
import "os"

func addMembers() {

	Members.Add(Members.Member{
		Address: "127.0.0.1:4000",
		Name: 1,
	})

	Members.Add(Members.Member{
		Address: "127.0.0.1:5000",
		Name: 2,
	})

	Members.Add(Members.Member{
		Address: "127.0.0.1:6000",
		Name: 3,
	})

	Members.Add(Members.Member{
		Address: "127.0.0.1:7000",
		Name: 4,
	})

}

func main() {

	if len(os.Args) < 2 {
		fmt.Println(fmt.Errorf("Member name missing!"))
		return
	}

	addMembers();
	i, _ := strconv.Atoi(os.Args[1])
	Members.SetSelf(i)

	PaxosMessages.Init(Members.GetSelf())
	mp := Paxos.NewMulti()

	if (i == 1) {
		mp.ProposeValue(100 + i)
		mp.ProposeValue(200 + i)
		mp.ProposeValue(300 + i)
		mp.ProposeValue(400 + i)
		mp.ProposeValue(500 + i)
		mp.ProposeValue(600 + i)
		mp.ProposeValue(700 + i)
		mp.ProposeValue(800 + i)
		mp.ProposeValue(900 + i)
		mp.ProposeValue(1000 + i)
	}

	a := make(chan int)
	for { <- a }
}
