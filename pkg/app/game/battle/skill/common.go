package skill

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

type chipNameDraw struct {
	name  string
	count int
	tm    int
}

func loadImages() error {
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

	return nil
}

func cleanupImages() {
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
}

func setChipNameDraw(name string) {
	battlecommon.AddSystem(&chipNameDraw{
		name:  name,
		count: 0,
		tm:    10,
	})
}

func (c *chipNameDraw) Draw() {
	r := float64(0)
	if c.count < c.tm {
		r = float64(c.count) / float64(c.tm)
	} else if c.count < c.tm*3 {
		r = 1
	} else if c.count < c.tm*4 {
		r = 1 - float64(c.count-c.tm*3)/float64(c.tm)
	}

	if r <= 0 {
		return
	}

	draw.ExtendString(50, 70, r, 0xffffff, c.name)
}

func (c *chipNameDraw) Process() bool {
	c.count++
	return c.count > c.tm*4
}