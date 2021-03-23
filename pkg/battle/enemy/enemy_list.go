package enemy

import (
	"fmt"

	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/anim"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/skill"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/draw"
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

const (
	delayMetallAtk = 3
)

func getObject(id int, initParam EnemyParam) enemyObject {
	switch id {
	case idMetall:
		return &enemyMetall{pm: initParam}
	}
	return nil
}

type metallAtk struct {
	ownerID string
	count   int
	images  []int32
}

type enemyMetall struct {
	pm        EnemyParam
	imgMove   []int32
	count     int
	moveCount int
	atkID     string
	atk       metallAtk
}

func (e *enemyMetall) Init(ID string) error {
	e.pm.ID = ID
	e.imgMove = make([]int32, 1)
	fname := common.ImagePath + "battle/character/メットール_move.png"
	e.imgMove[0] = dxlib.LoadGraph(fname)
	if e.imgMove[0] == -1 {
		return fmt.Errorf("Failed to load image: %s", fname)
	}
	e.atk.images = make([]int32, 15)
	fname = common.ImagePath + "battle/character/メットール_atk.png"
	if res := dxlib.LoadDivGraph(fname, 15, 15, 1, 100, 140, e.atk.images); res == -1 {
		return fmt.Errorf("Failed to load image: %s", fname)
	}

	return nil
}

func (e *enemyMetall) End() {
	dxlib.DeleteGraph(e.imgMove[0])
	for _, img := range e.atk.images {
		dxlib.DeleteGraph(img)
	}
}

func (e *enemyMetall) Process() (bool, error) {
	e.count++

	if e.pm.HP <= 0 {
		return true, nil
	}

	const waitCount = 1 * 60
	const actionInterval = 1 * 60

	if e.atkID != "" {
		// Anim end
		if !anim.IsProcessing(e.atkID) {
			e.atkID = ""
			e.count = waitCount + 1 // Skip initial wait
		}
	}

	// Metall Actions
	if e.count < waitCount {
		return false, nil
	}

	if e.count%actionInterval == 0 {
		_, py := field.GetPos(e.pm.PlayerID)
		if py == e.pm.PosY {
			// Attack
			e.atk.count = 0
			e.atk.ownerID = e.pm.ID
			e.atkID = anim.New(&e.atk)
		} else {
			// Move
			if py > e.pm.PosY {
				battlecommon.MoveObject(&e.pm.PosX, &e.pm.PosY, common.DirectDown, field.PanelTypeEnemy, true)
			} else {
				battlecommon.MoveObject(&e.pm.PosX, &e.pm.PosY, common.DirectUp, field.PanelTypeEnemy, true)
			}
		}
	}

	return false, nil
}

func (e *enemyMetall) Draw() {
	x, y := battlecommon.ViewPos(e.pm.PosX, e.pm.PosY)
	img := e.imgMove[0]
	if e.atkID != "" {
		img = e.atk.images[e.atk.GetImageNo()]
	}
	dxlib.DrawRotaGraph(x, y, 1, 0, img, dxlib.TRUE)

	// Show HP
	if e.pm.HP > 0 {
		draw.Number(x, y+field.PanelSizeY-10, int32(e.pm.HP), draw.NumberOption{
			Color:    draw.NumberColorWhiteSmall,
			Centered: true,
		})
	}
}

func (e *enemyMetall) DamageProc(dm *damage.Damage) {
	if dm == nil {
		return
	}
	if dm.TargetType|damage.TargetEnemy != 0 {
		e.pm.HP -= dm.Power
		anim.New(effect.Get(dm.HitEffectType, e.pm.PosX, e.pm.PosY))
	}
}

func (e *enemyMetall) GetParam() anim.Param {
	return anim.Param{
		PosX:     e.pm.PosX,
		PosY:     e.pm.PosY,
		AnimType: anim.TypeObject,
	}
}

func (a *metallAtk) Draw() {
	// Nothing to do
}

func (a *metallAtk) Process() (bool, error) {
	a.count++

	if a.count == delayMetallAtk*10 {
		anim.New(skill.Get(skill.SkillShockWave, skill.Argument{
			OwnerID:    a.ownerID,
			Power:      10, // TODO: ダメージ
			TargetType: damage.TargetPlayer,
		}))
	}

	return a.count >= (len(a.images) * delayMetallAtk), nil
}

func (a *metallAtk) DamageProc(dm *damage.Damage) {
}

func (a *metallAtk) GetParam() anim.Param {
	return anim.Param{
		AnimType: anim.TypeEffect,
	}
}

func (a *metallAtk) GetImageNo() int {
	n := a.count / delayMetallAtk
	if n >= len(a.images) {
		n = len(a.images) - 1
	}
	return n
}
