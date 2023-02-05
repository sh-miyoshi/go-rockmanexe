package session

import (
	"errors"
	"fmt"
	"time"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/oldnet/damage"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/oldnet/netconnpb"
)

const (
	statusConnectWait int = iota
	statusChipSelectWait
	statusActing
	statusGameEnd
)

const (
	publishInterval      = 100 * time.Millisecond // debug
	sessionExpireTime    = 30 * time.Minute
	sessionCheckInterval = 1 * time.Second
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
			status:    statusConnectWait,
			dmMgr:     &damage.Manager{},
			expiresAt: time.Now().Add(sessionExpireTime),
		}

		inst.sessions[sessionID].Run()
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

func ManagerExec() {
	for {
		before := time.Now().UnixNano() / (1000 * 1000)

		for key, s := range inst.sessions {
			if s.IsEnd() {
				s.End()
				delete(inst.sessions, key)
			}
		}

		after := time.Now().UnixNano() / (1000 * 1000)
		time.Sleep(sessionCheckInterval - time.Duration(after-before))
	}
}
