package player

import (
	"fmt"
	"math/rand"
	"sort"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/chip"
	chipimage "github.com/sh-miyoshi/go-rockmanexe/pkg/app/chip/image"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/manager"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common/deleteanim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/player"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/system"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/inputs"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/math"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
	"github.com/stretchr/stew/slice"
)

const (
	NextActNone int = iota
	NextActChipSelect
	NextActLose
)

type SelectChip struct {
	ID        int
	Code      string
	PlusPower int
}

type BattlePlayerAct struct {
	MoveDirect int
	Charged    bool
	ShotPower  uint

	typ        int
	count      int
	endCount   int
	pPos       *point.Point
	skillObjID string
	skillInst  skill.SkillAnim
	animMgr    *manager.Manager
}

type BattlePlayer struct {
	ID            string
	Pos           point.Point
	HP            uint
	HPMax         uint
	ChargeCount   uint
	ChargeTime    uint
	GaugeCount    uint
	ShotPower     uint
	ChipFolder    []player.ChipInfo
	SelectedChips []SelectChip
	NextAction    int
	EnableAct     bool
	MoveNum       int
	DamageNum     int
	MindStatus    int
	IsUnderShirt  bool
	ChipSelectMax int

	act             BattlePlayerAct
	invincibleCount int
	visible         bool
	isShiftFrame    bool
	isShowGauge     bool
	barrierHP       int
	animMgr         *manager.Manager
}

var (
	imgPlayers    [battlecommon.PlayerActMax][]int
	imgHPFrame    int
	imgGaugeFrame int
	imgGaugeMax   []int
	imgCharge     [2][]int
	imgMinds      []int
	imgMindFrame  int
	imgBarrier    []int
)

