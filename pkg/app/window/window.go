package window

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/inputs"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
)

type MessageWindow struct {
	Image    int
	Messages []string
}

func New(msg string) (MessageWindow, error) {
	fname := common.ImagePath + "msg_frame.png"
	res := MessageWindow{
		Image: dxlib.LoadGraph(fname),
	}
	if res.Image == -1 {
		return res, fmt.Errorf("failed to load message frame image %s", fname)
	}
	res.Messages = common.SplitMsg(msg, 19)
	return res, nil
}

func (w *MessageWindow) End() {
	dxlib.DeleteGraph(w.Image)
}

func (w *MessageWindow) Draw() {
	dxlib.DrawGraph(40, 205, w.Image, true)
	for i, msg := range w.Messages {
		draw.MessageText(120, 220+i*30, 0x000000, msg)
	}
}

func (w *MessageWindow) Process() bool {
	// TODO: 改ページ, 流れるように表示
	if inputs.CheckKey(inputs.KeyEnter) == 1 {
		logger.Debug("end window process")
		return true
	}
	return false
}
