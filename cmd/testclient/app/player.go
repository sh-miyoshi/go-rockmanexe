package app

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/cmd/testclient/skill"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	appskill "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/netconn"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/object"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/net/routerpb"
)

type player struct {
	Object     object.Object
	Count      int
	ActNo      int
	Act        *Act
	HitDamages map[string]bool
}

func newPlayer(clientID string) *player {
	res := &player{
		Object: object.Object{
			ID:             uuid.New().String(),
			ClientID:       clientID,
			Type:           object.TypeRockmanStand,
			HP:             150,
			X:              1,
			Y:              1,
			Hittable:       true,
			UpdateBaseTime: true,
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
	netconn.SendObject(p.Object)

	if err := netconn.SendSignal(pb.Action_CHIPSEND); err != nil {
		return fmt.Errorf("failed to get data stream: %w", err)
	}

	// TODO cleanup hit damages?

	return nil
}

func (p *player) Action() bool {
	if p.Object.HP <= 0 {
		// Player deleted
		return true
	}

	if p.damageProc() {
		return false
	}

	if p.Act.Process() {
		return false
	}

	actTable := []int{2, 2}
	// Wait, Move, Cannon, Buster
	actInterval := []int{60, 60, 120, 60}

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
			}, p.Object.ClientID)
		case 3:
			p.Act.Set(battlecommon.PlayerActBuster, nil)
		}
	}

	return false
}

func (p *player) damageProc() bool {
	finfo, _ := netconn.GetFieldInfo()
	if finfo.HitDamage.ID == "" {
		return false
	}

	if p.Object.Invincible {
		return false
	}

	log.Printf("got damage: %+v", finfo.HitDamage)

	if _, exists := p.HitDamages[finfo.HitDamage.ID]; exists {
		return false
	} else {
		p.HitDamages[finfo.HitDamage.ID] = true
	}

	p.Object.HP -= finfo.HitDamage.Power
	if p.Object.HP < 0 {
		p.Object.HP = 0
	}

	// TODO Add damage animation
	netconn.SendObject(p.Object)

	log.Printf("add hit damage effect: %d", finfo.HitDamage.HitEffectType)
	netconn.SendEffect(effect.Effect{
		ID:       uuid.New().String(),
		ClientID: p.Object.ClientID,
		Type:     finfo.HitDamage.HitEffectType,
		X:        p.Object.X,
		Y:        p.Object.Y,
		ViewOfsX: finfo.HitDamage.ViewOfsX,
		ViewOfsY: finfo.HitDamage.ViewOfsY,
	})

	netconn.RemoveDamage()
	return true
}