func New(plyr *player.Player, animMgr *manager.Manager) (*BattlePlayer, error) {
	logger.Info("Initialize battle player data")

	res := BattlePlayer{
		ID:            uuid.New().String(),
		HP:            plyr.HP,
		HPMax:         plyr.HP, // TODO HPは引き継がない
		Pos:           point.Point{X: 1, Y: 1},
		ShotPower:     plyr.ShotPower,
		ChargeTime:    plyr.ChargeTime,
		EnableAct:     true,
		MindStatus:    battlecommon.PlayerMindStatusNormal, // TODO playerにstatusを持つ
		visible:       true,
		IsUnderShirt:  plyr.IsUnderShirt(),
		ChipSelectMax: plyr.ChipSelectMax,
		barrierHP:     0,
		animMgr:       animMgr,
	}
	res.act.Init(&res.Pos, animMgr)

	if config.Get().Debug.AlwaysInvisible {
		logger.Debug("enable inbisible mode")
		res.invincibleCount = 1
	}

	for _, c := range plyr.ChipFolder {
		res.ChipFolder = append(res.ChipFolder, c)
	}
	if !config.Get().Debug.UseDebugFolder {
		// Shuffle folder
		for i := 0; i < 10; i++ {
			for j := 0; j < len(res.ChipFolder); j++ {
				n := rand.Intn(len(res.ChipFolder))
				res.ChipFolder[j], res.ChipFolder[n] = res.ChipFolder[n], res.ChipFolder[j]
			}
		}
	}

	logger.Debug("Player info: %+v", res)

	fname := config.ImagePath + "battle/character/player_move.png"
	imgPlayers[battlecommon.PlayerActMove] = make([]int, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 100, 100, imgPlayers[battlecommon.PlayerActMove]); res == -1 {
		return nil, errors.Newf("failed to load player move image: %s", fname)
	}

	fname = config.ImagePath + "battle/character/player_damaged.png"
	imgPlayers[battlecommon.PlayerActDamage] = make([]int, 6)
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 100, 100, imgPlayers[battlecommon.PlayerActDamage]); res == -1 {
		return nil, errors.Newf("failed to load player damage image: %s", fname)
	}
	// 1 -> 2,3  2-4 3-5
	imgPlayers[battlecommon.PlayerActDamage][4] = imgPlayers[battlecommon.PlayerActDamage][2]
	imgPlayers[battlecommon.PlayerActDamage][5] = imgPlayers[battlecommon.PlayerActDamage][3]
	imgPlayers[battlecommon.PlayerActDamage][2] = imgPlayers[battlecommon.PlayerActDamage][1]
	imgPlayers[battlecommon.PlayerActDamage][3] = imgPlayers[battlecommon.PlayerActDamage][1]

	fname = config.ImagePath + "battle/character/player_shot.png"
	imgPlayers[battlecommon.PlayerActShot] = make([]int, 6)
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 180, 100, imgPlayers[battlecommon.PlayerActShot]); res == -1 {
		return nil, errors.Newf("failed to load player shot image: %s", fname)
	}

	fname = config.ImagePath + "battle/character/player_cannon.png"
	imgPlayers[battlecommon.PlayerActCannon] = make([]int, 6)
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 100, 100, imgPlayers[battlecommon.PlayerActCannon]); res == -1 {
		return nil, errors.Newf("failed to load player cannon image: %s", fname)
	}

	fname = config.ImagePath + "battle/character/player_sword.png"
	imgPlayers[battlecommon.PlayerActSword] = make([]int, 7)
	if res := dxlib.LoadDivGraph(fname, 7, 7, 1, 128, 128, imgPlayers[battlecommon.PlayerActSword]); res == -1 {
		return nil, errors.Newf("failed to load player sword image: %s", fname)
	}

	fname = config.ImagePath + "battle/character/player_bomb.png"
	imgPlayers[battlecommon.PlayerActBomb] = make([]int, 7)
	if res := dxlib.LoadDivGraph(fname, 5, 5, 1, 100, 114, imgPlayers[battlecommon.PlayerActBomb]); res == -1 {
		return nil, errors.Newf("failed to load player bomb image: %s", fname)
	}
	imgPlayers[battlecommon.PlayerActBomb][5] = imgPlayers[battlecommon.PlayerActBomb][4]
	imgPlayers[battlecommon.PlayerActBomb][6] = imgPlayers[battlecommon.PlayerActBomb][4]

	fname = config.ImagePath + "battle/character/player_buster.png"
	imgPlayers[battlecommon.PlayerActBuster] = make([]int, 6)
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 180, 100, imgPlayers[battlecommon.PlayerActBuster]); res == -1 {
		return nil, errors.Newf("failed to load player buster image: %s", fname)
	}

	fname = config.ImagePath + "battle/character/player_pick.png"
	imgPlayers[battlecommon.PlayerActPick] = make([]int, 6)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 96, 124, imgPlayers[battlecommon.PlayerActPick]); res == -1 {
		return nil, errors.Newf("failed to load player pick image: %s", fname)
	}
	imgPlayers[battlecommon.PlayerActPick][4] = imgPlayers[battlecommon.PlayerActPick][3]
	imgPlayers[battlecommon.PlayerActPick][5] = imgPlayers[battlecommon.PlayerActPick][3]

	fname = config.ImagePath + "battle/character/player_throw.png"
	imgPlayers[battlecommon.PlayerActThrow] = make([]int, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 97, 115, imgPlayers[battlecommon.PlayerActThrow]); res == -1 {
		return nil, errors.Newf("failed to load player throw image: %s", fname)
	}

	imgPlayers[battlecommon.PlayerActParalyzed] = make([]int, 4)
	for i := 0; i < 4; i++ {
		imgPlayers[battlecommon.PlayerActParalyzed][i] = imgPlayers[battlecommon.PlayerActDamage][i]
	}

	fname = config.ImagePath + "battle/hp_frame.png"
	imgHPFrame = dxlib.LoadGraph(fname)
	if imgHPFrame < 0 {
		return nil, errors.Newf("failed to read hp frame image %s", fname)
	}
	fname = config.ImagePath + "battle/gauge.png"
	imgGaugeFrame = dxlib.LoadGraph(fname)
	if imgGaugeFrame < 0 {
		return nil, errors.Newf("failed to read gauge frame image %s", fname)
	}
	fname = config.ImagePath + "battle/gauge_max.png"
	imgGaugeMax = make([]int, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 1, 4, 288, 30, imgGaugeMax); res == -1 {
		return nil, errors.Newf("failed to read gauge max image %s", fname)
	}

	fname = config.ImagePath + "battle/skill/charge.png"
	tmp := make([]int, 16)
	if res := dxlib.LoadDivGraph(fname, 16, 8, 2, 158, 150, tmp); res == -1 {
		return nil, errors.Newf("failed to load image %s", fname)
	}
	for i := 0; i < 8; i++ {
		imgCharge[0] = append(imgCharge[0], tmp[i])
		imgCharge[1] = append(imgCharge[1], tmp[i+8])
	}

	fname = config.ImagePath + "battle/mind_window_frame.png"
	if imgMindFrame = dxlib.LoadGraph(fname); imgMindFrame == -1 {
		return nil, errors.Newf("failed to read mind frame image %s", fname)
	}

	fname = config.ImagePath + "battle/mind_status.png"
	imgMinds = make([]int, battlecommon.PlayerMindStatusMax)
	if res := dxlib.LoadDivGraph(fname, battlecommon.PlayerMindStatusMax, 6, 3, 88, 32, imgMinds); res == -1 {
		return nil, errors.Newf("failed to load image %s", fname)
	}

	fname = config.ImagePath + "battle/skill/バリア.png"
	imgBarrier = make([]int, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 144, 136, imgBarrier); res == -1 {
		return nil, errors.Newf("failed to load image %s", fname)
	}

	logger.Info("Successfully initialized battle player data")
	return &res, nil
}

