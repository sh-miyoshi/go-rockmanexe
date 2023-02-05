package session

import (
	"fmt"
	"time"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/fps"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/net/netconnpb"
)

// TODO
const (
	statusConnectWait int = iota
	statusChipSelectWait
	statusActing
	statusGameEnd
)

type sessionError struct {
	generatorClientID string
	reason            error
}

type SessionClient struct {
	clientID   string
	dataStream pb.NetConn_TransDataServer
}

type Session struct {
	id        string
	clients   [2]SessionClient
	expiresAt time.Time
	// WIP gameInfo  GameInfo
	exitErr *sessionError
	fpsMgr  fps.Fps
	status  int

	// TODO
	/*
		go clientからデータを受け取る -> gameInfoに反映
		go fpsごとにgameInfoの計算
		go 一定時間ごとにclientに送信(※)
	*/
}

func newSession(sessionID string) *Session {
	res := &Session{
		id:        sessionID,
		expiresAt: time.Now().Add(sessionExpireTime),
		fpsMgr:    fps.Fps{TargetFPS: 60},
		status:    statusConnectWait,
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
	for {
		if s.exitErr != nil {
			return
		}

		// TODO 処理

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
