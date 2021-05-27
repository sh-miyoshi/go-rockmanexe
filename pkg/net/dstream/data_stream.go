package dstream

import (
	"context"
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/db"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/net/routerpb"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/session"
)

type RouterStream struct {
	sendQueue chan *pb.Data
}

func New() *RouterStream {
	return &RouterStream{
		sendQueue: make(chan *pb.Data),
	}
}

func (s *RouterStream) SendAction(ctx context.Context, action *pb.Action) (*pb.Result, error) {
	if err := session.ActionProc(action); err != nil {
		return &pb.Result{
			Success: false,
			ErrMsg:  fmt.Sprintf("request failed: %v", err),
		}, nil
	}

	return &pb.Result{
		Success: true,
	}, nil
}

func (s *RouterStream) PublishData(authReq *pb.AuthRequest, dataStream pb.Router_PublishDataServer) error {
	// Verify auth request
	authRes := &pb.AuthResponse{}

	// TODO validate version
	c, err := db.GetInst().ClientGetByID(authReq.Id)
	if err != nil {
		logger.Info("Failed to get client: %v", err)
		authRes.Success = false
		authRes.ErrMsg = "authenticate failed"
		dataStream.Send(makeAuthRes(authRes))
		return nil
	}
	if c.Key != authReq.Key {
		logger.Info("got invalid key from user")
		authRes.Success = false
		authRes.ErrMsg = "authenticate failed"
		dataStream.Send(makeAuthRes(authRes))
		return nil
	}

	// Add to sessionList
	sid, err := session.Add(c.ID, s.sendQueue)
	if err != nil {
		logger.Error("Failed to add session: %v", err)
		return fmt.Errorf("add session failed: %w", err)
	}
	logger.Info("add to session %s", sid)

	authRes.Success = true
	authRes.SessionID = sid
	dataStream.Send(makeAuthRes(authRes))

	session.Run(sid)

	// Publish data
	for {
		data := <-s.sendQueue
		logger.Debug("Send to client %s: %+v", c.ID, data)

		dataStream.Send(data)
	}
}

func makeAuthRes(authRes *pb.AuthResponse) *pb.Data {
	return &pb.Data{
		Type: pb.Data_AUTHRESPONSE,
		Data: &pb.Data_AuthRes{
			AuthRes: authRes,
		},
	}
}
