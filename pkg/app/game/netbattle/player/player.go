package player

import (
	"fmt"
	"math/rand"
	"sort"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/net"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/player"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/inputs"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/newnet/action"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/newnet/netconnpb"
)

type BattlePlayer struct {
	objectID      string
	chipFolder    []player.ChipInfo
	selectedChips []player.ChipInfo
	imgHPFrame    int
	imgGaugeFrame int
	imgGaugeMax   []int
	imgMinds      []int
	imgMindFrame  int
	imgCharge     [2][]int

	chargeCount int
	shotPower   int
}

func New(plyr *player.Player) (*BattlePlayer, error) {
	res := &BattlePlayer{
		objectID:    uuid.New().String(),
		chargeCount: 0,
		shotPower:   1,
	}
	for _, c := range plyr.ChipFolder {
		res.chipFolder = append(res.chipFolder, c)
	}
	// Shuffle folder
	for i := 0; i < 10; i++ {
		for j := 0; j < len(res.chipFolder); j++ {
			n := rand.Intn(len(res.chipFolder))
			res.chipFolder[j], res.chipFolder[n] = res.chipFolder[n], res.chipFolder[j]
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
	// TODO: draw hp

	// Show Mind Status
	dxlib.DrawGraph(x, 40, p.imgMindFrame, true)
	dxlib.DrawGraph(x, 40, p.imgMinds[battlecommon.PlayerMindStatusNormal], true) // TODO set mind status

	// TODO: Show Custom Gauge
}

func (p *BattlePlayer) LocalDraw() {
	// TODO
}

func (p *BattlePlayer) Process() (bool, error) {
	// TODO

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
			sound.On(sound.SEBusterCharging)
		}
		if p.chargeCount == battlecommon.ChargeTime {
			sound.On(sound.SEBusterCharged)
		}
	} else if p.chargeCount > 0 {
		sound.On(sound.SEBusterShot)
		charged := p.chargeCount > battlecommon.ChargeTime
		power := p.shotPower
		if charged {
			power *= 10
		}

		buster := action.Buster{
			Power: power,
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
	p.selectedChips = []player.ChipInfo{}
	for _, s := range selected {
		p.selectedChips = append(p.selectedChips, p.chipFolder[s])
	}

	// Remove selected chips from folder
	sort.Sort(sort.Reverse(sort.IntSlice(selected)))
	for _, s := range selected {
		p.chipFolder = append(p.chipFolder[:s], p.chipFolder[s+1:]...)
	}
}

func (p *BattlePlayer) GetSelectedChips() []player.ChipInfo {
	return p.selectedChips
}

func (p *BattlePlayer) UpdatePA() {
	// Check program advance
	// TODO
}

func (p *BattlePlayer) IsDead() bool {
	// TODO
	return false
}
