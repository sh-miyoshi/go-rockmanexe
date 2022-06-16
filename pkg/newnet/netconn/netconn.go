package netconn

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	api "github.com/sh-miyoshi/go-rockmanexe/pkg/newnet/apiclient"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/newnet/netconnpb"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/newnet/object"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/newnet/session"
)

type NetConn struct {
}

func New() *NetConn {
	return &NetConn{}
}

func (s *NetConn) TransData(stream pb.NetConn_TransDataServer) error {
	for {
		action, err := stream.Recv()
		if err != nil {
			logger.Info("Failed to recv from client: %w", err)
			return nil
		}

		if action.GetType() == pb.Action_AUTHENTICATE {
			res := authClient(action.GetReq(), stream)
			stream.Send(makeAuthRes(&res))
			if !res.Success {
				return nil
			}
			continue
		}

		sessionID := action.GetSessionID()
		s := session.GetSession(sessionID)
		if s == nil {
			logger.Info("No such session: %s", sessionID)
			return fmt.Errorf("failed to get session info for %s", sessionID)
		}

		switch action.GetType() {
		case pb.Action_UPDATEOBJECT:
			var obj object.Object
			object.Unmarshal(&obj, action.GetObjectInfo())
			s.UpdateObject(obj)
		case pb.Action_REMOVEOBJECT:
			var obj object.Object
			object.Unmarshal(&obj, action.GetObjectInfo())
			s.RemoveObject(obj.ID)
		case pb.Action_ADDSKILL:
			s.AddSkill()
		case pb.Action_ADDDAMAGE:
			s.AddDamage()
		case pb.Action_ADDEFFECT:
			s.AddEffect()
		case pb.Action_SENDSIGNAL:
			s.SendSignal(action.GetClientID(), action.GetSignal())
		default:
			return fmt.Errorf("invalid action type: %v", action.GetType())
		}
	}
}

func authClient(authReq *pb.Action_AuthRequest, stream pb.NetConn_TransDataServer) pb.Data_AuthResponse {
	if err := api.VersionCheck(authReq.Version); err != nil {
		logger.Info("Got missmatched version: %v", err)
		return pb.Data_AuthResponse{
			Success: false,
			ErrMsg:  err.Error(),
		}
	}

	caRes, err := api.ClientAuth(authReq.Id, authReq.Key)
	if err != nil {
		logger.Info("Failed to authenticate client: %v", err)
		return pb.Data_AuthResponse{
			Success: false,
			ErrMsg:  "authenticate failed",
		}
	}

	sinfo, err := api.GetSessionInfo(caRes.SessionID)
	if err != nil {
		logger.Error("Failed to get session: %v", err)
		return pb.Data_AuthResponse{
			Success: false,
			ErrMsg:  "internal server error",
		}
	}

	if err := session.Add(sinfo.ID, authReq.Id, stream); err != nil {
		logger.Error("Failed to add to session manager: %v", err)
		return pb.Data_AuthResponse{
			Success: false,
			ErrMsg:  "internal server error",
		}
	}

	return pb.Data_AuthResponse{
		Success:   true,
		SessionID: sinfo.ID,
		AllUserIDs: []string{
			fmt.Sprintf("%s:%s", sinfo.OwnerClientID, sinfo.OwnerUserID),
			fmt.Sprintf("%s:%s", sinfo.GuestClientID, sinfo.GuestUserID),
		},
	}
}

func makeAuthRes(authRes *pb.Data_AuthResponse) *pb.Data {
	return &pb.Data{
		Type: pb.Data_AUTHRESPONSE,
		Data: &pb.Data_AuthRes{
			AuthRes: authRes,
		},
	}
}
