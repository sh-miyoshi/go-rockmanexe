package skilldraw

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

var (
	imgCannonAtk         [resources.SkillTypeCannonMax][]int
	imgCannonBody        [resources.SkillTypeCannonMax][]int
	imgSword             [resources.SkillTypeSwordMax][]int
	imgBombThrow         []int
	imgShockWave         []int
	imgRecover           []int
	imgSpreadGunAtk      []int
	imgSpreadGunBody     []int
	imgVulcan            []int
	imgPick              []int
	imgThunderBall       []int
	imgWideShotBody      []int
	imgWideShotBegin     []int
	imgWideShotMove      []int
	imgBoomerang         []int
	imgAquamanShot       []int
	imgBambooLance       []int
	imgDreamSword        []int
	imgGarooBreath       []int
	imgFlamePillar       []int
	imgFlameLineBody     []int
	imgHeatShotBody      []int
	imgHeatShotAtk       []int
	imgAreaStealMain     []int
	imgAreaStealPanel    []int
	imgAquamanCharStand  []int
	imgAquamanCharCreate []int
	imgSpreadHit         []int
	imgCountBomb         []int
)

func LoadImages() error {
	path := common.ImagePath + "battle/skill/"

	tmp := make([]int, 24)
	fname := path + "キャノン_atk.png"
	if res := dxlib.LoadDivGraph(fname, 24, 8, 3, 120, 140, tmp); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	for i := 0; i < 8; i++ {
		imgCannonAtk[0] = append(imgCannonAtk[0], tmp[i])
		imgCannonAtk[1] = append(imgCannonAtk[1], tmp[i+8])
		imgCannonAtk[2] = append(imgCannonAtk[2], tmp[i+16])
	}
	fname = path + "キャノン_body.png"
	if res := dxlib.LoadDivGraph(fname, 15, 5, 3, 46, 40, tmp); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	for i := 0; i < 5; i++ {
		imgCannonBody[0] = append(imgCannonBody[0], tmp[i])
		imgCannonBody[1] = append(imgCannonBody[1], tmp[i+5])
		imgCannonBody[2] = append(imgCannonBody[2], tmp[i+10])
	}

	fname = path + "ミニボム.png"
	if res := dxlib.LoadDivGraph(fname, 5, 5, 1, 40, 30, tmp); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	for i := 0; i < 5; i++ {
		imgBombThrow = append(imgBombThrow, tmp[i])
	}

	fname = path + "ソード.png"
	if res := dxlib.LoadDivGraph(fname, 12, 4, 3, 160, 150, tmp); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	for i := 0; i < 4; i++ {
		// Note: In the image, the order of wide sword and long sword is swapped.
		imgSword[0] = append(imgSword[0], tmp[i])
		imgSword[1] = append(imgSword[1], tmp[i+8])
		imgSword[2] = append(imgSword[2], tmp[i+4])
	}

	fname = path + "ショックウェーブ.png"
	if res := dxlib.LoadDivGraph(fname, 7, 7, 1, 100, 140, tmp); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	for i := 0; i < 7; i++ {
		imgShockWave = append(imgShockWave, tmp[i])
	}

	fname = path + "リカバリー.png"
	if res := dxlib.LoadDivGraph(fname, 8, 8, 1, 84, 144, tmp); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	for i := 0; i < 8; i++ {
		imgRecover = append(imgRecover, tmp[i])
	}

	fname = path + "スプレッドガン_atk.png"
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 75, 76, tmp); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	for i := 0; i < 4; i++ {
		imgSpreadGunAtk = append(imgSpreadGunAtk, tmp[i])
	}
	fname = path + "スプレッドガン_body.png"
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 56, 76, tmp); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	for i := 0; i < 4; i++ {
		imgSpreadGunBody = append(imgSpreadGunBody, tmp[i])
	}

	fname = path + "バルカン.png"
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 66, 50, tmp); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	for i := 0; i < 4; i++ {
		imgVulcan = append(imgVulcan, tmp[i])
	}

	fname = path + "ウェーブ_body.png"
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 128, 136, tmp); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	for i := 0; i < 4; i++ {
		imgPick = append(imgPick, tmp[i])
	}
	imgPick = append(imgPick, tmp[3])
	imgPick = append(imgPick, tmp[3])

	fname = path + "サンダーボール.png"
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 64, 80, tmp); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	for i := 0; i < 4; i++ {
		imgThunderBall = append(imgThunderBall, tmp[i])
	}

	fname = path + "ワイドショット_body.png"
	if res := dxlib.LoadDivGraph(fname, 3, 3, 1, 56, 66, tmp); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	for i := 0; i < 3; i++ {
		imgWideShotBody = append(imgWideShotBody, tmp[i])
	}
	fname = path + "ワイドショット_begin.png"
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 90, 147, tmp); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	for i := 0; i < 4; i++ {
		imgWideShotBegin = append(imgWideShotBegin, tmp[i])
	}
	fname = path + "ワイドショット_move.png"
	if res := dxlib.LoadDivGraph(fname, 3, 3, 1, 90, 148, tmp); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	for i := 0; i < 3; i++ {
		imgWideShotMove = append(imgWideShotMove, tmp[i])
	}

	fname = path + "ブーメラン.png"
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 100, 80, tmp); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	for i := 0; i < 4; i++ {
		imgBoomerang = append(imgBoomerang, tmp[i])
	}
	fname = path + "aquaman_shot.png"
	imgAquamanShot = make([]int, 1)
	if imgAquamanShot[0] = dxlib.LoadGraph(fname); imgAquamanShot[0] == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	fname = path + "バンブーランス.png"
	imgBambooLance = make([]int, 1)
	if imgBambooLance[0] = dxlib.LoadGraph(fname); imgBambooLance[0] == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	fname = path + "ドリームソード.png"
	imgDreamSword = make([]int, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 200, 188, imgDreamSword); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	fname = path + "ガルー_atk.png"
	imgGarooBreath = make([]int, 3)
	if res := dxlib.LoadDivGraph(fname, 3, 3, 1, 108, 62, imgGarooBreath); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	fname = path + "フレイムライン_火柱.png"
	imgFlamePillar = make([]int, 5)
	if res := dxlib.LoadDivGraph(fname, 5, 5, 1, 72, 120, imgFlamePillar); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	fname = path + "フレイムライン_body.png"
	imgFlameLineBody = make([]int, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 64, 64, imgFlameLineBody); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	fname = path + "ヒートショット_body.png"
	imgHeatShotBody = make([]int, 5)
	if res := dxlib.LoadDivGraph(fname, 5, 5, 1, 60, 40, imgHeatShotBody); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	fname = path + "ヒートショット_atk.png"
	imgHeatShotAtk = make([]int, 3)
	if res := dxlib.LoadDivGraph(fname, 3, 3, 1, 60, 45, imgHeatShotAtk); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	fname = path + "エリアスチール_main.png"
	imgAreaStealMain = make([]int, 5)
	if res := dxlib.LoadDivGraph(fname, 5, 5, 1, 74, 69, imgAreaStealMain); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	fname = path + "エリアスチール_panel.png"
	imgAreaStealPanel = make([]int, 6)
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 110, 76, imgAreaStealPanel); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = path + "カウントボム.png"
	imgCountBomb = make([]int, 1)
	if imgCountBomb[0] = dxlib.LoadGraph(fname); imgCountBomb[0] == -1 {
		return fmt.Errorf("failed to load image: %s", fname)
	}

	fname = common.ImagePath + "battle/character/アクアマン_stand.png"
	imgAquamanCharStand = make([]int, 9)
	if res := dxlib.LoadDivGraph(fname, 9, 9, 1, 62, 112, imgAquamanCharStand); res == -1 {
		return fmt.Errorf("failed to load image: %s", fname)
	}

	fname = common.ImagePath + "battle/character/アクアマン_create.png"
	imgAquamanCharCreate = make([]int, 1)
	if res := dxlib.LoadDivGraph(fname, 1, 1, 1, 80, 92, imgAquamanCharCreate); res == -1 {
		return fmt.Errorf("failed to load image: %s", fname)
	}

	fname = common.ImagePath + "battle/effect/spread_and_bamboo_hit.png"
	imgSpreadHit = make([]int, 6)
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 92, 88, imgSpreadHit); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	return nil
}

