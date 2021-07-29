package draw

import (
	"fmt"

	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	appdraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/object"
)

type Option struct {
	Reverse bool
	ViewHP  int
}

var (
	imgObjs        [object.TypeMax][]int32
	imgEffs        [effect.TypeMax][]int32
	imgUnknownIcon int32
)

func Init() error {
	if err := loadObjs(); err != nil {
		return fmt.Errorf("load objects failed: %w", err)
	}

	if err := loadEffects(); err != nil {
		return fmt.Errorf("load effects failed: %w", err)
	}

	fname := common.ImagePath + "chipInfo/unknown_icon.png"
	if imgUnknownIcon = dxlib.LoadGraph(fname); imgUnknownIcon == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	return nil
}

func End() {
	for _, image := range imgObjs {
		for _, img := range image {
			dxlib.DeleteGraph(img)
		}
	}
	dxlib.DeleteGraph(imgUnknownIcon)
}

func Object(obj object.Object, opt Option) {
	vx, vy := battlecommon.ViewPos(obj.X, obj.Y)
	imgNo := obj.Count / object.ImageDelays[obj.Type]
	dxopts := dxlib.DrawRotaGraphOption{}

	if opt.Reverse {
		flag := int32(dxlib.TRUE)
		dxopts.ReverseXFlag = &flag
		obj.ViewOfsX *= -1
	}

	vx += obj.ViewOfsX
	vy += obj.ViewOfsY

	// Special object draw
	switch obj.Type {
	case object.TypeVulcan:
		objectVulcan(vx, vy, imgNo, dxopts)
	case object.TypeWideShotMove:
		objectWideShotMove(vx, vy, obj, dxopts)
	case object.TypeThunderBall:
		objectThunderBall(vx, vy, obj, dxopts)
	case object.TypeMiniBomb:
		objectMiniBomb(vx, vy, obj, dxopts)
	default:
		if obj.Invincible {
			if cnt := obj.Count / 5 % 2; cnt == 0 {
				return
			}
		}

		if imgNo >= len(imgObjs[obj.Type]) {
			imgNo = len(imgObjs[obj.Type]) - 1
		}
		dxlib.DrawRotaGraph(vx, vy, 1, 0, imgObjs[obj.Type][imgNo], dxlib.TRUE, dxopts)
	}

	// Show HP
	if opt.ViewHP > 0 {
		appdraw.Number(vx, vy+40, int32(opt.ViewHP), appdraw.NumberOption{
			Color:    appdraw.NumberColorWhiteSmall,
			Centered: true,
		})
	}

	if len(obj.Chips) > 0 {
		x := field.PanelSizeX*obj.X + field.PanelSizeX/2 - 18
		y := field.DrawPanelTopY + field.PanelSizeY*obj.Y - 83
		dxlib.DrawBox(int32(x-1), int32(y-1), int32(x+29), int32(y+29), 0x000000, dxlib.FALSE)
		dxlib.DrawGraph(int32(x), int32(y), imgUnknownIcon, dxlib.TRUE)
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
	return len(imgObjs[objType]), object.ImageDelays[objType]
}

func GetEffectImageInfo(effType int) (imageNum, delay int) {
	return len(imgEffs[effType]), effect.Delays[effType]
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
	tmp := make([]int32, 24)
	if res := dxlib.LoadDivGraph(fname, 24, 8, 3, 120, 140, tmp); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	imgObjs[object.TypeNormalCannonAtk] = make([]int32, 8)
	imgObjs[object.TypeHighCannonAtk] = make([]int32, 8)
	imgObjs[object.TypeMegaCannonAtk] = make([]int32, 8)
	for i := 0; i < 8; i++ {
		imgObjs[object.TypeNormalCannonAtk][i] = tmp[i]
		imgObjs[object.TypeHighCannonAtk][i] = tmp[i+8]
		imgObjs[object.TypeMegaCannonAtk][i] = tmp[i+16]
	}

	fname = skillPath + "キャノン_body.png"
	tmp = make([]int32, 15)
	if res := dxlib.LoadDivGraph(fname, 15, 5, 3, 46, 40, tmp); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	imgObjs[object.TypeNormalCannonBody] = make([]int32, 5)
	imgObjs[object.TypeHighCannonBody] = make([]int32, 5)
	imgObjs[object.TypeMegaCannonBody] = make([]int32, 5)
	for i := 0; i < 5; i++ {
		imgObjs[object.TypeNormalCannonBody][i] = tmp[i]
		imgObjs[object.TypeHighCannonBody][i] = tmp[i+5]
		imgObjs[object.TypeMegaCannonBody][i] = tmp[i+10]
	}

	fname = skillPath + "ミニボム.png"
	imgObjs[object.TypeMiniBomb] = make([]int32, 5)
	if res := dxlib.LoadDivGraph(fname, 5, 5, 1, 40, 30, imgObjs[object.TypeMiniBomb]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = skillPath + "ソード.png"
	tmp = make([]int32, 12)
	if res := dxlib.LoadDivGraph(fname, 12, 4, 3, 160, 150, tmp); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	imgObjs[object.TypeSword] = make([]int32, 4)
	imgObjs[object.TypeWideSword] = make([]int32, 4)
	imgObjs[object.TypeLongSword] = make([]int32, 4)
	for i := 0; i < 4; i++ {
		imgObjs[object.TypeSword][i] = tmp[i]
		imgObjs[object.TypeWideSword][i] = tmp[i+8]
		imgObjs[object.TypeLongSword][i] = tmp[i+4]
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

	fname = skillPath + "ワイドショット_body.png"
	imgObjs[object.TypeWideShotBody] = make([]int32, 3)
	if res := dxlib.LoadDivGraph(fname, 3, 3, 1, 56, 66, imgObjs[object.TypeWideShotBody]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = skillPath + "ワイドショット_begin.png"
	imgObjs[object.TypeWideShotBegin] = make([]int32, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 90, 147, imgObjs[object.TypeWideShotBegin]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = skillPath + "ワイドショット_move.png"
	imgObjs[object.TypeWideShotMove] = make([]int32, 3)
	if res := dxlib.LoadDivGraph(fname, 3, 3, 1, 90, 148, imgObjs[object.TypeWideShotMove]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

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

func objectVulcan(vx, vy int32, imgNo int, dxopts dxlib.DrawRotaGraphOption) {
	if imgNo > 2 {
		imgNo /= 5 // slow down animation
	}

	ofsBody := int32(50)
	ofsAtk := int32(100)
	if dxopts.ReverseXFlag != nil && *dxopts.ReverseXFlag == dxlib.TRUE {
		ofsBody *= -1
		ofsAtk *= -1
	}

	// Show body
	no := imgNo
	if no > 2 {
		no = no % 2
	}

	dxlib.DrawRotaGraph(vx+ofsBody, vy-18, 1, 0, imgObjs[object.TypeVulcan][no], dxlib.TRUE)
	// Show attack
	if imgNo != 0 {
		if imgNo%2 == 0 {
			dxlib.DrawRotaGraph(vx+ofsAtk, vy-10, 1, 0, imgObjs[object.TypeVulcan][3], dxlib.TRUE, dxopts)
		} else {
			dxlib.DrawRotaGraph(vx+ofsAtk, vy-15, 1, 0, imgObjs[object.TypeVulcan][3], dxlib.TRUE, dxopts)
		}
	}
}

func objectWideShotMove(vx, vy int32, obj object.Object, dxopts dxlib.DrawRotaGraphOption) {
	if obj.Speed == 0 {
		panic("ワイドショット描画のためのSpeedが0です")
	}

	imgNo := (obj.Count / object.ImageDelays[obj.Type]) % len(imgObjs[obj.Type])
	ofsx := int32(field.PanelSizeX * obj.Count / obj.Speed)
	if dxopts.ReverseXFlag != nil && *dxopts.ReverseXFlag == dxlib.TRUE {
		ofsx *= -1
	}

	dxlib.DrawRotaGraph(vx+ofsx, vy, 1, 0, imgObjs[obj.Type][imgNo], dxlib.TRUE, dxopts)
}

func objectThunderBall(vx, vy int32, obj object.Object, dxopts dxlib.DrawRotaGraphOption) {
	imgNo := (obj.Count / object.ImageDelays[obj.Type]) % len(imgObjs[obj.Type])
	if obj.Count < obj.Speed {
		ofsx := field.PanelSizeX * (obj.TargetX - obj.X) * obj.Count / obj.Speed
		ofsy := field.PanelSizeY * (obj.TargetY - obj.Y) * obj.Count / obj.Speed
		dxlib.DrawRotaGraph(vx+int32(ofsx), vy+25+int32(ofsy), 1, 0, imgObjs[obj.Type][imgNo], dxlib.TRUE)
	}
}

func objectMiniBomb(vx, vy int32, obj object.Object, dxopts dxlib.DrawRotaGraphOption) {
	imgNo := (obj.Count / object.ImageDelays[obj.Type]) % len(imgObjs[obj.Type])

	// y = ax^2 + bx +c
	// (0,0), (d/2, ymax), (d, 0)
	size := field.PanelSizeX * 3
	ofsx := size * obj.Count / obj.Speed
	ymax := 100
	ofsy := ymax*4*ofsx*ofsx/(size*size) - ymax*4*ofsx/size

	dxlib.DrawRotaGraph(vx+int32(ofsx), vy+int32(ofsy), 1, 0, imgObjs[obj.Type][imgNo], dxlib.TRUE, dxopts)
}
