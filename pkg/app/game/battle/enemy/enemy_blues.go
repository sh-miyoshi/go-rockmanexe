package enemy

import (
	"math/rand"

	"github.com/cockroachdb/errors"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common/deleteanim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/system"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	bluesActTypeStand = iota
	bluesActTypeMove
	bluesActTypeWideSword
	bluesActTypeFighterSword
	bluesActTypeSonicBoom
	bluesActTypeDeltaRayEdge
	bluesActTypeDeltaRayEdgeEnd
	bluesActTypeBehindSlash
	bluesActTypeDamage

	bluesActTypeMax
)

const (
	bluesShieldTime = 20
)

var (
	bluesDelays      = [bluesActTypeMax]int{1, 2, 4, 4, 4, 1, 4, 1}
	bluesAtkPower    = [bluesActTypeMax]uint{0, 0, 60, 60, 60, 90, 0, 60}
	bluesShieldDelay = 4
)

type enemyBlues struct {
	pm             EnemyParam
	state          int
	count          int
	waitCount      int
	nextState      int
	targetPoses    [4]point.Point
	isTargetMoved  bool
	moveNum        int
	images         [bluesActTypeMax][]int
	atkCount       int
	atkID          string
	isCharReverse  bool
	isEdgeEffectOn bool
	edgeEndPos     point.Point
	imgShields     []int
	shieldCount    int
}

func (e *enemyBlues) Init(objID string) error {
	e.pm.ObjectID = objID
	e.state = bluesActTypeStand
	e.count = 0
	e.waitCount = 20
	e.nextState = bluesActTypeMove
	e.targetPoses = [4]point.Point{emptyPos, emptyPos, emptyPos, emptyPos}
	e.moveNum = 2
	e.atkCount = 0
	e.isCharReverse = false
	e.isEdgeEffectOn = false
	e.atkID = ""
	e.isTargetMoved = false
	e.shieldCount = 0

	// Load Images
	name, ext := GetStandImageFile(IDBlues)

	fname := name + "_all" + ext
	tmp := make([]int, 36)
	if res := dxlib.LoadDivGraph(fname, 36, 7, 6, 170, 156, tmp); res == -1 {
		return errors.Newf("failed to load image %s", fname)
	}

	releases := [36]int{}
	for i := 0; i < 36; i++ {
		releases[i] = i
	}

	e.images[bluesActTypeStand] = make([]int, 1)
	e.images[bluesActTypeStand][0] = tmp[0]
	releases[0] = -1

	e.images[bluesActTypeMove] = make([]int, 4)
	for i := 0; i < 4; i++ {
		e.images[bluesActTypeMove][i] = tmp[i]
		releases[i] = -1
	}

	e.images[bluesActTypeWideSword] = make([]int, 6)
	e.images[bluesActTypeFighterSword] = make([]int, 6)
	e.images[bluesActTypeDeltaRayEdge] = make([]int, 6)
	for i := 0; i < 6; i++ {
		e.images[bluesActTypeWideSword][i] = tmp[i+7]
		e.images[bluesActTypeFighterSword][i] = tmp[i+7]
		e.images[bluesActTypeDeltaRayEdge][i] = tmp[i+7]
		releases[i+7] = -1
	}

	e.images[bluesActTypeDeltaRayEdgeEnd] = make([]int, 1)
	e.images[bluesActTypeDeltaRayEdgeEnd][0] = tmp[7]
	e.images[bluesActTypeBehindSlash] = make([]int, 1)
	e.images[bluesActTypeBehindSlash][0] = tmp[7]

	// 使わないイメージを削除
	for i, r := range releases {
		if r != -1 {
			dxlib.DeleteGraph(tmp[i])
		}
	}

	e.imgShields = make([]int, 7)
	fname = config.ImagePath + "battle/skill/ブルース_シールド.png"
	if res := dxlib.LoadDivGraph(fname, 7, 7, 1, 64, 116, e.imgShields); res == -1 {
		return errors.Newf("failed to load image %s", fname)
	}

	return nil
}

func (e *enemyBlues) End() {
	// Delete Images
	for _, imgs := range e.images {
		for _, img := range imgs {
			dxlib.DeleteGraph(img)
		}
	}
	for _, img := range e.imgShields {
		dxlib.DeleteGraph(img)
	}
}

