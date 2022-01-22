package player

import (
	"fmt"
	"math/rand"
	"sort"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/chip"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/inputs"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/player"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
)

const (
	NextActNone int = iota
	NextActChipSelect
	NextActLose
)

type act struct {
	ID         string
	MoveDirect int
	Charged    bool
	ShotPower  uint

	typ       int
	count     int
	keepCount int
	pPos      *common.Point
}

// BattlePlayer ...
type BattlePlayer struct {
	ID            string
	Pos           common.Point
	HP            uint
	HPMax         uint
	ChargeCount   uint
	GaugeCount    uint
	ShotPower     uint
	ChipFolder    []player.ChipInfo
	SelectedChips []player.ChipInfo
	NextAction    int
	EnableAct     bool
	MoveNum       int
	DamageNum     int

	act             act
	invincibleCount int
}

var (
	imgPlayers    [battlecommon.PlayerActMax][]int32
	imgDelays     = [battlecommon.PlayerActMax]int{1, 2, 2, 6, 3, 4, 1, 4, 3}
	imgHPFrame    int32
	imgGaugeFrame int32
	imgGaugeMax   []int32
	imgCharge     [2][]int32
)

// New ...
func New(plyr *player.Player) (*BattlePlayer, error) {
	logger.Info("Initialize battle player data")

	res := BattlePlayer{
		ID:        uuid.New().String(),
		HP:        plyr.HP,
		HPMax:     plyr.HP, // TODO HPは引き継がない
		Pos:       common.Point{X: 1, Y: 1},
		ShotPower: plyr.ShotPower,
		EnableAct: true,
	}
	res.act.typ = -1
	res.act.pPos = &res.Pos

	for _, c := range plyr.ChipFolder {
		res.ChipFolder = append(res.ChipFolder, c)
	}
	// Shuffle folder
	for i := 0; i < 10; i++ {
		for j := 0; j < len(res.ChipFolder); j++ {
			n := rand.Intn(len(res.ChipFolder))
			res.ChipFolder[j], res.ChipFolder[n] = res.ChipFolder[n], res.ChipFolder[j]
		}
	}

	logger.Debug("Player info: %+v", res)

	fname := common.ImagePath + "battle/character/player_move.png"
	imgPlayers[battlecommon.PlayerActMove] = make([]int32, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 100, 100, imgPlayers[battlecommon.PlayerActMove]); res == -1 {
		return nil, fmt.Errorf("failed to load player move image: %s", fname)
	}

	fname = common.ImagePath + "battle/character/player_damaged.png"
	imgPlayers[battlecommon.PlayerActDamage] = make([]int32, 6)
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 100, 100, imgPlayers[battlecommon.PlayerActDamage]); res == -1 {
		return nil, fmt.Errorf("failed to load player damage image: %s", fname)
	}
	// 1 -> 2,3  2-4 3-5
	imgPlayers[battlecommon.PlayerActDamage][4] = imgPlayers[battlecommon.PlayerActDamage][2]
	imgPlayers[battlecommon.PlayerActDamage][5] = imgPlayers[battlecommon.PlayerActDamage][3]
	imgPlayers[battlecommon.PlayerActDamage][2] = imgPlayers[battlecommon.PlayerActDamage][1]
	imgPlayers[battlecommon.PlayerActDamage][3] = imgPlayers[battlecommon.PlayerActDamage][1]

	fname = common.ImagePath + "battle/character/player_shot.png"
	imgPlayers[battlecommon.PlayerActShot] = make([]int32, 6)
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 180, 100, imgPlayers[battlecommon.PlayerActShot]); res == -1 {
		return nil, fmt.Errorf("failed to load player shot image: %s", fname)
	}

	fname = common.ImagePath + "battle/character/player_cannon.png"
	imgPlayers[battlecommon.PlayerActCannon] = make([]int32, 6)
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 100, 100, imgPlayers[battlecommon.PlayerActCannon]); res == -1 {
		return nil, fmt.Errorf("failed to load player cannon image: %s", fname)
	}

	fname = common.ImagePath + "battle/character/player_sword.png"
	imgPlayers[battlecommon.PlayerActSword] = make([]int32, 7)
	if res := dxlib.LoadDivGraph(fname, 7, 7, 1, 128, 128, imgPlayers[battlecommon.PlayerActSword]); res == -1 {
		return nil, fmt.Errorf("failed to load player sword image: %s", fname)
	}

	fname = common.ImagePath + "battle/character/player_bomb.png"
	imgPlayers[battlecommon.PlayerActBomb] = make([]int32, 7)
	if res := dxlib.LoadDivGraph(fname, 5, 5, 1, 100, 114, imgPlayers[battlecommon.PlayerActBomb]); res == -1 {
		return nil, fmt.Errorf("failed to load player bomb image: %s", fname)
	}
	imgPlayers[battlecommon.PlayerActBomb][5] = imgPlayers[battlecommon.PlayerActBomb][4]
	imgPlayers[battlecommon.PlayerActBomb][6] = imgPlayers[battlecommon.PlayerActBomb][4]

	fname = common.ImagePath + "battle/character/player_buster.png"
	imgPlayers[battlecommon.PlayerActBuster] = make([]int32, 6)
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 180, 100, imgPlayers[battlecommon.PlayerActBuster]); res == -1 {
		return nil, fmt.Errorf("failed to load player buster image: %s", fname)
	}

	fname = common.ImagePath + "battle/character/player_pick.png"
	imgPlayers[battlecommon.PlayerActPick] = make([]int32, 6)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 96, 124, imgPlayers[battlecommon.PlayerActPick]); res == -1 {
		return nil, fmt.Errorf("failed to load player pick image: %s", fname)
	}
	imgPlayers[battlecommon.PlayerActPick][4] = imgPlayers[battlecommon.PlayerActPick][3]
	imgPlayers[battlecommon.PlayerActPick][5] = imgPlayers[battlecommon.PlayerActPick][3]

	fname = common.ImagePath + "battle/character/player_throw.png"
	imgPlayers[battlecommon.PlayerActThrow] = make([]int32, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 97, 115, imgPlayers[battlecommon.PlayerActThrow]); res == -1 {
		return nil, fmt.Errorf("failed to load player throw image: %s", fname)
	}

	fname = common.ImagePath + "battle/hp_frame.png"
	imgHPFrame = dxlib.LoadGraph(fname)
	if imgHPFrame < 0 {
		return nil, fmt.Errorf("failed to read hp frame image %s", fname)
	}
	fname = common.ImagePath + "battle/gauge.png"
	imgGaugeFrame = dxlib.LoadGraph(fname)
	if imgGaugeFrame < 0 {
		return nil, fmt.Errorf("failed to read gauge frame image %s", fname)
	}
	fname = common.ImagePath + "battle/gauge_max.png"
	imgGaugeMax = make([]int32, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 1, 4, 288, 30, imgGaugeMax); res == -1 {
		return nil, fmt.Errorf("failed to read gauge max image %s", fname)
	}

	fname = common.ImagePath + "battle/skill/charge.png"
	tmp := make([]int32, 16)
	if res := dxlib.LoadDivGraph(fname, 16, 8, 2, 158, 150, tmp); res == -1 {
		return nil, fmt.Errorf("failed to load image %s", fname)
	}
	for i := 0; i < 8; i++ {
		imgCharge[0] = append(imgCharge[0], tmp[i])
		imgCharge[1] = append(imgCharge[1], tmp[i+8])
	}

	logger.Info("Successfully initialized battle player data")
	return &res, nil
}

