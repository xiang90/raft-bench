package main

import (
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
	"github.com/coreos/etcd/raft"
	"github.com/coreos/etcd/raft/raftpb"
	"github.com/coreos/etcd/wal"
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
	dir, err := ioutil.TempDir("", "raft-bench")
	if err != nil {
		log.Fatalf("raft-bench: cannot create dir for wal (%v)", err)
	}
	defer os.RemoveAll(dir)

	w, err := wal.Create(dir, nil)
	if err != nil {
		log.Fatalf("raft-bench: create wal error: %v", err)
	}
	defer w.Close()

	for {
		select {
		case <-n.ticker:
			n.Tick()
		case rd := <-n.Node.Ready():
			w.Save(rd.HardState, rd.Entries)
			n.raftStorage.Append(rd.Entries)

			n.trans.Send(rd.Messages)
			for i := range rd.CommittedEntries {
				if rd.CommittedEntries[i].Index == n.goal {
					n.reach <- struct{}{}
				}
			}
			n.Node.Advance()
			// batch for 0.5ms
			time.Sleep(500 * time.Microsecond)
		case <-n.done:
			return
		}
	}
}

func (n *raftNode) Process(ctx context.Context, m raftpb.Message) error {
	return n.Step(ctx, m)
}
