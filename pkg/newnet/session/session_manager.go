package session

import (
	"fmt"
	"time"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/fps"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/newnet/netconnpb"
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

type clientInfo struct {
	clientID   string
	dataStream pb.NetConn_TransDataServer
}

type session struct {
	id        string
	clients   [2]clientInfo
	info      *GameInfo
	status    int
	expiresAt time.Time
	cancel    chan struct{}
}

type SessionManager struct {
	sessions map[string]*session
}

var (
	inst = &SessionManager{
		sessions: make(map[string]*session),
	}
)

func Add(sessionID, clientID string, stream pb.NetConn_TransDataServer) error {
	s, ok := inst.sessions[sessionID]
	if ok {
		// check exists
		for i, c := range s.clients {
			if c.clientID == "" {
				s.clients[i] = clientInfo{
					clientID:   clientID,
					dataStream: stream,
				}
				s.start()
				logger.Info("set new client %s to session %s", clientID, sessionID)
				return nil
			}
			if c.clientID == clientID {
				return fmt.Errorf("the client already added")
			}
		}
	} else {
		c := clientInfo{
			clientID:   clientID,
			dataStream: stream,
		}
		inst.sessions[sessionID] = &session{
			id:      sessionID,
			clients: [2]clientInfo{c},
			info:    &GameInfo{},
			status:  statusConnectWait,
			cancel:  make(chan struct{}),
		}
		logger.Info("create new session %s for client %s", sessionID, clientID)
	}
	return nil
}

func GetGameInfo(sessionID string) *GameInfo {
	s, ok := inst.sessions[sessionID]
	if !ok {
		return nil
	}

	return s.info
}

func (s *session) start() {
	s.status = statusChipSelectWait
	s.expiresAt = time.Now().Add(sessionExpireTime)
	go s.frameProc()
	go s.gameInfoPublish()
}

func (s *session) frameProc() {
	fpsMgr := fps.Fps{TargetFPS: 60}
	for {
		select {
		case <-s.cancel:
			return
		default:
			// TODO

			fpsMgr.Wait()
		}
	}
}

func (s *session) gameInfoPublish() {
	for {
		select {
		case <-s.cancel:
			return
		default:
			now := time.Now()
			before := now.UnixNano() / (1000 * 1000)

			// check session expires
			if s.expiresAt.Before(now) {
				// TODO publish to clients
				s.cancel <- struct{}{}
				return
			}

			// TODO publish to clients

			after := time.Now().UnixNano() / (1000 * 1000)
			time.Sleep(publishInterval - time.Duration(after-before))
		}
	}
}
