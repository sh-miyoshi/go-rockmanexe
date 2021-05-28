package player

import (
	"fmt"
	"math/rand"
	"sort"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/player"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/field"
)

type BattlePlayer struct {
	Object      field.Object
	HPMax       uint
	ChargeCount uint
	GaugeCount  uint
	ShotPower   uint
	ChipFolder  []player.ChipInfo
}

var (
	// imgPlayers    [playerAnimMax][]int32
	// imgDelays     = [playerAnimMax]int{1, 2, 2, 6, 3, 4, 1, 4}
	imgHPFrame    int32
	imgGaugeFrame int32
	imgGaugeMax   []int32
	// imgCharge     [2][]int32
)

// New ...
func New(plyr *player.Player) (*BattlePlayer, error) {
	logger.Info("Initialize net battle player data")

	res := BattlePlayer{
		Object: field.Object{
			ID: uuid.New().String(),
			HP: int(plyr.HP),
			X:  1,
			Y:  1,
		},
		HPMax:     plyr.HP, // TODO HPは引き継がない
		ShotPower: plyr.ShotPower,
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
	imgGaugeMax = make([]int32, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 1, 4, 288, 30, imgGaugeMax); res == -1 {
		return nil, fmt.Errorf("failed to read gauge max image %s", fname)
	}

	// TODO

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
}

func (p *BattlePlayer) Draw() {
}

func (p *BattlePlayer) DrawFrame(xShift bool, showGauge bool) {
	x := int32(7)
	y := int32(5)
	if xShift {
		x += 235
	}

	// Show HP
	dxlib.DrawGraph(x, y, imgHPFrame, dxlib.TRUE)
	draw.Number(x+2, y+2, int32(p.Object.HP), draw.NumberOption{RightAligned: true, Length: 4})

	// Show Custom Gauge
	if showGauge {
		if p.GaugeCount < common.BattleGaugeMaxCount {
			dxlib.DrawGraph(96, 5, imgGaugeFrame, dxlib.TRUE)
			const gaugeMaxSize = 256
			size := int32(gaugeMaxSize * p.GaugeCount / common.BattleGaugeMaxCount)
			dxlib.DrawBox(112, 19, 112+size, 21, dxlib.GetColor(123, 154, 222), dxlib.TRUE)
			dxlib.DrawBox(112, 21, 112+size, 29, dxlib.GetColor(231, 235, 255), dxlib.TRUE)
			dxlib.DrawBox(112, 29, 112+size, 31, dxlib.GetColor(123, 154, 222), dxlib.TRUE)
		} else {
			i := (p.GaugeCount / 20) % 4
			dxlib.DrawGraph(96, 5, imgGaugeMax[i], dxlib.TRUE)
		}
	}
}

func (p *BattlePlayer) Process() (bool, error) {
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
