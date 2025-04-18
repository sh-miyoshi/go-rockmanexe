package enemy

import (
	"math/rand"

	"github.com/cockroachdb/errors"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/manager"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common/deleteanim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/math"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	forteActTypeStand = iota
	forteActTypeMove
	forteActTypeShooting
	forteActTypeHellsRolling
	forteActTypeDarkArmBlade1
	forteActTypeDarkArmBlade3
	forteActTypeDarknessOverload
	forteActTypeDamage

	forteActTypeMax
)

var (
	forteDelays   = [forteActTypeMax]int{1, 1, 6, 6, 1, 1, 6, 5}
	forteAtkPower = [forteActTypeMax]uint{0, 0, 50, 80, 100, 100, 280, 0}
)

// Note:
//   キャラがロックマンを起点に移動する場合、移動後ロックマンが移動すると追従してしまう
//   一回のみ移動させるためisTargetPosMovedを使用している

type enemyForte struct {
	pm               EnemyParam
	images           [forteActTypeMax][]int
	count            int
	state            int
	waitCount        int
	nextState        int
	moveNum          int
	targetPos        point.Point
	isTargetPosMoved bool
	bladeAtkCount    int
	atkIDs           []string
	isCharReverse    bool
	animMgr          *manager.Manager
}

func (e *enemyForte) Init(objID string, animMgr *manager.Manager) error {
	e.pm.ObjectID = objID
	e.state = forteActTypeStand
	e.waitCount = 20
	e.nextState = forteActTypeMove
	e.moveNum = 2
	e.targetPos = emptyPos
	e.isTargetPosMoved = false
	e.bladeAtkCount = 0
	e.isCharReverse = false
	e.animMgr = animMgr

	// Load Images
	name, ext := GetStandImageFile(IDForte)

	fname := name + "_all" + ext
	tmp := make([]int, 45)
	if res := dxlib.LoadDivGraph(fname, 45, 9, 5, 136, 172, tmp); res == -1 {
		return errors.Newf("failed to load image: %s", fname)
	}
	e.images[forteActTypeStand] = make([]int, 1)
	e.images[forteActTypeMove] = make([]int, 6)
	e.images[forteActTypeShooting] = make([]int, 9)
	e.images[forteActTypeHellsRolling] = make([]int, 9)
	e.images[forteActTypeDarkArmBlade1] = make([]int, 1)
	e.images[forteActTypeDarkArmBlade3] = make([]int, 3)
	e.images[forteActTypeDarknessOverload] = make([]int, 9)
	e.images[forteActTypeDamage] = make([]int, 1)

	e.images[forteActTypeStand][0] = tmp[0]
	for i := 0; i < 6; i++ {
		e.images[forteActTypeMove][i] = tmp[i]
	}
	for i := 0; i < 9; i++ {
		e.images[forteActTypeShooting][i] = tmp[9+i]
		e.images[forteActTypeHellsRolling][i] = tmp[18+i]
		e.images[forteActTypeDarknessOverload][i] = tmp[18+i]
	}
	e.images[forteActTypeDarkArmBlade1][0] = tmp[27]
	for i := 0; i < 3; i++ {
		e.images[forteActTypeDarkArmBlade3][i] = tmp[27+i]
	}
	e.images[forteActTypeDamage][0] = tmp[36]

	cleanup := []int{6, 7, 8}
	for i := 30; i < len(tmp); i++ {
		if i != 36 {
			cleanup = append(cleanup, i)
		}
	}

	for _, t := range cleanup {
		dxlib.DeleteGraph(t)
	}

	return nil
}

func (e *enemyForte) End() {
	// Delete Images
	for i := 0; i < forteActTypeMax; i++ {
		for j := 0; j < len(e.images[i]); j++ {
			dxlib.DeleteGraph(e.images[i][j])
		}
		e.images[i] = []int{}
	}
}