func (e *enemyBlues) Update() (bool, error) {
	if e.pm.HP <= 0 {
		// Delete Animation
		img := e.getCurrentImagePointer()
		deleteanim.New(*img, e.pm.Pos, false)
		localanim.AnimNew(effect.Get(resources.EffectTypeExplode, e.pm.Pos, 0))
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
	case bluesActTypeStand:
		e.waitCount--
		if e.waitCount <= 0 {
			return e.stateChange(e.nextState)
		}
	case bluesActTypeMove:
		if e.count == 3*bluesDelays[bluesActTypeMove] {
			if !e.targetPoses[e.atkCount].Equal(emptyPos) {
				if !battlecommon.MoveObjectDirect(
					&e.pm.Pos,
					e.targetPoses[e.atkCount],
					-1, // プレイヤーのパネルでも移動可能
					true,
					field.GetPanelInfo,
				) {
					// 移動に失敗したら、移動からやり直し
					logger.Debug("Forte move failed. retry")
					return e.clearState()
				}
				e.targetPoses[e.atkCount] = emptyPos
				if e.waitCount == 0 {
					e.waitCount = 20
				}
				e.isTargetMoved = true
				return e.stateChange(forteActTypeStand)
			}

			moveRandom(&e.pm.Pos)
			e.waitCount = 40

			e.moveNum--
			if e.moveNum <= 0 {
				if debugFlag {
					e.moveNum = 3
					e.nextState = bluesActTypeBehindSlash
					return e.stateChange(bluesActTypeStand)
				}

				e.moveNum = rand.Intn(2) + 3

				// Action process
				// 攻撃処理(HPの減り具合から乱数で攻撃決定)
				// ワイドソード(W),　ファイターソード(F),ビハインドスラッシュ(B),デルタレイエッジ(D)
				// HP: MAX～1/2 -> W(45%), F(45%), B(10%), D(0%)
				// HP: 1/2～1/4 -> W(35%), F(40%), B(20%), D(5%)
				// HP: 1/4～0   -> W(35%), F(35%), B(20%), D(10%)
				prob := rand.Intn(100)
				halfHP := e.pm.HPMax / 2
				quarterHP := e.pm.HPMax / 4
				var wLine, fLine, bLine int
				if e.pm.HP > halfHP {
					wLine = 45
					fLine = 90
					bLine = 100
				} else if e.pm.HP > quarterHP {
					wLine = 35
					fLine = 70
					bLine = 95
				} else {
					wLine = 35
					fLine = 70
					bLine = 90
				}
				if prob < wLine {
					e.nextState = bluesActTypeWideSword
				} else if prob < fLine {
					e.nextState = bluesActTypeFighterSword
				} else if prob < bLine {
					e.nextState = bluesActTypeBehindSlash
				} else {
					e.nextState = bluesActTypeDeltaRayEdge
				}
			}

			return e.stateChange(bluesActTypeStand)
		}
	case bluesActTypeWideSword:
		if e.count == 0 && !e.isTargetMoved {
			// Move to attack position
			objs := localanim.ObjAnimGetObjs(objanim.Filter{ObjType: objanim.ObjTypePlayer})
			if len(objs) == 0 {
				// エラー処理
				logger.Info("Failed to get player position")
				return e.clearState()
			}
			// 一旦右上だが、移動できなければ右下も候補にする
			targetPos := point.Point{X: objs[0].Pos.X + 1, Y: objs[0].Pos.Y - 1}
			if !battlecommon.MoveObjectDirect(&e.pm.Pos, targetPos, -1, false, field.GetPanelInfo) {
				targetPos = point.Point{X: objs[0].Pos.X + 1, Y: objs[0].Pos.Y + 1}
			}

			if !targetPos.Equal(e.pm.Pos) {
				e.targetPoses[0] = targetPos
				e.nextState = bluesActTypeWideSword
				return e.stateChange(bluesActTypeMove)
			}
		}

		if e.count == 1*bluesDelays[bluesActTypeWideSword] {
			logger.Debug("Blues Wide Sword Attack")
			localanim.AnimNew(skill.Get(resources.SkillWideSword, skillcore.Argument{
				OwnerID:    e.pm.ObjectID,
				Power:      bluesAtkPower[e.state],
				TargetType: damage.TargetPlayer,
				IsReverse:  !e.isCharReverse,
			}))
		}

		if e.count == 6*bluesDelays[bluesActTypeWideSword] {
			return e.clearState()
		}
	case bluesActTypeFighterSword:
		if e.count == 0 && !e.isTargetMoved {
			// Move to attack position
			objs := localanim.ObjAnimGetObjs(objanim.Filter{ObjType: objanim.ObjTypePlayer})
			if len(objs) == 0 {
				// エラー処理
				logger.Info("Failed to get player position")
				return e.clearState()
			}
			tx := 0
			for x := 0; x < battlecommon.FieldNum.X; x++ {
				if field.GetPanelInfo(point.Point{X: x, Y: objs[0].Pos.Y}).Type == battlecommon.PanelTypeEnemy {
					tx = x
					break
				}
			}

			targetPos := point.Point{X: tx, Y: objs[0].Pos.Y}
			if !targetPos.Equal(e.pm.Pos) {
				e.targetPoses[0] = targetPos
				e.nextState = bluesActTypeFighterSword
				return e.stateChange(bluesActTypeMove)
			}
		}

		if e.count == 1*bluesDelays[bluesActTypeFighterSword] {
			logger.Debug("Blues Fighter Sword Attack")
			localanim.AnimNew(skill.Get(resources.SkillFighterSword, skillcore.Argument{
				OwnerID:    e.pm.ObjectID,
				Power:      bluesAtkPower[e.state],
				TargetType: damage.TargetPlayer,
				IsReverse:  true,
			}))
		}

		if e.count == 6*bluesDelays[bluesActTypeWideSword] {
			return e.clearState()
		}
	case bluesActTypeDeltaRayEdge:
		if e.count == 0 && !e.isEdgeEffectOn {
			// 移動先を決める
			objs := localanim.ObjAnimGetObjs(objanim.Filter{ObjType: objanim.ObjTypePlayer})
			if len(objs) == 0 {
				// エラー処理
				logger.Info("Failed to get player position")
				return e.clearState()
			}
			e.targetPoses[0] = point.Point{X: objs[0].Pos.X + 1, Y: objs[0].Pos.Y - 1}
			e.targetPoses[1] = point.Point{X: objs[0].Pos.X - 1, Y: objs[0].Pos.Y}
			e.targetPoses[2] = point.Point{X: objs[0].Pos.X + 1, Y: objs[0].Pos.Y + 1}
			// 最後は最初の場所に戻る
			e.targetPoses[3] = e.pm.Pos
			for i := 0; i < 3; i++ {
				if !battlecommon.MoveObjectDirect(&e.pm.Pos, e.targetPoses[i], -1, false, field.GetPanelInfo) {
					logger.Debug("Blues DeltaRayEdge move failed at %d time.", i)
					return e.clearState()
				}
			}
			e.edgeEndPos = objs[0].Pos

			e.isEdgeEffectOn = true
			localanim.AnimNew(effect.Get(resources.EffectTypeSpecialStart, e.pm.Pos, 0))
			e.nextState = bluesActTypeDeltaRayEdge
			e.waitCount = 10
			return e.stateChange(bluesActTypeStand)
		}

		if e.count == 0 && !e.targetPoses[e.atkCount].Equal(emptyPos) {
			logger.Debug("Blues move to %s", e.targetPoses[e.atkCount].String())
			if e.atkCount == 1 {
				e.isCharReverse = true
			} else {
				e.isCharReverse = false
			}
			e.waitCount = 2
			e.nextState = bluesActTypeDeltaRayEdge
			return e.stateChange(bluesActTypeMove)
		}

		if e.count == 1*bluesDelays[bluesActTypeDeltaRayEdge] {
			logger.Debug("Blues DeltaRayEdge %d times Attack", e.atkCount+1)
			skillID := resources.SkillNonEffectWideSword
			isReverse := true
			if e.atkCount == 1 {
				isReverse = false
			} else if e.atkCount == 2 {
				skillID = resources.SkillWideSword
			}

			e.atkID = localanim.AnimNew(skill.Get(skillID, skillcore.Argument{
				OwnerID:    e.pm.ObjectID,
				Power:      bluesAtkPower[e.state],
				TargetType: damage.TargetPlayer,
				IsReverse:  isReverse,
			}))
		}

		if e.atkID != "" {
			if !localanim.AnimIsProcessing(e.atkID) {
				e.atkCount++
				if e.atkCount == 3 {
					e.waitCount = 1
					e.nextState = bluesActTypeDeltaRayEdgeEnd
					return e.stateChange(bluesActTypeMove)
				} else {
					e.nextState = bluesActTypeDeltaRayEdge
					e.waitCount = 2
					e.atkID = ""
					return e.stateChange(forteActTypeStand)
				}
			}
		}
	case bluesActTypeDeltaRayEdgeEnd:
		if e.count == 0 {
			e.atkID = localanim.AnimNew(effect.Get(resources.EffectTypeDeltaRayEdge, e.edgeEndPos, 0))
			sound.On(resources.SEDeltaRayEdgeEnd)
		}

		if !localanim.AnimIsProcessing(e.atkID) {
			return e.clearState()
		}
	case bluesActTypeSonicBoom:
		system.SetError("WIP: Blues SonicBoom is not implemented yet")
	case bluesActTypeBehindSlash:
		if e.count == 0 {
			e.shieldCount = bluesShieldTime
		}
		if e.shieldCount > 0 {
			e.shieldCount--
			if e.shieldCount == 0 {
				return e.clearState()
			}
		}
	}

	e.count++
	return false, nil
}

