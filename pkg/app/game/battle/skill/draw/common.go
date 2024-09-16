package skilldraw

import (
	"github.com/cockroachdb/errors"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	imageTypeCannonAtk int = iota
	imageTypeCannonBody
	imageTypeSword
	imageTypeBombThrow
	imageTypeShockWave
	imageTypeRecover
	imageTypeSpreadGunAtk
	imageTypeSpreadGunBody
	imageTypeVulcan
	imageTypePick
	imageTypeThunderBall
	imageTypeWideShotBody
	imageTypeWideShotBegin
	imageTypeWideShotMove
	imageTypeBoomerang
	imageTypeAquamanShot
	imageTypeBambooLance
	imageTypeDreamSword
	imageTypeGarooBreath
	imageTypeFlamePillar
	imageTypeFlameLineBody
	imageTypeHeatShotBody
	imageTypeHeatShotAtk
	imageTypeAreaStealMain
	imageTypeAreaStealPanel
	imageTypeAquamanCharStand
	imageTypeAquamanCharCreate
	imageTypeSpreadHit
	imageTypeCountBomb
	imageTypeTornadoAtk
	imageTypeTornadoBody
	imageTypeCirkillShot
	imageTypeCountBombNumber
	imageTypeShrimpyAtkBegin
	imageTypeShrimpyAtkMove
	imageTypeBubbleShotBody
	imageTypeBubbleShotAtk
	imageTypeForteHellsRolling
	imageTypeForteDarkArmBlade
	imageTypeForteShootingBuster
	imageTypeForteDarknessOverload

	imageTypeMax
)

var (
	images [imageTypeMax][]int
)

