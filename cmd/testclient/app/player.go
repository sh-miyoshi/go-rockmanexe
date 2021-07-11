package app

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/cmd/testclient/netconn"
	"github.com/sh-miyoshi/go-rockmanexe/cmd/testclient/skill"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	appskill "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/field"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/net/routerpb"
)

type player struct {
	Object     field.Object
	Count      int
	ActNo      int
	Act        *Act
	HitDamages map[string]bool // TODO
}

func newPlayer() *player {
	res := &player{
		Object: field.Object{
			ID:            uuid.New().String(),
			Type:          field.ObjectTypeRockmanStand,
			HP:            150,
			X:             1,
			Y:             1,
			DamageChecked: true,
		},
		Count:      0,
		ActNo:      0,
		HitDamages: make(map[string]bool),
	}
	res.Act = NewAct(&res.Object)

	return res
}

func (p *player) ChipSelect() error {
	n := rand.Intn(2) + 1
	time.Sleep(time.Duration(n) * time.Second)
	p.Object.Chips = []int{1, 3} // debug

	// Finished chip select, so send action
	if err := netconn.SendObject(p.Object); err != nil {
		return fmt.Errorf("failed to get data stream: %w", err)
	}

	if err := netconn.SendSignal(pb.Action_CHIPSEND); err != nil {
		return fmt.Errorf("failed to get data stream: %w", err)
	}

	// TODO cleanup hit damages?

	return nil
}

func (p *player) Action() {
	// Damage process
	finfo, _ := netconn.GetFieldInfo()
	for _, obj := range finfo.Objects {
		if obj.ID == p.Object.ID {
			if !obj.DamageChecked {
				if _, exists := p.HitDamages[obj.HitDamage.ID]; exists {
					break
				} else {
					p.HitDamages[obj.HitDamage.ID] = true
				}

				log.Printf("got damage: %+v", obj.HitDamage)

				p.Object.DamageChecked = true
				p.Object.HP -= obj.HitDamage.Power
				if p.Object.HP < 0 {
					p.Object.HP = 0
				}
				// TODO Add damage animation
				netconn.SendObject(p.Object)

				if obj.HitDamage.HitEffectType > 0 {
					log.Printf("add hit damage info: %d", obj.HitDamage.HitEffectType)
					ttl := 0 // TTL = len(images) * delay
					switch obj.HitDamage.HitEffectType {
					case field.ObjectTypeHitSmallEffect:
						ttl = 4
					default:
						panic("not implemented yet")
					}

					eff := field.Object{
						ID:   uuid.New().String(),
						Type: obj.HitDamage.HitEffectType,
						HP:   0,
						X:    p.Object.X,
						Y:    p.Object.Y,
						TTL:  ttl,
					}
					netconn.SendObject(eff)
				}

				return
			}
			break
		}
	}

	if p.Act.Process() {
		return
	}

	actTable := []int{0}
	// Wait, Move, Cannon
	actInterval := []int{60, 30, 120}

	current := actTable[p.ActNo]

	p.Count++
	if p.Count == actInterval[current] {
		p.Count = 0
		p.ActNo = (p.ActNo + 1) % len(actTable)

		// Add action
		log.Printf("Set action %d", current)
		p.Object.UpdateBaseTime = true
		switch current {
		case 1: // Move
			p.Act.Set(battlecommon.PlayerActMove, nil)
		case 2: // Cannon
			p.Act.Set(battlecommon.PlayerActCannon, nil)
			skill.Add(appskill.SkillCannon, skill.Argument{
				X: p.Object.X,
				Y: p.Object.Y,
			})
		}
	}
}
