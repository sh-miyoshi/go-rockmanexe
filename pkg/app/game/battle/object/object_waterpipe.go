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
)

const (
	delayWaterPipeSet    = 3
	delayWaterPipeAttack = 8
)

type WaterPipeAtk struct {
	count       int
	images      []int32
	isAttacking bool
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
	o.pm.ObjectID = uuid.New().String()

	// Load Images
	o.imgSet = make([]int32, 4)
	fname := common.ImagePath + "battle/character/水道管_set.png"
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 72, 88, o.imgSet); res == -1 {
		return fmt.Errorf("failed to load image: %s", fname)
	}

	if err := o.atk.Init(); err != nil {
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

	if o.count%150 == 0 {
		o.atk.Start()
		o.atkCount++

		if o.atkCount > 5 {
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

	n := o.count / delayWaterPipeSet
	if n > len(o.imgSet)-1 {
		n = len(o.imgSet) - 1
	}
	dxlib.DrawRotaGraph(x-8, y+16, 1, 0, o.imgSet[n], dxlib.TRUE)
}

func (o *WaterPipe) DamageProc(dm *damage.Damage) bool {
	// TODO: はじかれエフェクト
	return false
}

func (o *WaterPipe) GetParam() anim.Param {
	return anim.Param{
		ObjID:    o.pm.ObjectID,
		PosX:     o.pm.PosX,
		PosY:     o.pm.PosY,
		AnimType: anim.AnimTypeObject,
	}
}

func (o *WaterPipe) GetObjectType() int {
	return objanim.ObjTypeNone
}

func (a *WaterPipeAtk) Init() error {
	a.count = 0
	a.isAttacking = false
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

	dxlib.DrawRotaGraph(x-81, y+13, 1, 0, a.images[n], dxlib.TRUE)
}

func (a *WaterPipeAtk) Process() {
	a.count++

	// TODO(ダメージ)

	if a.count > len(a.images)*2*delayWaterPipeAttack {
		a.count = 0
		a.isAttacking = false
	}
}
