package event

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/window"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
)

type MessageHandler struct {
	win window.MessageWindow
}

func (h *MessageHandler) Init(values []byte) error {
	var err error
	h.win, err = window.New(string(values))
	if err != nil {
		return err
	}

	logger.Info("init message handler with %s", string(values))
	return nil
}

func (h *MessageHandler) End() {
	h.win.End()
}

func (h *MessageHandler) Draw() {
	h.win.Draw()
}

func (h *MessageHandler) Process() (bool, error) {
	return h.win.Process(), nil
}
