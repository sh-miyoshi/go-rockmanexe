package battle

import (
	"fmt"

	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/anim"
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
	posX int
	posY int
	hp   uint
	act  act
}

var (
	imgPlayers [playerAnimMax][]int32
	playerInfo battlePlayer
)

func playerInit(hp uint) error {
	logger.Info("Initialize battle player data")

	playerInfo.hp = hp
	playerInfo.posX = 1
	playerInfo.posY = 1
	playerInfo.act.typ = playerAnimMove

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

func playerEnd() {
	logger.Info("Cleanup battle player data")

	for i := 0; i < playerAnimMax; i++ {
		for j := 0; j < len(imgPlayers[i]); j++ {
			dxlib.DeleteGraph(imgPlayers[i][j])
			imgPlayers[i][j] = -1
		}
	}

	logger.Info("Successfully cleanuped battle player data")
}

func playerDraw() {
	x := panelSizeX*playerInfo.posX + panelSizeX/2
	y := drawPanelTopY + panelSizeY*playerInfo.posY - 10
	img := imgPlayers[playerInfo.act.typ][playerInfo.act.getImageNo()]
	dxlib.DrawRotaGraph(int32(x), int32(y), 1, 0, img, dxlib.TRUE)
}

func playerMainProcess() {
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
	// TODO: shot

	// TODO: move
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
		if moveObject(&playerInfo.posX, &playerInfo.posY, moveDirect, false) {
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

func (a *act) reset() {
	a.typ = playerAnimMove
	a.count = 0
	a.animID = ""
}

func (a *act) Process() (bool, error) {
	switch a.typ {
	case playerAnimMove:
		if a.count == 2 {
			moveObject(&playerInfo.posX, &playerInfo.posY, a.moveDirect, true)
		}
		if a.count > len(imgPlayers[playerAnimMove]) {
			return true, nil
		}
	default:
		return false, fmt.Errorf("Anim %d is not implemented yet", a.typ)
	}
	a.count++
	return false, nil
}

func (a *act) getImageNo() int {
	// TODO image delay
	return a.count % len(imgPlayers[a.typ])
}
