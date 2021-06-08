package common

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

	PlayerActMax
)
