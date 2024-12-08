package event

import "github.com/sh-miyoshi/go-rockmanexe/pkg/logger"

type EndHandler struct {
}

func (h *EndHandler) Init(values []byte) error {
	logger.Info("init end scenario")
	return nil
}

func (h *EndHandler) End() {
}

func (h *EndHandler) Draw() {
}

func (h *EndHandler) Update() (bool, error) {
	resultCode = ResultEnd
	return true, nil
}
