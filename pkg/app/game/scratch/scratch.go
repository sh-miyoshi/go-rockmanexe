package scratch

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/manager"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

var (
	imgPlayers [battlecommon.PlayerActMax][]int
	imgMetall  int
	imgAquaman int
	animMgr    *manager.Manager
)

func Init() {
	if animMgr == nil {
		animMgr = manager.NewManager()
	}
	field.Init(animMgr)

	fname := config.ImagePath + "battle/character/player_move.png"
	imgPlayers[battlecommon.PlayerActMove] = make([]int, 4)
	dxlib.LoadDivGraph(fname, 4, 4, 1, 100, 100, imgPlayers[battlecommon.PlayerActMove])

	fname = config.ImagePath + "battle/character/player_damaged.png"
	imgPlayers[battlecommon.PlayerActDamage] = make([]int, 6)
	dxlib.LoadDivGraph(fname, 6, 6, 1, 100, 100, imgPlayers[battlecommon.PlayerActDamage])
	// 1 -> 2,3  2-4 3-5
	imgPlayers[battlecommon.PlayerActDamage][4] = imgPlayers[battlecommon.PlayerActDamage][2]
	imgPlayers[battlecommon.PlayerActDamage][5] = imgPlayers[battlecommon.PlayerActDamage][3]
	imgPlayers[battlecommon.PlayerActDamage][2] = imgPlayers[battlecommon.PlayerActDamage][1]
	imgPlayers[battlecommon.PlayerActDamage][3] = imgPlayers[battlecommon.PlayerActDamage][1]

	fname = config.ImagePath + "battle/character/player_shot.png"
	imgPlayers[battlecommon.PlayerActShot] = make([]int, 6)
	dxlib.LoadDivGraph(fname, 6, 6, 1, 180, 100, imgPlayers[battlecommon.PlayerActShot])

	fname = config.ImagePath + "battle/character/player_cannon.png"
	imgPlayers[battlecommon.PlayerActCannon] = make([]int, 6)
	dxlib.LoadDivGraph(fname, 6, 6, 1, 100, 100, imgPlayers[battlecommon.PlayerActCannon])

	fname = config.ImagePath + "battle/character/player_sword.png"
	imgPlayers[battlecommon.PlayerActSword] = make([]int, 7)
	dxlib.LoadDivGraph(fname, 7, 7, 1, 128, 128, imgPlayers[battlecommon.PlayerActSword])

	fname = config.ImagePath + "battle/character/player_bomb.png"
	imgPlayers[battlecommon.PlayerActBomb] = make([]int, 7)
	dxlib.LoadDivGraph(fname, 5, 5, 1, 100, 114, imgPlayers[battlecommon.PlayerActBomb])
	imgPlayers[battlecommon.PlayerActBomb][5] = imgPlayers[battlecommon.PlayerActBomb][4]
	imgPlayers[battlecommon.PlayerActBomb][6] = imgPlayers[battlecommon.PlayerActBomb][4]

	fname = config.ImagePath + "battle/character/player_buster.png"
	imgPlayers[battlecommon.PlayerActBuster] = make([]int, 6)
	dxlib.LoadDivGraph(fname, 6, 6, 1, 180, 100, imgPlayers[battlecommon.PlayerActBuster])

	fname = config.ImagePath + "battle/character/player_pick.png"
	imgPlayers[battlecommon.PlayerActPick] = make([]int, 6)
	dxlib.LoadDivGraph(fname, 4, 4, 1, 96, 124, imgPlayers[battlecommon.PlayerActPick])
	imgPlayers[battlecommon.PlayerActPick][4] = imgPlayers[battlecommon.PlayerActPick][3]
	imgPlayers[battlecommon.PlayerActPick][5] = imgPlayers[battlecommon.PlayerActPick][3]

	fname = config.ImagePath + "battle/character/player_throw.png"
	imgPlayers[battlecommon.PlayerActThrow] = make([]int, 4)
	dxlib.LoadDivGraph(fname, 4, 4, 1, 97, 115, imgPlayers[battlecommon.PlayerActThrow])

	fname = config.ImagePath + "battle/character/メットール_move.png"
	imgMetall = dxlib.LoadGraph(fname)

	fname = config.ImagePath + "battle/character/アクアマン.png"
	imgAquaman = dxlib.LoadGraph(fname)
}

func Draw() {
	field.Draw()
	animMgr.Draw()

	view := battlecommon.ViewPos(point.Point{X: 1, Y: 1})
	dxlib.DrawRotaGraph(view.X, view.Y, 1, 0, imgPlayers[battlecommon.PlayerActMove][0], true)

	view = battlecommon.ViewPos(point.Point{X: 4, Y: 1})
	dxlib.DrawRotaGraph(view.X, view.Y, 1, 0, imgMetall, true)
}

func Update() {
	field.Update()
	animMgr.Update(true)
}
