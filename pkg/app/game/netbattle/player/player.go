package player

import (
	"fmt"
	"math/rand"
	"sort"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/chip"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	appfield "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill"
	netfield "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/field"
	netskill "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/skill"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/inputs"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/netconn"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/player"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/effect"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/net/netconnpb"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/object"
)

type BattlePlayer struct {
	Object        object.Object
	ChipFolder    []player.ChipInfo
	GaugeCount    uint
	ChargeCount   uint
	ShotPower     uint
	Act           *Act
	HPMax         uint
	HitDamages    map[string]bool
	ManagedSkills []string

	imgHPFrame    int
	imgGaugeFrame int
	imgGaugeMax   []int
	imgMinds      []int
	imgMindFrame  int
	imgCharge     [2][]int
}

func New(plyr *player.Player) (*BattlePlayer, error) {
	logger.Info("Initialize net battle player data")
	cfg := config.Get()

	res := BattlePlayer{
		Object: object.Object{
			ID:       uuid.New().String(),
			HP:       int(plyr.HP),
			X:        1,
			Y:        1,
			ClientID: cfg.Net.ClientID,
			Hittable: true,
		},
		ShotPower:     plyr.ShotPower,
		HPMax:         plyr.HP,
		HitDamages:    make(map[string]bool),
		ManagedSkills: []string{},
	}
	res.Act = NewAct(&res.Object)

	for _, c := range plyr.ChipFolder {
		res.ChipFolder = append(res.ChipFolder, c)
	}
	// Shuffle folder
	for i := 0; i < 10; i++ {
		for j := 0; j < len(res.ChipFolder); j++ {
			n := rand.Intn(len(res.ChipFolder))
			res.ChipFolder[j], res.ChipFolder[n] = res.ChipFolder[n], res.ChipFolder[j]
		}
	}

	logger.Debug("Player info: %+v", res)

	fname := common.ImagePath + "battle/hp_frame.png"
	res.imgHPFrame = dxlib.LoadGraph(fname)
	if res.imgHPFrame < 0 {
		return nil, fmt.Errorf("failed to read hp frame image %s", fname)
	}
	fname = common.ImagePath + "battle/gauge.png"
	res.imgGaugeFrame = dxlib.LoadGraph(fname)
	if res.imgGaugeFrame < 0 {
		return nil, fmt.Errorf("failed to read gauge frame image %s", fname)
	}
	fname = common.ImagePath + "battle/gauge_max.png"
	res.imgGaugeMax = make([]int, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 1, 4, 288, 30, res.imgGaugeMax); res == -1 {
		return nil, fmt.Errorf("failed to read gauge max image %s", fname)
	}

	fname = common.ImagePath + "battle/mind_window_frame.png"
	if res.imgMindFrame = dxlib.LoadGraph(fname); res.imgMindFrame == -1 {
		return nil, fmt.Errorf("failed to read mind frame image %s", fname)
	}

	fname = common.ImagePath + "battle/mind_status.png"
	res.imgMinds = make([]int, battlecommon.PlayerMindStatusMax)
	if res := dxlib.LoadDivGraph(fname, battlecommon.PlayerMindStatusMax, 6, 3, 88, 32, res.imgMinds); res == -1 {
		return nil, fmt.Errorf("failed to load image %s", fname)
	}

	fname = common.ImagePath + "battle/skill/charge.png"
	tmp := make([]int, 16)
	if res := dxlib.LoadDivGraph(fname, 16, 8, 2, 158, 150, tmp); res == -1 {
		return nil, fmt.Errorf("failed to load image %s", fname)
	}
	for i := 0; i < 8; i++ {
		res.imgCharge[0] = append(res.imgCharge[0], tmp[i])
		res.imgCharge[1] = append(res.imgCharge[1], tmp[i+8])
	}

	logger.Info("Successfully initialized net battle player data")
	return &res, nil
}

