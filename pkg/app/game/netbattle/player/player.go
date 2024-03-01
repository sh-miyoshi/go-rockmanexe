package player

import (
	"fmt"
	"math/rand"
	"sort"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/chip"
	chipimage "github.com/sh-miyoshi/go-rockmanexe/pkg/app/chip/image"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	battleplayer "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/player"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/net"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/player"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/inputs"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/action"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/net/netconnpb"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/router/gameinfo"
)

type BattlePlayer struct {
	objectID      string
	chipFolder    []player.ChipInfo
	selectedChips []battleplayer.SelectChip
	imgHPFrame    int
	imgGaugeFrame int
	imgGaugeMax   []int
	imgMinds      []int
	imgMindFrame  int
	imgCharge     [2][]int
	chipAnimID    string
	chargeCount   int
	shotPower     int
	chargeTime    int
	gaugeCount    int
}

func New(plyr *player.Player) (*BattlePlayer, error) {
	res := &BattlePlayer{
		objectID:    uuid.New().String(),
		chargeCount: 0,
		shotPower:   int(plyr.ShotPower),
		chargeTime:  int(plyr.ChargeTime),
		gaugeCount:  0,
	}
	for _, c := range plyr.ChipFolder {
		res.chipFolder = append(res.chipFolder, c)
	}
	if !config.Get().Debug.UseDebugFolder {
		// Shuffle folder
		for i := 0; i < 10; i++ {
			for j := 0; j < len(res.chipFolder); j++ {
				n := rand.Intn(len(res.chipFolder))
				res.chipFolder[j], res.chipFolder[n] = res.chipFolder[n], res.chipFolder[j]
			}
		}
	}
	logger.Debug("Player info: %+v", res)

	fname := config.ImagePath + "battle/hp_frame.png"
	res.imgHPFrame = dxlib.LoadGraph(fname)
	if res.imgHPFrame < 0 {
		return nil, fmt.Errorf("failed to read hp frame image %s", fname)
	}
	fname = config.ImagePath + "battle/gauge.png"
	res.imgGaugeFrame = dxlib.LoadGraph(fname)
	if res.imgGaugeFrame < 0 {
		return nil, fmt.Errorf("failed to read gauge frame image %s", fname)
	}
	fname = config.ImagePath + "battle/gauge_max.png"
	res.imgGaugeMax = make([]int, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 1, 4, 288, 30, res.imgGaugeMax); res == -1 {
		return nil, fmt.Errorf("failed to read gauge max image %s", fname)
	}

	fname = config.ImagePath + "battle/mind_window_frame.png"
	if res.imgMindFrame = dxlib.LoadGraph(fname); res.imgMindFrame == -1 {
		return nil, fmt.Errorf("failed to read mind frame image %s", fname)
	}

	fname = config.ImagePath + "battle/mind_status.png"
	res.imgMinds = make([]int, battlecommon.PlayerMindStatusMax)
	if res := dxlib.LoadDivGraph(fname, battlecommon.PlayerMindStatusMax, 6, 3, 88, 32, res.imgMinds); res == -1 {
		return nil, fmt.Errorf("failed to load image %s", fname)
	}

	fname = config.ImagePath + "battle/skill/charge.png"
	tmp := make([]int, 16)
	if res := dxlib.LoadDivGraph(fname, 16, 8, 2, 158, 150, tmp); res == -1 {
		return nil, fmt.Errorf("failed to load image %s", fname)
	}
	for i := 0; i < 8; i++ {
		res.imgCharge[0] = append(res.imgCharge[0], tmp[i])
		res.imgCharge[1] = append(res.imgCharge[1], tmp[i+8])
	}

	logger.Info("Successfully initialized net battle player data")

	return res, nil
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
	obj := p.GetObject()
	draw.Number(x+2, y+2, obj.HP, draw.NumberOption{RightAligned: true, Length: 4})

	// Show Mind Status
	dxlib.DrawGraph(x, 40, p.imgMindFrame, true)
	dxlib.DrawGraph(x, 40, p.imgMinds[battlecommon.PlayerMindStatusNormal], true) // TODO set mind status

	// Show Custom Gauge
	if showGauge {
		baseX := 5
		if field.Is4x4Area() {
			baseX = 80
		}

		if p.gaugeCount < battlecommon.GaugeMaxCount {
			dxlib.DrawGraph(96+baseX, y, p.imgGaugeFrame, true)
			const gaugeMaxSize = 256
			size := int(gaugeMaxSize * p.gaugeCount / battlecommon.GaugeMaxCount)
			dxlib.DrawBox(112+baseX, y+14, 112+baseX+size, y+16, dxlib.GetColor(123, 154, 222), true)
			dxlib.DrawBox(112+baseX, y+16, 112+baseX+size, y+24, dxlib.GetColor(231, 235, 255), true)
			dxlib.DrawBox(112+baseX, y+24, 112+baseX+size, y+26, dxlib.GetColor(123, 154, 222), true)
		} else {
			i := (p.gaugeCount / 40) % 4
			dxlib.DrawGraph(96+baseX, y, p.imgGaugeMax[i], true)
		}
	}
}

