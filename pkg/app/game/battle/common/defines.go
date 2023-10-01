package common

import "github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"

const (
	GaugeMaxCount   = 1200
	ChargeViewDelay = 20
)

var (
	PlayerDefaultInvincibleTime = 120
	DefaultParalyzedTime        = 120
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
	FieldNum      = common.Point{X: 6, Y: 3}
	PanelSize     = common.Point{X: 80, Y: 50}
	DrawPanelTopY = common.ScreenSize.Y - (PanelSize.Y * FieldNum.Y) - 30
)

const (
	PanelStatusNormal int = iota
	PanelStatusCrack
	PanelStatusHole

	PanelStatusMax
)

const (
	PanelTypePlayer int = iota
	PanelTypeEnemy

	PanelTypeMax
)

type PanelInfo struct {
	Type      int
	ObjectID  string
	Status    int
	HoleCount int
}
