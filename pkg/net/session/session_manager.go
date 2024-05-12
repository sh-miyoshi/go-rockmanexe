package session

import (
	"fmt"
	"sync"
	"time"

	"github.com/cockroachdb/errors"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/system"
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
	mu                sync.Mutex
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
	inst.mu.Lock()
	defer inst.mu.Unlock()

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

func ManagerExec() error {
	for {
		before := time.Now().UnixNano() / (1000 * 1000)

		for key, s := range inst.sessions {
			if s.IsEnd() {
				s.End()
				delete(inst.sessions, key)
			}
		}

		// できる限りsessionを終了させてからrouterは死ぬ
		if len(inst.sessions) == 0 {
			if err := system.Error(); err != nil {
				return fmt.Errorf("failed to exec system: %+w", err)
			}
		}

		after := time.Now().UnixNano() / (1000 * 1000)
		time.Sleep(sessionCheckInterval - time.Duration(after-before))
	}
}
