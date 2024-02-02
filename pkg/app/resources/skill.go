package resources

const (
	SkillCannon int = iota
	SkillHighCannon
	SkillMegaCannon
	SkillMiniBomb
	SkillSword
	SkillWideSword
	SkillLongSword
	SkillEnemyShockWave
	SkillRecover
	SkillSpreadGun
	SkillVulcan1
	SkillPlayerShockWave
	SkillThunderBall
	SkillPlayerWideShot
	SkillBoomerang
	SkillWaterBomb
	SkillAquamanShot
	SkillAquaman
	SkillCrackout
	SkillDoubleCrack
	SkillTripleCrack
	SkillBambooLance
	SkillDreamSword
	SkillInvisible
	SkillGarooBreath
	SkillFlamePillarRandom
	SkillFlamePillarTracking
	SkillHeatShot
	SkillHeatV
	SkillHeatSide
	SkillFlamePillarLine
	SkillAreaSteal
	SkillPanelSteal
	SkillCountBomb
	SkillTornado
	SkillQuickGauge
	SkillCirkillShot
	SkillEnemyWideShot

	SkillFailed
)

const (
	SkillWaterBombEndCount    = 60
	SkillAreaStealHitEndCount = 12
	SkillHeatShotEndCount     = 15
	SkillFlamePillarEndCount  = 20
)

const (
	SkillBoomerangNextStepCount   = 6
	SkillGarooBreathNextStepCount = 10
	SkillThunderBallNextStepCount = 80
	SkillCirkillShotNextStepCount = 20

	SkillFlamePillarDelay = 4
)

const (
	SkillWideShotStateBegin int = iota
	SkillWideShotStateMove
)

const (
	SkillAquamanStateInit int = iota
	SkillAquamanStateAppear
	SkillAquamanStateCreatePipe
	SkillAquamanStateAttack
)

const (
	SkillAreaStealStateBlackout int = iota
	SkillAreaStealStateActing
	SkillAreaStealStateHit
)

const (
	SkillFlamePillarStateWakeup int = iota
	SkillFlamePillarStateDoing
	SkillFlamePillarStateEnd
	SkillFlamePillarStateDeleted
)

const (
	SkillFlamePillarTypeRandom int = iota
	SkillFlamePillarTypeTracking
	SkillFlamePillarTypeLine
)
