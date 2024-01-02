package resources

const (
	SkillTypeNormalCannon int = iota
	SkillTypeHighCannon
	SkillTypeMegaCannon

	SkillTypeCannonMax
)

const (
	SkillTypeSword int = iota
	SkillTypeWideSword
	SkillTypeLongSword

	SkillTypeSwordMax
)

const (
	SkillRecoverEndCount      = 8
	SkillSpreadGunEndCount    = 8 // imgAtkNum*delay
	SkillSwordEndCount        = 12
	SkillWideShotEndCount     = 16
	SkillWaterBombEndCount    = 60
	SkillAreaStealHitEndCount = 12
	SkillHeatShotEndCount     = 15
	SkillFlamePillarEndCount  = 20
)

const (
	SkillShockWaveInitWait    = 9
	SkillShockWavePlayerSpeed = 3
	SkillShockWaveImageNum    = 7

	SkillWideShotPlayerNextStepCount = 8
	SkillBoomerangNextStepCount      = 6
	SkillGarooBreathNextStepCount    = 10
	SkillThunderBallNextStepCount    = 80
	SkillCirkillShotNextStepCount    = 20

	SkillSwordDelay       = 3
	SkillVulcanDelay      = 2
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
