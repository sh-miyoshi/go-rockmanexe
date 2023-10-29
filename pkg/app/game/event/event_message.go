package event

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/inputs"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
)

type MessageHandler struct {
	messages    []string
	imgMsgFrame int
}

func (h *MessageHandler) Init(values []byte) error {
	h.messages = common.SplitMsg(string(values), 19)
	fname := common.ImagePath + "msg_frame.png"
	h.imgMsgFrame = dxlib.LoadGraph(fname)
	if h.imgMsgFrame == -1 {
		return fmt.Errorf("failed to load image: %s", fname)
	}

	logger.Info("init message handler with %s", string(values))
	return nil
}

func (h *MessageHandler) End() {
	dxlib.DeleteGraph(h.imgMsgFrame)
}

func (h *MessageHandler) Draw() {
	dxlib.DrawGraph(40, 205, h.imgMsgFrame, true)
	for i, msg := range h.messages {
		draw.MessageText(120, 220+i*30, 0x000000, msg)
	}
}

func (h *MessageHandler) Process() (bool, error) {
	if inputs.CheckKey(inputs.KeyEnter) == 1 {
		return true, nil
	}
	return false, nil
}
