package main

import (
	"../Members"
	"../Omega"
	"../PAXOS"
	"strconv"
	"fmt"
	"os"
)

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

	if (len(os.Args) < 2) {
		fmt.Println(fmt.Errorf("Member name missing!"))
		return
	}

	addMembers();

	i, _ := strconv.Atoi(os.Args[1])
	Members.SetSelf(i)

	PAXOS.Init()

	var omega = Omega.Init()
	for { <- omega.Ind; }

}