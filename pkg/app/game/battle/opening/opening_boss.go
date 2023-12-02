package opening

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/enemy"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/math"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	horizontal int = iota
	vertical
)

type Boss struct {
	enemyImages []int
	enemies     []enemy.EnemyParam
	playerImage int
	count       int
}

func NewWithBoss(enemyList []enemy.EnemyParam) (*Boss, error) {
	res := &Boss{}

	res.enemies = enemyList
	res.count = 0

	for _, e := range enemyList {
		name, ext := enemy.GetStandImageFile(e.CharID)
		fname := name + ext
		res.enemyImages = append(res.enemyImages, dxlib.LoadGraph(fname))
	}

	res.playerImage = dxlib.LoadGraph(common.ImagePath + "battle/character/ロックマン_player_side.png")
	if res.playerImage == -1 {
		return nil, fmt.Errorf("failed to load player image")
	}

	return res, nil
}

func (b *Boss) End() {
	for _, img := range b.enemyImages {
		dxlib.DeleteGraph(img)
	}
	b.enemyImages = []int{}
}

func (b *Boss) Process() bool {
	if config.Get().Debug.SkipBattleOpening {
		return true
	}

	b.count++

	return b.count > 70
}

func (b *Boss) Draw() {
	dxlib.DrawBox(0, 0, common.ScreenSize.X, common.ScreenSize.Y, 0x000000, true)

	// debug(初期位置)
	view := battlecommon.ViewPos(point.Point{X: 1, Y: 1})

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
	for i := 0; i < battlecommon.FieldNum.Y+1; i++ {
		y := battlecommon.DrawPanelTopY + i*battlecommon.PanelSize.Y
		len := (b.count - i*10) * 40
		if len > common.ScreenSize.X {
			len = common.ScreenSize.X
		}

		drawLine(0, y, len, horizontal, color)
	}

	// vertical lines
	for i := 0; i < battlecommon.FieldNum.X-1; i++ {
		x := (i + 1) * battlecommon.PanelSize.X
		len := (b.count - 40) * 40
		s := 0
		delay := 45 + math.MountainIndex(i, battlecommon.FieldNum.X-1)*5
		if b.count >= delay {
			s = (b.count - delay) * 20
			if s > battlecommon.DrawPanelTopY {
				s = battlecommon.DrawPanelTopY
			}
		}

		maxLen := common.ScreenSize.Y - battlecommon.DrawPanelTopY
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
