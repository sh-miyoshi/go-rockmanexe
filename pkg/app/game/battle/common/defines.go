package common

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	GaugeMaxCount              = 1200
	ChargeViewDelay            = 20
	DefaultCustomGaugeSpeed    = 4
	DefaultPanelTypeEndCount   = 480
	DefaultPanelStatusEndCount = 480
	PanelReturnAnimCount       = 60
)

var (
	PlayerDefaultInvincibleTime = 120
	DefaultParalyzedTime        = 60
	CustomGaugeSpeed            = 4
)

const (
	PlayerActMove int = iota
	PlayerActDamage
	PlayerActShot
	PlayerActCannon
	PlayerActSword
	PlayerActBomb
	PlayerActBuster
	PlayerActPick
	PlayerActThrow
	PlayerActParalyzed
	PlayerActBShot

	PlayerActMax
)

const (
	PlayerMindStatusFullSync int = iota
	PlayerMindStatusAnger
	PlayerMindStatusNormal
	PlayerMindStatusFear
	PlayerMindStatusDark
	PlayerMindStatusRollSoul
	PlayerMindStatusAquaSoul
	PlayerMindStatusWoodSoul
	PlayerMindStatusJunkSoul
	PlayerMindStatusBluesSoul
	PlayerMindStatusMetalSoul
	PlayerMindStatusGutsSoul
	PlayerMindStatusSearchSoul
	PlayerMindStatusNumberSoul
	PlayerMindStatusFireSoul
	PlayerMindStatusWindSoul
	PlayerMindStatusThunderSoul

	PlayerMindStatusMax
)

var (
	FieldNum      = point.Point{X: 6, Y: 3}
	PanelSize     = point.Point{X: 80, Y: 50}
	DrawPanelTopY = config.ScreenSize.Y - (PanelSize.Y * FieldNum.Y) - 30
)

const (
	PanelStatusNormal int = iota
	PanelStatusCrack
	PanelStatusHole
	PanelStatusPoison
	PanelStatusEmpty

	PanelStatusMax
)

const (
	PanelTypePlayer int = iota
	PanelTypeEnemy

	PanelTypeMax
)

type PanelInfo struct {
	Type     int
	ObjectID string
	Status   int
}
