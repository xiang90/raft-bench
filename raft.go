package main

import (
	"time"

	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
	"github.com/coreos/etcd/raft"
	"github.com/coreos/etcd/raft/raftpb"
)

var (
	tickDuration = 100 * time.Millisecond
)

type raftNode struct {
	raft.Node

	ticker      <-chan time.Time
	raftStorage *raft.MemoryStorage
	trans       transport

	goal  uint64
	reach chan struct{}
	done  chan struct{}
}

func (n *raftNode) run() {
	for {
		select {
		case <-n.ticker:
			n.Tick()
		case rd := <-n.Node.Ready():
			// save to WAL
			time.Sleep(time.Millisecond)
			n.raftStorage.Append(rd.Entries)
			n.trans.Send(rd.Messages)
			for i := range rd.CommittedEntries {
				// var r etcdserverpb.Request
				// err := r.Unmarshal(rd.CommittedEntries[i].Data)
				// if err != nil {
				// 	log.Fatal("err unmarshal ", err)
				// }
				if rd.CommittedEntries[i].Index == n.goal {
					n.reach <- struct{}{}
				}
			}
			n.Node.Advance()
		case <-n.done:
			return
		}
	}
}

func (n *raftNode) Process(ctx context.Context, m raftpb.Message) error {
	return n.Step(ctx, m)
}