func (p *BattlePlayer) End() {
	dxlib.DeleteGraph(p.imgHPFrame)
	dxlib.DeleteGraph(p.imgGaugeFrame)
	dxlib.DeleteGraph(p.imgMindFrame)

	for _, img := range p.imgGaugeMax {
		dxlib.DeleteGraph(img)
	}
	p.imgGaugeMax = []int{}

	for _, img := range p.imgMinds {
		dxlib.DeleteGraph(img)
	}
	p.imgMinds = []int{}

	for i := 0; i < 2; i++ {
		for _, img := range p.imgCharge[i] {
			dxlib.DeleteGraph(img)
		}
		p.imgCharge[i] = []int{}
	}
}

func (p *BattlePlayer) DrawFrame(xShift bool, showGauge bool) {
	x := 7
	y := 5
	if xShift {
		x += 235
	}

	// Show HP
	dxlib.DrawGraph(x, y, p.imgHPFrame, true)
	draw.Number(x+2, y+2, p.Object.HP, draw.NumberOption{RightAligned: true, Length: 4})

	// Show Mind Status
	dxlib.DrawGraph(x, 40, p.imgMindFrame, true)
	dxlib.DrawGraph(x, 40, p.imgMinds[battlecommon.PlayerMindStatusNormal], true) // TODO set mind status

	// Show Custom Gauge
	if showGauge {
		if p.GaugeCount < battlecommon.GaugeMaxCount {
			dxlib.DrawGraph(96, 5, p.imgGaugeFrame, true)
			const gaugeMaxSize = 256
			size := int(gaugeMaxSize * p.GaugeCount / battlecommon.GaugeMaxCount)
			dxlib.DrawBox(112, 19, 112+size, 21, dxlib.GetColor(123, 154, 222), true)
			dxlib.DrawBox(112, 21, 112+size, 29, dxlib.GetColor(231, 235, 255), true)
			dxlib.DrawBox(112, 29, 112+size, 31, dxlib.GetColor(123, 154, 222), true)
		} else {
			i := (p.GaugeCount / 40) % 4
			dxlib.DrawGraph(96, 5, p.imgGaugeMax[i], true)
		}
	}
}

func (p *BattlePlayer) LocalDraw() {
	view := battlecommon.ViewPos(common.Point{
		X: p.Object.X,
		Y: p.Object.Y,
	})

	// Show charge image
	if p.ChargeCount > battlecommon.ChargeViewDelay {
		n := 0
		if p.ChargeCount > battlecommon.ChargeTime {
			n = 1
		}
		imgNo := int(p.ChargeCount/4) % len(p.imgCharge[n])
		dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_ALPHA, 224)
		dxlib.DrawRotaGraph(view.X, view.Y, 1, 0, p.imgCharge[n][imgNo], true)
		dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_NOBLEND, 0)
	}
}

func (p *BattlePlayer) Process() (bool, error) {
	p.GaugeCount += 4 // TODO GaugeSpeed

	if p.Object.HP <= 0 {
		return true, nil
	}

	if p.damageProc() {
		return false, nil
	}

	if p.Act.Process() {
		return false, nil
	}

	// Go to chip folder
	if p.GaugeCount >= battlecommon.GaugeMaxCount {
		if p.GaugeCount == battlecommon.GaugeMaxCount {
			sound.On(sound.SEGaugeMax)
		}

		// State change to chip select
		if inputs.CheckKey(inputs.KeyLButton) == 1 || inputs.CheckKey(inputs.KeyRButton) == 1 {
			netconn.GetInst().SendSignal(pb.Action_GOCHIPSELECT)

			return false, nil
		}
	}

	// Move
	moveDirect := -1
	if inputs.CheckKey(inputs.KeyUp) == 1 {
		moveDirect = common.DirectUp
	} else if inputs.CheckKey(inputs.KeyDown) == 1 {
		moveDirect = common.DirectDown
	} else if inputs.CheckKey(inputs.KeyRight) == 1 {
		moveDirect = common.DirectRight
	} else if inputs.CheckKey(inputs.KeyLeft) == 1 {
		moveDirect = common.DirectLeft
	}

	if moveDirect >= 0 {
		t := common.Point{X: p.Object.X, Y: p.Object.Y}
		if battlecommon.MoveObject(&t, moveDirect, appfield.PanelTypePlayer, false, netfield.GetPanelInfo) {
			p.Act.Set(battlecommon.PlayerActMove, &ActOption{
				MoveDirect: moveDirect,
			})
		}
	}

	// Chip use
	if inputs.CheckKey(inputs.KeyEnter) == 1 {
		if len(p.Object.Chips) > 0 {
			c := chip.Get(p.Object.Chips[0].ID)
			if c.PlayerAct != -1 {
				p.Act.Set(c.PlayerAct, &ActOption{
					KeepCount: c.KeepCount,
				})
			}

			sid := skill.GetSkillID(c.ID)
			id := netskill.GetInst().Add(sid, netskill.Argument{
				X:     p.Object.X,
				Y:     p.Object.Y,
				Power: int(c.Power),
			})
			p.ManagedSkills = append(p.ManagedSkills, id)

			p.Object.Chips = p.Object.Chips[1:]
			return false, nil
		}
	}

	// Rock buster
	if inputs.CheckKey(inputs.KeyCancel) > 0 {
		p.ChargeCount++
		if p.ChargeCount == battlecommon.ChargeViewDelay {
			sound.On(sound.SEBusterCharging)
		}
		if p.ChargeCount == battlecommon.ChargeTime {
			sound.On(sound.SEBusterCharged)
		}
	} else if p.ChargeCount > 0 {
		sound.On(sound.SEBusterShot)
		p.Act.Set(battlecommon.PlayerActBuster, &ActOption{
			Charged:   p.ChargeCount > battlecommon.ChargeTime,
			ShotPower: int(p.ShotPower),
		})
		p.ChargeCount = 0
	}

	return false, nil
}

