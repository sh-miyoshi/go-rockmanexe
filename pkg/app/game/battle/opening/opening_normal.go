package opening

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/enemy"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

const (
	normalViewDelay = 8
	normalViewCount = 256 / normalViewDelay
)

type Normal struct {
	count     int
	showCount int
	enemies   []enemy.EnemyParam
	images    map[int]int
}

func NewWithNormal(enemyList []enemy.EnemyParam) (*Normal, error) {
	res := &Normal{}

	res.count = 0
	res.showCount = 0
	res.enemies = enemyList
	res.images = make(map[int]int)

	for _, e := range enemyList {
		_, ok := res.images[e.CharID]
		if !ok {
			name, ext := enemy.GetStandImageFile(e.CharID)
			fname := name + ext
			res.images[e.CharID] = dxlib.LoadGraph(fname)
		}
	}

	return res, nil
}
func (n *Normal) End() {
	for _, img := range n.images {
		dxlib.DeleteGraph(img)
	}
	n.images = make(map[int]int)
}
func (n *Normal) Update() bool {
	if config.Get().Debug.SkipBattleOpening {
		return true
	}

	if n.count == 0 {
		sound.On(resources.SEEnemyAppear)
	}

	n.count++
	if n.count > normalViewCount {
		n.count = 0
		n.showCount++
		if n.showCount >= len(n.enemies) {
			return true
		}
	}
	return false
}
func (n *Normal) Draw() {
	// Show animationed enemies
	for i := 0; i < n.showCount; i++ {
		view := battlecommon.ViewPos(n.enemies[i].Pos)
		dxlib.DrawRotaGraph(view.X, view.Y, 1, 0, n.images[n.enemies[i].CharID], true)
	}

	// Show current enemy
	if n.showCount < len(n.enemies) {
		pm := n.count * normalViewDelay
		if pm >= 256 {
			pm = 255
		}

		dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_ALPHA, pm)
		view := battlecommon.ViewPos(n.enemies[n.showCount].Pos)
		dxlib.DrawRotaGraph(view.X, view.Y, 1, 0, n.images[n.enemies[n.showCount].CharID], true)
		dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_NOBLEND, 0)
	}
}
