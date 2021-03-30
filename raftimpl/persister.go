package raftimpl

import (
	"github.com/bitcapybara/raft"
	"sync"
)

// SnapshotPersister 接口的内存实现
type snapshotPersister struct {
	snapshot raft.Snapshot
	mu       sync.Mutex
}

func NewSnapshotPersister() *snapshotPersister {
	return &snapshotPersister{}
}

func (ps *snapshotPersister) SaveSnapshot(snapshot raft.Snapshot) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.snapshot = snapshot
	return nil
}

func (ps *snapshotPersister) LoadSnapshot() (raft.Snapshot, error) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	return ps.snapshot, nil
}

// RaftStatePersister 接口的内存实现，开发测试用
type raftStatePersister struct {
	raftState raft.RaftState
	mu        sync.Mutex
}

func NewRaftStatePersister() *raftStatePersister {
	return &raftStatePersister{
		raftState: raft.RaftState{
			Term:     0,
			VotedFor: "",
			Entries:  make([]raft.Entry, 0),
		},
	}
}

func (ps *raftStatePersister) SaveRaftState(state raft.RaftState) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.raftState = state
	return nil
}

func (ps *raftStatePersister) LoadRaftState() (raft.RaftState, error) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	return ps.raftState, nil
}
