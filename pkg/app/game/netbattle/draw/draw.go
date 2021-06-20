package draw

import (
	"fmt"

	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/field"
)

type Option struct {
	Reverse   bool
	SkillType int
	Speed     int
	ViewOfsX  int32
	ViewOfsY  int32
}

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

	skillPath := common.ImagePath + "battle/skill/"
	fname = skillPath + "キャノン_atk.png"
	images[field.ObjectTypeCannonAtk] = make([]int32, 24)
	if res := dxlib.LoadDivGraph(fname, 24, 8, 3, 120, 140, images[field.ObjectTypeCannonAtk]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = skillPath + "キャノン_body.png"
	images[field.ObjectTypeCannonBody] = make([]int32, 15)
	if res := dxlib.LoadDivGraph(fname, 15, 5, 3, 46, 40, images[field.ObjectTypeCannonBody]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = skillPath + "ミニボム.png"
	images[field.ObjectTypeMiniBomb] = make([]int32, 5)
	if res := dxlib.LoadDivGraph(fname, 5, 5, 1, 40, 30, images[field.ObjectTypeMiniBomb]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = skillPath + "ソード.png"
	images[field.ObjectTypeSword] = make([]int32, 12)
	if res := dxlib.LoadDivGraph(fname, 12, 4, 3, 160, 150, images[field.ObjectTypeSword]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = skillPath + "リカバリー.png"
	images[field.ObjectTypeRecover] = make([]int32, 8)
	if res := dxlib.LoadDivGraph(fname, 8, 8, 1, 84, 144, images[field.ObjectTypeRecover]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = skillPath + "スプレッドガン_atk.png"
	images[field.ObjectTypeSpreadGunAtk] = make([]int32, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 75, 76, images[field.ObjectTypeSpreadGunAtk]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = skillPath + "スプレッドガン_body.png"
	images[field.ObjectTypeSpreadGunBody] = make([]int32, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 56, 76, images[field.ObjectTypeSpreadGunBody]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = skillPath + "バルカン.png"
	images[field.ObjectTypeVulcan] = make([]int32, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 66, 50, images[field.ObjectTypeVulcan]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = skillPath + "ウェーブ_body.png"
	images[field.ObjectTypePick] = make([]int32, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 128, 136, images[field.ObjectTypePick]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = skillPath + "サンダーボール.png"
	images[field.ObjectTypeThunderBall] = make([]int32, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 64, 80, images[field.ObjectTypeThunderBall]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	// fname = skillPath + "ワイドショット_body.png"
	// if res := dxlib.LoadDivGraph(fname, 3, 3, 1, 56, 66, tmp); res == -1 {
	// 	return fmt.Errorf("failed to load image %s", fname)
	// }
	// for i := 0; i < 3; i++ {
	// 	imgWideShotBody = append(imgWideShotBody, tmp[i])
	// }
	// fname = skillPath + "ワイドショット_begin.png"
	// if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 90, 147, tmp); res == -1 {
	// 	return fmt.Errorf("failed to load image %s", fname)
	// }
	// for i := 0; i < 4; i++ {
	// 	imgWideShotBegin = append(imgWideShotBegin, tmp[i])
	// }
	// fname = skillPath + "ワイドショット_move.png"
	// if res := dxlib.LoadDivGraph(fname, 3, 3, 1, 90, 148, tmp); res == -1 {
	// 	return fmt.Errorf("failed to load image %s", fname)
	// }
	// for i := 0; i < 3; i++ {
	// 	imgWideShotMove = append(imgWideShotMove, tmp[i])
	// }

	fname = skillPath + "ショックウェーブ.png"
	images[field.ObjectTypeShockWave] = make([]int32, 7)
	if res := dxlib.LoadDivGraph(fname, 7, 7, 1, 100, 140, images[field.ObjectTypeShockWave]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	return nil
}

func End() {
	for _, image := range images {
		for _, img := range image {
			dxlib.DeleteGraph(img)
		}
	}
}

func Object(objType int, imgNo int, x, y int, opts ...Option) {
	if imgNo >= len(images[objType])/getTypeNum(objType) {
		imgNo = len(images[objType])/getTypeNum(objType) - 1
	}

	vx, vy := battlecommon.ViewPos(x, y)
	dxopts := dxlib.DrawRotaGraphOption{}

	if len(opts) > 0 {
		n := getTypeNum(objType)
		imgNo += opts[0].SkillType * n

		if opts[0].Reverse {
			flag := int32(dxlib.TRUE)
			dxopts.ReverseXFlag = &flag
		}

		// TODO reverse
		vx += opts[0].ViewOfsX
		vy += opts[0].ViewOfsY
	}

	dxlib.DrawRotaGraph(vx, vy, 1, 0, images[objType][imgNo], dxlib.TRUE, dxopts)
}

func GetImageInfo(objType int) (imageNum, delay int) {
	return len(images[objType]) / getTypeNum(objType), field.ImageDelays[objType]
}

func getTypeNum(objType int) int {
	switch objType {
	case field.ObjectTypeCannonAtk, field.ObjectTypeCannonBody:
		return skill.TypeCannonMax
	case field.ObjectTypeSword:
		return skill.TypeSwordMax
	}

	return 1
}
