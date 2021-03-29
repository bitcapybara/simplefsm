package simplefsm

import (
	"log"
	"sync"
)

type state uint8

const (
	Running state = iota
	Stopped
)

type command uint8

const (
	Start command = iota
	Stop
)

type fsm struct {
	state  state
	mu     sync.Mutex
	cmdCh  chan command
	exitCh chan struct{}
}

func (f *fsm) start() {
	for {
		select {
		case cmd := <-f.cmdCh:
			f.transfer(cmd)
		case <-f.exitCh:
			return
		}
	}
}

func (f *fsm) transfer(cmd command) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.state == Running && cmd == Stop {
		f.state = Stopped
	} else if f.state == Stopped && cmd == Start {
		f.state = Running
	} else {
		log.Printf("不支持的命令：%d\n", cmd)
	}
}
