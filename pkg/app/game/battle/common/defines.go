package common

import "github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"

const (
	GaugeMaxCount               = 1200
	ChargeTime                  = 180 // TODO 変数化
	PlayerDefaultInvincibleTime = 120
	ChargeViewDelay             = 20
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

var FieldNum = common.Point{X: 6, Y: 3}
