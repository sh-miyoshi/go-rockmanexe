package player

import (
	"fmt"
	"sort"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/anim"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/skill"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/chip"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/inputs"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/player"
)

const (
	playerAnimMove int = iota
	playerAnimDamage
	playerAnimShot
	playerAnimCannon
	playerAnimSword
	playerAnimBomb
	playerAnimBuster
	playerAnimMax
)

const (
	NextActNone int = iota
	NextActChipSelect
	MextActLose
)

type act struct {
	MoveDirect int
	Charged    bool
	ShotPower  uint

	typ   int
	count int
	pPosX *int
	pPosY *int
}

// BattlePlayer ...
type BattlePlayer struct {
	ID            string
	PosX          int
	PosY          int
	HP            uint
	HPMax         uint
	ChargeCount   uint
	GaugeCount    uint
	ShotPower     uint
	ChipFolder    []player.ChipInfo
	SelectedChips []player.ChipInfo
	NextAction    int
	EnableAct     bool

	act             act
	invincibleCount int
}

const (
	gaugeMaxCount  = 1200
	chargeTime     = 180 // TODO 変数化
	invincibleTime = 120
)

var (
	imgPlayers    [playerAnimMax][]int32
	imgDelays     = [playerAnimMax]int{1, 2, 2, 6, 3, 4, 1}
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
		PosX:      1,
		PosY:      1,
		ShotPower: plyr.ShotPower,
		EnableAct: true,
	}
	res.act.typ = -1
	res.act.pPosX = &res.PosX
	res.act.pPosY = &res.PosY

	for _, c := range plyr.ChipFolder {
		res.ChipFolder = append(res.ChipFolder, c)
	}
	// TODO: srand
	// Shuffle
	// for i := 0; i < 10; i++ {
	// 	for j := 0; j < len(res.ChipFolder); j++ {
	// 		n := rand.Intn(len(res.ChipFolder))
	// 		res.ChipFolder[j], res.ChipFolder[n] = res.ChipFolder[n], res.ChipFolder[j]
	// 	}
	// }

	fname := common.ImagePath + "battle/character/player_move.png"
	imgPlayers[playerAnimMove] = make([]int32, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 100, 100, imgPlayers[playerAnimMove]); res == -1 {
		return nil, fmt.Errorf("Failed to load player move image: %s", fname)
	}

	fname = common.ImagePath + "battle/character/player_damaged.png"
	imgPlayers[playerAnimDamage] = make([]int32, 6)
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 100, 100, imgPlayers[playerAnimDamage]); res == -1 {
		return nil, fmt.Errorf("Failed to load player damage image: %s", fname)
	}
	// 1 -> 2,3  2-4 3-5
	imgPlayers[playerAnimDamage][4] = imgPlayers[playerAnimDamage][2]
	imgPlayers[playerAnimDamage][5] = imgPlayers[playerAnimDamage][3]
	imgPlayers[playerAnimDamage][2] = imgPlayers[playerAnimDamage][1]
	imgPlayers[playerAnimDamage][3] = imgPlayers[playerAnimDamage][1]

	fname = common.ImagePath + "battle/character/player_shot.png"
	imgPlayers[playerAnimShot] = make([]int32, 6)
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 180, 100, imgPlayers[playerAnimShot]); res == -1 {
		return nil, fmt.Errorf("Failed to load player shot image: %s", fname)
	}

	fname = common.ImagePath + "battle/character/player_cannon.png"
	imgPlayers[playerAnimCannon] = make([]int32, 6)
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 100, 100, imgPlayers[playerAnimCannon]); res == -1 {
		return nil, fmt.Errorf("Failed to load player cannon image: %s", fname)
	}

	fname = common.ImagePath + "battle/character/player_sword.png"
	imgPlayers[playerAnimSword] = make([]int32, 7)
	if res := dxlib.LoadDivGraph(fname, 7, 7, 1, 128, 128, imgPlayers[playerAnimSword]); res == -1 {
		return nil, fmt.Errorf("Failed to load player sword image: %s", fname)
	}

	fname = common.ImagePath + "battle/character/player_bomb.png"
	imgPlayers[playerAnimBomb] = make([]int32, 5)
	if res := dxlib.LoadDivGraph(fname, 5, 5, 1, 100, 114, imgPlayers[playerAnimBomb]); res == -1 {
		return nil, fmt.Errorf("Failed to load player bomb image: %s", fname)
	}

	fname = common.ImagePath + "battle/character/player_buster.png"
	imgPlayers[playerAnimBuster] = make([]int32, 6)
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 180, 100, imgPlayers[playerAnimBuster]); res == -1 {
		return nil, fmt.Errorf("Failed to load player buster image: %s", fname)
	}

	fname = common.ImagePath + "battle/hp_frame.png"
	imgHPFrame = dxlib.LoadGraph(fname)
	if imgHPFrame < 0 {
		return nil, fmt.Errorf("Failed to read hp frame image %s", fname)
	}
	fname = common.ImagePath + "battle/gauge.png"
	imgGaugeFrame = dxlib.LoadGraph(fname)
	if imgGaugeFrame < 0 {
		return nil, fmt.Errorf("Failed to read gauge frame image %s", fname)
	}
	fname = common.ImagePath + "battle/gauge_max.png"
	imgGaugeMax = make([]int32, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 1, 4, 288, 30, imgGaugeMax); res == -1 {
		return nil, fmt.Errorf("Failed to read gauge max image %s", fname)
	}

	fname = common.ImagePath + "battle/skill/charge.png"
	tmp := make([]int32, 16)
	if res := dxlib.LoadDivGraph(fname, 16, 8, 2, 158, 150, tmp); res == -1 {
		return nil, fmt.Errorf("Failed to load image %s", fname)
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

	for i := 0; i < playerAnimMax; i++ {
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

	logger.Info("Successfully cleanuped battle player data")
}

// Draw ...
func (p *BattlePlayer) Draw() {
	if p.invincibleCount/5%2 != 0 {
		return
	}

	x, y := battlecommon.ViewPos(p.PosX, p.PosY)
	img := p.act.GetImage()
	dxlib.DrawRotaGraph(x, y, 1, 0, img, dxlib.TRUE)

	// Show charge image
	if p.ChargeCount > 20 {
		n := 0
		if p.ChargeCount > chargeTime {
			n = 1
		}
		imgNo := int(p.ChargeCount/4) % len(imgCharge[n])
		dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_ALPHA, 224)
		dxlib.DrawRotaGraph(x, y, 1, 0, imgCharge[n][imgNo], dxlib.TRUE)
		dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_NOBLEND, 0)
	}

	// Show selected chip icons
	n := len(p.SelectedChips)
	if n > 0 {
		// TODO Show chip info

		const px = 3
		max := n * px
		for i := 0; i < n; i++ {
			x := field.PanelSizeX*p.PosX + field.PanelSizeX/2 - 2 + (i * px) - max
			y := field.DrawPanelTopY + field.PanelSizeY*p.PosY - 10 - 81 + (i * px) - max
			dxlib.DrawBox(int32(x-1), int32(y-1), int32(x+29), int32(y+29), 0x000000, dxlib.FALSE)
			// draw from the end
			dxlib.DrawGraph(int32(x), int32(y), chip.GetIcon(p.SelectedChips[n-1-i].ID, true), dxlib.TRUE)
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
		if p.GaugeCount < gaugeMaxCount {
			dxlib.DrawGraph(96, 5, imgGaugeFrame, dxlib.TRUE)
			const gaugeMaxSize = 256
			size := int32(gaugeMaxSize * p.GaugeCount / gaugeMaxCount)
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

	if p.invincibleCount > 0 {
		p.invincibleCount++
		if p.invincibleCount > invincibleTime {
			p.invincibleCount = 0
		}
	}

	p.GaugeCount += 4 // TODO GaugeSpeed

	// TODO delete process
	// TODO return true, nil and set NextActLose if dead

	if p.act.Process() {
		return false, nil
	}

	if p.GaugeCount >= gaugeMaxCount {
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
		if battlecommon.MoveObject(&p.PosX, &p.PosY, moveDirect, field.PanelTypePlayer, false) {
			p.act.MoveDirect = moveDirect
			p.act.SetAnim(playerAnimMove)
			return false, nil
		}
	}

	// Chip use
	if inputs.CheckKey(inputs.KeyEnter) == 1 {
		if len(p.SelectedChips) > 0 {
			c := chip.Get(p.SelectedChips[0].ID)
			if c.PlayerAct != -1 {
				p.act.SetAnim(c.PlayerAct)
			}
			target := damage.TargetEnemy
			if c.ForMe {
				target = damage.TargetPlayer
			}

			anim.New(skill.GetByChip(c.ID, skill.Argument{
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
	} else if p.ChargeCount > 0 {
		p.act.Charged = p.ChargeCount > chargeTime
		p.act.ShotPower = p.ShotPower
		p.act.SetAnim(playerAnimBuster)
		p.ChargeCount = 0
	}

	return false, nil
}

func (p *BattlePlayer) DamageProc(dm *damage.Damage) {
	if dm == nil {
		return
	}

	if p.invincibleCount > 0 {
		return
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
		anim.New(effect.Get(dm.HitEffectType, p.PosX, p.PosY))
		if dm.Power > 0 {
			p.act.SetAnim(playerAnimDamage)
			p.invincibleCount = 1
		}
		logger.Debug("Player damaged: %+v", *dm)
	}
}

func (p *BattlePlayer) GetParam() anim.Param {
	return anim.Param{
		PosX:     p.PosX,
		PosY:     p.PosY,
		AnimType: anim.TypeObject,
	}
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
	case playerAnimBuster:
		if a.count == 1 {
			s := a.ShotPower
			eff := effect.TypeHitSmall
			if a.Charged {
				s *= 10
				eff = effect.TypeHitBig
			}

			y := *a.pPosY
			for x := *a.pPosX + 1; x < field.FieldNumX; x++ {
				// logger.Debug("Rock buster damage set %d to (%d, %d)", s, x, *a.pPosY)
				if field.GetPanelInfo(x, y).ObjectID != "" {
					damage.New(damage.Damage{
						PosX:          x,
						PosY:          y,
						Power:         int(s),
						TTL:           1,
						TargetType:    damage.TargetEnemy,
						HitEffectType: eff,
					})
					break
				}
			}
		}
	case playerAnimMove:
		if a.count == 2 {
			battlecommon.MoveObject(a.pPosX, a.pPosY, a.MoveDirect, field.PanelTypePlayer, true)
		}
	case playerAnimCannon, playerAnimSword, playerAnimBomb, playerAnimDamage, playerAnimShot:
		// No special action
	default:
		panic(fmt.Sprintf("Invalid player anim type %d was specified.", a.typ))
	}

	a.count++

	if a.count > len(imgPlayers[a.typ])*imgDelays[a.typ] {
		// Reset params
		a.typ = -1
		a.count = 0
		return false // finished
	}
	return true // processing now
}

func (a *act) SetAnim(animType int) {
	a.count = 0
	a.typ = animType
}

func (a *act) GetImage() int32 {
	if a.typ == -1 {
		// return stand image
		return imgPlayers[playerAnimMove][0]
	}

	imgNo := (a.count / imgDelays[a.typ]) % len(imgPlayers[a.typ])
	return imgPlayers[a.typ][imgNo]
}