// End ...
func (p *BattlePlayer) End() {
	logger.Info("Cleanup battle player data")

	for i := 0; i < battlecommon.PlayerActMax; i++ {
		for j := 0; j < len(imgPlayers[i]); j++ {
			dxlib.DeleteGraph(imgPlayers[i][j])
			imgPlayers[i][j] = -1
		}
	}
	dxlib.DeleteGraph(imgHPFrame)
	imgHPFrame = -1
	dxlib.DeleteGraph(imgGaugeFrame)
	imgGaugeFrame = -1
	for _, img := range imgGaugeMax {
		dxlib.DeleteGraph(img)
	}
	for i := 0; i < 2; i++ {
		for _, img := range imgCharge[i] {
			dxlib.DeleteGraph(img)
		}
		imgCharge[i] = []int32{}
	}

	logger.Info("Successfully cleanuped battle player data")
}

// Draw ...
func (p *BattlePlayer) Draw() {
	if p.invincibleCount/5%2 != 0 {
		return
	}

	view := battlecommon.ViewPos(p.Pos)
	img := p.act.GetImage()
	dxlib.DrawRotaGraph(view.X, view.Y, 1, 0, img, dxlib.TRUE)

	// Show charge image
	if p.ChargeCount > battlecommon.ChargeViewDelay {
		n := 0
		if p.ChargeCount > battlecommon.ChargeTime {
			n = 1
		}
		imgNo := int(p.ChargeCount/4) % len(imgCharge[n])
		dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_ALPHA, 224)
		dxlib.DrawRotaGraph(view.X, view.Y, 1, 0, imgCharge[n][imgNo], dxlib.TRUE)
		dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_NOBLEND, 0)
	}

	// Show selected chip icons
	n := len(p.SelectedChips)
	if n > 0 {
		// TODO Show chip info

		const px = 3
		max := n * px
		for i := 0; i < n; i++ {
			x := field.PanelSize.X*p.Pos.X + field.PanelSize.X/2 - 2 + int32((i*px)-max)
			y := field.DrawPanelTopY + field.PanelSize.Y*p.Pos.Y - 10 - 81 + int32((i*px)-max)
			dxlib.DrawBox(x-1, y-1, x+29, y+29, 0x000000, dxlib.FALSE)
			// draw from the end
			dxlib.DrawGraph(x, y, chip.GetIcon(p.SelectedChips[n-1-i].ID, true), dxlib.TRUE)
		}
	}
}

