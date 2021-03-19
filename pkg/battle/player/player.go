package player

import (
	"errors"
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
	playerAnimMax
)

type act struct {
	typ        int
	count      int
	animID     string
	moveDirect int
}

// BattlePlayer ...
type BattlePlayer struct {
	ID            string
	PosX          int
	PosY          int
	HP            uint
	ChargeCount   uint
	GaugeCount    uint
	ChipFolder    []player.ChipInfo
	SelectedChips []player.ChipInfo

	act act
}

const (
	gaugeMaxCount = 256 // debug(1200)
)

var (
	ErrPlayerDead = errors.New("player dead")
	ErrChipSelect = errors.New("chip select")

	imgPlayers    [playerAnimMax][]int32
	imgDelays     = [playerAnimMax]int{1, 1, 1, 6, 3, 1} // TODO: set correct value
	imgHPFrame    int32
	imgGaugeFrame int32
	imgGaugeMax   []int32
	playerInfo    BattlePlayer
)

// Init ...
func Init(hp uint, chipFolder [player.FolderSize]player.ChipInfo) error {
	logger.Info("Initialize battle player data")

	if playerInfo.ID == "" {
		playerInfo.ID = uuid.New().String()
	}

	playerInfo.HP = hp
	playerInfo.PosX = 1
	playerInfo.PosY = 1
	playerInfo.act.typ = playerAnimMove
	playerInfo.ChargeCount = 0
	playerInfo.GaugeCount = 0

	for _, c := range chipFolder {
		playerInfo.ChipFolder = append(playerInfo.ChipFolder, c)
	}
	// TODO: Shuffle

	fname := common.ImagePath + "battle/character/player_move.png"
	imgPlayers[playerAnimMove] = make([]int32, 4)
	res := dxlib.LoadDivGraph(fname, 4, 4, 1, 100, 100, imgPlayers[playerAnimMove])
	if res == -1 {
		return fmt.Errorf("Failed to load player move image: %s", fname)
	}

	fname = common.ImagePath + "battle/character/player_damaged.png"
	imgPlayers[playerAnimDamage] = make([]int32, 6)
	res = dxlib.LoadDivGraph(fname, 6, 6, 1, 100, 100, imgPlayers[playerAnimDamage])
	if res == -1 {
		return fmt.Errorf("Failed to load player damage image: %s", fname)
	}
	// 1 -> 2,3  2-4 3-5
	imgPlayers[playerAnimDamage][4] = imgPlayers[playerAnimDamage][2]
	imgPlayers[playerAnimDamage][5] = imgPlayers[playerAnimDamage][3]
	imgPlayers[playerAnimDamage][2] = imgPlayers[playerAnimDamage][1]
	imgPlayers[playerAnimDamage][3] = imgPlayers[playerAnimDamage][1]

	fname = common.ImagePath + "battle/character/player_shot.png"
	imgPlayers[playerAnimShot] = make([]int32, 6)
	res = dxlib.LoadDivGraph(fname, 6, 6, 1, 180, 100, imgPlayers[playerAnimShot])
	if res == -1 {
		return fmt.Errorf("Failed to load player shot image: %s", fname)
	}

	fname = common.ImagePath + "battle/character/player_cannon.png"
	imgPlayers[playerAnimCannon] = make([]int32, 6)
	res = dxlib.LoadDivGraph(fname, 6, 6, 1, 100, 100, imgPlayers[playerAnimCannon])
	if res == -1 {
		return fmt.Errorf("Failed to load player cannon image: %s", fname)
	}

	fname = common.ImagePath + "battle/character/player_sword.png"
	imgPlayers[playerAnimSword] = make([]int32, 7)
	res = dxlib.LoadDivGraph(fname, 7, 7, 1, 128, 128, imgPlayers[playerAnimSword])
	if res == -1 {
		return fmt.Errorf("Failed to load player sword image: %s", fname)
	}

	fname = common.ImagePath + "battle/character/player_bomb.png"
	imgPlayers[playerAnimBomb] = make([]int32, 5)
	res = dxlib.LoadDivGraph(fname, 5, 5, 1, 100, 114, imgPlayers[playerAnimBomb])
	if res == -1 {
		return fmt.Errorf("Failed to load player bomb image: %s", fname)
	}

	fname = common.ImagePath + "battle/hp_frame.png"
	imgHPFrame = dxlib.LoadGraph(fname)
	if imgHPFrame < 0 {
		return fmt.Errorf("Failed to read hp frame image %s", fname)
	}
	fname = common.ImagePath + "battle/gauge.png"
	imgGaugeFrame = dxlib.LoadGraph(fname)
	if imgGaugeFrame < 0 {
		return fmt.Errorf("Failed to read gauge frame image %s", fname)
	}
	fname = common.ImagePath + "battle/gauge_max.png"
	imgGaugeMax = make([]int32, 4)
	res = dxlib.LoadDivGraph(fname, 4, 1, 4, 288, 30, imgGaugeMax)
	if res == -1 {
		return fmt.Errorf("Failed to read gauge max image %s", fname)
	}

	logger.Info("Successfully initialized battle player data")
	return nil
}

