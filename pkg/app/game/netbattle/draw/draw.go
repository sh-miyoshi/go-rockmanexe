package draw

import (
	"fmt"

	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	appdraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/effect"
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
	imgObjs [object.TypeMax][]int32
	imgEffs [effect.TypeMax][]int32
)

func Init() error {
	if err := loadObjs(); err != nil {
		return fmt.Errorf("load objects failed: %w", err)
	}

	if err := loadEffects(); err != nil {
		return fmt.Errorf("load effects failed: %w", err)
	}

	return nil
}

func End() {
	for _, image := range imgObjs {
		for _, img := range image {
			dxlib.DeleteGraph(img)
		}
	}
}

func Object(objType int, imgNo int, x, y int, opts ...Option) {
	if imgNo >= len(imgObjs[objType])/getTypeNum(objType) {
		imgNo = len(imgObjs[objType])/getTypeNum(objType) - 1
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

	dxlib.DrawRotaGraph(vx, vy, 1, 0, imgObjs[objType][imgNo], dxlib.TRUE, dxopts)

	// Show HP
	if len(opts) > 0 && opts[0].ViewHP > 0 {
		appdraw.Number(vx, vy+40, int32(opts[0].ViewHP), appdraw.NumberOption{
			Color:    appdraw.NumberColorWhiteSmall,
			Centered: true,
		})
	}
}

func Effect(effType int, imgNo int, x, y int, ofsX, ofsY int32) {
	if imgNo >= len(imgEffs[effType]) {
		imgNo = len(imgEffs[effType]) - 1
	}

	vx, vy := battlecommon.ViewPos(x, y)
	vx += ofsX
	vy += ofsY

	dxlib.DrawRotaGraph(vx, vy, 1, 0, imgEffs[effType][imgNo], dxlib.TRUE)
}

func GetImageInfo(objType int) (imageNum, delay int) {
	return len(imgObjs[objType]) / getTypeNum(objType), object.ImageDelays[objType]
}

func GetEffectImageInfo(effType int) (imageNum, delay int) {
	return len(imgEffs[effType]), effect.Delays[effType]
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

func loadObjs() error {
	fname := common.ImagePath + "battle/character/player_move.png"
	imgObjs[object.TypeRockmanMove] = make([]int32, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 100, 100, imgObjs[object.TypeRockmanMove]); res == -1 {
		return fmt.Errorf("failed to load player move image: %s", fname)
	}

	fname = common.ImagePath + "battle/character/player_damaged.png"
	imgObjs[object.TypeRockmanDamage] = make([]int32, 6)
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 100, 100, imgObjs[object.TypeRockmanDamage]); res == -1 {
		return fmt.Errorf("failed to load player damage image: %s", fname)
	}
	// 1 -> 2,3  2-4 3-5
	imgObjs[object.TypeRockmanDamage][4] = imgObjs[object.TypeRockmanDamage][2]
	imgObjs[object.TypeRockmanDamage][5] = imgObjs[object.TypeRockmanDamage][3]
	imgObjs[object.TypeRockmanDamage][2] = imgObjs[object.TypeRockmanDamage][1]
	imgObjs[object.TypeRockmanDamage][3] = imgObjs[object.TypeRockmanDamage][1]

	fname = common.ImagePath + "battle/character/player_shot.png"
	imgObjs[object.TypeRockmanShot] = make([]int32, 6)
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 180, 100, imgObjs[object.TypeRockmanShot]); res == -1 {
		return fmt.Errorf("failed to load player shot image: %s", fname)
	}

	fname = common.ImagePath + "battle/character/player_cannon.png"
	imgObjs[object.TypeRockmanCannon] = make([]int32, 6)
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 100, 100, imgObjs[object.TypeRockmanCannon]); res == -1 {
		return fmt.Errorf("failed to load player cannon image: %s", fname)
	}

	fname = common.ImagePath + "battle/character/player_sword.png"
	imgObjs[object.TypeRockmanSword] = make([]int32, 7)
	if res := dxlib.LoadDivGraph(fname, 7, 7, 1, 128, 128, imgObjs[object.TypeRockmanSword]); res == -1 {
		return fmt.Errorf("failed to load player sword image: %s", fname)
	}

	fname = common.ImagePath + "battle/character/player_bomb.png"
	imgObjs[object.TypeRockmanBomb] = make([]int32, 7)
	if res := dxlib.LoadDivGraph(fname, 5, 5, 1, 100, 114, imgObjs[object.TypeRockmanBomb]); res == -1 {
		return fmt.Errorf("failed to load player bomb image: %s", fname)
	}
	imgObjs[object.TypeRockmanBomb][5] = imgObjs[object.TypeRockmanBomb][4]
	imgObjs[object.TypeRockmanBomb][6] = imgObjs[object.TypeRockmanBomb][4]

	fname = common.ImagePath + "battle/character/player_buster.png"
	imgObjs[object.TypeRockmanBuster] = make([]int32, 6)
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 180, 100, imgObjs[object.TypeRockmanBuster]); res == -1 {
		return fmt.Errorf("failed to load player buster image: %s", fname)
	}

	fname = common.ImagePath + "battle/character/player_pick.png"
	imgObjs[object.TypeRockmanPick] = make([]int32, 6)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 96, 124, imgObjs[object.TypeRockmanPick]); res == -1 {
		return fmt.Errorf("failed to load player pick image: %s", fname)
	}
	imgObjs[object.TypeRockmanPick][4] = imgObjs[object.TypeRockmanPick][3]
	imgObjs[object.TypeRockmanPick][5] = imgObjs[object.TypeRockmanPick][3]

	imgObjs[object.TypeRockmanStand] = make([]int32, 1)
	imgObjs[object.TypeRockmanStand][0] = imgObjs[object.TypeRockmanMove][0]

	skillPath := common.ImagePath + "battle/skill/"
	fname = skillPath + "キャノン_atk.png"
	imgObjs[object.TypeCannonAtk] = make([]int32, 24)
	if res := dxlib.LoadDivGraph(fname, 24, 8, 3, 120, 140, imgObjs[object.TypeCannonAtk]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = skillPath + "キャノン_body.png"
	imgObjs[object.TypeCannonBody] = make([]int32, 15)
	if res := dxlib.LoadDivGraph(fname, 15, 5, 3, 46, 40, imgObjs[object.TypeCannonBody]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = skillPath + "ミニボム.png"
	imgObjs[object.TypeMiniBomb] = make([]int32, 5)
	if res := dxlib.LoadDivGraph(fname, 5, 5, 1, 40, 30, imgObjs[object.TypeMiniBomb]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = skillPath + "ソード.png"
	imgObjs[object.TypeSword] = make([]int32, 12)
	if res := dxlib.LoadDivGraph(fname, 12, 4, 3, 160, 150, imgObjs[object.TypeSword]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = skillPath + "リカバリー.png"
	imgObjs[object.TypeRecover] = make([]int32, 8)
	if res := dxlib.LoadDivGraph(fname, 8, 8, 1, 84, 144, imgObjs[object.TypeRecover]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = skillPath + "スプレッドガン_atk.png"
	imgObjs[object.TypeSpreadGunAtk] = make([]int32, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 75, 76, imgObjs[object.TypeSpreadGunAtk]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = skillPath + "スプレッドガン_body.png"
	imgObjs[object.TypeSpreadGunBody] = make([]int32, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 56, 76, imgObjs[object.TypeSpreadGunBody]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = skillPath + "バルカン.png"
	imgObjs[object.TypeVulcan] = make([]int32, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 66, 50, imgObjs[object.TypeVulcan]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = skillPath + "ウェーブ_body.png"
	imgObjs[object.TypePick] = make([]int32, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 128, 136, imgObjs[object.TypePick]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = skillPath + "サンダーボール.png"
	imgObjs[object.TypeThunderBall] = make([]int32, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 64, 80, imgObjs[object.TypeThunderBall]); res == -1 {
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
	imgObjs[object.TypeShockWave] = make([]int32, 7)
	if res := dxlib.LoadDivGraph(fname, 7, 7, 1, 100, 140, imgObjs[object.TypeShockWave]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	return nil
}

func loadEffects() error {
	fname := common.ImagePath + "battle/effect/hit_small.png"
	imgEffs[effect.TypeHitSmallEffect] = make([]int32, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 40, 44, imgEffs[effect.TypeHitSmallEffect]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = common.ImagePath + "battle/effect/hit_big.png"
	imgEffs[effect.TypeHitBigEffect] = make([]int32, 6)
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 90, 76, imgEffs[effect.TypeHitBigEffect]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = common.ImagePath + "battle/effect/explode.png"
	imgEffs[effect.TypeExplodeEffect] = make([]int32, 16)
	if res := dxlib.LoadDivGraph(fname, 16, 8, 2, 110, 124, imgEffs[effect.TypeExplodeEffect]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = common.ImagePath + "battle/effect/cannon_hit.png"
	imgEffs[effect.TypeCannonHitEffect] = make([]int32, 7)
	if res := dxlib.LoadDivGraph(fname, 7, 7, 1, 110, 136, imgEffs[effect.TypeCannonHitEffect]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = common.ImagePath + "battle/effect/spread_hit.png"
	imgEffs[effect.TypeSpreadHitEffect] = make([]int32, 6)
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 92, 88, imgEffs[effect.TypeSpreadHitEffect]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	tmp := make([]int32, 8)
	fname = common.ImagePath + "battle/effect/vulcan_hit.png"
	if res := dxlib.LoadDivGraph(fname, 8, 8, 1, 50, 58, tmp); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	imgEffs[effect.TypeVulcanHit1Effect] = []int32{}
	imgEffs[effect.TypeVulcanHit2Effect] = []int32{}
	for i := 0; i < 4; i++ {
		imgEffs[effect.TypeVulcanHit1Effect] = append(imgEffs[effect.TypeVulcanHit1Effect], tmp[i])
		imgEffs[effect.TypeVulcanHit2Effect] = append(imgEffs[effect.TypeVulcanHit2Effect], tmp[i+4])
	}

	return nil
}
