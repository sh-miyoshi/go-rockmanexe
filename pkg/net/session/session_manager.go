package session

import (
	"errors"
	"fmt"
	"time"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/net/netconnpb"
)

const (
	sessionExpireTime    = 30 * time.Minute
	sessionCheckInterval = 1 * time.Second
)

type SessionManager struct {
	sessions          map[string]*Session
	gameLogicGenerate func() GameLogic
}

var (
	inst = &SessionManager{
		sessions: make(map[string]*Session),
	}
	errSendFailed = errors.New("data send failed")
)

func SetLogicGenerator(g func() GameLogic) {
	inst.gameLogicGenerate = g
}

func Add(sessionID, clientID string, stream pb.NetConn_TransDataServer) error {
	s, ok := inst.sessions[sessionID]
	if ok {
		if err := s.SetClient(clientID, stream); err != nil {
			return fmt.Errorf("failed to set client: %w", err)
		}
		logger.Info("set new client %s to session %s", clientID, sessionID)
	} else {
		handler := inst.gameLogicGenerate()
		inst.sessions[sessionID] = newSession(sessionID, handler)
		inst.sessions[sessionID].SetClient(clientID, stream)

		go inst.sessions[sessionID].Run()
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