func (p *BattlePlayer) DrawFrame(xShift bool, showGauge bool) {
	x := int32(7)
	y := int32(5)
	if xShift {
		x += 235
	}

	// Show HP
	dxlib.DrawGraph(x, y, imgHPFrame, dxlib.TRUE)
	draw.Number(x+2, y+2, int32(p.HP), draw.NumberOption{RightAligned: true, Length: 4})

	// Show Custom Gauge
	if showGauge {
		if p.GaugeCount < battlecommon.GaugeMaxCount {
			dxlib.DrawGraph(96, 5, imgGaugeFrame, dxlib.TRUE)
			const gaugeMaxSize = 256
			size := int32(gaugeMaxSize * p.GaugeCount / battlecommon.GaugeMaxCount)
			dxlib.DrawBox(112, 19, 112+size, 21, dxlib.GetColor(123, 154, 222), dxlib.TRUE)
			dxlib.DrawBox(112, 21, 112+size, 29, dxlib.GetColor(231, 235, 255), dxlib.TRUE)
			dxlib.DrawBox(112, 29, 112+size, 31, dxlib.GetColor(123, 154, 222), dxlib.TRUE)
		} else {
			i := (p.GaugeCount / 20) % 4
			dxlib.DrawGraph(96, 5, imgGaugeMax[i], dxlib.TRUE)
		}
	}
}

func (p *BattlePlayer) Process() (bool, error) {
	if !p.EnableAct {
		return false, nil
	}

	if p.HP <= 0 {
		// Player deleted
		img := &imgPlayers[battlecommon.PlayerActDamage][1]
		battlecommon.NewDelete(*img, p.Pos, true)
		*img = -1 // DeleteGraph at delete animation
		p.NextAction = NextActLose
		p.EnableAct = false
		p.invincibleCount = 15 // do not show player image
		return false, nil
	}

	if p.invincibleCount > 0 {
		p.invincibleCount++
		if p.invincibleCount > battlecommon.PlayerDefaultInvincibleTime {
			p.invincibleCount = 0
		}
	}

	p.GaugeCount += 4 // TODO GaugeSpeed

	if p.act.Process() {
		return false, nil
	}

	if p.GaugeCount >= battlecommon.GaugeMaxCount {
		if p.GaugeCount == battlecommon.GaugeMaxCount {
			sound.On(sound.SEGaugeMax)
		}

		// State change to chip select
		if inputs.CheckKey(inputs.KeyLButton) == 1 || inputs.CheckKey(inputs.KeyRButton) == 1 {
			p.GaugeCount = 0
			p.NextAction = NextActChipSelect
			return false, nil
		}
	}

	// Move
	moveDirect := -1
	if inputs.CheckKey(inputs.KeyUp) == 1 {
		moveDirect = common.DirectUp
	} else if inputs.CheckKey(inputs.KeyDown) == 1 {
		moveDirect = common.DirectDown
	} else if inputs.CheckKey(inputs.KeyRight) == 1 {
		moveDirect = common.DirectRight
	} else if inputs.CheckKey(inputs.KeyLeft) == 1 {
		moveDirect = common.DirectLeft
	}

	if moveDirect >= 0 {
		if battlecommon.MoveObject(&p.Pos, moveDirect, field.PanelTypePlayer, false, field.GetPanelInfo) {
			p.act.MoveDirect = moveDirect
			p.act.SetAnim(battlecommon.PlayerActMove, 0)
			p.MoveNum++
			return false, nil
		}
	}

	// Chip use
	if inputs.CheckKey(inputs.KeyEnter) == 1 {
		if len(p.SelectedChips) > 0 {
			c := chip.Get(p.SelectedChips[0].ID)
			if c.PlayerAct != -1 {
				p.act.SetAnim(c.PlayerAct, c.KeepCount)
			}
			target := damage.TargetEnemy
			if c.ForMe {
				target = damage.TargetPlayer
			}

			sid := skill.GetSkillID(c.ID)
			p.act.ID = anim.New(skill.Get(sid, skill.Argument{
				OwnerID:    p.ID,
				Power:      c.Power,
				TargetType: target,
			}))

			p.SelectedChips = p.SelectedChips[1:]
			return false, nil
		}
	}

	// Rock buster
	if inputs.CheckKey(inputs.KeyCancel) > 0 {
		p.ChargeCount++
		if p.ChargeCount == battlecommon.ChargeViewDelay {
			sound.On(sound.SEBusterCharging)
		}
		if p.ChargeCount == battlecommon.ChargeTime {
			sound.On(sound.SEBusterCharged)
		}
	} else if p.ChargeCount > 0 {
		sound.On(sound.SEBusterShot)
		p.act.Charged = p.ChargeCount > battlecommon.ChargeTime
		p.act.ShotPower = p.ShotPower
		p.act.SetAnim(battlecommon.PlayerActBuster, 0)
		p.ChargeCount = 0
	}

	return false, nil
}

