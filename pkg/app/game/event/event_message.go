package event

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/window"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
)

type MessageHandler struct {
	win window.MessageWindow
}

func (h *MessageHandler) Init(values []byte) error {
	if err := h.win.Init(); err != nil {
		return err
	}
	// TODO: 適切なface typeを設定できるようにする
	h.win.SetMessage(string(values), window.FaceTypeNone)

	logger.Info("init message handler with %s", string(values))
	return nil
}

func (h *MessageHandler) End() {
	h.win.End()
}

func (h *MessageHandler) Draw() {
	h.win.Draw()
}

func (h *MessageHandler) Update() (bool, error) {
	return h.win.Update(), nil
}
