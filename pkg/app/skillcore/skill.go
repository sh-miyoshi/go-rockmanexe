package skillcore

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/chip"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/system"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type DamageManager interface {
	New(dm damage.Damage) string
	Exists(id string) bool
}

type Argument struct {
	OwnerID       string
	OwnerClientID string
	Power         uint
	TargetType    int
	IsReverse     bool

	DamageMgr         DamageManager
	GetPanelInfo      func(pos point.Point) battlecommon.PanelInfo
	GetObjectPos      func(objID string) point.Point
	SoundOn           func(typ resources.SEType)
	GetObjects        func(filter objanim.Filter) []objanim.Param
	Cutin             func(skillName string, count int)
	ChangePanelStatus func(pos point.Point, pnStatus int, endCount int)
	ChangePanelType   func(pos point.Point, pnType int, endCount int)
	MakeInvisible     func(objID string, count int)
	AddBarrier        func(objID string, hp int)
}

type SkillCore interface {
	Update() (bool, error)
	GetCount() int
}

func GetIDByChipID(chipID int) int {
	switch chipID {
	case chip.IDCannon:
		return resources.SkillCannon
	case chip.IDHighCannon:
		return resources.SkillHighCannon
	case chip.IDMegaCannon:
		return resources.SkillMegaCannon
	case chip.IDSword:
		return resources.SkillSword
	case chip.IDWideSword:
		return resources.SkillWideSword
	case chip.IDLongSword:
		return resources.SkillLongSword
	case chip.IDMiniBomb:
		return resources.SkillMiniBomb
	case chip.IDRecover10, chip.IDRecover30, chip.IDRecover50, chip.IDRecover80, chip.IDRecover120, chip.IDRecover150, chip.IDRecover200, chip.IDRecover300:
		return resources.SkillRecover
	case chip.IDSpreadGun:
		return resources.SkillSpreadGun
	case chip.IDVulcan1:
		return resources.SkillVulcan1
	case chip.IDVulcan2:
		return resources.SkillVulcan2
	case chip.IDVulcan3:
		return resources.SkillVulcan3
	case chip.IDShockWave:
		return resources.SkillPlayerShockWave
	case chip.IDThunderBall1, chip.IDThunderBall2, chip.IDThunderBall3:
		return resources.SkillThunderBall
	case chip.IDWideShot1, chip.IDWideShot2, chip.IDWideShot3:
		return resources.SkillPlayerWideShot
	case chip.IDBoomerang1:
		return resources.SkillBoomerang
	case chip.IDAquaman:
		return resources.SkillAquaman
	case chip.IDCrackout:
		return resources.SkillCrackout
	case chip.IDDoubleCrack:
		return resources.SkillDoubleCrack
	case chip.IDTripleCrack:
		return resources.SkillTripleCrack
	case chip.IDBambooLance:
		return resources.SkillBambooLance
	case chip.IDDreamSword:
		return resources.SkillDreamSword
	case chip.IDInvisible:
		return resources.SkillInvisible
	case chip.IDHeatShot:
		return resources.SkillHeatShot
	case chip.IDHeatV:
		return resources.SkillHeatV
	case chip.IDHeatSide:
		return resources.SkillHeatSide
	case chip.IDFlameLine1, chip.IDFlameLine2, chip.IDFlameLine3:
		return resources.SkillFlamePillarLine
	case chip.IDAreaSteal:
		return resources.SkillAreaSteal
	case chip.IDPanelSteal:
		return resources.SkillPanelSteal
	case chip.IDCountBomb:
		return resources.SkillCountBomb
	case chip.IDTornado:
		return resources.SkillTornado
	case chip.IDAttack10:
		return resources.SkillFailed
	case chip.IDQuickGauge:
		return resources.SkillQuickGauge
	case chip.IDBubbleShot:
		return resources.SkillBubbleShot
	case chip.IDBubbleV:
		return resources.SkillBubbleV
	case chip.IDBubbleSide:
		return resources.SkillBubbleSide
	case chip.IDForteAnother:
		return resources.SkillChipForteAnother
	case chip.IDDeathMatch1:
		return resources.SkillDeathMatch1
	case chip.IDDeathMatch2:
		return resources.SkillDeathMatch2
	case chip.IDDeathMatch3:
		return resources.SkillDeathMatch3
	case chip.IDPanelReturn:
		return resources.SkillPanelReturn
	case chip.IDBarrier:
		return resources.SkillBarrier
	case chip.IDBarrier100:
		return resources.SkillBarrier100
	case chip.IDBarrier200:
		return resources.SkillBarrier200
	}

	system.SetError(fmt.Sprintf("Skill for Chip %d is not implemented yet", chipID))
	return 0
}
