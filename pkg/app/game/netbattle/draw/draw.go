package draw

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/net"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/router/object"
)

var (
	images    [object.TypeMax][]int
	imgDelays = [object.TypeMax]int{1, 2, 2, 6, 3, 4, 1, 4, 3} // debug
)

func Init() error {
	fname := common.ImagePath + "battle/character/player_move.png"
	images[object.TypePlayerMove] = make([]int, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 100, 100, images[object.TypePlayerMove]); res == -1 {
		return fmt.Errorf("failed to load player move image: %s", fname)
	}

	images[object.TypePlayerStand] = make([]int, 1)
	images[object.TypePlayerStand][0] = images[object.TypePlayerMove][0]

	fname = common.ImagePath + "battle/character/player_damaged.png"
	images[object.TypePlayerDamaged] = make([]int, 6)
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 100, 100, images[object.TypePlayerDamaged]); res == -1 {
		return fmt.Errorf("failed to load player damage image: %s", fname)
	}
	// 1 -> 2,3  2-4 3-5
	images[object.TypePlayerDamaged][4] = images[object.TypePlayerDamaged][2]
	images[object.TypePlayerDamaged][5] = images[object.TypePlayerDamaged][3]
	images[object.TypePlayerDamaged][2] = images[object.TypePlayerDamaged][1]
	images[object.TypePlayerDamaged][3] = images[object.TypePlayerDamaged][1]

	fname = common.ImagePath + "battle/character/player_shot.png"
	images[object.TypePlayerShot] = make([]int, 6)
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 180, 100, images[object.TypePlayerShot]); res == -1 {
		return fmt.Errorf("failed to load player shot image: %s", fname)
	}

	fname = common.ImagePath + "battle/character/player_cannon.png"
	images[object.TypePlayerCannon] = make([]int, 6)
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 100, 100, images[object.TypePlayerCannon]); res == -1 {
		return fmt.Errorf("failed to load player cannon image: %s", fname)
	}

	fname = common.ImagePath + "battle/character/player_sword.png"
	images[object.TypePlayerSword] = make([]int, 7)
	if res := dxlib.LoadDivGraph(fname, 7, 7, 1, 128, 128, images[object.TypePlayerSword]); res == -1 {
		return fmt.Errorf("failed to load player sword image: %s", fname)
	}

	fname = common.ImagePath + "battle/character/player_bomb.png"
	images[object.TypePlayerBomb] = make([]int, 7)
	if res := dxlib.LoadDivGraph(fname, 5, 5, 1, 100, 114, images[object.TypePlayerBomb]); res == -1 {
		return fmt.Errorf("failed to load player bomb image: %s", fname)
	}
	images[object.TypePlayerBomb][5] = images[object.TypePlayerBomb][4]
	images[object.TypePlayerBomb][6] = images[object.TypePlayerBomb][4]

	fname = common.ImagePath + "battle/character/player_buster.png"
	images[object.TypePlayerBuster] = make([]int, 6)
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 180, 100, images[object.TypePlayerBuster]); res == -1 {
		return fmt.Errorf("failed to load player buster image: %s", fname)
	}

	fname = common.ImagePath + "battle/character/player_pick.png"
	images[object.TypePlayerPick] = make([]int, 6)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 96, 124, images[object.TypePlayerPick]); res == -1 {
		return fmt.Errorf("failed to load player pick image: %s", fname)
	}
	images[object.TypePlayerPick][4] = images[object.TypePlayerPick][3]
	images[object.TypePlayerPick][5] = images[object.TypePlayerPick][3]

	fname = common.ImagePath + "battle/character/player_throw.png"
	images[object.TypePlayerThrow] = make([]int, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 97, 115, images[object.TypePlayerThrow]); res == -1 {
		return fmt.Errorf("failed to load player throw image: %s", fname)
	}

	return nil
}

func End() {
	for i := 0; i < object.TypeMax; i++ {
		for j := 0; j < len(images[i]); j++ {
			dxlib.DeleteGraph(images[i][j])
		}
		images[i] = []int{}
	}
}

func Draw() {
	ginfo := net.GetInst().GetGameInfo()
	for _, obj := range ginfo.Objects {
		pos := battlecommon.ViewPos(obj.Pos)
		ino := 0 // debug
		// TODO offset

		opts := dxlib.DrawRotaGraphOption{}
		if obj.IsReverse {
			rev := int32(dxlib.TRUE)
			opts.ReverseXFlag = &rev
		}

		dxlib.DrawRotaGraph(pos.X, pos.Y, 1, 0, images[obj.Type][ino], true, opts)
	}
}
