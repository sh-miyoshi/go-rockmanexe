package enemy

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/common"
)

/*
Enemy template

type enemy struct {
	pm EnemyParam
}

func (e *enemy) Init() error {
	return nil
}
func (e *enemy) End() {}

func (e *enemy) Process() (bool, error) {
	return false, nil
}
func (e *enemy) Draw() {}
func (e *enemy) Get() *EnemyParam {
	return &e.pm
}

*/

const (
	idMetall int = iota

	idMax
)

const (
	animMove int = iota
	animAtk1

	animMax
)

func getObject(id int, initParam EnemyParam) enemyObject {
	initParam.ID = uuid.New().String()

	switch id {
	case idMetall:
		return &enemyMetall{pm: initParam}
	}
	return nil
}

type enemyMetall struct {
	pm      EnemyParam
	imgMove []int32
}

func (e *enemyMetall) Init() error {
	e.imgMove = make([]int32, 1)
	fname := common.ImagePath + "battle/character/メットール_move.png"
	e.imgMove[0] = dxlib.LoadGraph(fname)
	if e.imgMove[0] == -1 {
		return fmt.Errorf("Failed to load image: %s", fname)
	}

	return nil
}
func (e *enemyMetall) End() {
	dxlib.DeleteGraph(e.imgMove[0])
}

func (e *enemyMetall) Process() (bool, error) {
	// TODO

	if e.pm.HP <= 0 {
		return true, nil
	}
	return false, nil
}
func (e *enemyMetall) Draw() {
	x := field.PanelSizeX*e.pm.PosX + field.PanelSizeX/2
	y := field.DrawPanelTopY + field.PanelSizeY*e.pm.PosY - 10
	img := e.imgMove[0] // TODO
	dxlib.DrawRotaGraph(int32(x), int32(y), 1, 0, img, dxlib.TRUE)
}

func (e *enemyMetall) Get() *EnemyParam {
	return &e.pm
}