func (e *enemyForte) Update() (bool, error) {
	if e.pm.HP <= 0 {
		// Delete Animation
		img := e.getCurrentImagePointer()
		deleteanim.New(*img, e.pm.Pos, false, e.animMgr)
		e.animMgr.EffectAnimNew(effect.Get(resources.EffectTypeExplode, e.pm.Pos, 0))
		*img = -1 // DeleteGraph at delete animation
		return true, nil
	}

	if e.pm.ParalyzedCount > 0 {
		e.pm.ParalyzedCount--
		return false, nil
	}

	// Enemy Logic
	if e.pm.InvincibleCount > 0 {
		e.pm.InvincibleCount--
	}

	switch e.state {
	case forteActTypeStand:
		e.waitCount--
		if e.waitCount <= 0 {
			return e.stateChange(e.nextState)
		}
	case forteActTypeMove:
		if e.count == 6*forteDelays[forteActTypeMove] {
			if !e.targetPos.Equal(emptyPos) {
				if !battlecommon.MoveObjectDirect(
					&e.pm.Pos,
					e.targetPos,
					-1, // プレイヤーのパネルでも移動可能
					true,
					field.GetPanelInfo,
				) {
					// 移動に失敗したら、移動からやり直し
					logger.Debug("Forte move failed. retry")
					return e.clearState()
				}
				e.targetPos = emptyPos
				if e.waitCount == 0 {
					e.waitCount = 20
				}
				return e.stateChange(forteActTypeStand)
			}

			moveRandom(&e.pm.Pos)
			e.waitCount = 40
			e.moveNum--
			if e.moveNum <= 0 {
				if debugFlag {
					e.moveNum = 3
					e.nextState = forteActTypeDarknessOverload
					return e.stateChange(forteActTypeStand)
				}

				e.moveNum = rand.Intn(3) + 3

				// Action process
				// 攻撃処理(HPの減り具合から乱数で攻撃決定)
				// シューティングバスター(S),ヘルズローリング(H),ダークアームブレード(D1or3),ダークネスオーバーロード(O)
				// HP: MAX～1/2 -> D1(60%), D3(20%), H(10%), S(10%)
				// HP: 1/2～1/4 -> D1(30%), D3(20%), H(20%), S(25%), O(5%)
				// HP: 1/4～0   -> D1(10%), D3(25%), H(20%), S(35%), O(10%)
				prob := rand.Intn(100)
				halfHP := e.pm.HPMax / 2
				quarterHP := e.pm.HPMax / 4
				var d1Line, d3Line, hLine, sLine int
				if e.pm.HP > halfHP {
					d1Line = 60
					d3Line = 80
					hLine = 90
					sLine = 100
				} else if e.pm.HP > quarterHP {
					d1Line = 30
					d3Line = 50
					hLine = 70
					sLine = 95
				} else {
					d1Line = 10
					d3Line = 35
					hLine = 55
					sLine = 90
				}
				if prob < d1Line {
					e.nextState = forteActTypeDarkArmBlade1
				} else if prob < d3Line {
					e.nextState = forteActTypeDarkArmBlade3
				} else if prob < hLine {
					e.nextState = forteActTypeHellsRolling
				} else if prob < sLine {
					e.nextState = forteActTypeShooting
				} else {
					e.nextState = forteActTypeDarknessOverload
				}
			}
			return e.stateChange(forteActTypeStand)
		}
	case forteActTypeShooting:
		initWait := 4 * forteDelays[forteActTypeShooting]
		if e.count < 120+initWait {
			// 攻撃フェーズ
			// 3回攻撃する
			if e.count%40 == initWait {
				logger.Debug("Forte Shooting Buster Attack at %d", e.count)
				for i := 0; i < 3; i++ {
					e.atkIDs = append(e.atkIDs, e.animMgr.SkillAnimNew(skill.Get(resources.SkillForteShootingBuster, skillcore.Argument{
						OwnerID:    e.pm.ObjectID,
						Power:      forteAtkPower[e.state],
						TargetType: damage.TargetPlayer,
					}, e.animMgr)))
				}
			}
		} else {
			// 終了チェックフェーズ
			end := true
			for _, id := range e.atkIDs {
				if e.animMgr.IsAnimProcessing(id) {
					end = false
					break
				}
			}
			if end {
				return e.clearState()
			}
		}
	case forteActTypeHellsRolling:
		if e.count == 0 {
			e.atkIDs = []string{}

			// Move to attack position
			targetPos := point.Point{X: 5, Y: 1}
			if !targetPos.Equal(e.pm.Pos) {
				e.targetPos = targetPos
				e.nextState = forteActTypeHellsRolling
				return e.stateChange(forteActTypeMove)
			}
		}

		if e.count == 7*forteDelays[forteActTypeHellsRolling] {
			logger.Debug("Forte Hells Rolling Attack 1st")
			e.atkIDs = append(e.atkIDs, e.animMgr.SkillAnimNew(skill.Get(resources.SkillForteHellsRollingUp, skillcore.Argument{
				OwnerID:    e.pm.ObjectID,
				Power:      forteAtkPower[e.state],
				TargetType: damage.TargetPlayer,
			}, e.animMgr)))
		}

		if e.count == 7*forteDelays[forteActTypeHellsRolling]+30 {
			logger.Debug("Forte Hells Rolling Attack 2st")
			e.atkIDs = append(e.atkIDs, e.animMgr.SkillAnimNew(skill.Get(resources.SkillForteHellsRollingDown, skillcore.Argument{
				OwnerID:    e.pm.ObjectID,
				Power:      forteAtkPower[e.state],
				TargetType: damage.TargetPlayer,
			}, e.animMgr)))
		}

		if len(e.atkIDs) > 0 {
			end := true
			for _, id := range e.atkIDs {
				if e.animMgr.IsAnimProcessing(id) {
					end = false
					break
				}
			}
			if end {
				return e.clearState()
			}
		}
	case forteActTypeDarkArmBlade1:
		if e.count == 0 && !e.isTargetPosMoved {
			e.isTargetPosMoved = true

			// Move to attack position
			objs := e.animMgr.ObjAnimGetObjs(objanim.Filter{ObjType: objanim.ObjTypePlayer})
			if len(objs) == 0 {
				// エラー処理
				logger.Info("Failed to get player position")
				return e.clearState()
			}
			targetPos := point.Point{X: objs[0].Pos.X + 1, Y: objs[0].Pos.Y}
			if !targetPos.Equal(e.pm.Pos) {
				e.targetPos = targetPos
				e.nextState = forteActTypeDarkArmBlade1
				return e.stateChange(forteActTypeMove)
			}
		}
		if e.count == 2*forteDelays[forteActTypeDarkArmBlade1] {
			logger.Debug("Forte Dark Arm Blade 1st Attack")
			e.animMgr.SkillAnimNew(skill.Get(resources.SkillForteDarkArmBladeType1, skillcore.Argument{
				OwnerID:    e.pm.ObjectID,
				Power:      forteAtkPower[e.state],
				TargetType: damage.TargetPlayer,
			}, e.animMgr))
		}

		if e.count == 5*forteDelays[forteActTypeDarkArmBlade1] {
			return e.clearState()
		}
	case forteActTypeDarkArmBlade3:
		if e.count == 0 && !e.isTargetPosMoved {
			e.isTargetPosMoved = true

			// Move to attack position
			objs := e.animMgr.ObjAnimGetObjs(objanim.Filter{ObjType: objanim.ObjTypePlayer})
			if len(objs) == 0 {
				// エラー処理
				logger.Info("Failed to get player position")
				return e.clearState()
			}
			var targetPos point.Point
			switch e.bladeAtkCount {
			case 0, 2:
				e.isCharReverse = false
				targetPos = point.Point{X: objs[0].Pos.X + 1, Y: objs[0].Pos.Y}
			case 1:
				e.isCharReverse = true
				targetPos = point.Point{X: objs[0].Pos.X - 1, Y: objs[0].Pos.Y}
			}

			if !targetPos.Equal(e.pm.Pos) {
				e.targetPos = targetPos
				e.nextState = forteActTypeDarkArmBlade3
				e.waitCount = 15
				return e.stateChange(forteActTypeMove)
			}
		}

		if e.count == 1*forteDelays[forteActTypeDarkArmBlade3] {
			logger.Debug("Forte Dark Arm Blade %d times Attack", e.bladeAtkCount+1)
			skillType := resources.SkillForteDarkArmBladeType1
			if e.bladeAtkCount == 1 {
				skillType = resources.SkillForteDarkArmBladeType2
			}
			e.atkIDs = []string{
				e.animMgr.SkillAnimNew(
					skill.Get(
						skillType,
						skillcore.Argument{
							OwnerID:    e.pm.ObjectID,
							Power:      forteAtkPower[e.state],
							TargetType: damage.TargetPlayer,
						},
						e.animMgr,
					),
				),
			}
		}

		if len(e.atkIDs) > 0 {
			if !e.animMgr.IsAnimProcessing(e.atkIDs[0]) {
				e.bladeAtkCount++
				if e.bladeAtkCount == 3 {
					// 終了
					return e.clearState()
				} else {
					e.nextState = forteActTypeDarkArmBlade3
					e.waitCount = 5
					e.isTargetPosMoved = false
					e.atkIDs = []string{}
					return e.stateChange(forteActTypeStand)
				}
			}
		}
	case forteActTypeDarknessOverload:
		if e.count == 0 && !e.isTargetPosMoved {
			e.isTargetPosMoved = true

			// Move to attack position
			targetPos := point.Point{X: 2, Y: 1}
			if !targetPos.Equal(e.pm.Pos) {
				e.targetPos = targetPos
				e.nextState = forteActTypeDarknessOverload
				e.waitCount = 5
				return e.stateChange(forteActTypeMove)
			}
		}

		if e.count == 7*forteDelays[forteActTypeDarknessOverload] {
			logger.Debug("Forte Darkness Overload Attack")
			e.atkIDs = []string{
				e.animMgr.SkillAnimNew(skill.Get(resources.SkillForteDarknessOverload, skillcore.Argument{
					OwnerID:    e.pm.ObjectID,
					Power:      forteAtkPower[e.state],
					TargetType: damage.TargetPlayer,
				}, e.animMgr)),
			}
		}

		if len(e.atkIDs) > 0 {
			if !e.animMgr.IsAnimProcessing(e.atkIDs[0]) {
				return e.clearState()
			}
		}
	case forteActTypeDamage:
		if e.count == 2*forteDelays[forteActTypeDamage] {
			return e.clearState()
		}
	}

	e.count++
	return false, nil
}

