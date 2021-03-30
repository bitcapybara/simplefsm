package simplefsm

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/bitcapybara/raft"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type server struct {
	node raft.Node
}

func (s *server) Run(me raft.NodeId, peers map[raft.NodeId]raft.NodeAddr) {

	// 启动 raft
	config := raft.Config{
		Fsm:                newFsm(),
		RaftStatePersister: newRaftStatePersister(),
		SnapshotPersister:  newSnapshotPersister(),
		Transport:          newHttpTransport(),
		Peers:              peers,
		Me:                 me,
		ElectionMaxTimeout: 10000,
		ElectionMinTimeout: 5000,
		HeartbeatTimeout:   1000,
		MaxLogLength:       2000,
	}
	node := raft.NewNode(config)
	go node.Run()

	// 启动 echo
	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.POST("/appendEntries", s.appendEntries)
	e.POST("/requestVote", s.requestVote)
	e.POST("/installSnapshot", s.installSnapshot)
	e.POST("/applyCommand", s.applyCommand)
	e.POST("/changeConfig", s.changeConfig)
	e.POST("/transferLeadership", s.transferLeadership)
	e.POST("/addNewNode", s.addNewNode)

	// Start server
	e.Logger.Fatal(e.Start(string(peers[me])))
}

func (s *server) appendEntries(ctx echo.Context) error {
	// 反序列化获取请求参数
	decoder := gob.NewDecoder(ctx.Request().Body)
	var args raft.AppendEntry
	deErr := decoder.Decode(&args)
	if decoder != nil {
		return fmt.Errorf("反序列化参数失败！%w", deErr)
	}
	// 调用 raft 逻辑
	var res raft.AppendEntryReply
	raftErr := s.node.AppendEntries(args, &res)
	if raftErr != nil {
		return fmt.Errorf("raft 操作失败！%w", raftErr)
	}
	// 序列化并返回结果
	data := new(bytes.Buffer)
	encoder := gob.NewEncoder(data)
	enErr := encoder.Encode(res)
	if enErr != nil {
		return fmt.Errorf("序列化结果失败！%w", enErr)
	}
	ctxErr := ctx.Blob(200, "application/octet-stream", data.Bytes())
	if ctxErr != nil {
		return fmt.Errorf("处理返回值失败！%w", ctxErr)
	}
	return nil
}

func (s *server) requestVote(ctx echo.Context) error {
	// 反序列化获取请求参数
	decoder := gob.NewDecoder(ctx.Request().Body)
	var args raft.RequestVote
	deErr := decoder.Decode(&args)
	if decoder != nil {
		return fmt.Errorf("反序列化参数失败！%w", deErr)
	}
	// 调用 raft 逻辑
	var res raft.RequestVoteReply
	raftErr := s.node.RequestVote(args, &res)
	if raftErr != nil {
		return fmt.Errorf("raft 操作失败！%w", raftErr)
	}
	// 序列化并返回结果
	data := new(bytes.Buffer)
	encoder := gob.NewEncoder(data)
	enErr := encoder.Encode(res)
	if enErr != nil {
		return fmt.Errorf("序列化结果失败！%w", enErr)
	}
	ctxErr := ctx.Blob(200, "application/octet-stream", data.Bytes())
	if ctxErr != nil {
		return fmt.Errorf("处理返回值失败！%w", ctxErr)
	}
	return nil
}

func (s *server) installSnapshot(ctx echo.Context) error {
	// 反序列化获取请求参数
	decoder := gob.NewDecoder(ctx.Request().Body)
	var args raft.RequestVote
	deErr := decoder.Decode(&args)
	if decoder != nil {
		return fmt.Errorf("反序列化参数失败！%w", deErr)
	}
	// 调用 raft 逻辑
	var res raft.RequestVoteReply
	raftErr := s.node.RequestVote(args, &res)
	if raftErr != nil {
		return fmt.Errorf("raft 操作失败！%w", raftErr)
	}
	// 序列化并返回结果
	data := new(bytes.Buffer)
	encoder := gob.NewEncoder(data)
	enErr := encoder.Encode(res)
	if enErr != nil {
		return fmt.Errorf("序列化结果失败！%w", enErr)
	}
	ctxErr := ctx.Blob(200, "application/octet-stream", data.Bytes())
	if ctxErr != nil {
		return fmt.Errorf("处理返回值失败！%w", ctxErr)
	}
	return nil
}