func (p *BattlePlayer) End() {
	logger.Info("Cleanup battle player data")

	for i := 0; i < battlecommon.PlayerActMax; i++ {
		for j := 0; j < len(imgPlayers[i]); j++ {
			dxlib.DeleteGraph(imgPlayers[i][j])
			imgPlayers[i][j] = -1
		}
	}
	dxlib.DeleteGraph(imgHPFrame)
	imgHPFrame = -1
	dxlib.DeleteGraph(imgGaugeFrame)
	imgGaugeFrame = -1
	for _, img := range imgGaugeMax {
		dxlib.DeleteGraph(img)
	}
	for i := 0; i < 2; i++ {
		for _, img := range imgCharge[i] {
			dxlib.DeleteGraph(img)
		}
		imgCharge[i] = []int{}
	}
	dxlib.DeleteGraph(imgMindFrame)
	imgMindFrame = -1
	for _, img := range imgMinds {
		dxlib.DeleteGraph(img)
	}
	imgMinds = []int{}
	for _, img := range imgBarrier {
		dxlib.DeleteGraph(img)
	}
	imgBarrier = []int{}

	logger.Info("Successfully cleanuped battle player data")
}

func (p *BattlePlayer) Draw() {
	if !p.visible {
		return
	}

	frameX := 7
	frameY := 5
	if field.Is4x4Area() {
		frameY = 25
	}
	if p.isShiftFrame {
		frameX += 235
	}

	// Show Mind Status
	dxlib.DrawGraph(frameX, frameY+35, imgMindFrame, true)
	dxlib.DrawGraph(frameX, frameY+35, imgMinds[p.MindStatus], true)

	// Show selected chip icons
	n := len(p.SelectedChips)
	if n > 0 {
		// Show current chip info
		c := chip.Get(p.SelectedChips[0].ID)
		powTxt := ""
		if c.Power > 0 && !c.ForMe {
			powTxt = fmt.Sprintf("%d", c.Power)
			if p.SelectedChips[0].PlusPower > 0 {
				powTxt += fmt.Sprintf("＋ %d", p.SelectedChips[0].PlusPower)
			}
		}
		draw.String(5, config.ScreenSize.Y-20, 0xffffff, "%s %s", c.Name, powTxt)

		const px = 3
		max := n * px
		for i := 0; i < n; i++ {
			x := battlecommon.PanelSize.X*p.Pos.X + battlecommon.PanelSize.X/2 - 2 + (i * px) - max
			y := battlecommon.DrawPanelTopY + battlecommon.PanelSize.Y*p.Pos.Y - 10 - 81 + (i * px) - max
			dxlib.DrawBox(x-1, y-1, x+29, y+29, 0x000000, false)
			// draw from the end
			dxlib.DrawGraph(x, y, chipimage.GetIcon(p.SelectedChips[n-1-i].ID, true), true)
		}
	}

	playerVisible := p.invincibleCount/5%2 == 0

	if playerVisible {
		view := battlecommon.ViewPos(p.Pos)

		if p.barrierHP > 0 {
			cnt := system.GetGlobalCount()
			n := math.MountainIndex(int(cnt/15)%(len(imgBarrier)*2), len(imgBarrier)*2)
			dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_ALPHA, 156)
			dxlib.DrawRotaGraph(view.X, view.Y, 1, 0, imgBarrier[n], true)
			dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_NOBLEND, 0)
		}

		img := p.act.GetImage()
		dxlib.DrawRotaGraph(view.X, view.Y, 1, 0, img, true)
		if p.act.IsParalyzed() {
			dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_ADD, 255)
			// 黄色と白を点滅させる
			pm := 0
			if p.act.count/10%2 == 0 {
				pm = 255
			}
			dxlib.SetDrawBright(255, 255, pm)
			dxlib.DrawRotaGraph(view.X, view.Y, 1, 0, img, true)
			dxlib.SetDrawBright(255, 255, 255)
			dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_NOBLEND, 0)
		}

		// Show charge image
		if p.ChargeCount > battlecommon.ChargeViewDelay {
			n := 0
			if p.ChargeCount > p.ChargeTime {
				n = 1
			}
			imgNo := int(p.ChargeCount/4) % len(imgCharge[n])
			dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_ALPHA, 224)
			dxlib.DrawRotaGraph(view.X, view.Y, 1, 0, imgCharge[n][imgNo], true)
			dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_NOBLEND, 0)
		}
	}

	// Show HP
	dxlib.DrawGraph(frameX, frameY, imgHPFrame, true)
	col := draw.NumberColorWhite
	if p.HP*3 < p.HPMax {
		// HPが1/3未満の時はオレンジ色にする
		col = draw.NumberColorRed
	}
	draw.Number(frameX+2, frameY+2, int(p.HP), draw.NumberOption{RightAligned: true, Length: 4, Color: col})

	// Show Custom Gauge
	if p.isShowGauge {
		baseX := 5
		if field.Is4x4Area() {
			baseX = 80
		}

		if p.GaugeCount < battlecommon.GaugeMaxCount {
			dxlib.DrawGraph(96+baseX, frameY, imgGaugeFrame, true)
			const gaugeMaxSize = 256
			size := int(gaugeMaxSize * p.GaugeCount / battlecommon.GaugeMaxCount)
			dxlib.DrawBox(112+baseX, frameY+14, 112+baseX+size, frameY+16, dxlib.GetColor(123, 154, 222), true)
			dxlib.DrawBox(112+baseX, frameY+16, 112+baseX+size, frameY+24, dxlib.GetColor(231, 235, 255), true)
			dxlib.DrawBox(112+baseX, frameY+24, 112+baseX+size, frameY+26, dxlib.GetColor(123, 154, 222), true)
		} else {
			i := (p.GaugeCount / 40) % 4
			dxlib.DrawGraph(96+baseX, frameY, imgGaugeMax[i], true)
		}
	}
}

