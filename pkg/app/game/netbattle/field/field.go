package field

import (
	"fmt"

	battlefield "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
)

type Field struct {
	bgInst battlefield.Background
}

func New() (*Field, error) {
	logger.Info("Initialize battle field data")

	res := &Field{}

	// TODO: Serverから取得する
	if err := res.bgInst.Init(battlefield.BGType秋原町); err != nil {
		return nil, fmt.Errorf("failed to load background: %w", err)
	}

	logger.Info("Successfully initialized battle field data")
	return res, nil
}

func (f *Field) End() {
	f.bgInst.End()
}

func (f *Field) Draw() {
	f.bgInst.Draw()
}

func (f *Field) Update() {
	f.bgInst.Process()
}