func (e *enemyBlues) Draw() {
	// Show Enemy Images
	view := battlecommon.ViewPos(e.pm.Pos)
	img := e.getCurrentImagePointer()
	ofs := [bluesActTypeMax]point.Point{
		{X: -5, Y: -20},  // Stand
		{X: -5, Y: -20},  // Move
		{X: -20, Y: -20}, // WideSword
		{X: -20, Y: -20}, // FighterSword
		{X: -20, Y: -20}, // SonicBoom
		{X: -20, Y: -20}, // DeltaRayEdge
		{X: -20, Y: -20}, // DeltaRayEdgeEnd
		{X: -20, Y: -20}, // BehindSlash
		{X: 0, Y: 0},     // Damage
	}

	// デフォルトは逆向き
	flag := int32(dxlib.TRUE)
	opt := dxlib.DrawRotaGraphOption{
		ReverseXFlag: &flag,
	}
	if e.isCharReverse {
		opt = dxlib.DrawRotaGraphOption{}
	}

	state := e.state
	if e.count == 0 {
		state = bluesActTypeStand
	}

	dxlib.DrawRotaGraph(view.X+ofs[state].X, view.Y+ofs[state].Y, 1, 0, *img, true, opt)

	drawParalysis(view.X+ofs[state].X, view.Y+ofs[state].Y, *img, e.pm.ParalyzedCount, opt)

	if e.shieldCount > 0 {
		n := (bluesShieldTime - e.shieldCount) / bluesShieldDelay
		if n >= 0 && n < len(e.imgShields) {
			dxlib.DrawRotaGraph(view.X-30, view.Y+8, 1, 0, e.imgShields[n], true, opt)
		}
	}

	// Show HP
	if e.pm.HP > 0 {
		draw.Number(view.X, view.Y+40, e.pm.HP, draw.NumberOption{
			Color:    draw.NumberColorWhiteSmall,
			Centered: true,
		})
	}
}

