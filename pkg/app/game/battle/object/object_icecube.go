package object

import (
	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/manager"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	delayIceCubeCreate = 2
)

type IceCube struct {
	pm      ObjectParam
	images  []int
	count   int
	animMgr *manager.Manager
}

func (o *IceCube) Init(ownerID string, initParam ObjectParam, animMgr *manager.Manager) error {
	o.pm = initParam
	o.pm.objectID = uuid.New().String()
	o.pm.xFlip = o.pm.OnwerCharType == objanim.ObjTypePlayer
	o.animMgr = animMgr

	// Load Images
	o.images = make([]int, 6)
	fname := config.ImagePath + "battle/character/アイスキューブ.png"
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 76, 90, o.images); res == -1 {
		return errors.Newf("failed to load image: %s", fname)
	}

	return nil
}

func (o *IceCube) End() {
	// Delete Images
	for _, img := range o.images {
		dxlib.DeleteGraph(img)
	}
}

func (o *IceCube) Update() (bool, error) {
	if o.pm.HP <= 0 {
		o.animMgr.EffectAnimNew(effect.Get(resources.EffectTypeIceBreak, o.pm.Pos, 0))
		return true, nil
	}

	o.count++

	return false, nil
}

func (o *IceCube) Draw() {
	view := battlecommon.ViewPos(o.pm.Pos)

	n := o.count / delayIceCubeCreate
	if n > len(o.images)-1 {
		n = len(o.images) - 1
	}
	dxlib.DrawRotaGraph(view.X, view.Y+16, 1, 0, o.images[n], true, dxlib.OptXReverse(o.pm.xFlip))
}

func (o *IceCube) DamageProc(dm *damage.Damage) bool {
	if dm == nil {
		return false
	}

	target := damage.TargetEnemy
	if o.pm.OnwerCharType == objanim.ObjTypePlayer {
		target = damage.TargetPlayer
	}

	hit := false
	if dm.DamageType == damage.TypePosition && dm.TargetObjType&target != 0 {
		hit = true
	} else if dm.DamageType == damage.TypeObject {
		hit = true
	}

	if hit {
		o.pm.HP -= dm.Power

		target = ^target

		for i := 0; i < dm.PushLeft; i++ {
			if !battlecommon.MoveObject(&o.pm.Pos, config.DirectLeft, -1, true, field.GetPanelInfo) {
				pos := point.Point{X: o.pm.Pos.X - 1, Y: o.pm.Pos.Y}
				if pos.X < 0 {
					o.pm.HP = 0 // 画面外のため終了
					return false
				}

				// もし目の前が敵キャラならダメージを与える
				objType := objanim.ObjTypePlayer
				if o.pm.OnwerCharType == objanim.ObjTypePlayer {
					objType = objanim.ObjTypeEnemy
				}

				objs := o.animMgr.ObjAnimGetObjs(objanim.Filter{Pos: &pos, ObjType: objType})
				if len(objs) > 0 {
					// Add damage
					o.animMgr.DamageManager().New(damage.Damage{
						DamageType:    damage.TypeObject,
						Power:         10,
						HitEffectType: resources.EffectTypeNone,
						StrengthType:  damage.StrengthHigh,
						Element:       damage.ElementNone,
						TargetObjID:   objs[0].ObjID,
						TargetObjType: target,
					})
					o.pm.HP = 0 // 自身は死ぬ
					return false
				}
				break
			}
		}
		for i := 0; i < dm.PushRight; i++ {
			if !battlecommon.MoveObject(&o.pm.Pos, config.DirectRight, -1, true, field.GetPanelInfo) {
				pos := point.Point{X: o.pm.Pos.X + 1, Y: o.pm.Pos.Y}
				if pos.X >= battlecommon.FieldNum.X {
					o.pm.HP = 0 // 画面外のため終了
					return false
				}
				// もし目の前が敵キャラならダメージを与える
				objType := objanim.ObjTypePlayer
				if o.pm.OnwerCharType == objanim.ObjTypePlayer {
					objType = objanim.ObjTypeEnemy
				}

				objs := o.animMgr.ObjAnimGetObjs(objanim.Filter{Pos: &pos, ObjType: objType})
				if len(objs) > 0 {
					// Add damage
					o.animMgr.DamageManager().New(damage.Damage{
						DamageType:    damage.TypeObject,
						Power:         10,
						HitEffectType: resources.EffectTypeNone,
						StrengthType:  damage.StrengthHigh,
						Element:       damage.ElementNone,
						TargetObjID:   objs[0].ObjID,
						TargetObjType: target,
					})
					o.pm.HP = 0 // 自身は死ぬ
					return false
				}
				break
			}
		}

		o.animMgr.EffectAnimNew(effect.Get(dm.HitEffectType, o.pm.Pos, 5))
		return true
	}

	return false
}

func (o *IceCube) GetParam() objanim.Param {
	return objanim.Param{
		Param: anim.Param{
			ObjID: o.pm.objectID,
			Pos:   o.pm.Pos,
		},
		HP: o.pm.HP,
	}
}

func (o *IceCube) GetObjectType() int {
	return objanim.ObjTypeNone
}

func (o *IceCube) MakeInvisible(count int) {}

func (o *IceCube) AddBarrier(hp int) {}

func (o *IceCube) SetCustomGaugeMax() {}
