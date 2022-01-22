package opening

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/enemy"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

const (
	horizontal int = iota
	vertical
)

type boss struct {
	enemyImages []int
	enemies     []enemy.EnemyParam
	playerImage int
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
	b.enemyImages = []int{}
}

func (b *boss) Process() bool {
	b.count++

	return b.count > 70
}

func (b *boss) Draw() {
	dxlib.DrawBox(0, 0, common.ScreenSize.X, common.ScreenSize.Y, 0x000000, true)

	// debug(初期位置)
	view := battlecommon.ViewPos(common.Point{X: 1, Y: 1})

	dxlib.SetDrawBright(17, 168, 10)
	dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_INVSRC, 255)
	dxlib.DrawRotaGraph(view.X, view.Y, 1, 0, b.playerImage, true)

	for i, e := range b.enemies {
		view := battlecommon.ViewPos(e.Pos)
		dxlib.DrawRotaGraph(view.X, view.Y, 1, 0, b.enemyImages[i], true)
	}

	dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_ADD, 255)
	dxlib.DrawRotaGraph(view.X, view.Y, 1, 0, b.playerImage, true)

	for i, e := range b.enemies {
		view := battlecommon.ViewPos(e.Pos)
		dxlib.DrawRotaGraph(view.X, view.Y, 1, 0, b.enemyImages[i], true)
	}

	dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_NOBLEND, 255)
	dxlib.SetDrawBright(255, 255, 255)

	color := dxlib.GetColor(17, 168, 10)

	// horizontal lines
	for i := 0; i < field.FieldNum.Y+1; i++ {
		y := field.DrawPanelTopY + i*field.PanelSize.Y
		len := b.count - i*10*40
		if len > common.ScreenSize.X {
			len = common.ScreenSize.X
		}

		drawLine(0, y, len, horizontal, color)
	}

	// vertical lines
	for i := 0; i < field.FieldNum.X-1; i++ {
		x := (i + 1) * field.PanelSize.X
		len := b.count - 40*40
		s := 0
		delay := 45 + common.MountainIndex(i, field.FieldNum.X-1)*5
		if b.count >= delay {
			s = b.count - delay*20
			if s > field.DrawPanelTopY {
				s = field.DrawPanelTopY
			}
		}

		maxLen := common.ScreenSize.Y - field.DrawPanelTopY
		if len > maxLen {
			len = maxLen
		}

		drawLine(x, s, len, vertical, color)
	}
}

func drawLine(x, y int, length int, direct int, color uint) {
	if length <= 0 {
		return
	}

	const s = 1

	switch direct {
	case horizontal:
		dxlib.DrawBox(x, y-s, x+length, y+s, color, true)
	case vertical:
		dxlib.DrawBox(x-s, y, x+s, y+length, color, true)
	}
}
