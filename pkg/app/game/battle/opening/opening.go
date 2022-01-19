package opening

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/enemy"
)

type opening interface {
	Init(enemyList []enemy.EnemyParam) error
	End()
	Process() bool
	Draw()
}

var (
	openingInst opening
)

func Init(enemyList []enemy.EnemyParam) error {
	if enemy.IsBoss(enemyList[0].CharID) {
		openingInst = &boss{}
	} else {
		openingInst = &normal{}
	}

	return openingInst.Init(enemyList)
}

func End() {
	if openingInst != nil {
		openingInst.End()
		openingInst = nil
	}
}

func Process() bool {
	return openingInst.Process()
}

func Draw() {
	openingInst.Draw()
}
