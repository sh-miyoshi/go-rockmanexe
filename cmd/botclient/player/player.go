package player

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill"
	netconn "github.com/sh-miyoshi/go-rockmanexe/pkg/app/netconn"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/effect"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/net/netconnpb"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/object"
)

type Player struct {
	Object             object.Object
	HitDamages         map[string]bool
	currentActNo       int
	currentActInterval int
	actTable           []Act
}

func New(clientID string) *Player {
	res := &Player{
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
		currentActNo:       0,
		currentActInterval: 0,
		HitDamages:         make(map[string]bool),
	}
	res.initActTable()

	return res
}

func (p *Player) ChipSelect() error {
	n := rand.Intn(2) + 1
	time.Sleep(time.Duration(n) * time.Second)
	p.Object.Chips = []object.ChipInfo{
		{ID: 1, Code: "*"},
		{ID: 3, Code: "a"},
	}

	// Finished chip select, so send action
	netconn.GetInst().SendObject(p.Object)

	if err := netconn.GetInst().SendSignal(pb.Action_CHIPSEND); err != nil {
		return fmt.Errorf("failed to get data stream: %w", err)
	}

	return nil
}

func (p *Player) Action() bool {
	if p.Object.HP <= 0 {
		// Player deleted
		return true
	}

	if p.damageProc() {
		return false
	}

	if p.currentActInterval > 0 {
		p.currentActInterval--
		return false
	}

	if p.actTable[p.currentActNo].Process() {
		logger.Info("finished process %d", p.currentActNo)
		p.Object.UpdateBaseTime = true
		p.Object.Type = object.TypeRockmanStand
		netconn.GetInst().SendObject(p.Object)

		p.currentActNo++
		if p.currentActNo >= len(p.actTable) {
			p.initActTable()
			return false
		}
		p.currentActInterval = p.actTable[p.currentActNo].Interval()
	}

	return false
}

func (p *Player) initActTable() {
	logger.Info("initialize player act table")

	p.actTable = []Act{
		NewActWait(30),
		NewActSkill(skill.SkillPlayerShockWave, &p.Object),
		NewActSkill(skill.SkillSpreadGun, &p.Object),
		NewActSkill(skill.SkillSword, &p.Object),
		NewActSkill(skill.SkillThunderBall, &p.Object),
		NewActSkill(skill.SkillVulcan1, &p.Object),
		NewActSkill(skill.SkillWideShot, &p.Object),
		// NewActMove(&p.Object, 0, 1),
		// NewActBuster(&p.Object),
		// NewActSkill(skill.SkillRecover, &p.Object),
	}
	p.currentActNo = 0
	p.currentActInterval = p.actTable[0].Interval()
}

func (p *Player) damageProc() bool {
	ginfo := netconn.GetInst().GetGameInfo()
	if len(ginfo.HitDamages) == 0 {
		return false
	}

	dm := ginfo.HitDamages[0]
	defer netconn.GetInst().RemoveDamage(dm.ID)

	if _, exists := p.HitDamages[dm.ID]; exists {
		return false
	} else {
		p.HitDamages[dm.ID] = true
	}

	// Recover系は使えるようにする
	if p.Object.Invincible && dm.Power >= 0 {
		return false
	}

	logger.Debug("Got damage: %+v", dm)

	p.Object.HP -= dm.Power
	if p.Object.HP < 0 {
		p.Object.HP = 0
	}
	// debug(回復は考慮しない)

	if dm.BigDamage {
		// p.Object.Invincible = true // インビジ未実装
		// TODO Skill関係
		// for _, sid := range p.ManagedSkills {
		// 	netskill.StopByPlayer(sid)
		// }
		// p.ManagedSkills = []string{}
		// netconn.GetInst().AddSound(sound.SEDamaged)
	} else {
		netconn.GetInst().SendObject(p.Object)
	}

	if dm.HitEffectType > 0 {
		logger.Info("Add hit effect %d", dm.HitEffectType)
		netconn.GetInst().SendEffect(effect.Effect{
			ID:       uuid.New().String(),
			Type:     dm.HitEffectType,
			X:        p.Object.X,
			Y:        p.Object.Y,
			ViewOfsX: dm.ViewOfsX,
			ViewOfsY: dm.ViewOfsY,
		})
	}

	return true
}
