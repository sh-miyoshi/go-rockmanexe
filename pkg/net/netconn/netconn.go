package netconn

import (
	"fmt"

	"github.com/cockroachdb/errors"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	api "github.com/sh-miyoshi/go-rockmanexe/pkg/net/api/client"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/net/netconnpb"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/session"
)

type NetConn struct {
}

func New() *NetConn {
	return &NetConn{}
}

func (n *NetConn) TransData(stream pb.NetConn_TransDataServer) error {
	for {
		msg, err := stream.Recv()
		if err != nil {
			logger.Info("Failed to recv from client: %+v", err)
			return nil
		}
		logger.Debug("Got msg: %+v", msg)

		if msg.GetType() == pb.Request_AUTHENTICATE {
			res := authClient(msg.GetReq(), stream)
			stream.Send(makeAuthRes(&res))
			if !res.Success {
				return nil
			}
			continue
		}

		s := session.GetSession(msg.GetSessionID())
		if s == nil {
			logger.Info("No such session: %s", msg.GetSessionID())
			return errors.Newf("failed to get session info for %s", msg.GetSessionID())
		}

		switch msg.GetType() {
		case pb.Request_SENDSIGNAL:
			if err := s.HandleSignal(msg.GetClientID(), msg.GetSignal()); err != nil {
				logger.Error("Failed to send signal %v: %+v", msg.GetSignal(), err)
				return errors.Wrap(err, "failed to send signal")
			}
		case pb.Request_ACTION:
			if err := s.HandleAction(msg.GetClientID(), msg.GetAct()); err != nil {
				logger.Error("Failed to handle action %v: %+v", msg.GetAct(), err)
				return errors.Wrap(err, "failed to handle action")
			}
		default:
			return errors.Newf("invalid message type: %v", msg.GetType())
		}

		// TODO_Next return current status
	}
}

func authClient(authReq *pb.Request_AuthRequest, stream pb.NetConn_TransDataServer) pb.Response_AuthResponse {
	if err := api.VersionCheck(authReq.Version); err != nil {
		logger.Info("Got missmatched version: %v", err)
		return pb.Response_AuthResponse{
			Success: false,
			ErrMsg:  err.Error(),
		}
	}

	caRes, err := api.ClientAuth(authReq.Id, authReq.Key)
	if err != nil {
		logger.Info("Failed to authenticate client: %v", err)
		return pb.Response_AuthResponse{
			Success: false,
			ErrMsg:  "authenticate failed",
		}
	}

	sinfo, err := api.GetSessionInfo(caRes.SessionID)
	if err != nil {
		logger.Error("Failed to get session: %v", err)
		return pb.Response_AuthResponse{
			Success: false,
			ErrMsg:  "internal server error",
		}
	}

	if err := session.Add(sinfo.ID, authReq.Id, stream); err != nil {
		logger.Error("Failed to add to session manager: %v", err)
		return pb.Response_AuthResponse{
			Success: false,
			ErrMsg:  "internal server error",
		}
	}

	return pb.Response_AuthResponse{
		Success:   true,
		SessionID: sinfo.ID,
		AllUserIDs: []string{
			fmt.Sprintf("%s:%s", sinfo.OwnerClientID, sinfo.OwnerUserID),
			fmt.Sprintf("%s:%s", sinfo.GuestClientID, sinfo.GuestUserID),
		},
	}
}

func makeAuthRes(authRes *pb.Response_AuthResponse) *pb.Response {
	return &pb.Response{
		Type: pb.Response_AUTHRESPONSE,
		Data: &pb.Response_AuthRes{
			AuthRes: authRes,
		},
	}
}