func (e *enemyForte) Draw() {
	if e.pm.InvincibleCount/5%2 != 0 {
		return
	}

	// Show Enemy Images
	view := battlecommon.ViewPos(e.pm.Pos)
	img := e.getCurrentImagePointer()

	ofsY := -20
	if e.state == forteActTypeStand {
		ofsY -= math.MountainIndex(e.count/10%5, 5)
	}

	opt := dxlib.DrawRotaGraphOption{}
	if e.isCharReverse {
		t := int32(dxlib.TRUE)
		opt.ReverseXFlag = &t
	}

	dxlib.DrawRotaGraph(view.X, view.Y+ofsY, 1, 0, *img, true, opt)

	drawParalysis(view.X, view.Y+ofsY, *img, e.pm.ParalyzedCount)

	// Show HP
	if e.pm.HP > 0 {
		draw.Number(view.X, view.Y+40, e.pm.HP, draw.NumberOption{
			Color:    draw.NumberColorWhiteSmall,
			Centered: true,
		})
	}
}

func (e *enemyForte) DamageProc(dm *damage.Damage) bool {
	if damageProc(dm, &e.pm) {
		if dm.StrengthType == damage.StrengthNone {
			return true
		}

		e.state = forteActTypeDamage
		if dm.StrengthType == damage.StrengthHigh {
			e.pm.InvincibleCount = battlecommon.PlayerDefaultInvincibleTime
		}
		e.count = 0
		return true
	}

	return false
}