func (p *BattlePlayer) SetFrameInfo(xShift bool, showGauge bool) {
	p.isShiftFrame = xShift
	p.isShowGauge = showGauge
}

func (p *BattlePlayer) Update() (bool, error) {
	if !p.EnableAct {
		return false, nil
	}

	if p.HP <= 0 {
		// Player deleted
		img := &imgPlayers[battlecommon.PlayerActDamage][1]
		deleteanim.New(*img, p.Pos, true, p.animMgr)
		*img = -1 // DeleteGraph at delete animation
		p.NextAction = NextActLose
		p.EnableAct = false
		p.visible = false
		return false, nil
	}

	if p.invincibleCount > 0 && !config.Get().Debug.AlwaysInvisible {
		p.invincibleCount--
	}

	p.GaugeCount += uint(battlecommon.CustomGaugeSpeed)

	if p.act.Update() {
		return false, nil
	}

	if p.GaugeCount >= battlecommon.GaugeMaxCount {
		if p.GaugeCount == battlecommon.GaugeMaxCount {
			sound.On(resources.SEGaugeMax)
		}

		// State change to chip select
		if inputs.CheckKey(inputs.KeyLButton) == 1 || inputs.CheckKey(inputs.KeyRButton) == 1 {
			p.GaugeCount = 0
			p.NextAction = NextActChipSelect
			return false, nil
		}
	}

	// Move
	moveDirect := -1
	if inputs.CheckKey(inputs.KeyUp)%10 == 1 {
		moveDirect = config.DirectUp
	} else if inputs.CheckKey(inputs.KeyDown)%10 == 1 {
		moveDirect = config.DirectDown
	} else if inputs.CheckKey(inputs.KeyRight)%10 == 1 {
		moveDirect = config.DirectRight
	} else if inputs.CheckKey(inputs.KeyLeft)%10 == 1 {
		moveDirect = config.DirectLeft
	}

	if moveDirect >= 0 {
		if battlecommon.MoveObject(&p.Pos, moveDirect, battlecommon.PanelTypePlayer, false, field.GetPanelInfo) {
			p.act.MoveDirect = moveDirect
			p.act.SetAnim(battlecommon.PlayerActMove, 0)
			p.MoveNum++
			return false, nil
		}
	}

	// Chip use
	if inputs.CheckKey(inputs.KeyEnter) == 1 {
		if len(p.SelectedChips) > 0 {
			c := chip.Get(p.SelectedChips[0].ID)
			if c.PlayerAct != -1 {
				p.act.SetAnim(c.PlayerAct, c.KeepCount)
			}
			target := damage.TargetEnemy
			if c.ForMe {
				target = damage.TargetPlayer
			}

			sid := skillcore.GetIDByChipID(c.ID)
			p.act.SetSkill(sid, skillcore.Argument{
				OwnerID:    p.ID,
				Power:      c.Power + uint(p.SelectedChips[0].PlusPower),
				TargetType: target,
			})
			logger.Info("Use chip %d", sid)

			p.SelectedChips = p.SelectedChips[1:]
			return false, nil
		}
	}

	// Rock buster
	if inputs.CheckKey(inputs.KeyCancel) > 0 {
		p.ChargeCount++
		if p.ChargeCount == battlecommon.ChargeViewDelay {
			sound.On(resources.SEBusterCharging)
		}
		if p.ChargeCount == p.ChargeTime {
			sound.On(resources.SEBusterCharged)
		}
	} else if p.ChargeCount > 0 {
		sound.On(resources.SEBusterShot)
		p.act.Charged = p.ChargeCount > p.ChargeTime
		p.act.ShotPower = p.ShotPower
		p.act.SetAnim(battlecommon.PlayerActBuster, 0)
		p.ChargeCount = 0
	}

	return false, nil
}

