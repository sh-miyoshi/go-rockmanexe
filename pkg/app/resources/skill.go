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
	SkillCannonEndCount    = 31 // imgAtkNum*delayAtk + 15
	SkillMiniBombEndCount  = 60
	SkillRecoverEndCount   = 8
	SkillSpreadGunEndCount = 8 // imgAtkNum*delay
	SkillSwordEndCount     = 12
	SkillWideShotEndCount  = 16
)

const (
	SkillShockWaveInitWait           = 9
	SkillShockWavePlayerSpeed        = 3
	SkillShockWaveImageNum           = 7
	SkillSwordDelay                  = 3
	SkillVulcanDelay                 = 2
	SkillWideShotDelay               = 4
	SkillWideShotPlayerNextStepCount = 8
)

const (
	SkillWideShotStateBegin int = iota
	SkillWideShotStateMove
)
