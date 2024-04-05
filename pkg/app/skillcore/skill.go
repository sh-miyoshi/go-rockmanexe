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

	DamageMgr       DamageManager
	GetPanelInfo    func(pos point.Point) battlecommon.PanelInfo
	GetObjectPos    func(objID string) point.Point
	SoundOn         func(typ resources.SEType)
	GetObjects      func(filter objanim.Filter) []objanim.Param
	PanelBreak      func(pos point.Point)
	Cutin           func(skillName string, count int)
	ChangePanelType func(pos point.Point, pnType int)
}

type SkillCore interface {
	Process() (bool, error)
	GetCount() int
	GetEndCount() int
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
	case chip.IDRecover10:
		return resources.SkillRecover
	case chip.IDRecover30:
		return resources.SkillRecover
	case chip.IDSpreadGun:
		return resources.SkillSpreadGun
	case chip.IDVulcan1:
		return resources.SkillVulcan1
	case chip.IDShockWave:
		return resources.SkillPlayerShockWave
	case chip.IDThunderBall:
		return resources.SkillThunderBall
	case chip.IDWideShot:
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
	}

	system.SetError(fmt.Sprintf("Skill for Chip %d is not implemented yet", chipID))
	return 0
}
