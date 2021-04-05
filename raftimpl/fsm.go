package raftimpl

import (
	"fmt"
	"github.com/bitcapybara/raft"
	"github.com/vmihailenco/msgpack"
	"sync"
)

type State uint8

const (
	Running State = iota
	Stopped
)

func StateFromString(st string) (state State) {
	switch st {
	case "Running":
		state = Running
	case "Stopped":
		state = Stopped
	}
	return
}

func StateToString(state State) (st string) {
	switch state {
	case Running:
		st = "Running"
	case Stopped:
		st = "Stopped"
	}
	return
}

type Command uint8

const (
	Start Command = iota
	Stop
)

func CommandFromString(cmd string) (command Command) {
	switch cmd {
	case "Start":
		command = Start
	case "Stop":
		command = Stop
	}
	return
}

func CommandToString(command Command) (cmd string) {
	switch command {
	case Start:
		cmd = "Start"
	case Stop:
		cmd = "Stop"
	}
	return
}

type Fsm struct {
	state  State
	logger raft.Logger
	mu     sync.Mutex
}

func NewFsm(logger raft.Logger) *Fsm {
	return &Fsm{
		state: Stopped,
		logger: logger,
	}
}

// 应用状态机命令
func (f *Fsm) Apply(data []byte) error {
	var cmd Command
	umsErr := msgpack.Unmarshal(data, &cmd)
	if umsErr != nil {
		err := fmt.Errorf("反序列化状态机命令失败！%w", umsErr)
		f.logger.Error(err.Error())
		return err
	}
	return f.transfer(cmd)
}

// 生成状态机快照
func (f *Fsm) Serialize() ([]byte, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	stBytes, err := msgpack.Marshal(f.state)
	if err != nil {
		return nil, err
	}
	return stBytes, nil
}

// 安装快照数据
func (f *Fsm) Install(data []byte) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	return msgpack.Unmarshal(data, &f.state)
}

func (f *Fsm) transfer(cmd Command) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	var err error
	f.logger.Trace(fmt.Sprintf("接收到 %s 命令", CommandToString(cmd)))
	if f.state == Running && cmd == Stop {
		f.state = Stopped
		f.logger.Trace(fmt.Sprintf("状态切换为 %s", StateToString(Stopped)))
	} else if f.state == Stopped && cmd == Start {
		f.state = Running
		f.logger.Trace(fmt.Sprintf("状态切换为 %s", StateToString(Running)))
	} else {
		f.logger.Error(fmt.Errorf("不支持的命令：state=%s, command=%s\n", StateToString(f.state), CommandToString(cmd)).Error())
	}
	return err
}
