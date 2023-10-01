package resources

const (
	SkillTypeNormalCannon int = iota
	SkillTypeHighCannon
	SkillTypeMegaCannon

	SkillTypeCannonMax
)

const (
	SkillCannonEndCount   = 31 // imgAtkNum*delayAtk + 15
	SkillMiniBombEndCount = 60
	SkillRecoverEndCount  = 8
)

const (
	SkillShockWaveInitWait    = 9
	SkillShockWavePlayerSpeed = 3
	SkillShockWaveImageNum    = 7
)
