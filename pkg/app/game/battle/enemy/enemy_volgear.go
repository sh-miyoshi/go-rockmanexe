package enemy

import (
	"fmt"
	"math/rand"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

type volgearAtk struct {
	id      string
	ownerID string
	count   int
	images  []int
}

type enemyVolgear struct {
	pm        EnemyParam
	imgMove   []int
	count     int
	atkID     string
	atk       volgearAtk
	moveNum   int
	waitCount int
}

const (
	delayVolgearMove = 16
	delayVolgearAtk  = 2
	volgearInitWait  = 20
)

func (e *enemyVolgear) Init(objID string) error {
	e.pm.ObjectID = objID
	e.moveNum = 5
	e.waitCount = volgearInitWait

	// Load Images
	name, ext := GetStandImageFile(IDVolgear)
	e.atk.id = uuid.New().String()
	e.imgMove = make([]int, 4)
	fname := name + "_move" + ext
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 100, 100, e.imgMove); res == -1 {
		return fmt.Errorf("failed to load image: %s", fname)
	}
	e.atk.images = make([]int, 10)
	fname = name + "_atk" + ext
	if res := dxlib.LoadDivGraph(fname, 10, 10, 1, 122, 100, e.atk.images); res == -1 {
		return fmt.Errorf("failed to load image: %s", fname)
	}

	return nil
}

func (e *enemyVolgear) End() {
	// Delete Images
	for _, img := range e.imgMove {
		dxlib.DeleteGraph(img)
	}
	for _, img := range e.atk.images {
		dxlib.DeleteGraph(img)
	}
}

func (e *enemyVolgear) Process() (bool, error) {
	if e.pm.HP <= 0 {
		// Delete Animation
		img := e.getCurrentImagePointer()
		battlecommon.NewDelete(*img, e.pm.Pos, false)
		anim.New(effect.Get(effect.TypeExplode, e.pm.Pos, 0))
		*img = -1 // DeleteGraph at delete animation
		return true, nil
	}

	if e.waitCount > 0 {
		e.waitCount--
		return false, nil
	}

	if e.atkID != "" {
		// Anim end
		if !anim.IsProcessing(e.atkID) {
			e.atkID = ""
			e.waitCount = volgearInitWait
		}
		return false, nil
	}

	e.count++

	const actionInterval = 70

	if e.count%actionInterval == 0 {
		y := rand.Intn(field.FieldNum.Y)
		for i := 0; i < field.FieldNum.Y+1; i++ {
			next := common.Point{
				X: e.pm.Pos.X,
				Y: y,
			}
			if battlecommon.MoveObjectDirect(
				&e.pm.Pos,
				next,
				field.PanelTypeEnemy,
				true,
				field.GetPanelInfo,
			) {
				break
			}
			y = (y + 1) % field.FieldNum.Y
		}
		e.moveNum--
		if e.moveNum <= 0 {
			// TODO set attack
			e.moveNum = 4 + rand.Intn(2)
		}
	}

	return false, nil
}

func (e *enemyVolgear) Draw() {
	if e.pm.InvincibleCount/5%2 != 0 {
		return
	}

	view := battlecommon.ViewPos(e.pm.Pos)
	xflip := int32(dxlib.TRUE)
	img := e.getCurrentImagePointer()
	dxlib.DrawRotaGraph(view.X, view.Y+10, 1, 0, *img, true, dxlib.DrawRotaGraphOption{ReverseXFlag: &xflip})

	// Show HP
	if e.pm.HP > 0 {
		draw.Number(view.X, view.Y+40, e.pm.HP, draw.NumberOption{
			Color:    draw.NumberColorWhiteSmall,
			Centered: true,
		})
	}
}

func (e *enemyVolgear) DamageProc(dm *damage.Damage) bool {
	return damageProc(dm, &e.pm)
}

func (e *enemyVolgear) GetParam() anim.Param {
	return anim.Param{
		ObjID:    e.pm.ObjectID,
		Pos:      e.pm.Pos,
		AnimType: anim.AnimTypeObject,
	}
}

func (e *enemyVolgear) GetObjectType() int {
	return objanim.ObjTypeEnemy
}

func (e *enemyVolgear) MakeInvisible(count int) {
	e.pm.InvincibleCount = count
}

func (e *enemyVolgear) getCurrentImagePointer() *int {
	if e.atkID != "" {
		n := (e.atk.count / delayVolgearAtk)
		if n >= len(e.atk.images) {
			n = len(e.atk.images) - 1
		}
		return &e.atk.images[n]
	}

	n := (e.count / delayVolgearMove) % len(e.imgMove)
	return &e.imgMove[n]
}
