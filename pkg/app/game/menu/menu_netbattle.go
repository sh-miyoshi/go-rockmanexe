package menu

import (
	"fmt"

	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/inputs"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/netconn"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
)

type menuNetBattle struct {
	imgMsgFrame int32
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

	return res, nil
}

func (m *menuNetBattle) End() {
	dxlib.DeleteGraph(m.imgMsgFrame)
}

func (m *menuNetBattle) Process() bool {
	if !m.isConnect {
		m.isConnect = true
		c := config.Get()
		if err := netconn.Connect(netconn.Config{
			StreamAddr:     c.Net.StreamAddr,
			ClientID:       c.Net.ClientID,
			ClientKey:      c.Net.ClientKey,
			ProgramVersion: common.ProgramVersion,
			Insecure:       c.Net.Insecure,
		}); err != nil {
			logger.Error("Failed to connect server: %v", err)
			m.messages = []string{
				"サーバーへの接続に失敗しました。",
				"設定を見直してください。",
			}
			return false
		}
		return true
	}

	if inputs.CheckKey(inputs.KeyCancel) == 1 {
		// Data init for next access
		netconn.Disconnect()
		m.isConnect = false
		m.messages = []string{"通信待機中です・・・"}

		stateChange(stateTop)
	}
	return false
}

func (m *menuNetBattle) Draw() {
	dxlib.DrawGraph(40, 205, m.imgMsgFrame, dxlib.TRUE)
	for i, msg := range m.messages {
		draw.MessageText(120, 220+int32(i)*30, 0x000000, msg)
	}
}
