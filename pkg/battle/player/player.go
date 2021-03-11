package player

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/anim"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/inputs"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
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

type battlePlayer struct {
	id          string
	posX        int
	posY        int
	hp          uint
	act         act
	chargeCount uint
}

var (
	imgPlayers [playerAnimMax][]int32
	imgDelays  = [playerAnimMax]int{1, 1, 1, 1, 1, 1} // TODO: set correct value
	playerInfo battlePlayer
)

// Init ...
func Init(hp uint) error {
	logger.Info("Initialize battle player data")

	if playerInfo.id == "" {
		playerInfo.id = uuid.New().String()
	}

	playerInfo.hp = hp
	playerInfo.posX = 1
	playerInfo.posY = 1
	playerInfo.act.typ = playerAnimMove
	playerInfo.chargeCount = 0

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

// Draw ...
func Draw() {
	x := field.PanelSizeX*playerInfo.posX + field.PanelSizeX/2
	y := field.DrawPanelTopY + field.PanelSizeY*playerInfo.posY - 10
	img := imgPlayers[playerInfo.act.typ][playerInfo.act.getImageNo()]
	dxlib.DrawRotaGraph(int32(x), int32(y), 1, 0, img, dxlib.TRUE)
}

// GetID ...
func GetID() string {
	return playerInfo.id
}

// GetPos ...
func GetPos() (x, y int) {
	return playerInfo.posX, playerInfo.posY
}

// MainProcess ...
func MainProcess() {
	if playerInfo.act.animID != "" {
		// still in animation
		if !anim.IsProcessing(playerInfo.act.animID) {
			// end animation
			playerInfo.act.reset()
		}
		return
	}

	// TODO: stateChange(chipSelect)
	// TODO: chip use

	// Rock buster
	if inputs.CheckKey(inputs.KeyCancel) > 0 {
		playerInfo.chargeCount++
	} else if playerInfo.chargeCount > 0 {
		playerInfo.act.setShot(playerInfo.chargeCount)
		playerInfo.chargeCount = 0
		return
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
		if battlecommon.MoveObject(&playerInfo.posX, &playerInfo.posY, moveDirect, false) {
			playerInfo.act.setMove(moveDirect)
		}
	}
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
			battlecommon.MoveObject(&playerInfo.posX, &playerInfo.posY, a.moveDirect, true)
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
