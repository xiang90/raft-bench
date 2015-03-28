package main

import "github.com/coreos/etcd/raft/raftpb"

type transport interface {
	Send(msgs []raftpb.Message)
}

type noopTrans struct{}

func (nt *noopTrans) Send(msgs []raftpb.Message) {}
