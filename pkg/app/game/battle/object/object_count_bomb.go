package object

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
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
	explodeTime = 240
)

type CountBomb struct {
	pm        ObjectParam
	imgBody   int
	imgNumber []int
	count     int
}

func (o *CountBomb) Init(ownerID string, initParam ObjectParam) error {
	o.pm = initParam
	o.pm.objectID = uuid.New().String()
	o.count = 0

	fname := config.ImagePath + "battle/skill/カウントボム.png"
	o.imgBody = dxlib.LoadGraph(fname)
	if o.imgBody == -1 {
		return fmt.Errorf("failed to load image: %s", fname)
	}

	o.imgNumber = make([]int, 4)
	fname = config.ImagePath + "battle/skill/カウントボム_数字.png"
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 16, 16, o.imgNumber); res == -1 {
		return fmt.Errorf("failed to load image: %s", fname)
	}

	return nil
}

func (o *CountBomb) End() {
	dxlib.DeleteGraph(o.imgBody)
	for _, img := range o.imgNumber {
		dxlib.DeleteGraph(img)
	}
	o.imgNumber = []int{}
}

func (o *CountBomb) Draw() {
	view := battlecommon.ViewPos(o.pm.Pos)
	dxlib.DrawRotaGraph(view.X, view.Y+16, 1, 0, o.imgBody, true)

	cnt := 3 - o.count/60
	dxlib.DrawRotaGraph(view.X, view.Y+20, 1, 0, o.imgNumber[cnt], true)
}

func (o *CountBomb) Process() (bool, error) {
	o.count++

	if o.pm.HP <= 0 {
		return true, nil
	}

	if o.count == explodeTime {
		target := damage.TargetPlayer
		if o.pm.OnwerCharType == objanim.ObjTypePlayer {
			target = damage.TargetEnemy
		}

		targetObjType := objanim.ObjTypeAll ^ o.pm.OnwerCharType ^ objanim.ObjTypeNone
		objs := localanim.ObjAnimGetObjs(objanim.Filter{ObjType: targetObjType})

		for _, obj := range objs {
			dm := damage.Damage{
				DamageType:    damage.TypeObject,
				Power:         o.pm.Power,
				TargetObjType: target,
				HitEffectType: resources.EffectTypeExplode,
				BigDamage:     true,
				TargetObjID:   obj.ObjID,
			}
			localanim.DamageManager().New(dm)
		}

		logger.Info("explode count bomb with %+v", o.pm)
		sound.On(resources.SEExplode)
		return true, nil
	}

	if o.count%60 == 0 {
		if o.count/60 == 3 {
			// 最終カウント
			sound.On(resources.SECountBombEnd)
		} else {
			sound.On(resources.SECountBombCountdown)
		}
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
