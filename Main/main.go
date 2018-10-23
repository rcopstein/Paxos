package main

import (
	"../Members"
	"../Omega"
	"fmt"
)

func add_members() {

	Members.Add(Members.Member{
		IP:   "127.0.0.1",
		Port: 4000,
		Nome: "1",
	});

	Members.Add(Members.Member{
		IP:   "127.0.0.1",
		Port: 5000,
		Nome: "2",
	});

	Members.Add(Members.Member{
		IP:   "127.0.0.1",
		Port: 6000,
		Nome: "3",
	});

}

func main() {

	add_members();

	var omega = Omega.Init()
	omega.Start()

	for {
		var a = <- omega.Ind;
		fmt.Print(a.Member.Nome)
	}

}
