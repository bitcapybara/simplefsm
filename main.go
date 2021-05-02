package main

import (
	"flag"
	"fmt"
	"github.com/bitcapybara/raft"
	"github.com/bitcapybara/simplefsm/raftimpl"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/vmihailenco/msgpack"
	"log"
	"strings"
)

const NoneOption = ""

func main() {
	// 命令行参数定义
	var me string
	flag.StringVar(&me, "me", NoneOption, "当前节点nodeId")
	var peerStr string
	flag.StringVar(&peerStr, "peers", NoneOption, "指定所有节点地址，nodeId@nodeAddr，多个地址使用逗号间隔")
	var role string
	flag.StringVar(&role, "role", NoneOption, "当前节点角色")
	flag.Parse()

	// 命令行参数解析
	if me == "" {
		log.Fatal("未指定当前节点id！")
	}

	if peerStr == "" {
		log.Fatalln("未指定集群节点")
	}

	if role == "" {
		log.Fatalln("未指定节点角色")
	}

	peerSplit := strings.Split(peerStr, ",")
	peers := make(map[raft.NodeId]raft.NodeAddr, len(peerSplit))
	for _, peerInfo := range peerSplit {
		idAndAddr := strings.Split(peerInfo, "@")
		peers[raft.NodeId(idAndAddr[0])] = raft.NodeAddr(idAndAddr[1])
	}

	// 启动 server
	s := newServer(raft.RoleFromString(role), raft.NodeId(me), peers)
	s.Start()
}

type server struct {
	addr   string
	node   *raft.Node
	echo   *echo.Echo
	fsm    *raftimpl.Fsm
	logger *raftimpl.SimpleLogger
	enable bool
}

func newServer(role raft.RoleStage, me raft.NodeId, peers map[raft.NodeId]raft.NodeAddr) *server {
	// 启动 raft
	logger := raftimpl.NewLogger()
	config := raft.Config{
		Fsm:                raftimpl.NewFsm(logger),
		RaftStatePersister: raftimpl.NewRaftStatePersister(),
		SnapshotPersister:  raftimpl.NewSnapshotPersister(),
		Transport:          raftimpl.NewHttpTransport(logger),
		Logger:             logger,
		Peers:              peers,
		Me:                 me,
		Role:               role,
		ElectionMaxTimeout: 10000,
		ElectionMinTimeout: 5000,
		HeartbeatTimeout:   1000,
		MaxLogLength:       50,
	}
	node := raft.NewNode(config)

	// 启动 echo
	e := echo.New()
	return &server{
		addr:   string(peers[me]),
		node:   node,
		echo:   e,
		fsm:    raftimpl.NewFsm(logger),
		logger: logger,
		enable: true,
	}
}

func (s *server) Start() {
	go s.node.Run()

	e := s.echo
	// Middleware
	e.Use(middleware.Recover())

	// 由用户调用
	e.GET("/state", s.getState)
	e.POST("/applyCommand", s.applyCommand)
	e.POST("/sleep", s.sleep)
	e.POST("/awake", s.awake)

	// 成员变更
	e.POST("/changeConfig", s.changeConfig)
	// 添加 Learner 角色，可通过成员变更将其升级为 Follower 加入集群
	e.POST("/addLearner", s.addLearner)
	// 领导权转移
	e.POST("/transferLeadership", s.transferLeadership)

	// 由 raft 调用
	e.POST("/appendEntries", s.appendEntries)
	e.POST("/requestVote", s.requestVote)
	e.POST("/installSnapshot", s.installSnapshot)

	// Start server
	e.Logger.Fatal(e.Start(s.addr))
}

func (s *server) getState(ctx echo.Context) error {
	return ctx.JSON(200, struct{ State string }{s.fsm.GetState()})
}

func (s *server) appendEntries(ctx echo.Context) (err error) {
	defer func() {
		if err != nil {
			s.logger.Error(err.Error())
		}
	}()
	if !s.enable {
		err = fmt.Errorf("server not enable")
		return
	}
	// 反序列化获取请求参数
	var args raft.AppendEntry
	bindErr := ctx.Bind(&args)
	if bindErr != nil {
		return fmt.Errorf("反序列化参数失败！%w", bindErr)
	}
	// 调用 raft 逻辑
	var res raft.AppendEntryReply
	raftErr := s.node.AppendEntries(args, &res)
	if raftErr != nil {
		return fmt.Errorf("raft 操作失败！%w", raftErr)
	}
	// 序列化并返回结果
	return ctx.JSON(200, res)
}

func (s *server) requestVote(ctx echo.Context) (err error) {
	defer func() {
		if err != nil {
			s.logger.Error(err.Error())
		}
	}()
	if !s.enable {
		err = fmt.Errorf("server not enable")
		return
	}
	// 反序列化获取请求参数
	var args raft.RequestVote
	bindErr := ctx.Bind(&args)
	if bindErr != nil {
		return fmt.Errorf("反序列化参数失败！%w", bindErr)
	}
	// 调用 raft 逻辑
	var res raft.RequestVoteReply
	raftErr := s.node.RequestVote(args, &res)
	if raftErr != nil {
		return fmt.Errorf("raft 操作失败！%w", raftErr)
	}
	// 序列化并返回结果
	return ctx.JSON(200, res)
}

