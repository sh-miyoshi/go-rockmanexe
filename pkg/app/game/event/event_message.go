package event

import (
	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/inputs"
)

type MessageHandler struct {
	Message string
}

func (h *MessageHandler) Draw() {
	dxlib.DrawString(0, 100, h.Message, 0xff0000)
}

func (h *MessageHandler) Process() (bool, error) {
	if inputs.CheckKey(inputs.KeyEnter) == 1 {
		return true, nil
	}
	return false, nil
}