func (e *enemyBlues) DamageProc(dm *damage.Damage) bool {
	if e.shieldCount > 0 {
		// シールド中は反撃する
		if dm.StrengthType != damage.StrengthNone {
			// 背後に回ってワイドソード
			objs := localanim.ObjAnimGetObjs(objanim.Filter{ObjType: objanim.ObjTypePlayer})
			if len(objs) == 0 {
				// エラー処理
				logger.Info("Failed to get player position")
				e.clearState()
				return true
			}
			e.targetPoses[0] = point.Point{X: objs[0].Pos.X - 1, Y: objs[0].Pos.Y}
			e.nextState = bluesActTypeWideSword
			e.shieldCount = 0
			e.isCharReverse = true
			e.waitCount = 2
			e.stateChange(bluesActTypeMove)
			return true
		} else {
			localanim.AnimNew(effect.Get(resources.EffectTypeBlock, e.pm.Pos, 5))
			return true
		}
	}

	if damageProc(dm, &e.pm) {
		if dm.StrengthType == damage.StrengthNone {
			return true
		}

		e.state = bluesActTypeDamage
		if dm.StrengthType == damage.StrengthHigh {
			e.pm.InvincibleCount = battlecommon.PlayerDefaultInvincibleTime
		}
		e.count = 0
		return true
	}

	return false
}

func (e *enemyBlues) GetParam() objanim.Param {
	return objanim.Param{
		Param: anim.Param{
			ObjID:    e.pm.ObjectID,
			Pos:      e.pm.Pos,
			DrawType: anim.DrawTypeObject,
		},
		HP: e.pm.HP,
	}
}

func (e *enemyBlues) GetObjectType() int {
	return objanim.ObjTypeEnemy
}

func (e *enemyBlues) MakeInvisible(count int) {
	e.pm.InvincibleCount = count
}

func (e *enemyBlues) AddBarrier(hp int) {}

func (e *enemyBlues) getCurrentImagePointer() *int {
	if e.count == 0 {
		return &e.images[bluesActTypeStand][0]
	}

	n := (e.count / bluesDelays[e.state])
	if n >= len(e.images[e.state]) {
		n = len(e.images[e.state]) - 1
	}
	return &e.images[e.state][n]
}

func (e *enemyBlues) stateChange(next int) (bool, error) {
	logger.Info("change blues state to %d", next)
	e.state = next
	e.count = 0

	return false, nil
}

func (e *enemyBlues) clearState() (bool, error) {
	e.waitCount = 20
	e.nextState = forteActTypeMove
	e.moveNum = 3 + rand.Intn(3)
	e.targetPoses = [4]point.Point{emptyPos, emptyPos, emptyPos, emptyPos}
	e.atkCount = 0
	e.isCharReverse = false
	e.isEdgeEffectOn = false
	e.atkID = ""
	e.isTargetMoved = false
	e.shieldCount = 0

	return e.stateChange(forteActTypeStand)
}
