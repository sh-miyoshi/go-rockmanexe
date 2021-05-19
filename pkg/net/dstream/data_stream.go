package dstream

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/db"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/net/routerpb"
)

const (
	publishInterval = 500 * time.Millisecond // debug
)

type session struct {
	routeID string
	clients [2]string
	// TODO app data
}

type RouterStream struct {
	sessions map[string]*session
}

func New() *RouterStream {
	return &RouterStream{
		sessions: make(map[string]*session),
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
	sid, err := s.addSession(c.ID)
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

	// TODO Publish data
	for {
		// debug
		time.Sleep(publishInterval)
	}
}

func (s *RouterStream) addSession(clientID string) (string, error) {
	route, err := db.GetInst().RouteGetByClient(clientID)
	if err != nil {
		return "", fmt.Errorf("route get failed: %v", err)
	}

	for sid, se := range s.sessions {
		if se.routeID == route.ID {
			if se.clients[0] == "" {
				se.clients[0] = clientID
			} else if se.clients[1] == "" {
				se.clients[1] = clientID
			}
			return sid, nil
		}
	}

	// no session in the list
	// so create new session
	sessionID := uuid.New().String()
	s.sessions[sessionID] = &session{
		routeID: route.ID,
		clients: route.Clients,
	}

	return sessionID, nil
}
