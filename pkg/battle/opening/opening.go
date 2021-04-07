package opening

import (
	"github.com/sh-miyoshi/dxlib"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/enemy"
)

const (
	viewDelay = 8
	viewCount = 256 / viewDelay
)

var (
	count     int
	showCount int
	enemies   []enemy.EnemyParam
	images    = make(map[int]int32)
)

func Init(enemyList []enemy.EnemyParam) error {
	count = 0
	showCount = 0
	enemies = enemyList

	for _, e := range enemyList {
		_, ok := images[e.CharID]
		if !ok {
			name, ext := enemy.GetStandImageFile(e.CharID)
			fname := name + ext
			images[e.CharID] = dxlib.LoadGraph(fname)
		}
	}

	return nil
}

func End() {
	for _, img := range images {
		dxlib.DeleteGraph(img)
	}
	images = make(map[int]int32)
}

func Process() bool {
	count++
	if count > viewCount {
		count = 0
		showCount++
		if showCount >= len(enemies) {
			return true
		}
	}
	return false
	// return true // debug
}

func Draw() {
	// Show animationed enemies
	for i := 0; i < showCount; i++ {
		x, y := battlecommon.ViewPos(enemies[i].PosX, enemies[i].PosY)
		dxlib.DrawRotaGraph(x, y, 1, 0, images[enemies[i].CharID], dxlib.TRUE)
	}

	// Show current enemy
	if showCount < len(enemies) {
		pm := int32(count * viewDelay)
		if pm >= 256 {
			pm = 255
		}

		dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_ALPHA, pm)
		x, y := battlecommon.ViewPos(enemies[showCount].PosX, enemies[showCount].PosY)
		dxlib.DrawRotaGraph(x, y, 1, 0, images[enemies[showCount].CharID], dxlib.TRUE)
		dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_NOBLEND, 0)
	}
}
