package menu

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/net"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/netconn"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/inputs"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
)

type menuNetBattle struct {
	imgMsgFrame  int
	messages     []string
	isConnect    bool
	isErrorLoged bool
	result       Result
}

func netBattleNew() (*menuNetBattle, error) {
	res := &menuNetBattle{
		messages:     []string{"通信待機中です・・・"},
		isConnect:    false,
		isErrorLoged: false,
		result:       ResultContinue,
	}

	fname := config.ImagePath + "msg_frame.png"
	res.imgMsgFrame = dxlib.LoadGraph(fname)
	if res.imgMsgFrame == -1 {
		return nil, fmt.Errorf("failed to load menu message frame image %s", fname)
	}

	net.Init()
	return res, nil
}

func (m *menuNetBattle) End() {
	dxlib.DeleteGraph(m.imgMsgFrame)
}

func (m *menuNetBattle) Process() bool {
	if !m.isConnect {
		m.isConnect = true

		net.GetInst().ConnectRequest()
	}

	if inputs.CheckKey(inputs.KeyCancel) == 1 {
		net.GetInst().Disconnect()
		return true
	}

	status := net.GetInst().GetConnStatus()
	if status.Status == netconn.ConnStateOK {
		m.result = ResultGoNetBattle
		return true
	}
	if status.Status == netconn.ConnStateError {
		if !m.isErrorLoged {
			logger.Error("Failed to connect server: %v", status.Error)
			m.isErrorLoged = true
		}
		m.messages = []string{
			"サーバーへの接続に失敗しました。",
			"設定を見直してください。",
		}
	}

	return false
}

func (m *menuNetBattle) Draw() {
	dxlib.DrawGraph(40, 205, m.imgMsgFrame, true)
	for i, msg := range m.messages {
		draw.MessageText(120, 220+i*30, 0x000000, msg)
	}
}

func (m *menuNetBattle) GetResult() Result {
	return m.result
}
