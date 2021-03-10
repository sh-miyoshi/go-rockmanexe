package battle

import (
	"fmt"

	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/common"
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

type battlePlayer struct {
	posX int
	posY int
	hp   uint
	act  int
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
	playerInfo.act = playerAnimMove

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
	img := imgPlayers[playerInfo.act][0] // TODO
	dxlib.DrawRotaGraph(int32(x), int32(y), 1, 0, img, dxlib.TRUE)
}

func playerMainProcess() {
	// TODO: stateChange(chipSelect)
}