func (p *BattlePlayer) DamageProc(dm *damage.Damage) bool {
	if dm == nil {
		return false
	}

	// インビジ中は無効、ただしRecover系は使えるようにする
	if p.invincibleCount > 0 && dm.Power >= 0 {
		return false
	}

	if p.barrierHP > 0 && dm.Power > 0 {
		logger.Debug("Barrier HP: %d, Damage: %d", p.barrierHP, dm.Power)
		p.barrierHP -= int(dm.Power)
		if p.barrierHP < 0 {
			p.barrierHP = 0
		}
		return true
	}

	if dm.TargetObjType&damage.TargetPlayer != 0 {
		prevHP := p.HP
		hp := int(p.HP) - dm.Power
		if hp <= 0 {
			if p.IsUnderShirt && prevHP >= 2 {
				p.IsUnderShirt = false
				p.HP = 1
			} else {
				p.HP = 0
			}
		} else if hp > int(p.HPMax) {
			p.HP = p.HPMax
		} else {
			p.HP = uint(hp)
		}
		p.animMgr.EffectAnimNew(effect.Get(dm.HitEffectType, p.Pos, 5))

		for i := 0; i < dm.PushLeft; i++ {
			if !battlecommon.MoveObject(&p.Pos, config.DirectLeft, battlecommon.PanelTypePlayer, true, field.GetPanelInfo) {
				break
			}
		}
		for i := 0; i < dm.PushRight; i++ {
			if !battlecommon.MoveObject(&p.Pos, config.DirectRight, battlecommon.PanelTypePlayer, true, field.GetPanelInfo) {
				break
			}
		}

		if dm.Power <= 0 {
			// Not damage, maybe recover or special anim
			return true
		}

		if dm.StrengthType == damage.StrengthNone {
			return true
		}

		sound.On(resources.SEDamaged)

		// Stop current animation
		if p.animMgr.IsAnimProcessing(p.act.skillObjID) {
			p.act.skillInst.StopByOwner()
		}
		p.act.skillObjID = ""
		p.ChargeCount = 0

		if dm.IsParalyzed {
			p.act.SetAnim(battlecommon.PlayerActParalyzed, battlecommon.DefaultParalyzedTime)
		} else {
			p.act.SetAnim(battlecommon.PlayerActDamage, 0)
			if dm.StrengthType == damage.StrengthHigh {
				p.MakeInvisible(battlecommon.PlayerDefaultInvincibleTime)
			}
		}

		p.DamageNum++
		logger.Debug("Player damaged: %+v", *dm)
		return true
	}
	return false
}

