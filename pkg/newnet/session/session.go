package session

import (
	"fmt"
	"time"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/fps"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/newnet/action"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/newnet/netconnpb"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/newnet/object"
)

const (
	stateConnectWait int = iota
	stateChipSelectWait
	stateActing
	stateGameEnd
)

type GameLogic interface {
	Init(clientIDs [2]string) error
	AddObject(clientID string, param object.InitParam)
	MoveObject(moveInfo action.Move)
	AddBuster(clientID string, busterInfo action.Buster)
	UseChip(clientID string, chipInfo action.UseChip)
	GetInfo() []byte
	UpdateGameStatus()
}

type sessionError struct {
	generatorClientID string
	reason            error
}

type SessionClient struct {
	clientID   string
	dataStream pb.NetConn_TransDataServer
	chipSent   bool
}

type Session struct {
	id          string
	clients     [2]SessionClient
	expiresAt   time.Time
	gameHandler GameLogic
	exitErr     *sessionError
	fpsMgr      fps.Fps
	state       int

	// TODO
	/*
		go clientからデータを受け取る -> gameInfoに反映
		go fpsごとにgameInfoの計算
		go 一定時間ごとにclientに送信(※)
	*/
}

func newSession(sessionID string, gameHandler GameLogic) *Session {
	res := &Session{
		id:          sessionID,
		expiresAt:   time.Now().Add(sessionExpireTime),
		fpsMgr:      fps.Fps{TargetFPS: 60},
		state:       stateConnectWait,
		gameHandler: gameHandler,
	}
	return res
}

func (s *Session) SetClient(clientID string, stream pb.NetConn_TransDataServer) error {
	for i := 0; i < len(s.clients); i++ {
		if s.clients[i].clientID == "" {
			s.clients[i].clientID = clientID
			s.clients[i].dataStream = stream
			return nil
		} else if s.clients[i].clientID == clientID {
			return fmt.Errorf("the client already added")
		}
	}
	return fmt.Errorf("session is full")
}

func (s *Session) Run() {
MAIN_LOOP:
	for {
		if s.exitErr != nil {
			return
		}

		now := time.Now()

		// check session expires
		if s.expiresAt.Before(now) {
			s.exitErr = &sessionError{
				reason: fmt.Errorf("session expired"),
			}
			return
		}

		switch s.state {
		case stateConnectWait:
			for _, c := range s.clients {
				if c.clientID == "" {
					continue MAIN_LOOP
				}
			}

			clientIDs := [2]string{}
			for i := 0; i < len(s.clients); i++ {
				s.clients[i].chipSent = false
				clientIDs[i] = s.clients[i].clientID
			}
			if err := s.gameHandler.Init(clientIDs); err != nil {
				s.exitErr = &sessionError{
					reason: fmt.Errorf("failed to initialize game handler"),
				}
				return
			}
			s.changeState(stateChipSelectWait)
			s.publishStateToClient(pb.Response_CHIPSELECTWAIT)
		case stateChipSelectWait:
			for _, c := range s.clients {
				if !c.chipSent {
					continue MAIN_LOOP
				}
			}

			for i := 0; i < len(s.clients); i++ {
				s.clients[i].chipSent = false
			}
			s.changeState(stateActing)
			s.publishStateToClient(pb.Response_ACTING)
		case stateActing:
			s.gameHandler.UpdateGameStatus()
			s.publishGameInfo() // debug(送信頻度は要確認)

			// TODO(game info情報を見て必要に応じてstateGameEndへ)
		case stateGameEnd:
			// TODO(未実装)
		}

		s.fpsMgr.Wait()
	}
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
				c.dataStream.Send(&pb.Response{
					Type: pb.Response_UPDATESTATUS,
					Data: &pb.Response_Status_{
						Status: pb.Response_GAMEEND,
					},
				})
			}
		}
		logger.Error("Got error in session %s: %+v", s.id, s.exitErr.reason)
	}
}

func (s *Session) HandleSignal(clientID string, signal *pb.Request_Signal) error {
	switch signal.GetType() {
	case pb.Request_CHIPSELECT:
		// TODO(rawDataから選択したchipを取得)
		for i, c := range s.clients {
			if c.clientID == clientID {
				s.clients[i].chipSent = true
				return nil
			}
		}
		return fmt.Errorf("no such client %s", clientID)
	case pb.Request_GOCHIPSELECT:
		panic("TODO: 未実装")
	case pb.Request_INITPARAMS:
		var obj object.InitParam
		obj.Unmarshal(signal.GetRawData())
		s.gameHandler.AddObject(clientID, obj)
	}
	return nil
}

func (s *Session) HandleAction(clientID string, act *pb.Request_Action) error {
	switch act.GetType() {
	case pb.Request_MOVE:
		var move action.Move
		move.Unmarshal(act.GetRawData())
		s.gameHandler.MoveObject(move)
	case pb.Request_BUSTER:
		var buster action.Buster
		buster.Unmarshal(act.GetRawData())
		s.gameHandler.AddBuster(clientID, buster)
	case pb.Request_CHIPUSE:
		var chipInfo action.UseChip
		chipInfo.Unmarshal(act.GetRawData())
		s.gameHandler.UseChip(clientID, chipInfo)
	default:
		return fmt.Errorf("invalid action type %d is specified", act.GetType())
	}
	return nil
}

func (s *Session) changeState(next int) {
	logger.Info("Change state from %d to %d", s.state, next)
	s.state = next
}

func (s *Session) publishStateToClient(st pb.Response_Status) {
	for _, c := range s.clients {
		if c.dataStream == nil {
			continue
		}

		err := c.dataStream.Send(&pb.Response{
			Type: pb.Response_UPDATESTATUS,
			Data: &pb.Response_Status_{
				Status: st,
			},
		})
		if err != nil {
			logger.Error("failed to send status to client %s: %v", c.clientID, err)
			s.exitErr = &sessionError{
				generatorClientID: c.clientID,
				reason:            errSendFailed,
			}
			return
		}
	}
}

func (s *Session) publishGameInfo() {
	for _, c := range s.clients {
		if c.dataStream == nil {
			continue
		}

		err := c.dataStream.Send(&pb.Response{
			Type: pb.Response_DATA,
			Data: &pb.Response_RawData{
				RawData: s.gameHandler.GetInfo(),
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
	}
}
