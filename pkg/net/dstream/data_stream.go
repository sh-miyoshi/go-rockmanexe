package dstream

import (
	"context"
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	api "github.com/sh-miyoshi/go-rockmanexe/pkg/net/apiclient"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/net/routerpb"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/session"
)

type RouterStream struct {
}

func New() *RouterStream {
	return &RouterStream{}
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

	caRes, err := api.ClientAuth(authReq.Id, authReq.Key)
	if err != nil {
		logger.Info("Failed to authenticate client: %v", err)
		authRes.Success = false
		authRes.ErrMsg = "authenticate failed"
		dataStream.Send(makeAuthRes(authRes))
		return nil
	}

	sinfo, err := api.GetSessionInfo(caRes.SessionID)
	if err != nil {
		logger.Info("Failed to get session: %v", err)
		authRes.Success = false
		authRes.ErrMsg = "internal server error"
		dataStream.Send(makeAuthRes(authRes))
		return nil
	}

	// Add to sessionList
	sid, err := session.Add(authReq.Id, *sinfo, dataStream)
	if err != nil {
		logger.Error("Failed to add session: %v", err)
		return fmt.Errorf("add session failed: %w", err)
	}
	logger.Info("add to session %s for client %s", sid, authReq.Id)

	authRes.AllUserIDs = []string{
		fmt.Sprintf("%s:%s", sinfo.OwnerClientID, sinfo.OwnerUserID),
		fmt.Sprintf("%s:%s", sinfo.GuestClientID, sinfo.GuestUserID),
	}

	authRes.Success = true
	authRes.SessionID = sid
	dataStream.Send(makeAuthRes(authRes))

	session.Run(sid)

	return nil
}

func makeAuthRes(authRes *pb.AuthResponse) *pb.Data {
	return &pb.Data{
		Type: pb.Data_AUTHRESPONSE,
		Data: &pb.Data_AuthRes{
			AuthRes: authRes,
		},
	}
}
