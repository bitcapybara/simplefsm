package raftimpl

import (
	"fmt"
	"github.com/bitcapybara/raft"
	"github.com/go-resty/resty/v2"
)

type httpTransport struct {
	logger raft.Logger
	client *resty.Client
}

func NewHttpTransport(logger raft.Logger) *httpTransport {
	return &httpTransport{logger: logger, client: resty.New()}
}

func (h *httpTransport) AppendEntries(addr raft.NodeAddr, args raft.AppendEntry, res *raft.AppendEntryReply) error {
	url := fmt.Sprintf("%s%s%s", "http://", addr, "/appendEntries")
	// 发送请求
	response, resErr := h.client.R().SetHeader("Content-Type", "application/json").SetBody(args).SetResult(res).Post(url)
	if resErr != nil {
		return fmt.Errorf("发送请求失败！%w", resErr)
	}
	if response.StatusCode() != 200 {
		return fmt.Errorf("发送请求响应码异常：%d", response.StatusCode())
	}
	return nil
}

func (h *httpTransport) RequestVote(addr raft.NodeAddr, args raft.RequestVote, res *raft.RequestVoteReply) error {
	url := fmt.Sprintf("%s%s%s", "http://", addr, "/requestVote")
	// 发送请求
	response, resErr := h.client.R().SetHeader("Content-Type", "application/json").SetBody(args).SetResult(res).Post(url)
	if resErr != nil {
		return fmt.Errorf("发送请求失败！%w", resErr)
	}
	if response.StatusCode() != 200 {
		return fmt.Errorf("发送请求响应码异常：%d", response.StatusCode())
	}
	return nil
}

func (h *httpTransport) InstallSnapshot(addr raft.NodeAddr, args raft.InstallSnapshot, res *raft.InstallSnapshotReply) error {
	url := fmt.Sprintf("%s%s%s", "http://", addr, "/installSnapshot")
	// 发送请求
	response, resErr := h.client.R().SetHeader("Content-Type", "application/json").SetBody(args).SetResult(res).Post(url)
	if resErr != nil {
		return fmt.Errorf("发送请求失败！%w", resErr)
	}
	if response.StatusCode() != 200 {
		return fmt.Errorf("发送请求响应码异常：%d", response.StatusCode())
	}
	return nil
}



