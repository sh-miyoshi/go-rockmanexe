package object

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
)

const (
	delayWaterPipeSet    = 3
	delayWaterPipeAttack = 6
)

type WaterPipeAtk struct {
	count       int
	images      []int32
	isAttacking bool
	pm          ObjectParam
}

type WaterPipe struct {
	pm       ObjectParam
	atk      WaterPipeAtk
	imgSet   []int32
	count    int
	atkCount int
}

func (o *WaterPipe) Init(ownerID string, initParam ObjectParam) error {
	o.pm = initParam
	o.pm.objectID = uuid.New().String()
	o.pm.xFlip = o.pm.OnwerCharType == objanim.ObjTypePlayer

	// Load Images
	o.imgSet = make([]int32, 4)
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
	x, y := battlecommon.ViewPos(o.pm.PosX, o.pm.PosY)

	if o.atk.IsAttacking() {
		o.atk.Draw(x, y)
		return
	}

	ofsx := int32(-8)
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
	dxlib.DrawRotaGraph(x+ofsx, y+16, 1, 0, o.imgSet[n], dxlib.TRUE, opt)
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
		anim.New(effect.Get(effect.TypeBlock, o.pm.PosX, o.pm.PosY, 5))
	}

	return false
}

func (o *WaterPipe) GetParam() anim.Param {
	return anim.Param{
		ObjID:    o.pm.objectID,
		PosX:     o.pm.PosX,
		PosY:     o.pm.PosY,
		AnimType: anim.AnimTypeObject,
	}
}

func (o *WaterPipe) GetObjectType() int {
	return objanim.ObjTypeNone
}

func (a *WaterPipeAtk) Init(pm ObjectParam) error {
	a.count = 0
	a.isAttacking = false
	a.pm = pm
	a.images = make([]int32, 9)
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

func (a *WaterPipeAtk) Draw(x, y int32) {
	n := 0
	if a.isAttacking {
		c := (a.count / delayWaterPipeAttack) % (len(a.images) * 2)
		s := len(a.images)
		n = c - (c/s)*((c-s)*2+1)
	}

	ofsx := int32(-81)
	opt := dxlib.DrawRotaGraphOption{}
	if a.pm.xFlip {
		f := int32(dxlib.TRUE)
		opt.ReverseXFlag = &f
		ofsx *= -1
	}

	dxlib.DrawRotaGraph(x+ofsx, y+13, 1, 0, a.images[n], dxlib.TRUE, opt)
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
			PosY:          a.pm.PosY,
			Power:         a.pm.Power,
			TTL:           6 * delayWaterPipeAttack,
			TargetType:    target,
			HitEffectType: effect.TypeNone,
			BigDamage:     true,
		}

		if a.pm.xFlip {
			dm.PosX = a.pm.PosX + 1
			damage.New(dm)
			dm.PosX = a.pm.PosX + 2
			damage.New(dm)
		} else {
			dm.PosX = a.pm.PosX - 1
			damage.New(dm)
			dm.PosX = a.pm.PosX - 2
			damage.New(dm)
		}
	}

	if a.count > len(a.images)*2*delayWaterPipeAttack {
		a.count = 0
		a.isAttacking = false
	}
}