// End ...
func End() {
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

// DrawChar ...
func DrawChar() {
	x, y := battlecommon.ViewPos(playerInfo.PosX, playerInfo.PosY)
	img := imgPlayers[playerInfo.act.typ][playerInfo.act.getImageNo()]
	dxlib.DrawRotaGraph(x, y, 1, 0, img, dxlib.TRUE)
}

// DrawChipIcon ...
func DrawChipIcon() {
	n := len(playerInfo.SelectedChips)
	if n > 0 {
		// TODO Show chip info

		const px = 3
		max := n * px
		for i := 0; i < n; i++ {
			x := field.PanelSizeX*playerInfo.PosX + field.PanelSizeX/2 - 2 + (i * px) - max
			y := field.DrawPanelTopY + field.PanelSizeY*playerInfo.PosY - 10 - 81 + (i * px) - max
			dxlib.DrawBox(int32(x-1), int32(y-1), int32(x+29), int32(y+29), 0x000000, dxlib.FALSE)
			// draw from the end
			dxlib.DrawGraph(int32(x), int32(y), chip.GetIcon(playerInfo.SelectedChips[n-1-i].ID, true), dxlib.TRUE)
		}
	}
}

func DrawFrame(xShift bool, showGauge bool) {
	x := int32(7)
	y := int32(5)
	if xShift {
		x += 235
	}

	// Show HP
	dxlib.DrawGraph(x, y, imgHPFrame, dxlib.TRUE)
	draw.Number(x+2, y+2, int32(playerInfo.HP), draw.NumberOption{RightAligned: true, Length: 4})

	// Show Custom Gauge
	if showGauge {
		if playerInfo.GaugeCount < gaugeMaxCount {
			dxlib.DrawGraph(96, 5, imgGaugeFrame, dxlib.TRUE)
			const gaugeMaxSize = 256
			size := int32(gaugeMaxSize * playerInfo.GaugeCount / gaugeMaxCount)
			dxlib.DrawBox(112, 19, 112+size, 21, dxlib.GetColor(123, 154, 222), dxlib.TRUE)
			dxlib.DrawBox(112, 21, 112+size, 29, dxlib.GetColor(231, 235, 255), dxlib.TRUE)
			dxlib.DrawBox(112, 29, 112+size, 31, dxlib.GetColor(123, 154, 222), dxlib.TRUE)
		} else {
			i := (playerInfo.GaugeCount / 20) % 4
			dxlib.DrawGraph(96, 5, imgGaugeMax[i], dxlib.TRUE)
		}
	}
}

// Get ...
func Get() *BattlePlayer {
	return &playerInfo
}

// SetChipSelectResult ...
func SetChipSelectResult(selected []int) {
	playerInfo.SelectedChips = []player.ChipInfo{}
	for _, s := range selected {
		playerInfo.SelectedChips = append(playerInfo.SelectedChips, playerInfo.ChipFolder[s])
	}

	// Remove selected chips from folder
	sort.Sort(sort.Reverse(sort.IntSlice(selected)))
	for _, s := range selected {
		playerInfo.ChipFolder = append(playerInfo.ChipFolder[:s], playerInfo.ChipFolder[s+1:]...)
	}
}

// MainProcess ...
func MainProcess() error {
	playerInfo.GaugeCount++ // TODO GaugeSpeed

	if playerInfo.act.animID != "" {
		// still in animation
		if !anim.IsProcessing(playerInfo.act.animID) {
			// end animation
			playerInfo.act.reset()
		}
		return nil
	}

	if playerInfo.GaugeCount >= gaugeMaxCount {
		// State change to chip select
		if inputs.CheckKey(inputs.KeyLButton) == 1 || inputs.CheckKey(inputs.KeyRButton) == 1 {
			playerInfo.GaugeCount = 0
			return ErrChipSelect
		}
	}

	// Chip use
	if inputs.CheckKey(inputs.KeyEnter) == 1 {
		if len(playerInfo.SelectedChips) > 0 {
			c := chip.Get(playerInfo.SelectedChips[0].ID)
			if c.PlayerAct != -1 {
				playerInfo.act.set(c.PlayerAct)
			}
			anim.New(skill.Get(c.ID, skill.Argument{
				OwnerID:    playerInfo.ID,
				Power:      int(c.Power),
				TargetType: damage.TargetEnemy,
			}))

			playerInfo.SelectedChips = playerInfo.SelectedChips[1:]
			return nil
		}
	}

	// Rock buster
	if inputs.CheckKey(inputs.KeyCancel) > 0 {
		playerInfo.ChargeCount++
	} else if playerInfo.ChargeCount > 0 {
		// TODO set act.ShotPower by playerInfo.ChargeCount
		playerInfo.act.set(playerAnimShot)
		playerInfo.ChargeCount = 0
		return nil
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
		if battlecommon.MoveObject(&playerInfo.PosX, &playerInfo.PosY, moveDirect, field.PanelTypePlayer, false) {
			playerInfo.act.moveDirect = moveDirect
			playerInfo.act.set(playerAnimMove)
		}
	}

	return nil
}

func (a *act) set(typ int) {
	a.typ = typ
	a.count = 0
	a.animID = anim.New(a)
}

func (a *act) reset() {
	a.typ = playerAnimMove
	a.count = 0
	a.animID = ""
}

func (a *act) Draw() {
	// No common drawing process
}

func (a *act) Process() (bool, error) {
	switch a.typ {
	case playerAnimMove:
		if a.count == 2 {
			battlecommon.MoveObject(&playerInfo.PosX, &playerInfo.PosY, a.moveDirect, field.PanelTypePlayer, true)
		}
	case playerAnimShot:
		if a.count == 1 {
			for x := playerInfo.PosX + 1; x < field.FieldNumX; x++ {
				damage.New(damage.Damage{
					PosX:          x,
					PosY:          playerInfo.PosY,
					Power:         1, // debug
					TTL:           1,
					TargetType:    damage.TargetEnemy,
					HitEffectType: effect.TypeHitSmall, // TODO HitBig if charge shot
				})
			}
		}
	case playerAnimCannon, playerAnimSword:
		// nothing to do
	default:
		return false, fmt.Errorf("Anim %d is not implemented yet", a.typ)
	}

	a.count++

	if a.count > len(imgPlayers[a.typ])*imgDelays[a.typ] {
		return true, nil
	}
	return false, nil
}

func (a *act) getImageNo() int {
	n := a.count / imgDelays[a.typ]
	if n >= len(imgPlayers[a.typ]) {
		n = len(imgPlayers[a.typ]) - 1
	}
	return n
}
