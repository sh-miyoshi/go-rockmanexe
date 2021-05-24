package dstream

import (
	"context"
	"errors"

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
	return nil, nil
}

func (s *RouterStream) PublishData(authReq *pb.AuthRequest, dataStream pb.Router_PublishDataServer) error {
	// Verify auth request
	// TODO validate version
	c, err := db.GetInst().ClientGetByID(authReq.Id)
	if err != nil {
		logger.Info("Failed to get client: %v", err)
		return errors.New("authenticate failed")
	}
	if c.Key != authReq.Key {
		logger.Info("got invalid key from user")
		return errors.New("authenticate failed")
	}

	// Add to sessionList
	sid, err := session.Add(c.ID, s.sendQueue)
	if err != nil {
		logger.Error("Failed to add session: %v", err)
	}
	logger.Info("add to session %s", sid)

	dataStream.Send(&pb.Data{
		Type: pb.Data_AUTHRESPONSE,
		Data: &pb.Data_AuthRes{
			AuthRes: &pb.AuthResponse{
				Success:   true,
				SessionID: sid,
			},
		},
	})

	session.Run(sid)

	// Publish data
	for {
		data := <-s.sendQueue
		logger.Debug("Send to client: %+v", data)

		dataStream.Send(data)
	}
}
