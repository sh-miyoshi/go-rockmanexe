package resources

const (
	SkillCannon int = iota
	SkillHighCannon
	SkillMegaCannon
	SkillMiniBomb
	SkillSword
	SkillWideSword
	SkillLongSword
	SkillFighterSword
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
	SkillChipForteAnother
	SkillDeathMatch1
	SkillDeathMatch2
	SkillDeathMatch3
	SkillNonEffectWideSword
	SkillPanelReturn
	SkillBarrier
	SkillBarrier100
	SkillBarrier200

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

const (
	SkillChipForteAnotherStateInit int = iota
	SkillChipForteAnotherStateAppear
	SkillChipForteAnotherStateAttack
	SkillChipForteAnotherStateEnd
)
