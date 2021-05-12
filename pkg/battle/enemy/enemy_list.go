package enemy

import (
	"fmt"
	"math/rand"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/anim"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/skill"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/sound"
)

/*
Enemy template

type enemy struct {
	pm EnemyParam
}

func (e *enemy) Init(objID string) error {
	e.pm.ObjectID = objID

	// Load Images
	return nil
}

func (e *enemy) End() {
	// Delete Images
}

func (e *enemy) Process() (bool, error) {
	// Return true if finished(e.g. hp=0)
	// Enemy Logic
	return false, nil
}

func (e *enemy) Draw() {
	// Show Enemy Images
}

func (e *enemy) DamageProc(dm *damage.Damage) {
	if dm == nil {
		return
	}
	if dm.TargetType|damage.TargetEnemy != 0 {
		e.pm.HP -= dm.Power
		anim.New(effect.Get(dm.HitEffectType, e.pm.PosX, e.pm.PosY, 5))
	}
}

func (e *enemy) GetParam() anim.Param {
	return anim.Param{
		ObjID:    e.pm.ObjectID,
		PosX:     e.pm.PosX,
		PosY:     e.pm.PosY,
		AnimType: anim.TypeObject,
		ObjType:  anim.ObjTypeEnemy,
	}
}

*/

const (
	IDMetall int = iota
	IDTarget
	IDBilly
)

const (
	delayMetallAtk = 3
)

var (
	metallActQueue = []string{}
)

func getObject(id int, initParam EnemyParam) enemyObject {
	switch id {
	case IDMetall:
		return &enemyMetall{pm: initParam}
	case IDTarget:
		return &enemyTarget{pm: initParam}
	case IDBilly:
		return &enemyBilly{pm: initParam}
	}
	return nil
}

func GetStandImageFile(id int) (name, ext string) {
	ext = ".png"

	switch id {
	case IDMetall:
		name = common.ImagePath + "battle/character/メットール"
	case IDTarget:
		name = common.ImagePath + "battle/character/的"
	case IDBilly:
		name = common.ImagePath + "battle/character/ビリー"
	}
	return
}

//-----------------------------------
// Metall
//-----------------------------------

type metallAtk struct {
	id      string
	ownerID string
	count   int
	images  []int32
}

type enemyMetall struct {
	pm              EnemyParam
	imgMove         []int32
	count           int
	moveFailedCount int
	atkID           string
	atk             metallAtk
}

func (e *enemyMetall) Init(objID string) error {
	name, ext := GetStandImageFile(IDMetall)

	e.pm.ObjectID = objID
	e.atk.id = uuid.New().String()
	e.imgMove = make([]int32, 1)
	fname := name + "_move" + ext
	e.imgMove[0] = dxlib.LoadGraph(fname)
	if e.imgMove[0] == -1 {
		return fmt.Errorf("failed to load image: %s", fname)
	}
	e.atk.images = make([]int32, 15)
	fname = name + "_atk" + ext
	if res := dxlib.LoadDivGraph(fname, 15, 15, 1, 100, 140, e.atk.images); res == -1 {
		return fmt.Errorf("failed to load image: %s", fname)
	}

	metallActQueue = append(metallActQueue, objID)

	return nil
}

func (e *enemyMetall) End() {
	dxlib.DeleteGraph(e.imgMove[0])
	for _, img := range e.atk.images {
		dxlib.DeleteGraph(img)
	}
}

func (e *enemyMetall) Process() (bool, error) {
	if e.pm.HP <= 0 {
		// Delete Animation
		img := &e.imgMove[0]
		if e.atkID != "" {
			img = &e.atk.images[e.atk.GetImageNo()]
		}
		newDelete(*img, e.pm.PosX, e.pm.PosY)
		*img = -1 // DeleteGraph at delete animation

		// Delete from act queue
		for i, id := range metallActQueue {
			if e.pm.ObjectID == id {
				metallActQueue = append(metallActQueue[:i], metallActQueue[i+1:]...)
				break
			}
		}
		return true, nil
	}

	const waitCount = 1 * 60
	const actionInterval = 1 * 60
	const forceAttackCount = 3

	if e.atkID != "" {
		// Anim end
		if !anim.IsProcessing(e.atkID) {
			metallActQueue = metallActQueue[1:]
			metallActQueue = append(metallActQueue, e.pm.ObjectID)

			e.atkID = ""
			e.count = 0
		}
		return false, nil
	}

	if metallActQueue[0] != e.pm.ObjectID {
		// other metall is acting
		return false, nil
	}

	e.count++

	// Metall Actions
	if e.count < waitCount {
		return false, nil
	}

	if e.count%actionInterval == 0 {
		_, py := anim.GetObjPos(e.pm.PlayerID)
		if py == e.pm.PosY || e.moveFailedCount >= forceAttackCount {
			// Attack
			e.atk.count = 0
			e.atk.ownerID = e.pm.ObjectID
			e.atkID = anim.New(&e.atk)
			e.moveFailedCount = 0
		} else {
			// Move
			moved := false
			if py > e.pm.PosY {
				moved = battlecommon.MoveObject(&e.pm.PosX, &e.pm.PosY, common.DirectDown, field.PanelTypeEnemy, true)
			} else {
				moved = battlecommon.MoveObject(&e.pm.PosX, &e.pm.PosY, common.DirectUp, field.PanelTypeEnemy, true)
			}
			if moved {
				e.moveFailedCount = 0
			} else {
				e.moveFailedCount++
			}
		}
	}

	return false, nil
}