func ClearImages() {
	for i := 0; i < 3; i++ {
		for j := 0; j < len(imgCannonAtk[i]); j++ {
			dxlib.DeleteGraph(imgCannonAtk[i][j])
		}
		imgCannonAtk[i] = []int{}
		for j := 0; j < len(imgCannonBody[i]); j++ {
			dxlib.DeleteGraph(imgCannonBody[i][j])
		}
		imgCannonBody[i] = []int{}
	}
	for i := 0; i < 3; i++ {
		for j := 0; j < len(imgSword[i]); j++ {
			dxlib.DeleteGraph(imgSword[i][j])
		}
		imgSword[i] = []int{}
	}
	for i := 0; i < len(imgBombThrow); i++ {
		dxlib.DeleteGraph(imgBombThrow[i])
	}
	imgBombThrow = []int{}
	for i := 0; i < len(imgShockWave); i++ {
		dxlib.DeleteGraph(imgShockWave[i])
	}
	imgShockWave = []int{}
	for i := 0; i < len(imgSpreadGunAtk); i++ {
		dxlib.DeleteGraph(imgSpreadGunAtk[i])
	}
	imgSpreadGunAtk = []int{}
	for i := 0; i < len(imgSpreadGunBody); i++ {
		dxlib.DeleteGraph(imgSpreadGunBody[i])
	}
	imgSpreadGunBody = []int{}
	for i := 0; i < len(imgVulcan); i++ {
		dxlib.DeleteGraph(imgVulcan[i])
	}
	imgVulcan = []int{}
	for i := 0; i < len(imgRecover); i++ {
		dxlib.DeleteGraph(imgRecover[i])
	}
	imgRecover = []int{}
	for i := 0; i < len(imgPick); i++ {
		dxlib.DeleteGraph(imgPick[i])
	}
	imgPick = []int{}
	for i := 0; i < len(imgThunderBall); i++ {
		dxlib.DeleteGraph(imgThunderBall[i])
	}
	imgThunderBall = []int{}
	for i := 0; i < len(imgWideShotBody); i++ {
		dxlib.DeleteGraph(imgWideShotBody[i])
	}
	imgWideShotBody = []int{}
	for i := 0; i < len(imgWideShotBegin); i++ {
		dxlib.DeleteGraph(imgWideShotBegin[i])
	}
	imgWideShotBegin = []int{}
	for i := 0; i < len(imgWideShotMove); i++ {
		dxlib.DeleteGraph(imgWideShotMove[i])
	}
	imgWideShotMove = []int{}
	for i := 0; i < len(imgBoomerang); i++ {
		dxlib.DeleteGraph(imgBoomerang[i])
	}
	imgBoomerang = []int{}
	for i := 0; i < len(imgAquamanShot); i++ {
		dxlib.DeleteGraph(imgAquamanShot[i])
	}
	imgAquamanShot = []int{}
	for i := 0; i < len(imgBambooLance); i++ {
		dxlib.DeleteGraph(imgBambooLance[i])
	}
	imgBambooLance = []int{}
	for i := 0; i < len(imgDreamSword); i++ {
		dxlib.DeleteGraph(imgDreamSword[i])
	}
	imgDreamSword = []int{}
	for i := 0; i < len(imgGarooBreath); i++ {
		dxlib.DeleteGraph(imgGarooBreath[i])
	}
	imgGarooBreath = []int{}
	for i := 0; i < len(imgFlamePillar); i++ {
		dxlib.DeleteGraph(imgFlamePillar[i])
	}
	imgFlamePillar = []int{}
	for i := 0; i < len(imgFlameLineBody); i++ {
		dxlib.DeleteGraph(imgFlameLineBody[i])
	}
	imgFlameLineBody = []int{}
	for i := 0; i < len(imgHeatShotBody); i++ {
		dxlib.DeleteGraph(imgHeatShotBody[i])
	}
	imgHeatShotBody = []int{}
	for i := 0; i < len(imgHeatShotAtk); i++ {
		dxlib.DeleteGraph(imgHeatShotAtk[i])
	}
	imgHeatShotAtk = []int{}
	for i := 0; i < len(imgAreaStealMain); i++ {
		dxlib.DeleteGraph(imgAreaStealMain[i])
	}
	imgAreaStealMain = []int{}
	for i := 0; i < len(imgAreaStealPanel); i++ {
		dxlib.DeleteGraph(imgAreaStealPanel[i])
	}
	imgAreaStealPanel = []int{}
	for i := 0; i < len(imgAquamanCharStand); i++ {
		dxlib.DeleteGraph(imgAquamanCharStand[i])
	}
	imgAquamanCharStand = []int{}
	for i := 0; i < len(imgAquamanCharCreate); i++ {
		dxlib.DeleteGraph(imgAquamanCharCreate[i])
	}
	imgAquamanCharCreate = []int{}
	for i := 0; i < len(imgSpreadHit); i++ {
		dxlib.DeleteGraph(imgSpreadHit[i])
	}
	imgSpreadHit = []int{}
	for i := 0; i < len(imgCountBomb); i++ {
		dxlib.DeleteGraph(imgCountBomb[i])
	}
	imgCountBomb = []int{}
}
