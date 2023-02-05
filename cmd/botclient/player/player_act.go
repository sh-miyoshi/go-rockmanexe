package player

import (
	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/field"
	netskill "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/skill"
	netconn "github.com/sh-miyoshi/go-rockmanexe/pkg/app/netconn"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/oldnet/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/oldnet/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/oldnet/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/oldnet/object"
)

type Act interface {
	Process() bool
	Interval() int
}

/*

type ActTemplate struct {
}

func NewActTemplate() *ActTemplate {
	return &ActTemplate{}
}

func (a *ActTemplate) Process() bool {
	return false
}

func (a *ActTemplate) Interval() int {
	return 60
}

*/

type ActWait struct {
	waitFrame int
}

func NewActWait(waitFrame int) *ActWait {
	return &ActWait{
		waitFrame: waitFrame,
	}
}

func (a *ActWait) Process() bool {
	return true
}

func (a *ActWait) Interval() int {
	return a.waitFrame
}

type ActMove struct {
	targetX int
	targetY int
	obj     *object.Object
	conn    *netconn.NetConn
}

func NewActMove(obj *object.Object, targetX, targetY int, conn *netconn.NetConn) *ActMove {
	return &ActMove{
		targetX: targetX,
		targetY: targetY,
		obj:     obj,
		conn:    conn,
	}
}

func (a *ActMove) Process() bool {
	logger.Debug("Move to (%d, %d)", a.targetX, a.targetY)
	a.obj.X = a.targetX
	a.obj.Y = a.targetY
	a.conn.SendObject(*a.obj)
	return true
}

func (a *ActMove) Interval() int {
	return 30
}

type ActBuster struct {
	count     int
	obj       *object.Object
	conn      *netconn.NetConn
	shotPower uint
	charged   bool
}

func NewActBuster(obj *object.Object, conn *netconn.NetConn) *ActBuster {
	return &ActBuster{
		count: 0,
		obj:   obj,
		conn:  conn,

		// debug
		shotPower: 1,
		charged:   false,
	}
}

func (a *ActBuster) Process() bool {
	if a.count == 0 {
		a.obj.UpdateBaseTime = true
		a.obj.Type = object.TypeRockmanBuster
		a.conn.SendObject(*a.obj)
	}

	if a.count == 1 {
		s := a.shotPower
		eff := effect.TypeHitSmallEffect
		if a.charged {
			s *= 10
			eff = effect.TypeHitBigEffect
		}

		y := a.obj.Y
		for x := a.obj.X + 1; x < config.FieldNumX; x++ {
			pos := common.Point{X: x, Y: y}
			if field.GetPanelInfo(pos).ObjectID != "" {
				a.conn.AddDamage(damage.Damage{
					ID:            uuid.New().String(),
					PosX:          pos.X,
					PosY:          pos.Y,
					Power:         int(s),
					TTL:           1,
					TargetType:    damage.TargetOtherClient,
					BigDamage:     a.charged,
					ShowHitArea:   false,
					HitEffectType: eff,
				})
				break
			}
		}
	}

	a.count++
	delay := object.ImageDelays[a.obj.Type]
	num := 6
	return a.count > delay*num
}

func (a *ActBuster) Interval() int {
	return 30
}

type ActSkill struct {
	id        string
	skillType int
	obj       *object.Object
	conn      *netconn.NetConn
}

func NewActSkill(skillType int, obj *object.Object, conn *netconn.NetConn) *ActSkill {
	return &ActSkill{
		skillType: skillType,
		obj:       obj,
		conn:      conn,
	}
}

func (a *ActSkill) Process() bool {
	if a.id == "" {
		a.id = netskill.GetInst().Add(a.skillType, netskill.Argument{
			X:     a.obj.X,
			Y:     a.obj.Y,
			Power: 10, // debug(とりあえず全部10にする)
		})

		a.obj.UpdateBaseTime = true
		switch a.skillType {
		case skill.SkillCannon, skill.SkillHighCannon, skill.SkillMegaCannon:
			a.obj.Type = object.TypeRockmanCannon
		case skill.SkillSword, skill.SkillWideSword, skill.SkillLongSword:
			a.obj.Type = object.TypeRockmanSword
		case skill.SkillVulcan1, skill.SkillWideShot, skill.SkillSpreadGun:
			a.obj.Type = object.TypeRockmanShot
		case skill.SkillPlayerShockWave, skill.SkillShockWave:
			a.obj.Type = object.TypeRockmanPick
		case skill.SkillMiniBomb:
			a.obj.Type = object.TypeRockmanBomb
		default:
			a.obj.Type = object.TypeRockmanStand
		}
		a.conn.SendObject(*a.obj)
		return false
	}

	return !netskill.GetInst().Exists(a.id)
}

func (a *ActSkill) Interval() int {
	return 60
}