func (p *BattlePlayer) GetParam() objanim.Param {
	return objanim.Param{
		Param: anim.Param{
			ObjID: p.ID,
			Pos:   p.Pos,
		},
		HP: int(p.HP),
	}
}

func (p *BattlePlayer) GetObjectType() int {
	return objanim.ObjTypePlayer
}

func (p *BattlePlayer) MakeInvisible(count int) {
	p.invincibleCount = count
}

func (p *BattlePlayer) AddBarrier(hp int) {
	p.barrierHP = hp
}

func (p *BattlePlayer) SetCustomGaugeMax() {
	if p.GaugeCount < battlecommon.GaugeMaxCount {
		// 通常のゲージアップでSEがならないようにする
		p.GaugeCount = battlecommon.GaugeMaxCount + uint(battlecommon.CustomGaugeSpeed)
		sound.On(resources.SEGaugeMax)
	}
}

func (p *BattlePlayer) SetChipSelectResult(selected []int) {
	p.SelectedChips = []SelectChip{}
	for _, s := range selected {
		p.SelectedChips = append(
			p.SelectedChips,
			SelectChip{
				ID:   p.ChipFolder[s].ID,
				Code: p.ChipFolder[s].Code,
			},
		)
	}

	// Remove selected chips from folder
	sort.Sort(sort.Reverse(sort.IntSlice(selected)))
	for _, s := range selected {
		p.ChipFolder = append(p.ChipFolder[:s], p.ChipFolder[s+1:]...)
	}
}