func (p *BattlePlayer) DamageProc(dm *damage.Damage) bool {
	if dm == nil {
		return false
	}

	// Recover系は使えるようにする
	if p.invincibleCount > 0 && dm.Power >= 0 {
		return false
	}

	if dm.TargetType&damage.TargetPlayer != 0 {
		hp := int(p.HP) - dm.Power
		if hp < 0 {
			p.HP = 0
		} else if hp > int(p.HPMax) {
			p.HP = p.HPMax
		} else {
			p.HP = uint(hp)
		}
		anim.New(effect.Get(dm.HitEffectType, p.Pos, 5))

		if dm.Power <= 0 {
			// Not damage, maybe recover or special anim
			return true
		}

		if !dm.BigDamage {
			return true
		}

		sound.On(sound.SEDamaged)

		// Stop current animation
		if objanim.IsProcessing(p.act.ID) {
			objanim.Delete(p.act.ID)
			p.act.ID = ""
		}

		p.act.SetAnim(battlecommon.PlayerActDamage, 0)
		p.invincibleCount = 1
		p.DamageNum++
		logger.Debug("Player damaged: %+v", *dm)
		return true
	}
	return false
}

func (p *BattlePlayer) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		Pos:      p.Pos,
		AnimType: anim.AnimTypeObject,
	}
}

func (p *BattlePlayer) GetObjectType() int {
	return objanim.ObjTypePlayer
}

func (p *BattlePlayer) SetChipSelectResult(selected []int) {
	p.SelectedChips = []player.ChipInfo{}
	for _, s := range selected {
		p.SelectedChips = append(p.SelectedChips, p.ChipFolder[s])
	}

	// Remove selected chips from folder
	sort.Sort(sort.Reverse(sort.IntSlice(selected)))
	for _, s := range selected {
		p.ChipFolder = append(p.ChipFolder[:s], p.ChipFolder[s+1:]...)
	}
}

// Process method returns true if processing now
func (a *act) Process() bool {
	switch a.typ {
	case -1: // No animation
		return false
	case battlecommon.PlayerActBuster:
		if a.count == 1 {
			s := a.ShotPower
			eff := effect.TypeHitSmall
			if a.Charged {
				s *= 10
				eff = effect.TypeHitBig
			}

			y := a.pPos.Y
			for x := a.pPos.X + 1; x < field.FieldNum.X; x++ {
				// logger.Debug("Rock buster damage set %d to (%d, %d)", s, x, *a.pPosY)
				pos := common.Point{X: x, Y: y}
				if field.GetPanelInfo(pos).ObjectID != "" {
					damage.New(damage.Damage{
						Pos:           pos,
						Power:         int(s),
						TTL:           1,
						TargetType:    damage.TargetEnemy,
						HitEffectType: eff,
					})
					break
				}
			}
		}
	case battlecommon.PlayerActMove:
		if a.count == 2 {
			battlecommon.MoveObject(a.pPos, a.MoveDirect, field.PanelTypePlayer, true, field.GetPanelInfo)
		}
	case battlecommon.PlayerActCannon, battlecommon.PlayerActSword, battlecommon.PlayerActBomb, battlecommon.PlayerActDamage, battlecommon.PlayerActShot, battlecommon.PlayerActPick, battlecommon.PlayerActThrow:
		// No special action
	default:
		panic(fmt.Sprintf("Invalid player anim type %d was specified.", a.typ))
	}

	a.count++

	num := len(imgPlayers[a.typ]) + a.keepCount
	if a.count > num*imgDelays[a.typ] {
		// Reset params
		a.typ = -1
		a.count = 0
		a.keepCount = 0
		return false // finished
	}
	return true // processing now
}

func (a *act) SetAnim(animType int, keepCount int) {
	a.count = 0
	a.typ = animType
	a.keepCount = keepCount
}

func (a *act) GetImage() int32 {
	if a.typ == -1 {
		// return stand image
		return imgPlayers[battlecommon.PlayerActMove][0]
	}

	imgNo := (a.count / imgDelays[a.typ])
	if imgNo >= len(imgPlayers[a.typ]) {
		imgNo = len(imgPlayers[a.typ]) - 1
	}

	return imgPlayers[a.typ][imgNo]
}
