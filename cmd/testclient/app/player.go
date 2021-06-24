package app

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/cmd/testclient/netconn"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/field"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/net/routerpb"
)

type player struct {
	Object field.Object
	Count  int
	ActNo  int
	Act    *Act
}

func newPlayer() *player {
	res := &player{
		Object: field.Object{
			ID:   uuid.New().String(),
			Type: field.ObjectTypeRockmanStand,
			HP:   150,
			X:    0,
			Y:    1,
		},
		Count: 0,
		ActNo: 0,
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

	return nil
}

func (p *player) Action() {
	if p.Act.Process() {
		return
	}

	actTable := []int{0, 1, 1, 2}
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
			// TODO
		}
	}
}
