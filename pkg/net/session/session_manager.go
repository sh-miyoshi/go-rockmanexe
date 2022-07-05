package session

import (
	"fmt"
	"time"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/fps"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/effect"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/net/netconnpb"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/object"
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

type sessionError struct {
	generatorClientID string
	reason            error
}

type clientInfo struct {
	chipSent   bool
	clientID   string
	gameInfo   *GameInfo
	dataStream pb.NetConn_TransDataServer
}

type Session struct {
	id        string
	clients   [2]clientInfo
	status    int
	expiresAt time.Time
	dmMgr     *damage.Manager
	cancel    chan struct{}
	exitErr   chan sessionError
}

type SessionManager struct {
	sessions map[string]*Session
}

var (
	inst = &SessionManager{
		sessions: make(map[string]*Session),
	}
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

func (s *Session) UpdateObject(obj object.Object) {
	for i, c := range s.clients {
		isMyObj := c.clientID == obj.ClientID
		s.clients[i].gameInfo.UpdateObject(obj, isMyObj)
	}
}

func (s *Session) RemoveObject(id string) {
	for i := range s.clients {
		s.clients[i].gameInfo.RemoveObject(id)
	}
}

func (s *Session) AddSkill() {
	for i := range s.clients {
		s.clients[i].gameInfo.AddSkill()
	}
}

func (s *Session) AddDamage(dm []damage.Damage) error {
	return s.dmMgr.Add(dm)
}

func (s *Session) AddEffect(eff effect.Effect) {
	for i, c := range s.clients {
		isMyEff := c.clientID == eff.ClientID
		s.clients[i].gameInfo.AddEffect(eff, isMyEff)
	}
}

func (s *Session) SendSignal(clientID string, signal pb.Action_SignalType) error {
	for i, c := range s.clients {
		if c.clientID == clientID {
			switch signal {
			case pb.Action_CHIPSEND:
				s.clients[i].chipSent = true
			case pb.Action_GOCHIPSELECT:
				// TODO
			case pb.Action_PLAYERDEAD:
				// TODO
			}
			return nil
		}
	}

	return fmt.Errorf("no such client %s", clientID)
}

func (s *Session) start() {
	s.expiresAt = time.Now().Add(sessionExpireTime)
	go s.errorHandler()
	go s.frameProc()
	go s.gameInfoPublish()
}

func (s *Session) errorHandler() {
	err := <-s.exitErr
	// TODO publish to clients

	close(s.cancel)
	delete(inst.sessions, s.id)

	if err.reason != nil {
		logger.Error("Got error in session %s: %+v", s.id, err)
	}
}

func (s *Session) frameProc() {
	fpsMgr := fps.Fps{TargetFPS: 60}
	for {
		select {
		case <-s.cancel:
			return
		default:
			if s.status == statusActing {
				// damage process
				for i, c := range s.clients {
					for _, obj := range c.gameInfo.Objects {
						if !obj.Hittable {
							continue
						}

						dmList := []damage.Damage{}
						if dm := s.dmMgr.Hit(c.clientID, obj.ClientID, obj.X, obj.Y); dm != nil {
							dmList = append(dmList, *dm)
							logger.Debug("Hit damage for %s: %+v", c.clientID, dm)
						}
						s.clients[i].gameInfo.AddDamages(dmList)
					}
				}
				s.dmMgr.Update()
			}

			fpsMgr.Wait()
		}
	}
}

func (s *Session) gameInfoPublish() {
	for {
		select {
		case <-s.cancel:
			return
		default:
			now := time.Now()
			before := now.UnixNano() / (1000 * 1000)

			// check session expires
			if s.expiresAt.Before(now) {
				s.exitErr <- sessionError{
					reason: fmt.Errorf("session expired"),
				}
				return
			}

			if err := s.updateGameStatus(); err != nil {
				s.exitErr <- *err
				return
			}

			// publish game info to clients
			for _, c := range s.clients {
				if c.dataStream == nil {
					continue
				}
				c.gameInfo.CurrentTime = time.Now()

				gameInfoBin := c.gameInfo.Marshal()
				err := c.dataStream.Send(&pb.Data{
					Type: pb.Data_DATA,
					Data: &pb.Data_RawData{
						RawData: gameInfoBin,
					},
				})
				if err != nil {
					s.exitErr <- sessionError{
						generatorClientID: c.clientID,
						reason:            fmt.Errorf("failed to send game info: %v", err),
					}
					return
				}

				c.gameInfo.Cleanup()
			}

			after := time.Now().UnixNano() / (1000 * 1000)
			time.Sleep(publishInterval - time.Duration(after-before))
		}
	}
}

func (s *Session) updateGameStatus() *sessionError {
	switch s.status {
	case statusConnectWait:
		for _, c := range s.clients {
			if c.clientID == "" {
				return nil
			}
		}

		// Initialize panel info
		s.clients[0].gameInfo.InitPanel(s.clients[0].clientID, s.clients[1].clientID)
		s.clients[1].gameInfo.InitPanel(s.clients[1].clientID, s.clients[0].clientID)

		if err := s.sendStatusToClients(pb.Data_CHIPSELECTWAIT); err != nil {
			return err
		}
		s.changeStatus(statusChipSelectWait)
	case statusChipSelectWait:
		for _, c := range s.clients {
			if !c.chipSent {
				return nil
			}
		}

		if err := s.sendStatusToClients(pb.Data_ACTING); err != nil {
			return err
		}
		for i := range s.clients {
			s.clients[i].chipSent = false
		}
		s.changeStatus(statusActing)
	case statusActing:
		// TODO
	case statusGameEnd:
		// TODO
	}

	return nil
}

func (s *Session) changeStatus(next int) {
	logger.Info("Change state from %d to %d", s.status, next)
	s.status = next
}

func (s *Session) sendStatusToClients(st pb.Data_Status) *sessionError {
	for _, c := range s.clients {
		if c.dataStream == nil {
			continue
		}

		err := c.dataStream.Send(&pb.Data{
			Type: pb.Data_UPDATESTATUS,
			Data: &pb.Data_Status_{
				Status: st,
			},
		})
		if err != nil {
			return &sessionError{
				generatorClientID: c.clientID,
				reason:            fmt.Errorf("failed to send status to client %s: %v", c.clientID, err),
			}
		}
	}
	return nil
}
