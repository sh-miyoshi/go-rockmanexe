package enemy

import (
	"fmt"

	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
)

type enemyTarget struct {
	pm    EnemyParam
	image int32
}

func (e *enemyTarget) Init(objID string) error {
	e.pm.ObjectID = objID
	name, ext := GetStandImageFile(IDTarget)
	fname := name + ext
	e.image = dxlib.LoadGraph(fname)
	if e.image == -1 {
		return fmt.Errorf("failed to load enemy image %s", fname)
	}

	return nil
}

func (e *enemyTarget) End() {
	dxlib.DeleteGraph(e.image)
}

func (e *enemyTarget) Process() (bool, error) {
	if e.pm.HP <= 0 {
		battlecommon.NewDelete(e.image, e.pm.PosX, e.pm.PosY, false)
		anim.New(effect.Get(effect.TypeExplode, e.pm.PosX, e.pm.PosY, 0))
		e.image = -1 // DeleteGraph at delete animation
		return true, nil
	}
	return false, nil
}

func (e *enemyTarget) Draw() {
	x, y := battlecommon.ViewPos(e.pm.PosX, e.pm.PosY)
	dxlib.DrawRotaGraph(x, y, 1, 0, e.image, dxlib.TRUE)

	// Show HP
	if e.pm.HP > 0 {
		draw.Number(x, y+40, int32(e.pm.HP), draw.NumberOption{
			Color:    draw.NumberColorWhiteSmall,
			Centered: true,
		})
	}
}

func (e *enemyTarget) DamageProc(dm *damage.Damage) bool {
	if dm == nil {
		return false
	}
	logger.Debug("Enemy Target damaged: %+v", *dm)
	if dm.TargetType&damage.TargetEnemy != 0 {
		e.pm.HP -= dm.Power
		anim.New(effect.Get(dm.HitEffectType, e.pm.PosX, e.pm.PosY, 5))
		return true
	}
	return false
}

func (e *enemyTarget) GetParam() anim.Param {
	return anim.Param{
		ObjID:    e.pm.ObjectID,
		PosX:     e.pm.PosX,
		PosY:     e.pm.PosY,
		AnimType: anim.TypeObject,
		ObjType:  anim.ObjTypeEnemy,
	}
}
