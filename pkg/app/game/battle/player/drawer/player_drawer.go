package drawer

import (
	"github.com/cockroachdb/errors"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/system"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type PlayerDrawer struct {
	currentSoulUnison resources.SoulUnison
	imgNormals        [battlecommon.PlayerActMax][]int
	imgAquas          [battlecommon.PlayerActMax][]int
}

func (p *PlayerDrawer) Init() error {
	p.currentSoulUnison = resources.SoulUnisonNone

	// Load player normal images
	fname := config.ImagePath + "battle/character/player_move.png"
	p.imgNormals[battlecommon.PlayerActMove] = make([]int, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 100, 100, p.imgNormals[battlecommon.PlayerActMove]); res == -1 {
		return errors.Newf("failed to load player move image: %s", fname)
	}
	fname = config.ImagePath + "battle/character/player_damaged.png"
	p.imgNormals[battlecommon.PlayerActDamage] = make([]int, 6)
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 100, 100, p.imgNormals[battlecommon.PlayerActDamage]); res == -1 {
		return errors.Newf("failed to load player damage image: %s", fname)
	}
	// 1 -> 2,3  2-4 3-5
	p.imgNormals[battlecommon.PlayerActDamage][4] = p.imgNormals[battlecommon.PlayerActDamage][2]
	p.imgNormals[battlecommon.PlayerActDamage][5] = p.imgNormals[battlecommon.PlayerActDamage][3]
	p.imgNormals[battlecommon.PlayerActDamage][2] = p.imgNormals[battlecommon.PlayerActDamage][1]
	p.imgNormals[battlecommon.PlayerActDamage][3] = p.imgNormals[battlecommon.PlayerActDamage][1]

	fname = config.ImagePath + "battle/character/player_shot.png"
	p.imgNormals[battlecommon.PlayerActShot] = make([]int, 6)
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 180, 100, p.imgNormals[battlecommon.PlayerActShot]); res == -1 {
		return errors.Newf("failed to load player shot image: %s", fname)
	}

	fname = config.ImagePath + "battle/character/player_cannon.png"
	p.imgNormals[battlecommon.PlayerActCannon] = make([]int, 6)
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 100, 100, p.imgNormals[battlecommon.PlayerActCannon]); res == -1 {
		return errors.Newf("failed to load player cannon image: %s", fname)
	}

	fname = config.ImagePath + "battle/character/player_sword.png"
	p.imgNormals[battlecommon.PlayerActSword] = make([]int, 7)
	if res := dxlib.LoadDivGraph(fname, 7, 7, 1, 128, 128, p.imgNormals[battlecommon.PlayerActSword]); res == -1 {
		return errors.Newf("failed to load player sword image: %s", fname)
	}

	fname = config.ImagePath + "battle/character/player_bomb.png"
	p.imgNormals[battlecommon.PlayerActBomb] = make([]int, 7)
	if res := dxlib.LoadDivGraph(fname, 5, 5, 1, 100, 114, p.imgNormals[battlecommon.PlayerActBomb]); res == -1 {
		return errors.Newf("failed to load player bomb image: %s", fname)
	}
	p.imgNormals[battlecommon.PlayerActBomb][5] = p.imgNormals[battlecommon.PlayerActBomb][4]
	p.imgNormals[battlecommon.PlayerActBomb][6] = p.imgNormals[battlecommon.PlayerActBomb][4]

	fname = config.ImagePath + "battle/character/player_buster.png"
	p.imgNormals[battlecommon.PlayerActBuster] = make([]int, 6)
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 180, 100, p.imgNormals[battlecommon.PlayerActBuster]); res == -1 {
		return errors.Newf("failed to load player buster image: %s", fname)
	}

	fname = config.ImagePath + "battle/character/player_pick.png"
	p.imgNormals[battlecommon.PlayerActPick] = make([]int, 6)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 96, 124, p.imgNormals[battlecommon.PlayerActPick]); res == -1 {
		return errors.Newf("failed to load player pick image: %s", fname)
	}
	p.imgNormals[battlecommon.PlayerActPick][4] = p.imgNormals[battlecommon.PlayerActPick][3]
	p.imgNormals[battlecommon.PlayerActPick][5] = p.imgNormals[battlecommon.PlayerActPick][3]

	fname = config.ImagePath + "battle/character/player_throw.png"
	p.imgNormals[battlecommon.PlayerActThrow] = make([]int, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 97, 115, p.imgNormals[battlecommon.PlayerActThrow]); res == -1 {
		return errors.Newf("failed to load player throw image: %s", fname)
	}

	p.imgNormals[battlecommon.PlayerActParalyzed] = make([]int, 4)
	for i := 0; i < 4; i++ {
		p.imgNormals[battlecommon.PlayerActParalyzed][i] = p.imgNormals[battlecommon.PlayerActDamage][i]
	}

	// Load player aqua images
	fname = config.ImagePath + "battle/character/player_aqua_move.png"
	p.imgAquas[battlecommon.PlayerActMove] = make([]int, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 100, 100, p.imgAquas[battlecommon.PlayerActMove]); res == -1 {
		return errors.Newf("failed to load player aqua soul move image: %s", fname)
	}

	// WIP
	// Load player aqua damaged images
	// fname = config.ImagePath + "battle/character/player_aqua_damaged.png"
	p.imgAquas[battlecommon.PlayerActDamage] = make([]int, 6)
	// if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 100, 100, p.imgAquas[battlecommon.PlayerActDamage]); res == -1 {
	// 	return errors.Newf("failed to load player aqua damaged image: %s", fname)
	// }
	// p.imgAquas[battlecommon.PlayerActDamage][4] = p.imgAquas[battlecommon.PlayerActDamage][2]
	// p.imgAquas[battlecommon.PlayerActDamage][5] = p.imgAquas[battlecommon.PlayerActDamage][3]
	// p.imgAquas[battlecommon.PlayerActDamage][2] = p.imgAquas[battlecommon.PlayerActDamage][1]
	// p.imgAquas[battlecommon.PlayerActDamage][3] = p.imgAquas[battlecommon.PlayerActDamage][1]

	// Load player aqua shot images
	fname = config.ImagePath + "battle/character/player_aqua_shot.png"
	p.imgAquas[battlecommon.PlayerActShot] = make([]int, 6)
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 180, 100, p.imgAquas[battlecommon.PlayerActShot]); res == -1 {
		return errors.Newf("failed to load player aqua shot image: %s", fname)
	}

	// WIP
	// Load player aqua cannon images
	// fname = config.ImagePath + "battle/character/player_aqua_cannon.png"
	// p.imgAquas[battlecommon.PlayerActCannon] = make([]int, 6)
	// if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 100, 100, p.imgAquas[battlecommon.PlayerActCannon]); res == -1 {
	// 	return errors.Newf("failed to load player aqua cannon image: %s", fname)
	// }

	// Load player aqua sword images
	fname = config.ImagePath + "battle/character/player_aqua_sword.png"
	p.imgAquas[battlecommon.PlayerActSword] = make([]int, 7)
	if res := dxlib.LoadDivGraph(fname, 7, 7, 1, 128, 128, p.imgAquas[battlecommon.PlayerActSword]); res == -1 {
		return errors.Newf("failed to load player aqua sword image: %s", fname)
	}

	// Load player aqua bomb images
	fname = config.ImagePath + "battle/character/player_aqua_bomb.png"
	p.imgAquas[battlecommon.PlayerActBomb] = make([]int, 7)
	if res := dxlib.LoadDivGraph(fname, 5, 5, 1, 100, 114, p.imgAquas[battlecommon.PlayerActBomb]); res == -1 {
		return errors.Newf("failed to load player aqua bomb image: %s", fname)
	}
	p.imgAquas[battlecommon.PlayerActBomb][5] = p.imgAquas[battlecommon.PlayerActBomb][4]
	p.imgAquas[battlecommon.PlayerActBomb][6] = p.imgAquas[battlecommon.PlayerActBomb][4]

	// Load player aqua buster images
	fname = config.ImagePath + "battle/character/player_aqua_buster.png"
	p.imgAquas[battlecommon.PlayerActBuster] = make([]int, 6)
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 180, 100, p.imgAquas[battlecommon.PlayerActBuster]); res == -1 {
		return errors.Newf("failed to load player aqua buster image: %s", fname)
	}

	// Load player aqua pick images
	p.imgAquas[battlecommon.PlayerActPick] = make([]int, 6)
	for i := 0; i < 6; i++ {
		p.imgAquas[battlecommon.PlayerActPick][i] = p.imgAquas[battlecommon.PlayerActMove][0]
	}

	// Load player aqua throw images
	fname = config.ImagePath + "battle/character/player_aqua_throw.png"
	p.imgAquas[battlecommon.PlayerActThrow] = make([]int, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 97, 115, p.imgAquas[battlecommon.PlayerActThrow]); res == -1 {
		return errors.Newf("failed to load player aqua throw image: %s", fname)
	}

	// Set aqua paralyzed images same as damaged images
	p.imgAquas[battlecommon.PlayerActParalyzed] = make([]int, 4)
	for i := 0; i < 4; i++ {
		p.imgAquas[battlecommon.PlayerActParalyzed][i] = p.imgAquas[battlecommon.PlayerActDamage][i]
	}

	return nil
}

