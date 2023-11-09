package resources

type SEType int

const (
	SENone SEType = iota // required no SE as 0
	SETitleEnter
	SECursorMove
	SEMenuEnter
	SEDenied
	SECancel
	SEGoBattle
	SEEnemyAppear
	SEChipSelectOpen
	SESelect
	SEChipSelectEnd
	SEGaugeMax
	SECannon
	SEBusterCharging
	SEBusterCharged
	SEBusterShot
	SECannonHit
	SEExplode
	SEBusterHit
	SESword
	SERecover
	SEShockWave
	SEGun
	SESpreadHit
	SEBombThrow
	SEPlayerDeleted
	SEDamaged
	SEEnemyDeleted
	SEGotItem
	SEWindowChange
	SEThunderBall
	SEWideShot
	SEBoomerangThrow
	SEWaterLanding
	SEBlock
	SEObjectCreate
	SEWaterpipeAttack
	SEPanelBreak
	SEPAPrepare
	SEPACreated
	SEDreamSword
	SEFlameAttack
	SEAreaSteal
	SEAreaStealHit
	SERunOK
	SECountBombCountdown
	SECountBombEnd
	SETornado

	SEMax
)
