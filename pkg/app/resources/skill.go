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
	SkillVulcan2
	SkillVulcan3
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
	SkillShrimpyAttack
	SkillBubbleShot
	SkillBubbleV
	SkillBubbleSide
	SkillForteHellsRollingUp
	SkillForteHellsRollingDown
	SkillForteDarkArmBladeType1
	SkillForteDarkArmBladeType2
	SkillForteShootingBuster
	SkillForteDarknessOverload

	SkillFailed
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
	SkillShrimpyAttackStateBegin int = iota
	SkillShrimpyAttackStateMove
)