func (s *server) applyCommand(ctx echo.Context) error {
	// 反序列化获取请求参数
	decoder := gob.NewDecoder(ctx.Request().Body)
	var args raft.ApplyCommand
	deErr := decoder.Decode(&args)
	if decoder != nil {
		return fmt.Errorf("反序列化参数失败！%w", deErr)
	}
	// 调用 raft 逻辑
	var res raft.ApplyCommandReply
	raftErr := s.node.ApplyCommand(args, &res)
	if raftErr != nil {
		return fmt.Errorf("raft 操作失败！%w", raftErr)
	}
	// 序列化并返回结果
	data := new(bytes.Buffer)
	encoder := gob.NewEncoder(data)
	enErr := encoder.Encode(res)
	if enErr != nil {
		return fmt.Errorf("序列化结果失败！%w", enErr)
	}
	ctxErr := ctx.Blob(200, "application/octet-stream", data.Bytes())
	if ctxErr != nil {
		return fmt.Errorf("处理返回值失败！%w", ctxErr)
	}
	return nil
}

func (s *server) changeConfig(ctx echo.Context) error {
	// 反序列化获取请求参数
	decoder := gob.NewDecoder(ctx.Request().Body)
	var args raft.ChangeConfig
	deErr := decoder.Decode(&args)
	if decoder != nil {
		return fmt.Errorf("反序列化参数失败！%w", deErr)
	}
	// 调用 raft 逻辑
	var res raft.ChangeConfigReply
	raftErr := s.node.ChangeConfig(args, &res)
	if raftErr != nil {
		return fmt.Errorf("raft 操作失败！%w", raftErr)
	}
	// 序列化并返回结果
	data := new(bytes.Buffer)
	encoder := gob.NewEncoder(data)
	enErr := encoder.Encode(res)
	if enErr != nil {
		return fmt.Errorf("序列化结果失败！%w", enErr)
	}
	ctxErr := ctx.Blob(200, "application/octet-stream", data.Bytes())
	if ctxErr != nil {
		return fmt.Errorf("处理返回值失败！%w", ctxErr)
	}
	return nil
}

func (s *server) transferLeadership(ctx echo.Context) error {
	// 反序列化获取请求参数
	decoder := gob.NewDecoder(ctx.Request().Body)
	var args raft.TransferLeadership
	deErr := decoder.Decode(&args)
	if decoder != nil {
		return fmt.Errorf("反序列化参数失败！%w", deErr)
	}
	// 调用 raft 逻辑
	var res raft.TransferLeadershipReply
	raftErr := s.node.TransferLeadership(args, &res)
	if raftErr != nil {
		return fmt.Errorf("raft 操作失败！%w", raftErr)
	}
	// 序列化并返回结果
	data := new(bytes.Buffer)
	encoder := gob.NewEncoder(data)
	enErr := encoder.Encode(res)
	if enErr != nil {
		return fmt.Errorf("序列化结果失败！%w", enErr)
	}
	ctxErr := ctx.Blob(200, "application/octet-stream", data.Bytes())
	if ctxErr != nil {
		return fmt.Errorf("处理返回值失败！%w", ctxErr)
	}
	return nil
}

func (s *server) addNewNode(ctx echo.Context) error {
	// 反序列化获取请求参数
	decoder := gob.NewDecoder(ctx.Request().Body)
	var args raft.AddNewNode
	deErr := decoder.Decode(&args)
	if decoder != nil {
		return fmt.Errorf("反序列化参数失败！%w", deErr)
	}
	// 调用 raft 逻辑
	var res raft.AddNewNodeReply
	raftErr := s.node.AddNewNode(args, &res)
	if raftErr != nil {
		return fmt.Errorf("raft 操作失败！%w", raftErr)
	}
	// 序列化并返回结果
	data := new(bytes.Buffer)
	encoder := gob.NewEncoder(data)
	enErr := encoder.Encode(res)
	if enErr != nil {
		return fmt.Errorf("序列化结果失败！%w", enErr)
	}
	ctxErr := ctx.Blob(200, "application/octet-stream", data.Bytes())
	if ctxErr != nil {
		return fmt.Errorf("处理返回值失败！%w", ctxErr)
	}
	return nil
}
