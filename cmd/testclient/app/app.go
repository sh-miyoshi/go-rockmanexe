package app

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/cmd/testclient/netconn"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/field"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/net/routerpb"
)

const (
	statusWaiting int = iota
	statusChipSelect
	statusWaitActing
	statusActing
)

var (
	playerStatus = statusWaiting
	playerObject field.Object

	actCount = 0
)

func PlayerInit() error {
	playerObject = field.Object{
		ID:   uuid.New().String(),
		Type: field.ObjectTypeRockmanStand,
		HP:   150,
		X:    0,
		Y:    1,
	}

	return nil
}

func PlayerProc(exitErr chan error) {
	// set init data to router
	netconn.SendObject(playerObject)

	for {
		status, err := netconn.GetStatus()
		if err != nil {
			exitErr <- fmt.Errorf("got status failed: %v", err)
			return
		}
		playerStatusUpdate(status)

		switch playerStatus {
		case statusWaiting:
			// nothing to do
		case statusChipSelect:
			// Select using chip
			n := rand.Intn(2) + 1
			time.Sleep(time.Duration(n) * time.Second)
			playerObject.Chips = []int{1, 3} // debug

			// Finished chip select, so send action
			if err := netconn.SendObject(playerObject); err != nil {
				exitErr <- fmt.Errorf("failed to get data stream: %w", err)
				return
			}

			if err := netconn.SendSignal(pb.Action_CHIPSEND); err != nil {
				exitErr <- fmt.Errorf("failed to get data stream: %w", err)
				return
			}

			actCount = 0
			statusChange(statusWaitActing)
		case statusWaitActing:
			// 相手がselect完了になるのを待つ
		case statusActing:
			playerAct()
		}
		time.Sleep(16 * time.Millisecond)
	}
}

func playerStatusUpdate(status pb.Data_Status) {
	switch status {
	case pb.Data_CONNECTWAIT:
		// nothing to do
	case pb.Data_CHIPSELECTWAIT:
		if playerStatus == statusWaiting || playerStatus == statusActing {
			statusChange(statusChipSelect)
		}
	case pb.Data_ACTING:
		if playerStatus == statusWaitActing {
			statusChange(statusActing)
		}
	default:
		msg := fmt.Sprintf("unexpected status: player status %d, got status %d", playerStatus, status)
		panic(msg)
	}
}

func statusChange(next int) {
	log.Printf("player status change from %d to %d", playerStatus, next)
	playerStatus = next
}

func playerAct() {
	const waitCount = 120
	const moveInterval = 180

	actCount++
	if actCount < waitCount {
		return
	}

	if actCount%moveInterval == 0 {
		playerObject.UpdateBaseTime = true
		playerObject.Type = field.ObjectTypeRockmanMove
		netconn.SendObject(playerObject)
		log.Printf("Set to move")
	}

	if actCount%moveInterval == field.ImageDelays[field.ObjectTypeRockmanMove]*4 {
		playerObject.X = rand.Intn(3)
		playerObject.Y = rand.Intn(3)
		playerObject.Type = field.ObjectTypeRockmanStand
		log.Printf("Move to (%d, %d)", playerObject.X, playerObject.Y)
		netconn.SendObject(playerObject)
	}
}
