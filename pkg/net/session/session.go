package session

import (
	"fmt"
	"time"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/system"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/fps"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/net/netconnpb"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/object"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/sysinfo"
)

const (
	stateConnectWait int = iota
	stateChipSelectWait
	stateActing
)

type GameLogic interface {
	Init(clientIDs [2]string, sysReceiver chan sysinfo.SysInfo) error
	AddPlayerObject(clientID string, param object.InitParam) error
	HandleAction(clientID string, act *pb.Request_Action) error
	GetInfo(clientID string) []byte
	UpdateGameStatus()
	IsGameEnd() bool
	Cleanup()
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
	sysReceiver chan sysinfo.SysInfo
}

func newSession(sessionID string, gameHandler GameLogic) *Session {
	res := &Session{
		id:          sessionID,
		expiresAt:   time.Now().Add(sessionExpireTime),
		fpsMgr:      fps.Fps{},
		state:       stateConnectWait,
		gameHandler: gameHandler,
		sysReceiver: make(chan sysinfo.SysInfo),
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

		s.fpsMgr.Wait()

		now := time.Now()

		// check session expires
		if s.expiresAt.Before(now) {
			s.exitErr = &sessionError{
				reason: fmt.Errorf("session expired"),
			}
			return
		}

		// System error check
		if err := system.Error(); err != nil {
			s.exitErr = &sessionError{
				reason: err,
			}
			return
		}

		go func() {
			sysInfo := <-s.sysReceiver
			logger.Info("got system info %+v", sysInfo)
			switch sysInfo.Type {
			case sysinfo.TypeCutin:
				// カットイン状態にする
				s.publishStateToClient(pb.Response_CUTIN)
				// TODO: 情報を送る
			case sysinfo.TypeActing:
				// カットイン状態から戻す
				s.publishStateToClient(pb.Response_ACTING)
			default:
				system.SetError(fmt.Sprintf("System Info Type %d は現状扱えません", sysInfo.Type))
			}
		}()

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
			if err := s.gameHandler.Init(clientIDs, s.sysReceiver); err != nil {
				s.exitErr = &sessionError{
					reason: fmt.Errorf("failed to initialize game handler"),
				}
				return
			}
			s.changeState(stateChipSelectWait)
			s.publishStateToClient(pb.Response_CHIPSELECTWAIT)
		case stateChipSelectWait:
			s.publishGameInfo() // debug(送信頻度は要確認)
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

			// Game End
			if s.gameHandler.IsGameEnd() {
				s.publishStateToClient(pb.Response_GAMEEND)
				s.gameHandler.Cleanup()
				s.exitErr = &sessionError{}
				return
			}
		}
	}
}

func (s *Session) IsEnd() bool {
	return s.exitErr != nil
}

func (s *Session) End() {
	if s.exitErr.reason != nil {
		for _, c := range s.clients {
			if c.dataStream == nil || c.clientID == s.exitErr.generatorClientID {
				continue
			}

			// publish to alive clients
			logger.Debug("send error %v to client %s", s.exitErr.reason, c.clientID)
			c.dataStream.Send(&pb.Response{
				Type: pb.Response_UPDATESTATUS,
				Data: &pb.Response_Status_{
					Status: pb.Response_GAMEEND,
				},
			})
		}
		logger.Error("Got error in session %s: %+v", s.id, s.exitErr.reason)
	}
}

func (s *Session) HandleSignal(clientID string, signal *pb.Request_Signal) error {
	switch signal.GetType() {
	case pb.Request_CHIPSELECT:
		// TODO_Next rawDataから選択したchipを取得、検証用データにする
		for i, c := range s.clients {
			if c.clientID == clientID {
				s.clients[i].chipSent = true
				return nil
			}
		}
		return fmt.Errorf("no such client %s", clientID)
	case pb.Request_GOCHIPSELECT:
		for i := range s.clients {
			s.clients[i].chipSent = false
		}

		// Change game state to chpi select
		s.changeState(stateChipSelectWait)
		s.publishStateToClient(pb.Response_CHIPSELECTWAIT)
	case pb.Request_INITPARAMS:
		var obj object.InitParam
		obj.Unmarshal(signal.GetRawData())
		if err := s.gameHandler.AddPlayerObject(clientID, obj); err != nil {
			return fmt.Errorf("failed to add player object: %w", err)
		}
	}
	return nil
}

func (s *Session) HandleAction(clientID string, act *pb.Request_Action) error {
	return s.gameHandler.HandleAction(clientID, act)
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
				RawData: s.gameHandler.GetInfo(c.clientID),
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
