package main

import (
	"../Members"
	"../Omega"
	"fmt"
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

	var omega = Omega.Init()
	omega.Start()

	for {
		var a = <- omega.Ind;
		fmt.Print(a.Member.Name)
	}

}
