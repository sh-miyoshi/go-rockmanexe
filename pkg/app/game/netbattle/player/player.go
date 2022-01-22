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
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/object"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/net/routerpb"
)

type BattlePlayer struct {
	Object          object.Object
	HPMax           uint
	ChargeCount     uint
	GaugeCount      uint
	ShotPower       uint
	ChipFolder      []player.ChipInfo
	Act             *Act
	HitDamages      map[string]bool
	ManagedSkills   []string
	InvincibleCount uint
}

var (
	imgHPFrame    int
	imgGaugeFrame int
	imgGaugeMax   []int
	imgCharge     [2][]int
)

// New ...
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
		HPMax:      plyr.HP,
		ShotPower:  plyr.ShotPower,
		HitDamages: make(map[string]bool),
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
	imgHPFrame = dxlib.LoadGraph(fname)
	if imgHPFrame < 0 {
		return nil, fmt.Errorf("failed to read hp frame image %s", fname)
	}
	fname = common.ImagePath + "battle/gauge.png"
	imgGaugeFrame = dxlib.LoadGraph(fname)
	if imgGaugeFrame < 0 {
		return nil, fmt.Errorf("failed to read gauge frame image %s", fname)
	}
	fname = common.ImagePath + "battle/gauge_max.png"
	imgGaugeMax = make([]int, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 1, 4, 288, 30, imgGaugeMax); res == -1 {
		return nil, fmt.Errorf("failed to read gauge max image %s", fname)
	}

	fname = common.ImagePath + "battle/skill/charge.png"
	tmp := make([]int, 16)
	if res := dxlib.LoadDivGraph(fname, 16, 8, 2, 158, 150, tmp); res == -1 {
		return nil, fmt.Errorf("failed to load image %s", fname)
	}
	for i := 0; i < 8; i++ {
		imgCharge[0] = append(imgCharge[0], tmp[i])
		imgCharge[1] = append(imgCharge[1], tmp[i+8])
	}

	logger.Info("Successfully initialized net battle player data")
	return &res, nil
}

func (p *BattlePlayer) End() {
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
}

func (p *BattlePlayer) InitBattleFrame() {
	p.GaugeCount = 0
}

func (p *BattlePlayer) DrawOptions() {
	// Show Charge Shot
	if p.ChargeCount > battlecommon.ChargeViewDelay {
		n := 0
		if p.ChargeCount > battlecommon.ChargeTime {
			n = 1
		}
		view := battlecommon.ViewPos(common.Point{X: p.Object.X, Y: p.Object.Y})
		imgNo := int(p.ChargeCount/4) % len(imgCharge[n])
		dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_ALPHA, 224)
		dxlib.DrawRotaGraph(view.X, view.Y, 1, 0, imgCharge[n][imgNo], true)
		dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_NOBLEND, 0)
	}

	// Show selected chip icons
	n := len(p.Object.Chips)
	if n > 0 {
		// TODO Show chip info

		const px = 3
		max := n * px
		for i := 0; i < n; i++ {
			x := appfield.PanelSize.X*p.Object.X + appfield.PanelSize.X/2 - 2 + i*px - max
			y := appfield.DrawPanelTopY + appfield.PanelSize.Y*p.Object.Y - 10 - 81 + (i * px) - max
			dxlib.DrawBox(x-1, y-1, x+29, y+29, 0x000000, false)
			// draw from the end
			dxlib.DrawGraph(x, y, chip.GetIcon(p.Object.Chips[n-1-i], true), true)
		}
	}
}

func (p *BattlePlayer) DrawFrame(xShift bool, showGauge bool) {
	x := 7
	y := 5
	if xShift {
		x += 235
	}

	// Show HP
	dxlib.DrawGraph(x, y, imgHPFrame, true)
	draw.Number(x+2, y+2, p.Object.HP, draw.NumberOption{RightAligned: true, Length: 4})

	// Show Custom Gauge
	if showGauge {
		if p.GaugeCount < battlecommon.GaugeMaxCount {
			dxlib.DrawGraph(96, 5, imgGaugeFrame, true)
			const gaugeMaxSize = 256
			size := gaugeMaxSize * int(p.GaugeCount) / battlecommon.GaugeMaxCount
			dxlib.DrawBox(112, 19, 112+size, 21, dxlib.GetColor(123, 154, 222), true)
			dxlib.DrawBox(112, 21, 112+size, 29, dxlib.GetColor(231, 235, 255), true)
			dxlib.DrawBox(112, 29, 112+size, 31, dxlib.GetColor(123, 154, 222), true)
		} else {
			i := (p.GaugeCount / 20) % 4
			dxlib.DrawGraph(96, 5, imgGaugeMax[i], true)
		}
	}
}

func (p *BattlePlayer) Process() (bool, error) {
	p.GaugeCount += 4 // TODO GaugeSpeed

	if p.Object.HP <= 0 {
		netconn.SendSignal(pb.Action_PLAYERDEAD)
		return true, nil
	}

	if p.Object.Invincible {
		p.InvincibleCount++
		if p.InvincibleCount > battlecommon.PlayerDefaultInvincibleTime {
			p.InvincibleCount = 0
			p.Object.Invincible = false
			netconn.SendObject(p.Object)
		}
	}

	if p.damageProc() {
		return false, nil
	}

	if p.Act.Process() {
		return false, nil
	}

	if p.GaugeCount >= battlecommon.GaugeMaxCount {
		if p.GaugeCount == battlecommon.GaugeMaxCount {
			sound.On(sound.SEGaugeMax)
		}

		// State change to chip select
		if inputs.CheckKey(inputs.KeyLButton) == 1 || inputs.CheckKey(inputs.KeyRButton) == 1 {
			netconn.SendSignal(pb.Action_GOCHIPSELECT)

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
		return false, nil
	}

	// Chip use
	if inputs.CheckKey(inputs.KeyEnter) == 1 {
		if len(p.Object.Chips) > 0 {
			c := chip.Get(p.Object.Chips[0])
			if c.PlayerAct != -1 {
				p.Act.Set(c.PlayerAct, &ActOption{
					KeepCount: c.KeepCount,
				})
			}

			sid := skill.GetSkillID(c.ID)
			id := netskill.Add(sid, netskill.Argument{
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

func (p *BattlePlayer) SetChipSelectResult(selected []int) {
	p.Object.Chips = []int{}
	for _, s := range selected {
		p.Object.Chips = append(p.Object.Chips, p.ChipFolder[s].ID)
	}

	// Remove selected chips from folder
	sort.Sort(sort.Reverse(sort.IntSlice(selected)))
	for _, s := range selected {
		p.ChipFolder = append(p.ChipFolder[:s], p.ChipFolder[s+1:]...)
	}
}

func (p *BattlePlayer) damageProc() bool {
	finfo := netconn.GetFieldInfo()
	if len(finfo.HitDamages) == 0 {
		return false
	}

	dm := finfo.HitDamages[0]
	defer netconn.RemoveDamage(dm.ID)

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
		for _, sid := range p.ManagedSkills {
			netskill.StopByPlayer(sid)
		}
		p.ManagedSkills = []string{}
		p.Act.Set(battlecommon.PlayerActDamage, nil)
		netconn.AddSound(sound.SEDamaged)
	} else {
		netconn.SendObject(p.Object)
	}

	if dm.HitEffectType > 0 {
		netconn.SendEffect(effect.Effect{
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
