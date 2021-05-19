package dstream

import (
	"context"
	"errors"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/db"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/net/routerpb"
)

type RouterStream struct {
	// TODO 有効なセッションとか
	sessions map[string]string
}

func (s *RouterStream) SendAction(ctx context.Context, action *pb.Action) (*pb.Result, error) {
	return nil, nil
}

func (s *RouterStream) PublishData(authReq *pb.AuthRequest, dataStream pb.Router_PublishDataServer) error {
	// Validate auth request
	// TODO validate version
	c, err := db.GetInst().ClientGetByID(authReq.Id)
	if err != nil {
		logger.Info("Failed to get client: %v", err)
		return errors.New("authenticate failed")
	}
	if c.Key != authReq.Key {
		logger.Info("got invalid key from user")
		return errors.New("authenticate failed")
	}

	// Add to sessionList

	// やること
	// authReqを検証
	// もし失敗ならerrを返す
	// 成功ならconnectPoolに追加
	// 100msごとに状態を送信
	return nil
}
