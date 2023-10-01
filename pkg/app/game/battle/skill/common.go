package skill

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

type chipNameDraw struct {
	name         string
	count        int
	tm           int
	isUserPlayer bool
}

func loadImages() error {
	path := common.ImagePath + "battle/skill/"

	tmp := make([]int, 12)

	fname := path + "ミニボム.png"
	if res := dxlib.LoadDivGraph(fname, 5, 5, 1, 40, 30, tmp); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	for i := 0; i < 5; i++ {
		imgBombThrow = append(imgBombThrow, tmp[i])
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

	return nil
}

func cleanupImages() {
	for i := 0; i < len(imgBombThrow); i++ {
		dxlib.DeleteGraph(imgBombThrow[i])
	}
	imgBombThrow = []int{}
	for i := 0; i < len(imgVulcan); i++ {
		dxlib.DeleteGraph(imgVulcan[i])
	}
	imgVulcan = []int{}
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
}

func SetChipNameDraw(name string, isUserPlayer bool) {
	battlecommon.AddSystem(&chipNameDraw{
		name:         name,
		count:        0,
		tm:           10,
		isUserPlayer: isUserPlayer,
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

	x := 50
	if !c.isUserPlayer {
		x = 300
	}

	draw.ExtendString(x, 70, r, 0xffffff, c.name)
}

func (c *chipNameDraw) Process() bool {
	c.count++
	return c.count > c.tm*4
}
