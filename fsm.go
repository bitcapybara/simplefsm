package simplefsm

import (
	"bytes"
	"encoding/gob"
	"fmt"
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
}

func newFsm() *fsm {
	return &fsm{
		state: Stopped,
	}
}

// 应用状态机命令
func (f *fsm) Apply(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	var cmd command
	err := decoder.Decode(&cmd)
	if err != nil {
		return err
	}
	return f.transfer(cmd)
}

// 生成状态机快照
func (f *fsm) Serialize() ([]byte, error) {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(f.state)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (f *fsm) transfer(cmd command) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	var err error
	if f.state == Running && cmd == Stop {
		f.state = Stopped
	} else if f.state == Stopped && cmd == Start {
		f.state = Running
	} else {
		err = fmt.Errorf("不支持的命令：%d\n", cmd)
	}
	return err
}
