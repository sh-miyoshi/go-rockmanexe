package draw

import (
	"fmt"

	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/field"
)

var (
	images [field.ObjectTypeMax][]int32
)

func Init() error {
	fname := common.ImagePath + "battle/character/player_move.png"
	images[field.ObjectTypeRockmanMove] = make([]int32, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 100, 100, images[field.ObjectTypeRockmanMove]); res == -1 {
		return fmt.Errorf("failed to load player move image: %s", fname)
	}

	fname = common.ImagePath + "battle/character/player_damaged.png"
	images[field.ObjectTypeRockmanDamage] = make([]int32, 6)
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 100, 100, images[field.ObjectTypeRockmanDamage]); res == -1 {
		return fmt.Errorf("failed to load player damage image: %s", fname)
	}
	// 1 -> 2,3  2-4 3-5
	images[field.ObjectTypeRockmanDamage][4] = images[field.ObjectTypeRockmanDamage][2]
	images[field.ObjectTypeRockmanDamage][5] = images[field.ObjectTypeRockmanDamage][3]
	images[field.ObjectTypeRockmanDamage][2] = images[field.ObjectTypeRockmanDamage][1]
	images[field.ObjectTypeRockmanDamage][3] = images[field.ObjectTypeRockmanDamage][1]

	fname = common.ImagePath + "battle/character/player_shot.png"
	images[field.ObjectTypeRockmanShot] = make([]int32, 6)
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 180, 100, images[field.ObjectTypeRockmanShot]); res == -1 {
		return fmt.Errorf("failed to load player shot image: %s", fname)
	}

	fname = common.ImagePath + "battle/character/player_cannon.png"
	images[field.ObjectTypeRockmanCannon] = make([]int32, 6)
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 100, 100, images[field.ObjectTypeRockmanCannon]); res == -1 {
		return fmt.Errorf("failed to load player cannon image: %s", fname)
	}

	fname = common.ImagePath + "battle/character/player_sword.png"
	images[field.ObjectTypeRockmanSword] = make([]int32, 7)
	if res := dxlib.LoadDivGraph(fname, 7, 7, 1, 128, 128, images[field.ObjectTypeRockmanSword]); res == -1 {
		return fmt.Errorf("failed to load player sword image: %s", fname)
	}

	fname = common.ImagePath + "battle/character/player_bomb.png"
	images[field.ObjectTypeRockmanBomb] = make([]int32, 7)
	if res := dxlib.LoadDivGraph(fname, 5, 5, 1, 100, 114, images[field.ObjectTypeRockmanBomb]); res == -1 {
		return fmt.Errorf("failed to load player bomb image: %s", fname)
	}
	images[field.ObjectTypeRockmanBomb][5] = images[field.ObjectTypeRockmanBomb][4]
	images[field.ObjectTypeRockmanBomb][6] = images[field.ObjectTypeRockmanBomb][4]

	fname = common.ImagePath + "battle/character/player_buster.png"
	images[field.ObjectTypeRockmanBuster] = make([]int32, 6)
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 180, 100, images[field.ObjectTypeRockmanBuster]); res == -1 {
		return fmt.Errorf("failed to load player buster image: %s", fname)
	}

	fname = common.ImagePath + "battle/character/player_pick.png"
	images[field.ObjectTypeRockmanPick] = make([]int32, 6)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 96, 124, images[field.ObjectTypeRockmanPick]); res == -1 {
		return fmt.Errorf("failed to load player pick image: %s", fname)
	}
	images[field.ObjectTypeRockmanPick][4] = images[field.ObjectTypeRockmanPick][3]
	images[field.ObjectTypeRockmanPick][5] = images[field.ObjectTypeRockmanPick][3]

	images[field.ObjectTypeRockmanStand] = make([]int32, 1)
	images[field.ObjectTypeRockmanStand][0] = images[field.ObjectTypeRockmanMove][0]

	return nil
}

func End() {
	for _, image := range images {
		for _, img := range image {
			dxlib.DeleteGraph(img)
		}
	}
}

func Object(objType int, imgNo int, x, y int) {
	if imgNo >= len(images[objType]) {
		imgNo = len(images[objType]) - 1
	}

	vx, vy := battlecommon.ViewPos(x, y)
	dxlib.DrawRotaGraph(vx, vy, 1, 0, images[objType][imgNo], dxlib.TRUE)
}
