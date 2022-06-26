package player

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	netconn "github.com/sh-miyoshi/go-rockmanexe/pkg/app/netconn"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/net/netconnpb"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/object"
)

type Player struct {
	Object             object.Object
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
		NewActWait(150),
		NewActMove(&p.Object, 0, 1),
		NewActBuster(&p.Object),
	}
	p.currentActNo = 0
	p.currentActInterval = p.actTable[0].Interval()
}
