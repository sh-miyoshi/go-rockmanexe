package skill

import (
	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	routeranim "github.com/sh-miyoshi/go-rockmanexe/pkg/router/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/router/gameinfo"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/router/queue"
)

type vulcan struct {
	ID    string
	Arg   Argument
	Times int

	count    int
	atkCount int
	hit      bool
}

const (
	delayVulcan = 2
)

func newVulcan(times int, arg Argument) *vulcan {
	return &vulcan{
		ID:    arg.AnimObjID,
		Arg:   arg,
		Times: times,
	}
}

func (p *vulcan) Draw() {
	// nothing to do at router
}

func (p *vulcan) Process() (bool, error) {
	p.count++
	if p.count >= delayVulcan*1 {
		if p.count%(delayVulcan*5) == delayVulcan*1 {
			// Add damage
			pos := routeranim.ObjAnimGetObjPos(p.Arg.OwnerClientID, p.Arg.OwnerObjectID)
			hit := false
			p.atkCount++
			lastAtk := p.atkCount == p.Times
			for x := pos.X + 1; x < battlecommon.FieldNum.X; x++ {
				target := common.Point{X: x, Y: pos.Y}
				if p.Arg.GameInfo.GetPanelInfo(target).ObjectID != "" {
					routeranim.DamageNew(p.Arg.OwnerClientID, damage.Damage{
						Pos:           target,
						Power:         int(p.Arg.Power),
						TTL:           1,
						TargetType:    p.Arg.TargetType,
						HitEffectType: resources.EffectTypeSpreadHit,
						BigDamage:     lastAtk,
						DamageType:    damage.TypeNone,
					})
					queue.Push(p.Arg.QueueIDs[queue.TypeEffect], &gameinfo.Effect{
						ID:            uuid.New().String(),
						OwnerClientID: p.Arg.GameInfo.ClientID,
						Pos:           target,
						Type:          resources.EffectTypeVulcanHit1,
						RandRange:     20,
					})

					if p.hit && x < battlecommon.FieldNum.X-1 {
						target = common.Point{X: x + 1, Y: pos.Y}
						queue.Push(p.Arg.QueueIDs[queue.TypeEffect], &gameinfo.Effect{
							ID:            uuid.New().String(),
							OwnerClientID: p.Arg.GameInfo.ClientID,
							Pos:           target,
							Type:          resources.EffectTypeVulcanHit2,
							RandRange:     20,
						})

						routeranim.DamageNew(p.Arg.OwnerClientID, damage.Damage{
							Pos:           target,
							Power:         int(p.Arg.Power),
							TTL:           1,
							TargetType:    p.Arg.TargetType,
							HitEffectType: resources.EffectTypeNone,
							BigDamage:     lastAtk,
							DamageType:    damage.TypeNone,
						})
					}
					hit = true
					break
				}
			}
			p.hit = hit
			if lastAtk {
				return true, nil
			}
		}

	}

	return false, nil
}

func (p *vulcan) GetParam() anim.Param {
	info := routeranim.NetInfo{
		OwnerClientID: p.Arg.OwnerClientID,
		AnimType:      routeranim.TypeVulcan,
		ActCount:      p.count,
	}

	return anim.Param{
		ObjID:     p.ID,
		DrawType:  anim.DrawTypeSkill,
		Pos:       routeranim.ObjAnimGetObjPos(p.Arg.OwnerClientID, p.Arg.OwnerObjectID),
		ExtraInfo: info.Marshal(),
	}
}

func (p *vulcan) StopByOwner() {
	routeranim.AnimDelete(p.Arg.OwnerClientID, p.ID)
}

func (p *vulcan) GetEndCount() int {
	return delayVulcan*5*(p.Times-1) + delayVulcan*1
}