func LoadImages() error {
	path := config.ImagePath + "battle/skill/"

	images[imageTypeCannonAtk] = make([]int, 24)
	fname := path + "キャノン_atk.png"
	if res := dxlib.LoadDivGraph(fname, 24, 8, 3, 120, 140, images[imageTypeCannonAtk]); res == -1 {
		return errors.Newf("failed to load image %s", fname)
	}

	fname = path + "キャノン_body.png"
	images[imageTypeCannonBody] = make([]int, 15)
	if res := dxlib.LoadDivGraph(fname, 15, 5, 3, 46, 40, images[imageTypeCannonBody]); res == -1 {
		return errors.Newf("failed to load image %s", fname)
	}

	fname = path + "ミニボム.png"
	images[imageTypeBombThrow] = make([]int, 5)
	if res := dxlib.LoadDivGraph(fname, 5, 5, 1, 40, 30, images[imageTypeBombThrow]); res == -1 {
		return errors.Newf("failed to load image %s", fname)
	}

	fname = path + "ソード.png"
	images[imageTypeSword] = make([]int, 12)
	if res := dxlib.LoadDivGraph(fname, 12, 4, 3, 160, 150, images[imageTypeSword]); res == -1 {
		return errors.Newf("failed to load image %s", fname)
	}

	fname = path + "ショックウェーブ.png"
	images[imageTypeShockWave] = make([]int, 7)
	if res := dxlib.LoadDivGraph(fname, 7, 7, 1, 100, 140, images[imageTypeShockWave]); res == -1 {
		return errors.Newf("failed to load image %s", fname)
	}

	fname = path + "リカバリー.png"
	images[imageTypeRecover] = make([]int, 8)
	if res := dxlib.LoadDivGraph(fname, 8, 8, 1, 84, 144, images[imageTypeRecover]); res == -1 {
		return errors.Newf("failed to load image %s", fname)
	}

	fname = path + "スプレッドガン_atk.png"
	images[imageTypeSpreadGunAtk] = make([]int, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 75, 76, images[imageTypeSpreadGunAtk]); res == -1 {
		return errors.Newf("failed to load image %s", fname)
	}
	fname = path + "スプレッドガン_body.png"
	images[imageTypeSpreadGunBody] = make([]int, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 56, 76, images[imageTypeSpreadGunBody]); res == -1 {
		return errors.Newf("failed to load image %s", fname)
	}

	fname = path + "バルカン.png"
	images[imageTypeVulcan] = make([]int, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 66, 50, images[imageTypeVulcan]); res == -1 {
		return errors.Newf("failed to load image %s", fname)
	}

	fname = path + "ウェーブ_body.png"
	images[imageTypePick] = make([]int, 6)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 128, 136, images[imageTypePick]); res == -1 {
		return errors.Newf("failed to load image %s", fname)
	}
	images[imageTypePick][4] = images[imageTypePick][3]
	images[imageTypePick][5] = images[imageTypePick][3]

	fname = path + "サンダーボール.png"
	images[imageTypeThunderBall] = make([]int, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 64, 80, images[imageTypeThunderBall]); res == -1 {
		return errors.Newf("failed to load image %s", fname)
	}

	fname = path + "ワイドショット_body.png"
	images[imageTypeWideShotBody] = make([]int, 3)
	if res := dxlib.LoadDivGraph(fname, 3, 3, 1, 56, 66, images[imageTypeWideShotBody]); res == -1 {
		return errors.Newf("failed to load image %s", fname)
	}
	fname = path + "ワイドショット_begin.png"
	images[imageTypeWideShotBegin] = make([]int, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 90, 147, images[imageTypeWideShotBegin]); res == -1 {
		return errors.Newf("failed to load image %s", fname)
	}
	fname = path + "ワイドショット_move.png"
	images[imageTypeWideShotMove] = make([]int, 3)
	if res := dxlib.LoadDivGraph(fname, 3, 3, 1, 90, 148, images[imageTypeWideShotMove]); res == -1 {
		return errors.Newf("failed to load image %s", fname)
	}

	fname = path + "ブーメラン.png"
	images[imageTypeBoomerang] = make([]int, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 100, 80, images[imageTypeBoomerang]); res == -1 {
		return errors.Newf("failed to load image %s", fname)
	}

	fname = path + "aquaman_shot.png"
	images[imageTypeAquamanShot] = make([]int, 1)
	if images[imageTypeAquamanShot][0] = dxlib.LoadGraph(fname); images[imageTypeAquamanShot][0] == -1 {
		return errors.Newf("failed to load image %s", fname)
	}

	fname = path + "バンブーランス.png"
	images[imageTypeBambooLance] = make([]int, 1)
	if images[imageTypeBambooLance][0] = dxlib.LoadGraph(fname); images[imageTypeBambooLance][0] == -1 {
		return errors.Newf("failed to load image %s", fname)
	}

	fname = path + "ドリームソード.png"
	images[imageTypeDreamSword] = make([]int, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 200, 188, images[imageTypeDreamSword]); res == -1 {
		return errors.Newf("failed to load image %s", fname)
	}

	fname = path + "ガルー_atk.png"
	images[imageTypeGarooBreath] = make([]int, 3)
	if res := dxlib.LoadDivGraph(fname, 3, 3, 1, 108, 62, images[imageTypeGarooBreath]); res == -1 {
		return errors.Newf("failed to load image %s", fname)
	}

	fname = path + "フレイムライン_火柱.png"
	images[imageTypeFlamePillar] = make([]int, 5)
	if res := dxlib.LoadDivGraph(fname, 5, 5, 1, 72, 120, images[imageTypeFlamePillar]); res == -1 {
		return errors.Newf("failed to load image %s", fname)
	}
	fname = path + "フレイムライン_body.png"
	images[imageTypeFlameLineBody] = make([]int, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 64, 64, images[imageTypeFlameLineBody]); res == -1 {
		return errors.Newf("failed to load image %s", fname)
	}

	fname = path + "ヒートショット_body.png"
	images[imageTypeHeatShotBody] = make([]int, 5)
	if res := dxlib.LoadDivGraph(fname, 5, 5, 1, 60, 40, images[imageTypeHeatShotBody]); res == -1 {
		return errors.Newf("failed to load image %s", fname)
	}
	fname = path + "ヒートショット_atk.png"
	images[imageTypeHeatShotAtk] = make([]int, 3)
	if res := dxlib.LoadDivGraph(fname, 3, 3, 1, 60, 45, images[imageTypeHeatShotAtk]); res == -1 {
		return errors.Newf("failed to load image %s", fname)
	}

	fname = path + "エリアスチール_main.png"
	images[imageTypeAreaStealMain] = make([]int, 5)
	if res := dxlib.LoadDivGraph(fname, 5, 5, 1, 74, 69, images[imageTypeAreaStealMain]); res == -1 {
		return errors.Newf("failed to load image %s", fname)
	}
	fname = path + "エリアスチール_panel.png"
	images[imageTypeAreaStealPanel] = make([]int, 6)
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 110, 76, images[imageTypeAreaStealPanel]); res == -1 {
		return errors.Newf("failed to load image %s", fname)
	}

	fname = path + "カウントボム.png"
	images[imageTypeCountBomb] = make([]int, 1)
	if images[imageTypeCountBomb][0] = dxlib.LoadGraph(fname); images[imageTypeCountBomb][0] == -1 {
		return errors.Newf("failed to load image: %s", fname)
	}

	fname = path + "トルネード_atk.png"
	images[imageTypeTornadoAtk] = make([]int, 3)
	if res := dxlib.LoadDivGraph(fname, 3, 3, 1, 63, 96, images[imageTypeTornadoAtk]); res == -1 {
		return errors.Newf("failed to load image %s", fname)
	}
	fname = path + "トルネード_body.png"
	images[imageTypeTornadoBody] = make([]int, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 64, 64, images[imageTypeTornadoBody]); res == -1 {
		return errors.Newf("failed to load image %s", fname)
	}

	fname = config.ImagePath + "battle/character/アクアマン_stand.png"
	images[imageTypeAquamanCharStand] = make([]int, 9)
	if res := dxlib.LoadDivGraph(fname, 9, 9, 1, 62, 112, images[imageTypeAquamanCharStand]); res == -1 {
		return errors.Newf("failed to load image: %s", fname)
	}
	fname = config.ImagePath + "battle/character/アクアマン_create.png"
	images[imageTypeAquamanCharCreate] = make([]int, 1)
	if res := dxlib.LoadDivGraph(fname, 1, 1, 1, 80, 92, images[imageTypeAquamanCharCreate]); res == -1 {
		return errors.Newf("failed to load image: %s", fname)
	}

	fname = config.ImagePath + "battle/effect/spread_and_bamboo_hit.png"
	images[imageTypeSpreadHit] = make([]int, 6)
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 92, 88, images[imageTypeSpreadHit]); res == -1 {
		return errors.Newf("failed to load image %s", fname)
	}

	fname = path + "サーキラー_atk.png"
	images[imageTypeCirkillShot] = make([]int, 3)
	if res := dxlib.LoadDivGraph(fname, 3, 3, 1, 108, 62, images[imageTypeCirkillShot]); res == -1 {
		return errors.Newf("failed to load image %s", fname)
	}

	fname = config.ImagePath + "battle/skill/カウントボム_数字.png"
	images[imageTypeCountBombNumber] = make([]int, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 16, 16, images[imageTypeCountBombNumber]); res == -1 {
		return errors.Newf("failed to load image: %s", fname)
	}

	fname = config.ImagePath + "battle/skill/エビロン_atk_begin.png"
	images[imageTypeShrimpyAtkBegin] = make([]int, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 46, 44, images[imageTypeShrimpyAtkBegin]); res == -1 {
		return errors.Newf("failed to load image: %s", fname)
	}
	fname = config.ImagePath + "battle/skill/エビロン_atk_move.png"
	images[imageTypeShrimpyAtkMove] = make([]int, 8)
	if res := dxlib.LoadDivGraph(fname, 8, 8, 1, 74, 60, images[imageTypeShrimpyAtkMove]); res == -1 {
		return errors.Newf("failed to load image: %s", fname)
	}

	fname = config.ImagePath + "battle/skill/バブルショット_body.png"
	images[imageTypeBubbleShotBody] = make([]int, 5)
	if res := dxlib.LoadDivGraph(fname, 5, 5, 1, 70, 60, images[imageTypeBubbleShotBody]); res == -1 {
		return errors.Newf("failed to load image: %s", fname)
	}
	fname = config.ImagePath + "battle/skill/バブルショット_atk.png"
	images[imageTypeBubbleShotAtk] = make([]int, 2)
	if res := dxlib.LoadDivGraph(fname, 2, 2, 1, 40, 60, images[imageTypeBubbleShotAtk]); res == -1 {
		return errors.Newf("failed to load image: %s", fname)
	}

	fname = config.ImagePath + "battle/skill/フォルテ_ヘルズローリング.png"
	images[imageTypeForteHellsRolling] = make([]int, 5)
	if res := dxlib.LoadDivGraph(fname, 5, 5, 1, 96, 123, images[imageTypeForteHellsRolling]); res == -1 {
		return errors.Newf("failed to load image: %s", fname)
	}

	fname = config.ImagePath + "battle/skill/フォルテ_ダークアームブレード.png"
	images[imageTypeForteDarkArmBlade] = make([]int, 8)
	if res := dxlib.LoadDivGraph(fname, 8, 4, 2, 188, 150, images[imageTypeForteDarkArmBlade]); res == -1 {
		return errors.Newf("failed to load image: %s", fname)
	}

	fname = config.ImagePath + "battle/skill/シューティングバスター.png"
	images[imageTypeForteShootingBuster] = make([]int, 5)
	if res := dxlib.LoadDivGraph(fname, 5, 5, 1, 112, 96, images[imageTypeForteShootingBuster]); res == -1 {
		return errors.Newf("failed to load image: %s", fname)
	}

	fname = config.ImagePath + "battle/skill/フォルテ_ダークネスオーバーロード.png"
	images[imageTypeForteDarknessOverload] = make([]int, 8)
	if res := dxlib.LoadDivGraph(fname, 8, 4, 2, 240, 220, images[imageTypeForteDarknessOverload]); res == -1 {
		return errors.Newf("failed to load image: %s", fname)
	}

	return nil
}

func ClearImages() {
	for _, imgs := range images {
		for _, img := range imgs {
			dxlib.DeleteGraph(img)
		}
	}
}

func drawHitArea(panelPos point.Point) {
	view := battlecommon.ViewPos(panelPos)
	x1 := view.X - battlecommon.PanelSize.X/2
	y1 := view.Y
	x2 := view.X + battlecommon.PanelSize.X/2
	y2 := view.Y + battlecommon.PanelSize.Y
	const s = 5
	dxlib.DrawBox(x1+s, y1+s, x2-s, y2-s, 0xffff00, true)
}
