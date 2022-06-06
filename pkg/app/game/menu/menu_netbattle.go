package menu

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/inputs"
	netconn "github.com/sh-miyoshi/go-rockmanexe/pkg/app/newnetconn"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
)

type menuNetBattle struct {
	imgMsgFrame int
	messages    []string
	isConnect   bool
}

func netBattleNew() (*menuNetBattle, error) {
	res := &menuNetBattle{
		messages:  []string{"通信待機中です・・・"},
		isConnect: false,
	}

	fname := common.ImagePath + "menu/msg_frame.png"
	res.imgMsgFrame = dxlib.LoadGraph(fname)
	if res.imgMsgFrame == -1 {
		return nil, fmt.Errorf("failed to load menu message frame image %s", fname)
	}

	c := config.Get()
	netconn.Init(netconn.Config{
		StreamAddr:     c.Net.StreamAddr,
		ClientID:       c.Net.ClientID,
		ClientKey:      c.Net.ClientKey,
		ProgramVersion: common.ProgramVersion,
		Insecure:       c.Net.Insecure,
	})

	return res, nil
}

func (m *menuNetBattle) End() {
	dxlib.DeleteGraph(m.imgMsgFrame)
}

func (m *menuNetBattle) Process() bool {
	if !m.isConnect {
		m.isConnect = true

		netconn.GetInst().ConnectRequest()
	}

	if inputs.CheckKey(inputs.KeyCancel) == 1 {
		// Data init for next access
		netconn.GetInst().Disconnect()
		m.isConnect = false
		m.messages = []string{"通信待機中です・・・"}

		stateChange(stateTop)
	}

	status := netconn.GetInst().GetConnStatus()
	if status.Status == netconn.ConnStateOK {
		return true
	}
	if status.Status == netconn.ConnStateError {
		logger.Error("Failed to connect server: %v", status.Error)
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