func (p *PlayerDrawer) End() {
	for i := 0; i < battlecommon.PlayerActMax; i++ {
		for j := 0; j < len(p.imgNormals[i]); j++ {
			dxlib.DeleteGraph(p.imgNormals[i][j])
			p.imgNormals[i][j] = -1
		}
		for j := 0; j < len(p.imgAquas[i]); j++ {
			dxlib.DeleteGraph(p.imgAquas[i][j])
			p.imgAquas[i][j] = -1
		}
	}
}

func (p *PlayerDrawer) PopDeleteImage() int {
	img := p.images()[battlecommon.PlayerActDamage][1]
	p.images()[battlecommon.PlayerActDamage][1] = -1
	return img
}

func (p *PlayerDrawer) SetSoulUnison(soulUnison resources.SoulUnison) {
	p.currentSoulUnison = soulUnison
}

func (p *PlayerDrawer) Draw(count int, viewPos point.Point, actType int, isParalyzed bool) {
	img := p.getImage(count, actType)
	ofs := p.getOffset(actType)

	dxlib.DrawRotaGraph(viewPos.X+ofs.X, viewPos.Y+ofs.Y, 1, 0, img, true)
	if isParalyzed {
		dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_ADD, 255)
		// 黄色と白を点滅させる
		pm := 0
		if count/10%2 == 0 {
			pm = 255
		}
		dxlib.SetDrawBright(255, 255, pm)
		dxlib.DrawRotaGraph(viewPos.X+ofs.X, viewPos.Y+ofs.Y, 1, 0, img, true)
		dxlib.SetDrawBright(255, 255, 255)
		dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_NOBLEND, 0)
	}
}

