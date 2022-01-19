package opening

import (
	"fmt"

	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/enemy"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
)

const (
	horizontal int = iota
	vertical
)

type boss struct {
	enemyImages []int32
	enemies     []enemy.EnemyParam
	playerImage int32
	count       int
}

func (b *boss) Init(enemyList []enemy.EnemyParam) error {
	b.enemies = enemyList
	b.count = 0

	for _, e := range enemyList {
		name, ext := enemy.GetStandImageFile(e.CharID)
		fname := name + ext
		b.enemyImages = append(b.enemyImages, dxlib.LoadGraph(fname))
	}

	b.playerImage = dxlib.LoadGraph(common.ImagePath + "battle/character/ロックマン.png")
	if b.playerImage == -1 {
		return fmt.Errorf("failed to load player image")
	}

	return nil
}

func (b *boss) End() {
	for _, img := range b.enemyImages {
		dxlib.DeleteGraph(img)
	}
	b.enemyImages = []int32{}
}

func (b *boss) Process() bool {
	b.count++

	return b.count > 70
}

func (b *boss) Draw() {
	dxlib.DrawBox(0, 0, common.ScreenX, common.ScreenY, 0x000000, dxlib.TRUE)

	// debug(初期位置)
	x, y := battlecommon.ViewPos(1, 1)

	dxlib.SetDrawBright(17, 168, 10)
	dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_INVSRC, 255)
	dxlib.DrawRotaGraph(x, y, 1, 0, b.playerImage, dxlib.TRUE)

	for i, e := range b.enemies {
		x, y := battlecommon.ViewPos(e.PosX, e.PosY)
		dxlib.DrawRotaGraph(x, y, 1, 0, b.enemyImages[i], dxlib.TRUE)
	}

	dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_ADD, 255)
	dxlib.DrawRotaGraph(x, y, 1, 0, b.playerImage, dxlib.TRUE)

	for i, e := range b.enemies {
		x, y := battlecommon.ViewPos(e.PosX, e.PosY)
		dxlib.DrawRotaGraph(x, y, 1, 0, b.enemyImages[i], dxlib.TRUE)
	}

	dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_NOBLEND, 255)
	dxlib.SetDrawBright(255, 255, 255)

	color := dxlib.GetColor(17, 168, 10)

	// horizontal lines
	for i := 0; i < field.FieldNumY+1; i++ {
		y := field.DrawPanelTopY + i*field.PanelSizeY
		len := (b.count - i*10) * 40
		if len > common.ScreenX {
			len = common.ScreenX
		}

		drawLine(0, int32(y), int32(len), horizontal, color)
	}

	// vertical lines
	for i := 0; i < field.FieldNumX-1; i++ {
		x := (i + 1) * field.PanelSizeX
		len := (b.count - 40) * 40
		s := 0
		delay := 45 + common.MountainIndex(i, field.FieldNumX-1)*5
		if b.count >= delay {
			s = (b.count - delay) * 20
			if s > field.DrawPanelTopY {
				s = field.DrawPanelTopY
			}
		}

		maxLen := common.ScreenY - field.DrawPanelTopY
		if len > maxLen {
			len = maxLen
		}

		drawLine(int32(x), int32(s), int32(len), vertical, color)
	}
}

func drawLine(x, y int32, length int32, direct int, color uint32) {
	if length <= 0 {
		return
	}

	const s = 1

	switch direct {
	case horizontal:
		dxlib.DrawBox(x, y-s, x+length, y+s, color, dxlib.TRUE)
	case vertical:
		dxlib.DrawBox(x-s, y, x+s, y+length, color, dxlib.TRUE)
	}
}
