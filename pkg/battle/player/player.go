package player

import (
	"fmt"
	"sort"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/anim"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/chip"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/common"
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
	ChipFolder    []player.ChipInfo
	SelectedChips []player.ChipInfo

	act act
}

var (
	imgPlayers [playerAnimMax][]int32
	imgDelays  = [playerAnimMax]int{1, 1, 1, 1, 1, 1} // TODO: set correct value
	playerInfo BattlePlayer
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
	res = dxlib.LoadDivGraph(fname, 4, 4, 1, 100, 100, imgPlayers[playerAnimDamage])
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
	imgPlayers[playerAnimCannon] = make([]int32, 4)
	res = dxlib.LoadDivGraph(fname, 4, 4, 1, 100, 100, imgPlayers[playerAnimCannon])
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

	logger.Info("Successfully cleanuped battle player data")
}

// DrawChar ...
func DrawChar() {
	x := field.PanelSizeX*playerInfo.PosX + field.PanelSizeX/2
	y := field.DrawPanelTopY + field.PanelSizeY*playerInfo.PosY - 10
	img := imgPlayers[playerInfo.act.typ][playerInfo.act.getImageNo()]
	dxlib.DrawRotaGraph(int32(x), int32(y), 1, 0, img, dxlib.TRUE)
}

func DrawChipIcon() {
	n := len(playerInfo.SelectedChips)
	if n > 0 {
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
func MainProcess() bool {
	if playerInfo.act.animID != "" {
		// still in animation
		if !anim.IsProcessing(playerInfo.act.animID) {
			// end animation
			playerInfo.act.reset()
		}
		return false
	}

	// TODO: stateChange(chipSelect)
	// TODO: chip use

	// Rock buster
	if inputs.CheckKey(inputs.KeyCancel) > 0 {
		playerInfo.ChargeCount++
	} else if playerInfo.ChargeCount > 0 {
		playerInfo.act.setShot(playerInfo.ChargeCount)
		playerInfo.ChargeCount = 0
		return false
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
		if battlecommon.MoveObject(&playerInfo.PosX, &playerInfo.PosY, moveDirect, false) {
			playerInfo.act.setMove(moveDirect)
		}
	}

	if dxlib.CheckHitKey(dxlib.KEY_INPUT_A) == 1 {
		return true
	}
	return false
}

func (a *act) setMove(direct int) {
	a.typ = playerAnimMove
	a.count = 0
	a.animID = anim.New(a)
	a.moveDirect = direct
}

func (a *act) setShot(chargeCount uint) {
	a.typ = playerAnimShot
	a.count = 0
	a.animID = anim.New(a)

	// TODO: change show power by charge count
}

func (a *act) reset() {
	a.typ = playerAnimMove
	a.count = 0
	a.animID = ""
}

func (a *act) Process() (bool, error) {
	switch a.typ {
	case playerAnimMove:
		if a.count == 2 {
			battlecommon.MoveObject(&playerInfo.PosX, &playerInfo.PosY, a.moveDirect, true)
		}
		if a.count > len(imgPlayers[playerAnimMove]) {
			return true, nil
		}
	case playerAnimShot:
		// TODO
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
	return a.count % (len(imgPlayers[a.typ]) * imgDelays[a.typ])
}
