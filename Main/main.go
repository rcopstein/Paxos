package main

import (
	"../Members"
	"../Omega"
	"../PAXOS"
)

func addMembers() {

	Members.Add(Members.Member{
		Address: "127.0.0.1:4000",
		Name: "1",
	})

	Members.Add(Members.Member{
		Address: "127.0.0.1:5000",
		Name: "2",
	})

	Members.Add(Members.Member{
		Address: "127.0.0.1:6000",
		Name: "3",
	})

	Members.SetSelf("1")

}

func main() {

	addMembers();

	PAXOS.Init()

	var omega = Omega.Init()
	for { <- omega.Ind; }

}
