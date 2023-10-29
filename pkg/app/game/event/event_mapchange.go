package event

import (
	"bytes"
	"encoding/gob"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
)

type MapChangeArgs struct {
	MapID   int
	InitPos common.Point
}

type MapChangeHandler struct {
	args MapChangeArgs
}

func (h *MapChangeHandler) Init(values []byte) error {
	h.args.Unmarshal(values)
	logger.Info("init map change handler with %+v", h.args)
	return nil
}

func (h *MapChangeHandler) End() {
}

func (h *MapChangeHandler) Draw() {
}

func (h *MapChangeHandler) Process() (bool, error) {
	logger.Info("store map args %+v to event storedValues", h.args)
	storedValues = h.args.Marshal()
	resultCode = ResultMapChange
	return true, nil
}

func (p *MapChangeArgs) Marshal() []byte {
	buf := bytes.NewBuffer(nil)
	gob.NewEncoder(buf).Encode(p)
	return buf.Bytes()
}

func (p *MapChangeArgs) Unmarshal(data []byte) {
	buf := bytes.NewBuffer(data)
	_ = gob.NewDecoder(buf).Decode(p)
}
