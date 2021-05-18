package dstream

import (
	"context"

	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/net/routerpb"
)

type RouterStream struct {
}

func (s *RouterStream) SendAction(ctx context.Context, action *pb.Action) (*pb.Result, error) {
	return nil, nil
}

func (s *RouterStream) PublishData(*pb.AuthRequest, pb.Router_PublishDataServer) error {
	return nil
}
