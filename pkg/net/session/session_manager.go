package session

import (
	"errors"
	"fmt"
	"time"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/damage"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/net/netconnpb"
)

const (
	statusConnectWait int = iota
	statusChipSelectWait
	statusActing
	statusGameEnd
)

const (
	publishInterval   = 100 * time.Millisecond // debug
	sessionExpireTime = 30 * time.Minute
)

type SessionManager struct {
	sessions map[string]*Session
}

var (
	inst = &SessionManager{
		sessions: make(map[string]*Session),
	}
	errSendFailed = errors.New("data send failed")
)

func Add(sessionID, clientID string, stream pb.NetConn_TransDataServer) error {
	s, ok := inst.sessions[sessionID]
	if ok {
		// check exists
		for i, c := range s.clients {
			if c.clientID == "" {
				s.clients[i].clientID = clientID
				s.clients[i].dataStream = stream
				logger.Info("set new client %s to session %s", clientID, sessionID)
				return nil
			}
			if c.clientID == clientID {
				return fmt.Errorf("the client already added")
			}
		}
	} else {
		inst.sessions[sessionID] = &Session{
			id: sessionID,
			clients: [2]clientInfo{
				{
					chipSent:   false,
					clientID:   clientID,
					gameInfo:   NewGameInfo(),
					dataStream: stream,
				},
				{
					chipSent:   false,
					clientID:   "",
					gameInfo:   NewGameInfo(),
					dataStream: nil,
				},
			},
			status:  statusConnectWait,
			cancel:  make(chan struct{}),
			exitErr: make(chan sessionError),
			dmMgr:   &damage.Manager{},
		}

		inst.sessions[sessionID].start()
		logger.Info("create new session %s for client %s", sessionID, clientID)
	}
	return nil
}

func GetSession(sessionID string) *Session {
	s, ok := inst.sessions[sessionID]
	if !ok {
		return nil
	}
	return s
}
