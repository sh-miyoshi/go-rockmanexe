package session

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/db"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/net/routerpb"
)

const (
	publishInterval = 5 * time.Second // debug
)

type clientInfo struct {
	active    bool
	clientID  string
	sendQueue chan pb.Data
}

type session struct {
	sessionID string
	routeID   string
	clients   [2]clientInfo

	// status  int
	// // TODO app data
}

var (
	sessionList = []*session{}
)

func Add(clientID string, sendQueue chan pb.Data) (string, error) {
	route, err := db.GetInst().RouteGetByClient(clientID)
	if err != nil {
		return "", fmt.Errorf("route get failed: %v", err)
	}

	for _, se := range sessionList {
		if se.routeID == route.ID {
			if se.clients[0].clientID == clientID {
				se.clients[0].active = true
				se.clients[0].sendQueue = sendQueue
			} else if se.clients[1].clientID == clientID {
				se.clients[1].active = true
				se.clients[1].sendQueue = sendQueue
			}

			return se.sessionID, nil
		}
	}

	// no session in the list
	// so create new session
	sessionID := uuid.New().String()
	v := session{
		sessionID: sessionID,
		routeID:   route.ID,
	}
	v.clients[0] = clientInfo{
		active:    true,
		clientID:  route.Clients[0],
		sendQueue: sendQueue,
	}
	v.clients[1] = clientInfo{
		clientID: route.Clients[1],
	}
	logger.Debug("new session info: %+v", v)

	sessionList = append(sessionList, &v)

	go v.Process()

	return sessionID, nil
}

func (s *session) Process() {
	logger.Info("start new session for route %s", s.routeID)
	logger.Debug("client info: %+v", s.clients)

	// publish via data stream
	for {
		time.Sleep(publishInterval)
		for _, c := range s.clients {
			if c.active {
				logger.Debug("Send data to client %s", c.clientID)
				c.sendQueue <- pb.Data{
					Type: pb.Data_DATA,
					Data: &pb.Data_RawData{
						RawData: []byte("test"), // debug
					},
				}
			}
		}
	}
}
