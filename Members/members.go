package Members

import (
	"container/list"
	"math/rand"
)

type Member struct {
	IP string
	Port int
	Nome string
}

var l = list.New();

func Add(memb Member) {
	l.PushFront(memb);
}

func Find_Random() *Member {

	var num = rand.Intn(l.Len());

	for e := l.Front(); true; e = e.Next() {
		if (num == 0) {
			var item= e.Value.(Member);
			return &item;
		}
		num--;
	}

	return nil;
}

func Find(name string) *Member {
	for e := l.Front(); e != nil; e = e.Next() {
		var item = e.Value.(Member);
		if (item.Nome == name) { return &item; }
	}

	return nil;
}