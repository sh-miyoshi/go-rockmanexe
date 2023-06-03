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
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

const (
	delayWaterPipeSet    = 3
	delayWaterPipeAttack = 6
)

type WaterPipeAtk struct {
	count       int
	images      []int
	isAttacking bool
	pm          ObjectParam
}

type WaterPipe struct {
	pm       ObjectParam
	atk      WaterPipeAtk
	imgSet   []int
	count    int
	atkCount int
}

func (o *WaterPipe) Init(ownerID string, initParam ObjectParam) error {
	o.pm = initParam
	o.pm.objectID = uuid.New().String()
	o.pm.xFlip = o.pm.OnwerCharType == objanim.ObjTypePlayer

	// Load Images
	o.imgSet = make([]int, 4)
	fname := common.ImagePath + "battle/character/水道管_set.png"
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 72, 88, o.imgSet); res == -1 {
		return fmt.Errorf("failed to load image: %s", fname)
	}

	if err := o.atk.Init(o.pm); err != nil {
		return fmt.Errorf("failed to init water pipe attack %w", err)
	}

	return nil
}

func (o *WaterPipe) End() {
	// Delete Images
	for _, img := range o.imgSet {
		dxlib.DeleteGraph(img)
	}

	o.atk.End()
}

func (o *WaterPipe) Process() (bool, error) {
	if o.atk.IsAttacking() {
		o.atk.Process()
		return false, nil
	}

	o.count++

	if o.count == 1 {
		sound.On(sound.SEObjectCreate)

		pn := field.GetPanelInfo(o.pm.Pos)
		if pn.Status == battlecommon.PanelStatusHole {
			return true, nil
		}
	}

	if o.count%o.pm.Interval == 0 {
		o.atk.Start()
		o.atkCount++

		if o.atkCount > o.pm.AttackNum {
			return true, nil
		}
	}

	return false, nil
}

func (o *WaterPipe) Draw() {
	view := battlecommon.ViewPos(o.pm.Pos)

	if o.atk.IsAttacking() {
		o.atk.Draw(view)
		return
	}

	ofsx := -8
	opt := dxlib.DrawRotaGraphOption{}
	if o.pm.xFlip {
		f := int32(dxlib.TRUE)
		opt.ReverseXFlag = &f
		ofsx *= -1
	}

	n := o.count / delayWaterPipeSet
	if n > len(o.imgSet)-1 {
		n = len(o.imgSet) - 1
	}
	dxlib.DrawRotaGraph(view.X+ofsx, view.Y+16, 1, 0, o.imgSet[n], true, opt)
}

func (o *WaterPipe) DamageProc(dm *damage.Damage) bool {
	if dm == nil {
		return false
	}

	target := damage.TargetEnemy
	if o.pm.OnwerCharType == objanim.ObjTypePlayer {
		target = damage.TargetPlayer
	}

	if dm.TargetType&target != 0 {
		o.pm.HP--
		localanim.AnimNew(effect.Get(battlecommon.EffectTypeBlock, o.pm.Pos, 5))
	}

	return false
}

func (o *WaterPipe) GetParam() objanim.Param {
	return objanim.Param{
		Param: anim.Param{
			ObjID:    o.pm.objectID,
			Pos:      o.pm.Pos,
			DrawType: anim.DrawTypeObject,
		},
		HP: o.pm.HP,
	}
}

func (o *WaterPipe) GetObjectType() int {
	return objanim.ObjTypeNone
}

func (o *WaterPipe) MakeInvisible(count int) {}

func (a *WaterPipeAtk) Init(pm ObjectParam) error {
	a.count = 0
	a.isAttacking = false
	a.pm = pm
	a.images = make([]int, 9)
	fname := common.ImagePath + "battle/character/水道管_atk.png"
	if res := dxlib.LoadDivGraph(fname, 9, 9, 1, 234, 110, a.images); res == -1 {
		return fmt.Errorf("failed to load image: %s", fname)
	}
	return nil
}

func (a *WaterPipeAtk) End() {
	for _, img := range a.images {
		dxlib.DeleteGraph(img)
	}
}

func (a *WaterPipeAtk) Start() {
	a.isAttacking = true
}

func (a *WaterPipeAtk) IsAttacking() bool {
	return a.isAttacking
}

func (a *WaterPipeAtk) Draw(pos common.Point) {
	n := 0
	if a.isAttacking {
		c := (a.count / delayWaterPipeAttack) % (len(a.images) * 2)
		s := len(a.images)
		n = c - (c/s)*((c-s)*2+1)
	}

	ofsx := -81
	opt := dxlib.DrawRotaGraphOption{}
	if a.pm.xFlip {
		f := int32(dxlib.TRUE)
		opt.ReverseXFlag = &f
		ofsx *= -1
	}

	dxlib.DrawRotaGraph(pos.X+ofsx, pos.Y+13, 1, 0, a.images[n], true, opt)
}

func (a *WaterPipeAtk) Process() {
	a.count++

	if a.count == 1 {
		sound.On(sound.SEWaterpipeAttack)
	}

	if a.count == 7*delayWaterPipeAttack-2 {
		target := damage.TargetPlayer
		if a.pm.OnwerCharType == objanim.ObjTypePlayer {
			target = damage.TargetEnemy
		}

		dm := damage.Damage{
			Pos:           a.pm.Pos,
			Power:         a.pm.Power,
			TTL:           6 * delayWaterPipeAttack,
			TargetType:    target,
			HitEffectType: battlecommon.EffectTypeNone,
			BigDamage:     true,
			DamageType:    damage.TypeWater,
		}

		if a.pm.xFlip {
			dm.Pos.X = a.pm.Pos.X + 1
			localanim.DamageManager().New(dm)
			dm.Pos.X = a.pm.Pos.X + 2
			localanim.DamageManager().New(dm)
		} else {
			dm.Pos.X = a.pm.Pos.X - 1
			localanim.DamageManager().New(dm)
			dm.Pos.X = a.pm.Pos.X - 2
			localanim.DamageManager().New(dm)
		}
	}

	if a.count > len(a.images)*2*delayWaterPipeAttack {
		a.count = 0
		a.isAttacking = false
	}
}
