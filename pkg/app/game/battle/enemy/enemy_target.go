package enemy

import (
	"github.com/cockroachdb/errors"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/manager"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common/deleteanim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
)

type enemyTarget struct {
	pm      EnemyParam
	image   int
	animMgr *manager.Manager
}

func (e *enemyTarget) Init(objID string, animMgr *manager.Manager) error {
	e.pm.ObjectID = objID
	name, ext := GetStandImageFile(IDTarget)
	fname := name + ext
	e.image = dxlib.LoadGraph(fname)
	if e.image == -1 {
		return errors.Newf("failed to load enemy image %s", fname)
	}
	e.animMgr = animMgr

	return nil
}

func (e *enemyTarget) End() {
	dxlib.DeleteGraph(e.image)
}

func (e *enemyTarget) Update() (bool, error) {
	if e.pm.HP <= 0 {
		deleteanim.New(e.image, e.pm.Pos, false, e.animMgr)
		e.animMgr.EffectAnimNew(effect.Get(resources.EffectTypeExplode, e.pm.Pos, 0))
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

	drawParalysis(view.X, view.Y, e.image, e.pm.ParalyzedCount)

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
			ObjID: e.pm.ObjectID,
			Pos:   e.pm.Pos,
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

func (e *enemyTarget) AddBarrier(hp int) {}

func (e *enemyTarget) SetCustomGaugeMax() {}