func (e *enemyMetall) Draw() {
	x, y := battlecommon.ViewPos(e.pm.PosX, e.pm.PosY)
	img := e.imgMove[0]
	if e.atkID != "" {
		img = e.atk.images[e.atk.GetImageNo()]
	}
	dxlib.DrawRotaGraph(x, y, 1, 0, img, dxlib.TRUE)

	// Show HP
	if e.pm.HP > 0 {
		draw.Number(x, y-20, int32(e.pm.HP), draw.NumberOption{
			Color:    draw.NumberColorWhiteSmall,
			Centered: true,
		})
	}
}

func (e *enemyMetall) DamageProc(dm *damage.Damage) {
	if dm == nil {
		return
	}
	if dm.TargetType&damage.TargetEnemy != 0 {
		e.pm.HP -= dm.Power
		anim.New(effect.Get(dm.HitEffectType, e.pm.PosX, e.pm.PosY, 5))
	}
}

func (e *enemyMetall) GetParam() anim.Param {
	return anim.Param{
		ObjID:    e.pm.ObjectID,
		PosX:     e.pm.PosX,
		PosY:     e.pm.PosY,
		AnimType: anim.TypeObject,
		ObjType:  anim.ObjTypeEnemy,
	}
}

func (a *metallAtk) Draw() {
	// Nothing to do
}

func (a *metallAtk) Process() (bool, error) {
	a.count++

	if a.count == delayMetallAtk*10 {
		anim.New(skill.Get(skill.SkillShockWave, skill.Argument{
			OwnerID:    a.ownerID,
			Power:      10, // TODO: ダメージ
			TargetType: damage.TargetPlayer,
		}))
	}

	return a.count >= (len(a.images) * delayMetallAtk), nil
}

func (a *metallAtk) DamageProc(dm *damage.Damage) {
}

func (a *metallAtk) GetParam() anim.Param {
	return anim.Param{
		ObjID:    a.id,
		AnimType: anim.TypeObject,
		ObjType:  anim.ObjTypeNone,
	}
}

func (a *metallAtk) GetImageNo() int {
	n := a.count / delayMetallAtk
	if n >= len(a.images) {
		n = len(a.images) - 1
	}
	return n
}

//-----------------------------------
// Target
//-----------------------------------

type enemyTarget struct {
	pm    EnemyParam
	image int32
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
		newDelete(e.image, e.pm.PosX, e.pm.PosY)
		e.image = -1 // DeleteGraph at delete animation
		return true, nil
	}
	return false, nil
}

func (e *enemyTarget) Draw() {
	x, y := battlecommon.ViewPos(e.pm.PosX, e.pm.PosY)
	dxlib.DrawRotaGraph(x, y, 1, 0, e.image, dxlib.TRUE)

	// Show HP
	if e.pm.HP > 0 {
		draw.Number(x, y-55, int32(e.pm.HP), draw.NumberOption{
			Color:    draw.NumberColorWhiteSmall,
			Centered: true,
		})
	}
}

func (e *enemyTarget) DamageProc(dm *damage.Damage) {
	if dm == nil {
		return
	}
	logger.Debug("Enemy Target damaged: %+v", *dm)
	if dm.TargetType&damage.TargetEnemy != 0 {
		e.pm.HP -= dm.Power
		anim.New(effect.Get(dm.HitEffectType, e.pm.PosX, e.pm.PosY, 5))
	}
}

func (e *enemyTarget) GetParam() anim.Param {
	return anim.Param{
		ObjID:    e.pm.ObjectID,
		PosX:     e.pm.PosX,
		PosY:     e.pm.PosY,
		AnimType: anim.TypeObject,
		ObjType:  anim.ObjTypeEnemy,
	}
}

//-----------------------------------
// Billy
//-----------------------------------
const (
	billyActMove int = iota
	billyActAttack
)

const (
	delayBillyMove = 1
	delayBillyAtk  = 5
)

type billyAct struct {
	MoveDirect int

	ownerID  string
	count    int
	typ      int
	endCount int
	pPosX    *int
	pPosY    *int
}

type enemyBilly struct {
	pm        EnemyParam
	imgMove   []int32
	imgAttack []int32
	count     int
	moveCount int
	act       billyAct
}

