package enemy

import (
	"math/rand"

	"github.com/cockroachdb/errors"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/manager"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/system"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
	"github.com/stretchr/stew/slice"
)

const (
	IDMetall int = iota
	IDTarget
	IDBilly
	IDLark
	IDBoomer
	IDAquaman
	IDGaroo
	IDVolgear
	IDRockman
	IDSupportNPC
	IDColdman
	IDCirKill
	IDShrimpy
	IDForte
	IDBlues
)

type EnemyParam struct {
	CharID          int
	ObjectID        string
	PlayerID        string
	Pos             point.Point
	HP              int
	HPMax           int
	ActNo           int
	InvincibleCount int
	DamageElement   int
	ParalyzedCount  int
}

type enemyObject interface {
	objanim.Anim
	Init(ID string, animMgr *manager.Manager) error
	End()
}

var (
	ErrGameEnd  = errors.New("game end")
	enemies     = make(map[string]enemyObject)
	animManager *manager.Manager
	debugFlag   = false

	// 設定されてない場所であることがわかるような絶対にあり得ない座標
	emptyPos = point.Point{X: -100, Y: -100}
)

func Init(playerID string, enemyList []EnemyParam, animMgr *manager.Manager) error {
	animManager = animMgr
	for i, e := range enemyList {
		e.PlayerID = playerID
		e.ActNo = i
		e.HPMax = e.HP
		obj := getObject(e.CharID, e)
		objID := animManager.ObjAnimNew(obj)
		enemies[objID] = obj
	}

	// Init enemy data
	for id, e := range enemies {
		if err := e.Init(id, animManager); err != nil {
			return errors.Wrapf(err, "enemy %s init failed", id)
		}
	}
	return nil
}

func End() {
	// Cleanup existsting enemy data
	for _, e := range enemies {
		e.End()
	}
}

func MgrProcess() error {
	for id, e := range enemies {
		if !animManager.IsAnimProcessing(id) {
			e.End()
			delete(enemies, id)
		}
	}

	if len(enemies) == 0 {
		return ErrGameEnd
	}

	return nil
}

func GetStandImageFile(id int) (name, ext string) {
	ext = ".png"
	path := config.ImagePath + "battle/character/"
	name = path + GetName(id)
	return
}

func GetName(id int) string {
	switch id {
	case IDMetall:
		return "メットール"
	case IDTarget:
		return "的"
	case IDBilly:
		return "ビリー"
	case IDLark:
		return "ゲイラーク"
	case IDBoomer:
		return "ラウンダ"
	case IDAquaman:
		return "アクアマン"
	case IDGaroo:
		return "ガルー"
	case IDVolgear:
		return "ボルケギア"
	case IDRockman:
		return "ロックマン"
	case IDColdman:
		return "コールドマン"
	case IDCirKill:
		return "サーキラー"
	case IDShrimpy:
		return "エビロン"
	case IDForte:
		return "フォルテ"
	case IDBlues:
		return "ブルース"
	}
	return ""
}

func IsBoss(id int) bool {
	bossList := []int{IDAquaman, IDRockman, IDForte}
	return slice.Contains(bossList, id)
}

func getObject(id int, initParam EnemyParam) enemyObject {
	switch id {
	case IDMetall:
		return &enemyMetall{pm: initParam}
	case IDTarget:
		return &enemyTarget{pm: initParam}
	case IDBilly:
		return &enemyBilly{pm: initParam}
	case IDLark:
		return &enemyLark{pm: initParam}
	case IDBoomer:
		return &enemyBoomer{pm: initParam}
	case IDAquaman:
		return &enemyAquaman{pm: initParam}
	case IDGaroo:
		return &enemyGaroo{pm: initParam}
	case IDVolgear:
		return &enemyVolgear{pm: initParam}
	case IDRockman:
		system.SetError("enemy rockman is not implemented yet")
	case IDColdman:
		return &enemyColdman{pm: initParam}
	case IDCirKill:
		return &enemyCirKill{pm: initParam}
	case IDShrimpy:
		return &enemyShrimpy{pm: initParam}
	case IDForte:
		return &enemyForte{pm: initParam}
	case IDBlues:
		return &enemyBlues{pm: initParam}
	}
	return nil
}