func (p *BattlePlayer) UpdateChipInfo() {
	// Check program advance
	list := []chip.SelectParam{}
	for _, c := range p.SelectedChips {
		list = append(list, chip.SelectParam{ID: c.ID, Code: c.Code})
	}

	start, end, paID := chip.GetPAinList(list)
	if paID != -1 {
		before := append([]SelectChip{}, p.SelectedChips[:start]...)
		after := append([]SelectChip{}, p.SelectedChips[end:]...)
		p.SelectedChips = append(before, SelectChip{ID: paID})
		p.SelectedChips = append(p.SelectedChips, after...)
	}

	// アタック+10などの処理
	if len(p.SelectedChips) >= 2 {
		removes := []int{}
		target := 0
		for i := 1; i < len(p.SelectedChips); i++ {
			if p.SelectedChips[i].ID == chip.IDAttack10 {
				p.SelectedChips[target].PlusPower += 10
				removes = append(removes, i)
			} else {
				target = i
			}
		}
		tmp := append([]SelectChip{}, p.SelectedChips...)
		p.SelectedChips = []SelectChip{}
		for i := 0; i < len(tmp); i++ {
			if !slice.Contains(removes, i) {
				p.SelectedChips = append(p.SelectedChips, tmp[i])
			}
		}
	}

	logger.Info("selected player chips: %+v", p.SelectedChips)
}

func (a *BattlePlayerAct) Init(pPos *point.Point, animMgr *manager.Manager) {
	a.animMgr = animMgr
	a.typ = -1
	a.pPos = pPos
}

// Process method returns true if processing now
func (a *BattlePlayerAct) Update() bool {
	switch a.typ {
	case -1: // No animation
		return false
	case battlecommon.PlayerActBuster:
		if a.count == 1 {
			s := a.ShotPower
			eff := resources.EffectTypeHitSmall
			if a.Charged {
				s *= 10
				eff = resources.EffectTypeHitBig
			}

			y := a.pPos.Y
			for x := a.pPos.X + 1; x < battlecommon.FieldNum.X; x++ {
				// logger.Debug("Rock buster damage set %d to (%d, %d)", s, x, y)
				if objID := field.GetPanelInfo(point.Point{X: x, Y: y}).ObjectID; objID != "" {
					a.animMgr.DamageManager().New(damage.Damage{
						DamageType:    damage.TypeObject,
						TargetObjID:   objID,
						TargetObjType: damage.TargetEnemy,
						Power:         int(s),
						HitEffectType: eff,
						Element:       damage.ElementNone,
					})
					break
				}
			}
		}
	case battlecommon.PlayerActMove:
		if a.count == 2 {
			battlecommon.MoveObject(a.pPos, a.MoveDirect, battlecommon.PanelTypePlayer, true, field.GetPanelInfo)
		}
	case battlecommon.PlayerActCannon, battlecommon.PlayerActSword, battlecommon.PlayerActBomb, battlecommon.PlayerActDamage, battlecommon.PlayerActShot, battlecommon.PlayerActPick, battlecommon.PlayerActThrow, battlecommon.PlayerActParalyzed:
		// No special action
	default:
		system.SetError(fmt.Sprintf("Invalid player anim type %d was specified.", a.typ))
		return false
	}

	a.count++

	if a.count > a.endCount {
		// Reset params
		a.typ = -1
		a.count = 0
		a.endCount = 0
		return false // finished
	}
	return true // processing now
}

func (a *BattlePlayerAct) SetAnim(animType int, keepCount int) {
	a.count = 0
	a.typ = animType
	a.endCount = battlecommon.GetPlayerActCount(animType, keepCount)
}

func (a *BattlePlayerAct) GetImage() int {
	if a.typ == -1 {
		// return stand image
		return imgPlayers[battlecommon.PlayerActMove][0]
	}

	num, delay := battlecommon.GetPlayerImageInfo(a.typ)
	imgNo := (a.count / delay)
	if imgNo >= num {
		imgNo = num - 1
	}

	return imgPlayers[a.typ][imgNo]
}

func (a *BattlePlayerAct) IsParalyzed() bool {
	return a.typ == battlecommon.PlayerActParalyzed
}

func (a *BattlePlayerAct) SetSkill(id int, arg skillcore.Argument) {
	a.skillInst = skill.Get(id, arg, a.animMgr)
	a.skillObjID = a.animMgr.SkillAnimNew(a.skillInst)
}