func (e *enemyBilly) Init(objID string) error {
	e.pm.ObjectID = objID
	e.act.pPosX = &e.pm.PosX
	e.act.pPosY = &e.pm.PosY
	e.act.typ = -1
	e.act.ownerID = objID

	// Load Images
	name, ext := GetStandImageFile(IDBilly)
	e.imgMove = make([]int32, 6)
	fname := name + "_move" + ext
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 112, 114, e.imgMove); res == -1 {
		return fmt.Errorf("failed to load image: %s", fname)
	}

	e.imgAttack = make([]int32, 8)
	fname = name + "_atk" + ext
	if res := dxlib.LoadDivGraph(fname, 8, 8, 1, 124, 114, e.imgAttack); res == -1 {
		return fmt.Errorf("failed to load image: %s", fname)
	}

	return nil
}

func (e *enemyBilly) End() {
	// Delete Images
	for _, img := range e.imgMove {
		dxlib.DeleteGraph(img)
	}
	for _, img := range e.imgAttack {
		dxlib.DeleteGraph(img)
	}
}

func (e *enemyBilly) Process() (bool, error) {
	// Return true if finished(e.g. hp=0)
	// Enemy Logic
	if e.pm.HP <= 0 {
		// Delete Animation
		img := e.getCurrentImagePointer()
		newDelete(*img, e.pm.PosX, e.pm.PosY)
		*img = -1 // DeleteGraph at delete animation
		return true, nil
	}

	if e.act.Process() {
		return false, nil
	}

	const waitCount = 45
	const actionInterval = 75
	const moveNum = 4

	e.count++

	// Billy Actions
	if e.count < waitCount {
		return false, nil
	}

	if e.count%actionInterval == 0 {
		if e.moveCount < moveNum {
			// decide the direction to move
			// try 20 times to move
			for i := 0; i < 20; i++ {
				e.act.MoveDirect = 1 << rand.Intn(4)
				if battlecommon.MoveObject(&e.pm.PosX, &e.pm.PosY, e.act.MoveDirect, field.PanelTypeEnemy, false) {
					break
				}
			}

			e.act.SetAnim(billyActMove, len(e.imgMove)*delayBillyMove)
			e.moveCount++
		} else {
			// Attack
			e.act.SetAnim(billyActAttack, len(e.imgAttack)*delayBillyAtk)
			e.moveCount = 0
			e.count = 0
		}
	}

	return false, nil
}

func (e *enemyBilly) Draw() {
	x, y := battlecommon.ViewPos(e.pm.PosX, e.pm.PosY)
	img := e.getCurrentImagePointer()
	dxlib.DrawRotaGraph(x, y, 1, 0, *img, dxlib.TRUE)

	// Show HP
	if e.pm.HP > 0 {
		draw.Number(x, y-20, int32(e.pm.HP), draw.NumberOption{
			Color:    draw.NumberColorWhiteSmall,
			Centered: true,
		})
	}
}

func (e *enemyBilly) DamageProc(dm *damage.Damage) {
	if dm == nil {
		return
	}
	if dm.TargetType|damage.TargetEnemy != 0 {
		e.pm.HP -= dm.Power
		anim.New(effect.Get(dm.HitEffectType, e.pm.PosX, e.pm.PosY, 7))
	}
}

func (e *enemyBilly) GetParam() anim.Param {
	return anim.Param{
		ObjID:    e.pm.ObjectID,
		PosX:     e.pm.PosX,
		PosY:     e.pm.PosY,
		AnimType: anim.TypeObject,
		ObjType:  anim.ObjTypeEnemy,
	}
}

func (e *enemyBilly) getCurrentImagePointer() *int32 {
	img := &e.imgMove[0]
	if e.act.typ != -1 {
		imgs := e.imgMove
		delay := delayBillyMove
		if e.act.typ == billyActAttack {
			imgs = e.imgAttack
			delay = delayBillyAtk
		}

		cnt := e.act.count / delay
		if cnt >= len(imgs) {
			cnt = len(imgs) - 1
		}

		img = &imgs[cnt]
	}
	return img
}

func (a *billyAct) SetAnim(animType int, endCount int) {
	a.count = 0
	a.typ = animType
	a.endCount = endCount
}

func (a *billyAct) Process() bool {
	switch a.typ {
	case -1: // No animation
		return false
	case billyActAttack:
		if a.count == 5*delayBillyAtk {
			sound.On(sound.SEThunderBall)
			anim.New(skill.Get(skill.SkillThunderBall, skill.Argument{
				OwnerID:    a.ownerID,
				Power:      20, // TODO: ダメージ
				TargetType: damage.TargetPlayer,
			}))
		}
	case billyActMove:
		if a.count == 4*delayBillyMove {
			battlecommon.MoveObject(a.pPosX, a.pPosY, a.MoveDirect, field.PanelTypeEnemy, true)
		}
	}

	a.count++

	if a.count > a.endCount {
		// Reset params
		a.typ = -1
		a.count = 0
		return false // finished
	}
	return true // processing now
}
