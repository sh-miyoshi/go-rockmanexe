package drawer

import (
	"github.com/cockroachdb/errors"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type PlayerDrawer struct {
	images [battlecommon.PlayerActMax][]int
}

func (p *PlayerDrawer) Init() error {
	fname := config.ImagePath + "battle/character/player_move.png"
	p.images[battlecommon.PlayerActMove] = make([]int, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 100, 100, p.images[battlecommon.PlayerActMove]); res == -1 {
		return errors.Newf("failed to load player move image: %s", fname)
	}
	fname = config.ImagePath + "battle/character/player_damaged.png"
	p.images[battlecommon.PlayerActDamage] = make([]int, 6)
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 100, 100, p.images[battlecommon.PlayerActDamage]); res == -1 {
		return errors.Newf("failed to load player damage image: %s", fname)
	}
	// 1 -> 2,3  2-4 3-5
	p.images[battlecommon.PlayerActDamage][4] = p.images[battlecommon.PlayerActDamage][2]
	p.images[battlecommon.PlayerActDamage][5] = p.images[battlecommon.PlayerActDamage][3]
	p.images[battlecommon.PlayerActDamage][2] = p.images[battlecommon.PlayerActDamage][1]
	p.images[battlecommon.PlayerActDamage][3] = p.images[battlecommon.PlayerActDamage][1]

	fname = config.ImagePath + "battle/character/player_shot.png"
	p.images[battlecommon.PlayerActShot] = make([]int, 6)
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 180, 100, p.images[battlecommon.PlayerActShot]); res == -1 {
		return errors.Newf("failed to load player shot image: %s", fname)
	}

	fname = config.ImagePath + "battle/character/player_cannon.png"
	p.images[battlecommon.PlayerActCannon] = make([]int, 6)
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 100, 100, p.images[battlecommon.PlayerActCannon]); res == -1 {
		return errors.Newf("failed to load player cannon image: %s", fname)
	}

	fname = config.ImagePath + "battle/character/player_sword.png"
	p.images[battlecommon.PlayerActSword] = make([]int, 7)
	if res := dxlib.LoadDivGraph(fname, 7, 7, 1, 128, 128, p.images[battlecommon.PlayerActSword]); res == -1 {
		return errors.Newf("failed to load player sword image: %s", fname)
	}

	fname = config.ImagePath + "battle/character/player_bomb.png"
	p.images[battlecommon.PlayerActBomb] = make([]int, 7)
	if res := dxlib.LoadDivGraph(fname, 5, 5, 1, 100, 114, p.images[battlecommon.PlayerActBomb]); res == -1 {
		return errors.Newf("failed to load player bomb image: %s", fname)
	}
	p.images[battlecommon.PlayerActBomb][5] = p.images[battlecommon.PlayerActBomb][4]
	p.images[battlecommon.PlayerActBomb][6] = p.images[battlecommon.PlayerActBomb][4]

	fname = config.ImagePath + "battle/character/player_buster.png"
	p.images[battlecommon.PlayerActBuster] = make([]int, 6)
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 180, 100, p.images[battlecommon.PlayerActBuster]); res == -1 {
		return errors.Newf("failed to load player buster image: %s", fname)
	}

	fname = config.ImagePath + "battle/character/player_pick.png"
	p.images[battlecommon.PlayerActPick] = make([]int, 6)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 96, 124, p.images[battlecommon.PlayerActPick]); res == -1 {
		return errors.Newf("failed to load player pick image: %s", fname)
	}
	p.images[battlecommon.PlayerActPick][4] = p.images[battlecommon.PlayerActPick][3]
	p.images[battlecommon.PlayerActPick][5] = p.images[battlecommon.PlayerActPick][3]

	fname = config.ImagePath + "battle/character/player_throw.png"
	p.images[battlecommon.PlayerActThrow] = make([]int, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 97, 115, p.images[battlecommon.PlayerActThrow]); res == -1 {
		return errors.Newf("failed to load player throw image: %s", fname)
	}

	p.images[battlecommon.PlayerActParalyzed] = make([]int, 4)
	for i := 0; i < 4; i++ {
		p.images[battlecommon.PlayerActParalyzed][i] = p.images[battlecommon.PlayerActDamage][i]
	}

	return nil
}

func (p *PlayerDrawer) End() {
	for i := 0; i < battlecommon.PlayerActMax; i++ {
		for j := 0; j < len(p.images[i]); j++ {
			dxlib.DeleteGraph(p.images[i][j])
			p.images[i][j] = -1
		}
	}
}

func (p *PlayerDrawer) PopDeleteImage() int {
	img := p.images[battlecommon.PlayerActDamage][1]
	p.images[battlecommon.PlayerActDamage][1] = -1
	return img
}

func (p *PlayerDrawer) Draw(count int, viewPos point.Point, actType int, isParalyzed bool) {
	img := p.getImage(count, actType)

	dxlib.DrawRotaGraph(viewPos.X, viewPos.Y, 1, 0, img, true)
	if isParalyzed {
		dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_ADD, 255)
		// 黄色と白を点滅させる
		pm := 0
		if count/10%2 == 0 {
			pm = 255
		}
		dxlib.SetDrawBright(255, 255, pm)
		dxlib.DrawRotaGraph(viewPos.X, viewPos.Y, 1, 0, img, true)
		dxlib.SetDrawBright(255, 255, 255)
		dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_NOBLEND, 0)
	}
}

func (p *PlayerDrawer) getImage(count int, actType int) int {
	if actType == -1 {
		// return stand image
		return p.images[battlecommon.PlayerActMove][0]
	}

	num, delay := battlecommon.GetPlayerImageInfo(actType)
	imgNo := (count / delay)
	if imgNo >= num {
		imgNo = num - 1
	}

	return p.images[actType][imgNo]
}