func (s *server) installSnapshot(ctx echo.Context) (err error) {
	defer func() {
		if err != nil {
			s.logger.Error(err.Error())
		}
	}()
	if !s.enable {
		err = fmt.Errorf("server not enable")
		return
	}
	// 反序列化获取请求参数
	var args raft.InstallSnapshot
	bindErr := ctx.Bind(&args)
	if bindErr != nil {
		return fmt.Errorf("反序列化参数失败！%w", bindErr)
	}
	// 调用 raft 逻辑
	var res raft.InstallSnapshotReply
	raftErr := s.node.InstallSnapshot(args, &res)
	if raftErr != nil {
		return fmt.Errorf("raft 操作失败！%w", raftErr)
	}
	// 序列化并返回结果
	return ctx.JSON(200, res)
}

func (s *server) applyCommand(ctx echo.Context) (err error) {
	defer func() {
		if err != nil {
			s.logger.Error(err.Error())
		}
	}()
	if !s.enable {
		err = fmt.Errorf("server not enable")
		return
	}
	// 反序列化获取请求参数
	command := ctx.QueryParam("command")
	s.logger.Trace(fmt.Sprintf("获取到命令 %s", command))
	cmdBytes, msErr := msgpack.Marshal(raftimpl.CommandFromString(command))
	if msErr != nil {
		return fmt.Errorf("序列化命令失败！%w", msErr)
	}
	args := raft.ApplyCommand{
		Data: cmdBytes,
	}
	// 调用 raft 逻辑
	var res raft.ApplyCommandReply
	raftErr := s.node.ApplyCommand(args, &res)
	if raftErr != nil {
		return fmt.Errorf("raft 操作失败！%w", raftErr)
	}
	// 序列化并返回结果
	return ctx.JSON(200, res)
}

func (s *server) addLearner(ctx echo.Context) (err error) {
	defer func() {
		if err != nil {
			s.logger.Error(err.Error())
		}
	}()
	if !s.enable {
		err = fmt.Errorf("server not enable")
		return
	}
	// 反序列化获取请求参数
	var peers map[raft.NodeId]raft.NodeAddr
	argsErr := ctx.Bind(&peers)
	if argsErr != nil {
		return fmt.Errorf("读取请求参数失败！%w", argsErr)
	}
	args := raft.AddLearner{
		Learners: peers,
	}
	// 调用 raft 逻辑
	var res raft.AddLearnerReply
	raftErr := s.node.AddLearner(args, &res)
	if raftErr != nil {
		return fmt.Errorf("raft 操作失败！%w", raftErr)
	}
	// 序列化并返回结果
	return ctx.JSON(200, res)
}

func (s *server) changeConfig(ctx echo.Context) (err error) {
	defer func() {
		if err != nil {
			s.logger.Error(err.Error())
		}
	}()
	if !s.enable {
		err = fmt.Errorf("server not enable")
		return
	}
	// 反序列化获取请求参数
	var peers map[raft.NodeId]raft.NodeAddr
	argsErr := ctx.Bind(&peers)
	if argsErr != nil {
		return fmt.Errorf("读取请求参数失败！%w", argsErr)
	}
	args := raft.ChangeConfig{
		Peers: peers,
	}
	// 调用 raft 逻辑
	var res raft.ChangeConfigReply
	raftErr := s.node.ChangeConfig(args, &res)
	if raftErr != nil {
		return fmt.Errorf("raft 操作失败！%w", raftErr)
	}
	// 序列化并返回结果
	return ctx.JSON(200, res)
}

func (s *server) transferLeadership(ctx echo.Context) (err error) {
	defer func() {
		if err != nil {
			s.logger.Error(err.Error())
		}
	}()
	if !s.enable {
		err = fmt.Errorf("server not enable")
		return
	}
	// 反序列化获取请求参数
	var transferee raft.Server
	argsErr := ctx.Bind(&transferee)
	if argsErr != nil {
		return fmt.Errorf("读取请求参数失败！%w", argsErr)
	}
	args := raft.TransferLeadership{
		Transferee: transferee,
	}
	// 调用 raft 逻辑
	var res raft.TransferLeadershipReply
	raftErr := s.node.TransferLeadership(args, &res)
	if raftErr != nil {
		return fmt.Errorf("raft 操作失败！%w", raftErr)
	}
	// 序列化并返回结果
	return ctx.JSON(200, res)
}

func (s *server) sleep(ctx echo.Context) error {
	s.enable = false
	return ctx.JSON(200, &struct{ Res string }{Res: "OK"})
}

func (s *server) awake(ctx echo.Context) error {
	s.enable = true
	return ctx.JSON(200, &struct{ Res string }{Res: "OK"})
}
