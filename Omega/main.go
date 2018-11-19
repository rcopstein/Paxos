package Omega

import (
	"../Members"
	"math/rand"
	"time"
)

type Omega_Trust_Message struct {
	Member *Members.Member
}

type Omega_Module struct {

	Current *Members.Member
	Ind chan Omega_Trust_Message

}

func Init() *Omega_Module {

	var module = Omega_Module{ Ind : make(chan Omega_Trust_Message), }
	go module.Start()
	return &module

}

func (mod *Omega_Module) Start() {

	for {
		// var a = Members.Find_Random();
		var a = Members.GetSelf()
		var old = mod.Current
		mod.Current = a

		if old != a { mod.Ind <- Omega_Trust_Message { Member: a } }

		var b = rand.Intn(5)
		time.Sleep(time.Duration(b) * time.Second)
	}
}