func damageProc(dm *damage.Damage, pm *EnemyParam) bool {
	if dm == nil {
		return false
	}

	if dm.TargetObjType&damage.TargetEnemy != 0 {
		if pm.InvincibleCount > 0 && dm.Power > 0 {
			return false
		}

		if damage.IsWeakness(pm.DamageElement, *dm) {
			dm.Power *= 2
			animManager.EffectAnimNew(effect.Get(resources.EffectTypeExclamation, pm.Pos, 0))
		}

		pm.HP -= dm.Power

		for i := 0; i < dm.PushLeft; i++ {
			if !battlecommon.MoveObject(&pm.Pos, config.DirectLeft, battlecommon.PanelTypeEnemy, true, field.GetPanelInfo) {
				break
			}
		}
		for i := 0; i < dm.PushRight; i++ {
			if !battlecommon.MoveObject(&pm.Pos, config.DirectRight, battlecommon.PanelTypeEnemy, true, field.GetPanelInfo) {
				break
			}
		}

		animManager.EffectAnimNew(effect.Get(dm.HitEffectType, pm.Pos, 5))

		if dm.IsParalyzed {
			pm.ParalyzedCount = battlecommon.DefaultParalyzedTime
		}

		return true
	}
	return false
}

// 麻痺状態描画
func drawParalysis(x, y int, image int, count int, opt ...dxlib.DrawRotaGraphOption) {
	if count > 0 {
		dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_ADD, 255)
		// 黄色と白を点滅させる
		pm := 0
		if count/10%2 == 0 {
			pm = 255
		}
		dxlib.SetDrawBright(255, 255, pm)
		dxlib.DrawRotaGraph(x, y, 1, 0, image, true, opt...)
		dxlib.SetDrawBright(255, 255, 255)
		dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_NOBLEND, 0)
	}
}

func moveRandom(charPos *point.Point) {
	// 全エリアの中で移動可能な場所を探す
	movables := []point.Point{}
	for x := 0; x < battlecommon.FieldNum.X; x++ {
		for y := 0; y < battlecommon.FieldNum.Y; y++ {
			pos := point.Point{X: x, Y: y}
			if battlecommon.MoveObjectDirect(charPos, pos, battlecommon.PanelTypeEnemy, false, field.GetPanelInfo) {
				movables = append(movables, pos)
			}
		}
	}

	// 移動可能な場所があればランダムで移動
	if len(movables) > 0 {
		n := rand.Intn(len(movables))
		logger.Debug("enemy move to %v", movables[n])
		battlecommon.MoveObjectDirect(charPos, movables[n], battlecommon.PanelTypeEnemy, true, field.GetPanelInfo)
	}
}

/*
Enemy template

package enemy

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/manager"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
)

type enemy struct {
	pm EnemyParam
	animMgr *manager.Manager
}

func (e *enemy) Init(objID string, animMgr *manager.Manager) error {
	e.pm.ObjectID = objID
	e.animMgr = animMgr

	// Load Images
	return nil
}

func (e *enemy) End() {
	// Delete Images
}

func (e *enemy) Update() (bool, error) {
	// Return true if finished(e.g. hp=0)
	// Enemy Logic
	return false, nil
}

func (e *enemy) Draw() {
	// Show Enemy Images
}

func (e *enemy) DamageProc(dm *damage.Damage) bool {
	return damageProc(dm, &e.pm)
}

func (e *enemy) GetParam() objanim.Param {
	return objanim.Param{
		Param: anim.Param{
			ObjID:    e.pm.ObjectID,
			Pos:      e.pm.Pos,
		},
		HP: e.pm.HP,
	}
}

func (e *enemy) GetObjectType() int {
	return objanim.ObjTypeEnemy
}

func (e *enemy) MakeInvisible(count int) {
	e.pm.InvincibleCount = count
}

*/
