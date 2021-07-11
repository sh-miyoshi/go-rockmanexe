package app

import (
	"fmt"
	"log"
	"time"

	"github.com/sh-miyoshi/go-rockmanexe/cmd/testclient/netconn"
	"github.com/sh-miyoshi/go-rockmanexe/cmd/testclient/skill"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/net/routerpb"
)

const (
	statusWaiting int = iota
	statusChipSelect
	statusWaitActing
	statusActing
)

var (
	appStatus  = statusWaiting
	playerInst *player
)

func Init(clientID string) error {
	playerInst = newPlayer(clientID)

	return nil
}

func Process(exitErr chan error) {
	// set init data to router
	netconn.SendObject(playerInst.Object)

	// Main loop
	for {
		status, err := netconn.GetStatus()
		if err != nil {
			exitErr <- fmt.Errorf("got status failed: %v", err)
			return
		}
		statusUpdate(status)

		switch appStatus {
		case statusWaiting:
			// nothing to do
		case statusChipSelect:
			// Select using chip
			if err := playerInst.ChipSelect(); err != nil {
				exitErr <- err
				return
			}

			statusChange(statusWaitActing)
		case statusWaitActing:
			// 相手がselect完了になるのを待つ
		case statusActing:
			playerInst.Action()
			skill.Process()
		}
		time.Sleep(16 * time.Millisecond)
	}
}

func statusUpdate(status pb.Data_Status) {
	switch status {
	case pb.Data_CONNECTWAIT:
		// nothing to do
	case pb.Data_CHIPSELECTWAIT:
		if appStatus == statusWaiting || appStatus == statusActing {
			statusChange(statusChipSelect)
		}
	case pb.Data_ACTING:
		if appStatus == statusWaitActing {
			statusChange(statusActing)
		}
	default:
		msg := fmt.Sprintf("unexpected status: app status %d, got status %d", appStatus, status)
		panic(msg)
	}
}

func statusChange(next int) {
	log.Printf("app status change from %d to %d", appStatus, next)
	appStatus = next
}
