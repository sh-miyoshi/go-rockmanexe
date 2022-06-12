package app

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	netconn "github.com/sh-miyoshi/go-rockmanexe/pkg/app/newnetconn"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/newnet/netconnpb"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/newnet/object"
)

type player struct {
	Object object.Object
	Count  int
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
		Count: 0,
	}

	return res
}

func (p *player) ChipSelect() error {
	n := rand.Intn(2) + 1
	time.Sleep(time.Duration(n) * time.Second)
	p.Object.Chips = []int{1, 3} // debug

	// Finished chip select, so send action
	netconn.GetInst().SendObject(p.Object)

	if err := netconn.GetInst().SendSignal(pb.Action_CHIPSEND); err != nil {
		return fmt.Errorf("failed to get data stream: %w", err)
	}

	return nil
}

func (p *player) Action() bool {
	if p.Object.HP <= 0 {
		// Player deleted
		return true
	}

	p.Count++
	return false
}
