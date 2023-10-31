package object

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

const (
	explodeTimeSec = 3
)

type CountBomb struct {
	pm    ObjectParam
	image int
	count int
}

func (o *CountBomb) Init(ownerID string, initParam ObjectParam) error {
	return nil
}

func (o *CountBomb) End() {
	dxlib.DeleteGraph(o.image)
}

func (o *CountBomb) Draw() {
}

func (o *CountBomb) Process() (bool, error) {
	return false, nil
}
