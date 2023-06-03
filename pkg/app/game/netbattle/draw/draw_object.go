package draw

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/net"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/router/object"
)

type objectDraw struct {
	images         [object.TypeMax][]int
	imgDelays      [object.TypeMax]int
	playerObjectID string
}

func (d *objectDraw) Init(playerObjectID string) error {
	d.imgDelays = [object.TypeMax]int{1, 2, 2, 6, 3, 4, 1, 4, 3} // debug

	d.playerObjectID = playerObjectID

	fname := common.ImagePath + "battle/character/player_move.png"
	d.images[object.TypePlayerMove] = make([]int, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 100, 100, d.images[object.TypePlayerMove]); res == -1 {
		return fmt.Errorf("failed to load player move image: %s", fname)
	}

	d.images[object.TypePlayerStand] = make([]int, 1)
	d.images[object.TypePlayerStand][0] = d.images[object.TypePlayerMove][0]

	fname = common.ImagePath + "battle/character/player_damaged.png"
	d.images[object.TypePlayerDamaged] = make([]int, 6)
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 100, 100, d.images[object.TypePlayerDamaged]); res == -1 {
		return fmt.Errorf("failed to load player damage image: %s", fname)
	}
	// 1 -> 2,3  2-4 3-5
	d.images[object.TypePlayerDamaged][4] = d.images[object.TypePlayerDamaged][2]
	d.images[object.TypePlayerDamaged][5] = d.images[object.TypePlayerDamaged][3]
	d.images[object.TypePlayerDamaged][2] = d.images[object.TypePlayerDamaged][1]
	d.images[object.TypePlayerDamaged][3] = d.images[object.TypePlayerDamaged][1]

	fname = common.ImagePath + "battle/character/player_shot.png"
	d.images[object.TypePlayerShot] = make([]int, 6)
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 180, 100, d.images[object.TypePlayerShot]); res == -1 {
		return fmt.Errorf("failed to load player shot image: %s", fname)
	}

	fname = common.ImagePath + "battle/character/player_cannon.png"
	d.images[object.TypePlayerCannon] = make([]int, 6)
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 100, 100, d.images[object.TypePlayerCannon]); res == -1 {
		return fmt.Errorf("failed to load player cannon image: %s", fname)
	}

	fname = common.ImagePath + "battle/character/player_sword.png"
	d.images[object.TypePlayerSword] = make([]int, 7)
	if res := dxlib.LoadDivGraph(fname, 7, 7, 1, 128, 128, d.images[object.TypePlayerSword]); res == -1 {
		return fmt.Errorf("failed to load player sword image: %s", fname)
	}

	fname = common.ImagePath + "battle/character/player_bomb.png"
	d.images[object.TypePlayerBomb] = make([]int, 7)
	if res := dxlib.LoadDivGraph(fname, 5, 5, 1, 100, 114, d.images[object.TypePlayerBomb]); res == -1 {
		return fmt.Errorf("failed to load player bomb image: %s", fname)
	}
	d.images[object.TypePlayerBomb][5] = d.images[object.TypePlayerBomb][4]
	d.images[object.TypePlayerBomb][6] = d.images[object.TypePlayerBomb][4]

	fname = common.ImagePath + "battle/character/player_buster.png"
	d.images[object.TypePlayerBuster] = make([]int, 6)
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 180, 100, d.images[object.TypePlayerBuster]); res == -1 {
		return fmt.Errorf("failed to load player buster image: %s", fname)
	}

	fname = common.ImagePath + "battle/character/player_pick.png"
	d.images[object.TypePlayerPick] = make([]int, 6)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 96, 124, d.images[object.TypePlayerPick]); res == -1 {
		return fmt.Errorf("failed to load player pick image: %s", fname)
	}
	d.images[object.TypePlayerPick][4] = d.images[object.TypePlayerPick][3]
	d.images[object.TypePlayerPick][5] = d.images[object.TypePlayerPick][3]

	fname = common.ImagePath + "battle/character/player_throw.png"
	d.images[object.TypePlayerThrow] = make([]int, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 97, 115, d.images[object.TypePlayerThrow]); res == -1 {
		return fmt.Errorf("failed to load player throw image: %s", fname)
	}

	return nil
}

func (d *objectDraw) End() {
	for i := 0; i < object.TypeMax; i++ {
		for j := 0; j < len(d.images[i]); j++ {
			dxlib.DeleteGraph(d.images[i][j])
		}
		d.images[i] = []int{}
	}
}

func (d *objectDraw) Draw() {
	ginfo := net.GetInst().GetGameInfo()
	for _, obj := range ginfo.Objects {
		pos := battlecommon.ViewPos(obj.Pos)
		ino := obj.ActCount / d.imgDelays[obj.Type]
		if ino >= len(d.images[obj.Type]) {
			ino = len(d.images[obj.Type]) - 1
		}
		// TODO offset

		opts := dxlib.DrawRotaGraphOption{}
		if obj.IsReverse {
			rev := int32(dxlib.TRUE)
			opts.ReverseXFlag = &rev
		}

		dxlib.DrawRotaGraph(pos.X, pos.Y, 1, 0, d.images[obj.Type][ino], true, opts)

		// draw hp
		if obj.ID != d.playerObjectID {
			if obj.HP > 0 {
				draw.Number(pos.X, pos.Y+40, obj.HP, draw.NumberOption{
					Color:    draw.NumberColorWhiteSmall,
					Centered: true,
				})
			}
		}
	}
}
