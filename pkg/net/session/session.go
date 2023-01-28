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
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/sound"
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
	id         string
	clients    [2]clientInfo
	status     int
	nextStatus int
	expiresAt  time.Time
	dmMgr      *damage.Manager
	exitErr    *sessionError
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

func (s *Session) AddDamage(dm []damage.Damage) error {
	return s.dmMgr.Add(dm)
}

func (s *Session) AddEffect(eff effect.Effect) {
	for i, c := range s.clients {
		isMyEff := c.clientID == eff.ClientID
		s.clients[i].gameInfo.AddEffect(eff, isMyEff)
	}
}

func (s *Session) AddSound(se sound.Sound) {
	for i, c := range s.clients {
		// 自分のSEはローカルで流すので相手に流れてほしいSEのみを追加する
		if c.clientID != se.ClientID {
			s.clients[i].gameInfo.AddSound(se)
		}
	}
}

func (s *Session) SendSignal(clientID string, signal pb.Action_SignalType) error {
	for i, c := range s.clients {
		if c.clientID == clientID {
			switch signal {
			case pb.Action_CHIPSEND:
				s.clients[i].chipSent = true
			case pb.Action_GOCHIPSELECT:
				s.nextStatus = statusChipSelectWait
			case pb.Action_PLAYERDEAD:
				s.nextStatus = statusGameEnd
			default:
				return fmt.Errorf("got unexpected signal type: %d", signal)
			}
			return nil
		}
	}

	return fmt.Errorf("no such client %s", clientID)
}

func (s *Session) Run() {
	// TODO: 要調整
	go s.frameProc()
	go s.gameInfoPublish()
}

func (s *Session) IsEnd() bool {
	return s.exitErr != nil
}

func (s *Session) End() {
	if s.exitErr.reason != nil {
		if s.exitErr.reason == errSendFailed {
			for _, c := range s.clients {
				if c.dataStream == nil || c.clientID == s.exitErr.generatorClientID {
					continue
				}

				// publish to alive clients
				c.dataStream.Send(&pb.Data{
					Type: pb.Data_UPDATESTATUS,
					Data: &pb.Data_Status_{
						Status: pb.Data_GAMEEND,
					},
				})
			}
		}
		logger.Error("Got error in session %s: %+v", s.id, s.exitErr.reason)
	}
}

func (s *Session) frameProc() {
	fpsMgr := fps.Fps{TargetFPS: 60}
	for {
		if s.exitErr != nil {
			return
		}

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

func (s *Session) gameInfoPublish() {
	for {
		if s.exitErr != nil {
			return
		}

		now := time.Now()
		before := now.UnixNano() / (1000 * 1000)

		// check session expires
		if s.expiresAt.Before(now) {
			s.exitErr = &sessionError{
				reason: fmt.Errorf("session expired"),
			}
			return
		}

		if err := s.updateGameStatus(); err != nil {
			s.exitErr = err
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
				logger.Error("failed to send game info to client %s: %v", c.clientID, err)
				s.exitErr = &sessionError{
					generatorClientID: c.clientID,
					reason:            errSendFailed,
				}
				return
			}

			c.gameInfo.Cleanup()
		}

		after := time.Now().UnixNano() / (1000 * 1000)
		time.Sleep(publishInterval - time.Duration(after-before))
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
		s.nextStatus = -1
		s.changeStatus(statusActing)
	case statusActing:
		switch s.nextStatus {
		case statusChipSelectWait:
			if err := s.sendStatusToClients(pb.Data_CHIPSELECTWAIT); err != nil {
				return err
			}
			s.changeStatus(statusChipSelectWait)
		case statusGameEnd:
			if err := s.sendStatusToClients(pb.Data_GAMEEND); err != nil {
				return err
			}
			s.changeStatus(statusGameEnd)
		}
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
			logger.Error("failed to send status to client %s: %v", c.clientID, err)
			return &sessionError{
				generatorClientID: c.clientID,
				reason:            errSendFailed,
			}
		}
	}
	return nil
}
