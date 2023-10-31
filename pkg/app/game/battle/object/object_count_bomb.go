package object

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
)

const (
	explodeTime = 180
)

type CountBomb struct {
	pm    ObjectParam
	image int
	count int
}

func (o *CountBomb) Init(ownerID string, initParam ObjectParam) error {
	o.pm = initParam
	o.pm.objectID = uuid.New().String()
	o.count = 0

	fname := common.ImagePath + "battle/skill/カウントボム.png"
	o.image = dxlib.LoadGraph(fname)
	if o.image == -1 {
		return fmt.Errorf("failed to load image: %s", fname)
	}

	return nil
}

func (o *CountBomb) End() {
	dxlib.DeleteGraph(o.image)
}

func (o *CountBomb) Draw() {
	view := battlecommon.ViewPos(o.pm.Pos)
	dxlib.DrawRotaGraph(view.X, view.Y+16, 1, 0, o.image, true)
}

func (o *CountBomb) Process() (bool, error) {
	o.count++

	if o.pm.HP <= 0 {
		return true, nil
	}

	if o.count >= explodeTime {
		// TODO(ダメージ処理)
		logger.Info("explode count bomb with %+v", o.pm)
		return true, nil
	}

	if o.count%60 == 0 {
		sound.On(resources.SECountBombCountdown)
	}

	return false, nil
}

func (o *CountBomb) DamageProc(dm *damage.Damage) bool {
	if dm == nil {
		return false
	}

	o.pm.HP -= dm.Power
	localanim.AnimNew(effect.Get(dm.HitEffectType, o.pm.Pos, 5))
	return true
}

func (o *CountBomb) GetParam() objanim.Param {
	return objanim.Param{
		Param: anim.Param{
			ObjID:    o.pm.objectID,
			Pos:      o.pm.Pos,
			DrawType: anim.DrawTypeObject,
		},
		HP: o.pm.HP,
	}
}

func (o *CountBomb) GetObjectType() int {
	return objanim.ObjTypeNone
}

func (o *CountBomb) MakeInvisible(count int) {
}
