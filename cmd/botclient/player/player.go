package player

import (
	"math/rand"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/netconn"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/net/netconnpb"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type Player struct {
	ID  string
	HP  int
	Pos point.Point

	currentActNo       int
	currentActInterval int
	actTable           []Act
	conn               *netconn.NetConn
	clientID           string
}

func New(clientID string, conn *netconn.NetConn) *Player {
	res := &Player{
		ID:                 uuid.New().String(),
		HP:                 300,
		Pos:                point.Point{X: 1, Y: 1},
		currentActNo:       0,
		currentActInterval: 0,
		conn:               conn,
		clientID:           clientID,
	}
	res.initActTable()

	return res
}

func (p *Player) ChipSelect() error {
	n := rand.Intn(2) + 1
	time.Sleep(time.Duration(n) * time.Second)
	// TODO: 選択したチップを送る

	if err := p.conn.SendSignal(pb.Request_CHIPSELECT, nil); err != nil {
		return errors.Wrap(err, "failed to get data stream")
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
		// NewActSkill(chip.IDCannon, p.clientID, p.conn),
		// NewActSkill(chip.IDHighCannon, p.clientID, p.conn),
		// NewActSkill(chip.IDMegaCannon, p.clientID, p.conn),
		// NewActSkill(chip.IDMiniBomb, p.clientID, p.conn),
		// NewActSkill(chip.IDRecover10, p.clientID, p.conn),
		// NewActSkill(chip.IDRecover30, p.clientID, p.conn),
		// NewActSkill(chip.IDShockWave, p.clientID, p.conn),
		// NewActSkill(chip.IDSpreadGun, p.clientID, p.conn),
		// NewActSkill(chip.IDSword, p.clientID, p.conn),
		// NewActSkill(chip.IDWideSword, p.clientID, p.conn),
		// NewActSkill(chip.IDLongSword, p.clientID, p.conn),
		// NewActSkill(chip.IDVulcan1, p.clientID, p.conn),
		// NewActSkill(chip.IDWideShot, p.clientID, p.conn),
		// NewActSkill(chip.IDHeatShot, p.clientID, p.conn),
		// NewActSkill(chip.IDHeatV, p.clientID, p.conn),
		// NewActSkill(chip.IDHeatSide, p.clientID, p.conn),
		// NewActSkill(chip.IDFlameLine1, p.clientID, p.conn),
		// NewActSkill(chip.IDFlameLine2, p.clientID, p.conn),
		// NewActSkill(chip.IDFlameLine3, p.clientID, p.conn),
		// NewActSkill(chip.IDTornado, p.clientID, p.conn),
		// NewActSkill(chip.IDBoomerang1, p.clientID, p.conn),
		// NewActSkill(chip.IDBambooLance, p.clientID, p.conn),
		// NewActSkill(chip.IDCrackout, p.clientID, p.conn),
		// NewActSkill(chip.IDDoubleCrack, p.clientID, p.conn),
		// NewActSkill(chip.IDTripleCrack, p.clientID, p.conn),
		// NewActSkill(chip.IDAreaSteal, p.clientID, p.conn),
	}
	p.currentActNo = 0
	p.currentActInterval = p.actTable[0].Interval()
}
