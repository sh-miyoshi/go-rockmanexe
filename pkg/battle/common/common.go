package common

import (
	"github.com/google/uuid"
	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/sound"
)

type deleteAction struct {
	id    string
	image int32
	x, y  int
	count int
}

// MoveObject ...
func MoveObject(x, y *int, direct int, objPanelType int, isMove bool) bool {
	nx := *x
	ny := *y

	// Check field out
	switch direct {
	case common.DirectUp:
		if ny <= 0 {
			return false
		}
		ny--
	case common.DirectDown:
		if ny >= field.FieldNumY-1 {
			return false
		}
		ny++
	case common.DirectLeft:
		if nx <= 0 {
			return false
		}
		nx--
	case common.DirectRight:
		if nx >= field.FieldNumX-1 {
			return false
		}
		nx++
	}

	pn := field.GetPanelInfo(nx, ny)
	if pn.ObjectID != "" {
		return false
	}
	// Check panel type
	if objPanelType != pn.Type {
		return false
	}

	if isMove {
		*x = nx
		*y = ny
	}

	return true
}

func MoveObjectDirect(x, y *int, targetX, targetY int, objPanelType int, isMove bool) bool {
	pn := field.GetPanelInfo(targetX, targetY)
	if pn.ObjectID != "" {
		return false
	}
	// Check panel type
	if objPanelType != pn.Type {
		return false
	}

	if isMove {
		*x = targetX
		*y = targetY
	}

	return true
}

func ViewPos(x, y int) (viewX, viewY int32) {
	viewX = int32(field.PanelSizeX*x + field.PanelSizeX/2)
	viewY = int32(field.DrawPanelTopY + field.PanelSizeY*y - 10)
	return viewX, viewY
}

func GetOffset(nextPos, nowPos, beforePos int, cnt, totalCnt int, size int) int {
	// if cnt < total_count/2
	//   init_offset = (before - now) * size / 2
	//   offset = init_offset - (before - now)*(count*size/total_count))

	var res int
	if cnt < totalCnt/2 {
		res = (beforePos - nowPos)
	} else {
		res = (nowPos - nextPos)
	}

	return res * size * (totalCnt - 2*cnt) / (totalCnt * 2)
}

func NewDelete(image int32, x, y int, isPlayer bool) {
	if isPlayer {
		sound.On(sound.SEPlayerDeleted)
	} else {
		sound.On(sound.SEEnemyDeleted)
	}

	anim.New(&deleteAction{
		id:    uuid.New().String(),
		image: image,
		x:     x,
		y:     y,
	})
}

func (p *deleteAction) Process() (bool, error) {
	p.count++
	if p.count == 15 {
		dxlib.DeleteGraph(p.image)
		return true, nil
	}
	return false, nil
}

func (p *deleteAction) Draw() {
	x, y := ViewPos(p.x, p.y)

	dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_INVSRC, 255)
	dxlib.DrawRotaGraph(x, y, 1, 0, p.image, dxlib.TRUE)
	dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_ADD, 255)
	dxlib.DrawRotaGraph(x, y, 1, 0, p.image, dxlib.TRUE)
	dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_NOBLEND, 0)
}

func (p *deleteAction) DamageProc(dm *damage.Damage) {
}

func (p *deleteAction) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.id,
		PosX:     p.x,
		PosY:     p.y,
		AnimType: anim.TypeObject,
		ObjType:  anim.ObjTypeNone,
	}
}