func (e *enemyForte) GetParam() objanim.Param {
	return objanim.Param{
		Param: anim.Param{
			ObjID: e.pm.ObjectID,
			Pos:   e.pm.Pos,
		},
		HP: e.pm.HP,
	}
}

func (e *enemyForte) GetObjectType() int {
	return objanim.ObjTypeEnemy
}

func (e *enemyForte) MakeInvisible(count int) {
	e.pm.InvincibleCount = count
}

func (e *enemyForte) AddBarrier(hp int) {}

func (e *enemyForte) SetCustomGaugeMax() {}

func (e *enemyForte) getCurrentImagePointer() *int {
	if e.count == 0 {
		return &e.images[forteActTypeStand][0]
	}

	n := (e.count / forteDelays[e.state])
	if n >= len(e.images[e.state]) {
		n = len(e.images[e.state]) - 1
	}
	return &e.images[e.state][n]
}

func (e *enemyForte) stateChange(next int) (bool, error) {
	logger.Info("change forte state to %d", next)
	e.state = next
	e.count = 0

	return false, nil
}

func (e *enemyForte) clearState() (bool, error) {
	e.waitCount = 20
	e.nextState = forteActTypeMove
	e.moveNum = 3 + rand.Intn(3)
	e.targetPos = emptyPos
	e.isTargetPosMoved = false
	e.bladeAtkCount = 0
	e.atkIDs = []string{}
	e.isCharReverse = false

	return e.stateChange(forteActTypeStand)
}
