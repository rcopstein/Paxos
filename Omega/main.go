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
	Ind chan Omega_Trust_Message
}

func Init() Omega_Module {

	var module = Omega_Module{ Ind : make(chan Omega_Trust_Message), }
	go module.Start()
	return module;

}

func (self Omega_Module) Start() {

	for {
		var a = Members.Find_Random();
		self.Ind <- Omega_Trust_Message{Member: a}

		var b = rand.Intn(5)
		time.Sleep(time.Duration(b) * time.Second)
	}
}