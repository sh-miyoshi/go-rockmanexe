package enemy

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	deleteanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/delete"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
)

type enemyTarget struct {
	pm    EnemyParam
	image int
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
		deleteanim.New(e.image, e.pm.Pos, false)
		localanim.AnimNew(effect.Get(resources.EffectTypeExplode, e.pm.Pos, 0))
		e.image = -1 // DeleteGraph at delete animation
		return true, nil
	}

	if e.pm.ParalyzedCount > 0 {
		e.pm.ParalyzedCount--
		return false, nil
	}

	return false, nil
}

func (e *enemyTarget) Draw() {
	if e.pm.InvincibleCount/5%2 != 0 {
		return
	}

	view := battlecommon.ViewPos(e.pm.Pos)
	dxlib.DrawRotaGraph(view.X, view.Y, 1, 0, e.image, true)

	if e.pm.ParalyzedCount > 0 {
		// 麻痺状態描画
		dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_ADD, 255)
		pm := 0
		if e.pm.ParalyzedCount/10%2 == 0 {
			pm = 255
		}
		dxlib.SetDrawBright(255, 255, pm)
		dxlib.DrawRotaGraph(view.X, view.Y, 1, 0, e.image, true)
		dxlib.SetDrawBright(255, 255, 255)
		dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_NOBLEND, 0)
	}

	// Show HP
	if e.pm.HP > 0 {
		draw.Number(view.X, view.Y+40, e.pm.HP, draw.NumberOption{
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
	return damageProc(dm, &e.pm)
}

func (e *enemyTarget) GetParam() objanim.Param {
	return objanim.Param{
		Param: anim.Param{
			ObjID:    e.pm.ObjectID,
			Pos:      e.pm.Pos,
			DrawType: anim.DrawTypeObject,
		},
		HP: e.pm.HP,
	}
}

func (e *enemyTarget) GetObjectType() int {
	return objanim.ObjTypeEnemy
}

func (e *enemyTarget) MakeInvisible(count int) {
	e.pm.InvincibleCount = count
}
