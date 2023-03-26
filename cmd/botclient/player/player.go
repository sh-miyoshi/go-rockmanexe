package player

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	netconn "github.com/sh-miyoshi/go-rockmanexe/pkg/app/newnetconn"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/newnet/netconnpb"
)

type Player struct {
	ID  string
	HP  int
	Pos common.Point

	currentActNo       int
	currentActInterval int
	actTable           []Act
	conn               *netconn.NetConn
}

func New(clientID string, conn *netconn.NetConn) *Player {
	res := &Player{
		ID:                 uuid.New().String(),
		HP:                 10,
		Pos:                common.Point{X: 1, Y: 1},
		currentActNo:       0,
		currentActInterval: 0,
		conn:               conn,
	}
	res.initActTable()

	return res
}

func (p *Player) ChipSelect() error {
	n := rand.Intn(2) + 1
	time.Sleep(time.Duration(n) * time.Second)
	// TODO: 選択したチップを送る

	if err := p.conn.SendSignal(pb.Request_CHIPSELECT, nil); err != nil {
		return fmt.Errorf("failed to get data stream: %w", err)
	}

	return nil
}

func (p *Player) Action() bool {
	if p.HP <= 0 {
		// Player deleted
		return true
	}

	if p.currentActInterval > 0 {
		p.currentActInterval--
		return false
	}

	if p.actTable[p.currentActNo].Process() {
		logger.Info("finished process %d", p.currentActNo)
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
		NewActMove(0, 1, p.conn),
		// NewActSkill(skill.SkillPlayerShockWave, &p.Object),
		// NewActSkill(skill.SkillSpreadGun, &p.Object),
		// NewActSkill(skill.SkillSword, &p.Object),
		// NewActSkill(skill.SkillThunderBall, &p.Object),
		// NewActSkill(skill.SkillVulcan1, &p.Object),
		// NewActSkill(skill.SkillWideShot, &p.Object),
		// NewActMove(&p.Object, 0, 1),
		// NewActBuster(&p.Object),
		// NewActSkill(skill.SkillRecover, &p.Object),
	}
	p.currentActNo = 0
	p.currentActInterval = p.actTable[0].Interval()
}