func (p *BattlePlayer) damageProc() bool {
	ginfo := netconn.GetInst().GetGameInfo()
	if len(ginfo.HitDamages) == 0 {
		return false
	}

	dm := ginfo.HitDamages[0]
	defer netconn.GetInst().RemoveDamage(dm.ID)

	if _, exists := p.HitDamages[dm.ID]; exists {
		return false
	} else {
		p.HitDamages[dm.ID] = true
	}

	// Recover系は使えるようにする
	if p.Object.Invincible && dm.Power >= 0 {
		return false
	}

	logger.Debug("Got damage: %+v", dm)

	p.Object.HP -= dm.Power
	if p.Object.HP < 0 {
		p.Object.HP = 0
	}
	if p.Object.HP > int(p.HPMax) {
		p.Object.HP = int(p.HPMax)
	}

	if dm.BigDamage {
		p.Object.Invincible = true
		// TODO Skill関係
		// for _, sid := range p.ManagedSkills {
		// 	netskill.StopByPlayer(sid)
		// }
		// p.ManagedSkills = []string{}
		// netconn.GetInst().AddSound(sound.SEDamaged)
		p.Act.Set(battlecommon.PlayerActDamage, nil)
	} else {
		netconn.GetInst().SendObject(p.Object)
	}

	if dm.HitEffectType > 0 {
		netconn.GetInst().SendEffect(effect.Effect{
			ID:       uuid.New().String(),
			Type:     dm.HitEffectType,
			X:        p.Object.X,
			Y:        p.Object.Y,
			ViewOfsX: dm.ViewOfsX,
			ViewOfsY: dm.ViewOfsY,
		})
	}

	return true
}

func (p *BattlePlayer) SetChipSelectResult(selected []int) {
	p.Object.Chips = []object.ChipInfo{}
	for _, s := range selected {
		p.Object.Chips = append(p.Object.Chips, object.ChipInfo{ID: p.ChipFolder[s].ID, Code: p.ChipFolder[s].Code})
	}

	// Remove selected chips from folder
	sort.Sort(sort.Reverse(sort.IntSlice(selected)))
	for _, s := range selected {
		p.ChipFolder = append(p.ChipFolder[:s], p.ChipFolder[s+1:]...)
	}
}

func (p *BattlePlayer) GetSelectedChips() []player.ChipInfo {
	res := []player.ChipInfo{}
	for _, c := range p.Object.Chips {
		res = append(res, player.ChipInfo{ID: c.ID, Code: c.Code})
	}
	return res
}

func (p *BattlePlayer) UpdatePA() {
	// Check program advance
	// TODO
}