func (p *PlayerDrawer) getImage(count int, actType int) int {
	if actType == -1 {
		// return stand image
		return p.images()[battlecommon.PlayerActMove][0]
	}

	num, delay := battlecommon.GetPlayerImageInfo(actType)
	imgNo := (count / delay)
	if imgNo >= num {
		imgNo = num - 1
	}

	return p.images()[actType][imgNo]
}

func (p *PlayerDrawer) images() [battlecommon.PlayerActMax][]int {
	switch p.currentSoulUnison {
	case resources.SoulUnisonNone:
		return p.imgNormals
	case resources.SoulUnisonAqua:
		return p.imgAquas
	}
	system.SetError("Invalid soul unison type")
	return p.imgNormals
}

func (p *PlayerDrawer) getOffset(actType int) point.Point {
	switch p.currentSoulUnison {
	case resources.SoulUnisonNone:
		switch actType {
		case battlecommon.PlayerActSword:
			return point.Point{X: 20, Y: -3}
		}
	case resources.SoulUnisonAqua:
		switch actType {
		case battlecommon.PlayerActBuster, battlecommon.PlayerActShot:
			return point.Point{X: 45, Y: 5}
		case battlecommon.PlayerActSword:
			return point.Point{X: 25, Y: -10}
		}
	}
	return point.Point{X: 0, Y: 0}
}
