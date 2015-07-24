package gen

import (
	"fmt"
	"math/rand"
)

type Gen interface{}

type RandIntGen struct {
	C chan int64
}

func (r *RandIntGen) Init() {
	r.C = make(chan int64)

	go func() {
		for {
			<-r.C
			r.C <- rand.Int63()
		}
	}()
}

func MakeRandIntGen() *RandIntGen {
	g := RandIntGen{}
	g.Init()

	return &g
}

type RandStringGen struct {
	C chan string
}

func (r *RandStringGen) Init() {
	r.C = make(chan string)

	go func() {
		for {
			<-r.C
			r.C <- " this is msg " + fmt.Sprint(rand.Int63())
		}
	}()
}

func MakeRandMsgGen() *RandStringGen {
	g := RandStringGen{}
	g.Init()

	return &g
}
