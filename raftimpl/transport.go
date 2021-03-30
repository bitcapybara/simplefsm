package raftimpl

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/bitcapybara/raft"
	"net/http"
)

type httpTransport struct {
}

func NewHttpTransport() *httpTransport {
	return &httpTransport{}
}

func (h *httpTransport) AppendEntries(addr raft.NodeAddr, args raft.AppendEntry, res *raft.AppendEntryReply) error {
	return httpClient(string(addr)+"/appendEntries", args, res)
}

func (h *httpTransport) RequestVote(addr raft.NodeAddr, args raft.RequestVote, res *raft.RequestVoteReply) error {
	return httpClient(string(addr)+"/requestVote", args, res)
}

func (h *httpTransport) InstallSnapshot(addr raft.NodeAddr, args raft.InstallSnapshot, res *raft.InstallSnapshotReply) error {
	return httpClient(string(addr)+"/installSnapshot", args, res)
}

func httpClient(url string, args interface{}, res interface{}) error {
	// 对请求参数编码
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	enErr := encoder.Encode(args)
	if enErr != nil {
		return fmt.Errorf("序列化失败！%w", enErr)
	}
	// 发送请求
	postRes, httpErr := http.Post(url, "application/octet-stream", buffer)
	if httpErr != nil {
		return fmt.Errorf("发送请求失败！%w", httpErr)
	}
	if postRes.StatusCode != 200 {
		return fmt.Errorf("发送请求响应码异常：%d", postRes.StatusCode)
	}
	body := postRes.Body
	decoder := gob.NewDecoder(body)
	deErr := decoder.Decode(res)
	if deErr != nil {
		return fmt.Errorf("反序列化失败！%w", deErr)
	}
	return nil
}



