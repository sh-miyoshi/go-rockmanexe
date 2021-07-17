package draw

import (
	"fmt"

	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	appdraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/object"
)

type Option struct {
	Reverse   bool
	SkillType int
	Speed     int
	ViewOfsX  int32
	ViewOfsY  int32
	ViewHP    int
}

var (
	images [object.TypeMax][]int32
)

func Init() error {
	fname := common.ImagePath + "battle/character/player_move.png"
	images[object.TypeRockmanMove] = make([]int32, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 100, 100, images[object.TypeRockmanMove]); res == -1 {
		return fmt.Errorf("failed to load player move image: %s", fname)
	}

	fname = common.ImagePath + "battle/character/player_damaged.png"
	images[object.TypeRockmanDamage] = make([]int32, 6)
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 100, 100, images[object.TypeRockmanDamage]); res == -1 {
		return fmt.Errorf("failed to load player damage image: %s", fname)
	}
	// 1 -> 2,3  2-4 3-5
	images[object.TypeRockmanDamage][4] = images[object.TypeRockmanDamage][2]
	images[object.TypeRockmanDamage][5] = images[object.TypeRockmanDamage][3]
	images[object.TypeRockmanDamage][2] = images[object.TypeRockmanDamage][1]
	images[object.TypeRockmanDamage][3] = images[object.TypeRockmanDamage][1]

	fname = common.ImagePath + "battle/character/player_shot.png"
	images[object.TypeRockmanShot] = make([]int32, 6)
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 180, 100, images[object.TypeRockmanShot]); res == -1 {
		return fmt.Errorf("failed to load player shot image: %s", fname)
	}

	fname = common.ImagePath + "battle/character/player_cannon.png"
	images[object.TypeRockmanCannon] = make([]int32, 6)
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 100, 100, images[object.TypeRockmanCannon]); res == -1 {
		return fmt.Errorf("failed to load player cannon image: %s", fname)
	}

	fname = common.ImagePath + "battle/character/player_sword.png"
	images[object.TypeRockmanSword] = make([]int32, 7)
	if res := dxlib.LoadDivGraph(fname, 7, 7, 1, 128, 128, images[object.TypeRockmanSword]); res == -1 {
		return fmt.Errorf("failed to load player sword image: %s", fname)
	}

	fname = common.ImagePath + "battle/character/player_bomb.png"
	images[object.TypeRockmanBomb] = make([]int32, 7)
	if res := dxlib.LoadDivGraph(fname, 5, 5, 1, 100, 114, images[object.TypeRockmanBomb]); res == -1 {
		return fmt.Errorf("failed to load player bomb image: %s", fname)
	}
	images[object.TypeRockmanBomb][5] = images[object.TypeRockmanBomb][4]
	images[object.TypeRockmanBomb][6] = images[object.TypeRockmanBomb][4]

	fname = common.ImagePath + "battle/character/player_buster.png"
	images[object.TypeRockmanBuster] = make([]int32, 6)
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 180, 100, images[object.TypeRockmanBuster]); res == -1 {
		return fmt.Errorf("failed to load player buster image: %s", fname)
	}

	fname = common.ImagePath + "battle/character/player_pick.png"
	images[object.TypeRockmanPick] = make([]int32, 6)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 96, 124, images[object.TypeRockmanPick]); res == -1 {
		return fmt.Errorf("failed to load player pick image: %s", fname)
	}
	images[object.TypeRockmanPick][4] = images[object.TypeRockmanPick][3]
	images[object.TypeRockmanPick][5] = images[object.TypeRockmanPick][3]

	images[object.TypeRockmanStand] = make([]int32, 1)
	images[object.TypeRockmanStand][0] = images[object.TypeRockmanMove][0]

	skillPath := common.ImagePath + "battle/skill/"
	fname = skillPath + "キャノン_atk.png"
	images[object.TypeCannonAtk] = make([]int32, 24)
	if res := dxlib.LoadDivGraph(fname, 24, 8, 3, 120, 140, images[object.TypeCannonAtk]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = skillPath + "キャノン_body.png"
	images[object.TypeCannonBody] = make([]int32, 15)
	if res := dxlib.LoadDivGraph(fname, 15, 5, 3, 46, 40, images[object.TypeCannonBody]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = skillPath + "ミニボム.png"
	images[object.TypeMiniBomb] = make([]int32, 5)
	if res := dxlib.LoadDivGraph(fname, 5, 5, 1, 40, 30, images[object.TypeMiniBomb]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = skillPath + "ソード.png"
	images[object.TypeSword] = make([]int32, 12)
	if res := dxlib.LoadDivGraph(fname, 12, 4, 3, 160, 150, images[object.TypeSword]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = skillPath + "リカバリー.png"
	images[object.TypeRecover] = make([]int32, 8)
	if res := dxlib.LoadDivGraph(fname, 8, 8, 1, 84, 144, images[object.TypeRecover]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = skillPath + "スプレッドガン_atk.png"
	images[object.TypeSpreadGunAtk] = make([]int32, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 75, 76, images[object.TypeSpreadGunAtk]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = skillPath + "スプレッドガン_body.png"
	images[object.TypeSpreadGunBody] = make([]int32, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 56, 76, images[object.TypeSpreadGunBody]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = skillPath + "バルカン.png"
	images[object.TypeVulcan] = make([]int32, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 66, 50, images[object.TypeVulcan]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = skillPath + "ウェーブ_body.png"
	images[object.TypePick] = make([]int32, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 128, 136, images[object.TypePick]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = skillPath + "サンダーボール.png"
	images[object.TypeThunderBall] = make([]int32, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 64, 80, images[object.TypeThunderBall]); res == -1 {
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
	images[object.TypeShockWave] = make([]int32, 7)
	if res := dxlib.LoadDivGraph(fname, 7, 7, 1, 100, 140, images[object.TypeShockWave]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = common.ImagePath + "battle/effect/hit_small.png"
	images[object.TypeHitSmallEffect] = make([]int32, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 40, 44, images[object.TypeHitSmallEffect]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = common.ImagePath + "battle/effect/hit_big.png"
	images[object.TypeHitBigEffect] = make([]int32, 6)
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 90, 76, images[object.TypeHitBigEffect]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = common.ImagePath + "battle/effect/explode.png"
	images[object.TypeExplodeEffect] = make([]int32, 16)
	if res := dxlib.LoadDivGraph(fname, 16, 8, 2, 110, 124, images[object.TypeExplodeEffect]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = common.ImagePath + "battle/effect/cannon_hit.png"
	images[object.TypeCannonHitEffect] = make([]int32, 7)
	if res := dxlib.LoadDivGraph(fname, 7, 7, 1, 110, 136, images[object.TypeCannonHitEffect]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = common.ImagePath + "battle/effect/spread_hit.png"
	images[object.TypeSpreadHitEffect] = make([]int32, 6)
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 92, 88, images[object.TypeSpreadHitEffect]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	tmp := make([]int32, 8)
	fname = common.ImagePath + "battle/effect/vulcan_hit.png"
	if res := dxlib.LoadDivGraph(fname, 8, 8, 1, 50, 58, tmp); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	images[object.TypeVulcanHit1Effect] = []int32{}
	images[object.TypeVulcanHit2Effect] = []int32{}
	for i := 0; i < 4; i++ {
		images[object.TypeVulcanHit1Effect] = append(images[object.TypeVulcanHit1Effect], tmp[i])
		images[object.TypeVulcanHit2Effect] = append(images[object.TypeVulcanHit2Effect], tmp[i+4])
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
			opts[0].ViewOfsX *= -1
		}

		vx += opts[0].ViewOfsX
		vy += opts[0].ViewOfsY
	}

	dxlib.DrawRotaGraph(vx, vy, 1, 0, images[objType][imgNo], dxlib.TRUE, dxopts)

	// Show HP
	if len(opts) > 0 && opts[0].ViewHP > 0 {
		appdraw.Number(vx, vy+40, int32(opts[0].ViewHP), appdraw.NumberOption{
			Color:    appdraw.NumberColorWhiteSmall,
			Centered: true,
		})
	}
}

func GetImageInfo(objType int) (imageNum, delay int) {
	return len(images[objType]) / getTypeNum(objType), object.ImageDelays[objType]
}

func getTypeNum(objType int) int {
	switch objType {
	case object.TypeCannonAtk, object.TypeCannonBody:
		return skill.TypeCannonMax
	case object.TypeSword:
		return skill.TypeSwordMax
	}

	return 1
}
