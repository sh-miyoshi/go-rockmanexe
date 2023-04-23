package queue

import (
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/net/netconnpb"
)

var (
	allActionQueues = make(map[string][]*pb.Request_Action)
)

func Push(id string, act *pb.Request_Action) {
	allActionQueues[id] = append(allActionQueues[id], act)
}

func Pop(id string) *pb.Request_Action {
	acts, ok := allActionQueues[id]
	if !ok || len(acts) == 0 {
		return nil
	}
	res := acts[0]
	allActionQueues[id] = acts[1:]
	return res
}

func Delete(id string) {
	delete(allActionQueues, id)
}
