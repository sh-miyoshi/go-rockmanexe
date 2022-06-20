package player

import (
	"fmt"
	"math/rand"
	"sort"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	appfield "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	netdraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/draw"
	netfield "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/inputs"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/player"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/newnet/object"
)

type BattlePlayer struct {
	Object     object.Object
	ChipFolder []player.ChipInfo
	GaugeCount uint
	Act        *Act

	imgHPFrame    int
	imgGaugeFrame int
	imgGaugeMax   []int
	imgMinds      []int
	imgMindFrame  int
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
		// TODO
		// HPMax:      plyr.HP,
		// ShotPower:  plyr.ShotPower,
		// HitDamages: make(map[string]bool),
	}

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

	logger.Info("Successfully initialized net battle player data")
	return &res, nil
}

func (p *BattlePlayer) InitAct(drawMgr *netdraw.DrawManager) {
	p.Act = NewAct(drawMgr, &p.Object)
}

func (p *BattlePlayer) End() {
	// TODO imageの解放
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

func (p *BattlePlayer) Process() (bool, error) {
	p.GaugeCount += 4 // TODO GaugeSpeed

	if p.Object.HP <= 0 {
		return true, nil
	}

	if p.Act.Process() {
		return false, nil
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

	return false, nil
}

func (p *BattlePlayer) DamageProc(dm *damage.Damage) bool {
	return false
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