func (p *BattlePlayer) LocalDraw() {
	// Show selected chip icons
	n := len(p.selectedChips)
	if n > 0 {
		// Show current chip info
		c := chip.Get(p.selectedChips[0].ID)
		powTxt := ""
		if c.Power > 0 && !c.ForMe {
			powTxt = fmt.Sprintf("%d", c.Power)
			if p.selectedChips[0].PlusPower > 0 {
				powTxt += fmt.Sprintf("＋ %d", p.selectedChips[0].PlusPower)
			}
		}
		draw.String(5, config.ScreenSize.Y-20, 0xffffff, "%s %s", c.Name, powTxt)

		const px = 3
		max := n * px
		obj := p.GetObject()
		for i := 0; i < n; i++ {
			x := battlecommon.PanelSize.X*obj.Pos.X + battlecommon.PanelSize.X/2 - 2 + (i * px) - max
			y := battlecommon.DrawPanelTopY + battlecommon.PanelSize.Y*obj.Pos.Y - 10 - 81 + (i * px) - max
			dxlib.DrawBox(x-1, y-1, x+29, y+29, 0x000000, false)
			// draw from the end
			dxlib.DrawGraph(x, y, chipimage.GetIcon(p.selectedChips[n-1-i].ID, true), true)
		}
	}
}

func (p *BattlePlayer) Process() (bool, error) {
	p.gaugeCount += battlecommon.CustomGaugeSpeed

	info := net.GetInst().GetGameInfo()
	for _, anim := range info.Anims {
		if anim.ObjectID == p.chipAnimID {
			return false, nil // まだ処理中
		}
	}
	p.chipAnimID = ""

	if p.gaugeCount >= battlecommon.GaugeMaxCount {
		if p.gaugeCount == battlecommon.GaugeMaxCount {
			sound.On(resources.SEGaugeMax)
		}

		// State change to chip select
		if inputs.CheckKey(inputs.KeyLButton) == 1 || inputs.CheckKey(inputs.KeyRButton) == 1 {
			p.gaugeCount = 0
			net.GetInst().SendSignal(pb.Request_GOCHIPSELECT, nil)
			return false, nil
		}
	}

	// Chip Use
	if inputs.CheckKey(inputs.KeyEnter) == 1 {
		if len(p.selectedChips) > 0 {
			cid := p.selectedChips[0].ID
			p.chipAnimID = uuid.New().String()
			logger.Info("Use chip %d", cid)

			chipInfo := action.UseChip{
				AnimID:           p.chipAnimID,
				ChipUserClientID: config.Get().Net.ClientID,
				ChipID:           cid,
			}
			net.GetInst().SendAction(pb.Request_CHIPUSE, chipInfo.Marshal())

			p.selectedChips = p.selectedChips[1:]
			return false, nil
		}
	}

	// Move
	moveDirect := -1
	if inputs.CheckKey(inputs.KeyUp) == 1 {
		moveDirect = config.DirectUp
	} else if inputs.CheckKey(inputs.KeyDown) == 1 {
		moveDirect = config.DirectDown
	} else if inputs.CheckKey(inputs.KeyRight) == 1 {
		moveDirect = config.DirectRight
	} else if inputs.CheckKey(inputs.KeyLeft) == 1 {
		moveDirect = config.DirectLeft
	}

	if moveDirect >= 0 {
		move := action.Move{
			Type:   action.MoveTypeDirect,
			Direct: moveDirect,
		}
		net.GetInst().SendAction(pb.Request_MOVE, move.Marshal())
		return false, nil
	}

	// Rock buster
	if inputs.CheckKey(inputs.KeyCancel) > 0 {
		p.chargeCount++
		if p.chargeCount == battlecommon.ChargeViewDelay {
			sound.On(resources.SEBusterCharging)
		}
		if p.chargeCount == p.chargeTime {
			sound.On(resources.SEBusterCharged)
		}
	} else if p.chargeCount > 0 {
		sound.On(resources.SEBusterShot)
		charged := p.chargeCount > p.chargeTime
		power := p.shotPower
		if charged {
			power *= 10
		}

		buster := action.Buster{
			Power:     power,
			IsCharged: charged,
		}
		net.GetInst().SendAction(pb.Request_BUSTER, buster.Marshal())
		p.chargeCount = 0
	}

	return false, nil
}

func (p *BattlePlayer) GetChipFolder() []player.ChipInfo {
	return p.chipFolder
}

func (p *BattlePlayer) SetChipSelectResult(selected []int) {
	p.selectedChips = []battleplayer.SelectChip{}
	for _, s := range selected {
		p.selectedChips = append(
			p.selectedChips,
			battleplayer.SelectChip{
				ID:   p.chipFolder[s].ID,
				Code: p.chipFolder[s].Code,
			},
		)
	}

	// Remove selected chips from folder
	sort.Sort(sort.Reverse(sort.IntSlice(selected)))
	for _, s := range selected {
		p.chipFolder = append(p.chipFolder[:s], p.chipFolder[s+1:]...)
	}
}

func (p *BattlePlayer) GetSelectedChips() []battleplayer.SelectChip {
	return p.selectedChips
}

func (p *BattlePlayer) UpdatePA() {
	// Check program advance
	// TODO
}

func (p *BattlePlayer) IsDead() bool {
	obj := p.GetObject()
	return obj.HP <= 0
}

func (p *BattlePlayer) GetObjectID() string {
	return p.objectID
}

func (p *BattlePlayer) GetObject() gameinfo.Object {
	objs := net.GetInst().GetGameInfo().Objects
	for _, o := range objs {
		if o.ID == p.objectID {
			return o
		}
	}

	return gameinfo.Object{}
}
