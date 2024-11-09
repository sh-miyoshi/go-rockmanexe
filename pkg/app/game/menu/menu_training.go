package menu

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/inputs"
)

type menuTraining struct {
}

func trainingNew() (*menuTraining, error) {
	res := &menuTraining{}
	return res, nil
}

func (r *menuTraining) End() {
}

func (r *menuTraining) Update() bool {
	if inputs.CheckKey(inputs.KeyCancel) == 1 {
		sound.On(resources.SECancel)
		return true
	}
	return false
}

func (r *menuTraining) Draw() {
}

func (r *menuTraining) GetResult() Result {
	return ResultContinue
}
