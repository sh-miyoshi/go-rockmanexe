package menu

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/inputs"
)

type menuDevFeature struct {
	pointer int
}

func devFeatureNew() (*menuDevFeature, error) {
	return &menuDevFeature{
		pointer: 0,
	}, nil
}

func (t *menuDevFeature) End() {
}

func (t *menuDevFeature) Draw() {

}

func (t *menuDevFeature) Process() error {
	if inputs.CheckKey(inputs.KeyLButton) == 1 {
		return ErrGoMap
	}
	if inputs.CheckKey(inputs.KeyRButton) == 1 {
		field.Set4x4Area()
		return ErrGoBattle
	}
	return nil
}
