package app

import (
	"fmt"
	"time"

	"github.com/sh-miyoshi/go-rockmanexe/cmd/testclient/skill"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/netconn"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/net/routerpb"
)

const (
	statusWaiting int = iota
	statusChipSelect
	statusWaitActing
	statusActing
	statusGameEnd
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
		status := netconn.GetStatus()
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
			if playerInst.Action() {
				netconn.SendSignal(pb.Action_PLAYERDEAD)
				statusChange(statusGameEnd)
			}
			skill.Process()
		case statusGameEnd:
			// TODO
		}

		if err := netconn.BulkSendFieldInfo(); err != nil {
			exitErr <- fmt.Errorf("net send failed: %v", err)
			return
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
	case pb.Data_GAMEEND:
		if appStatus == statusChipSelect || appStatus == statusActing {
			statusChange(statusGameEnd)
		}
	default:
		msg := fmt.Sprintf("unexpected status: app status %d, got status %d", appStatus, status)
		panic(msg)
	}
}

func statusChange(next int) {
	logger.Info("app status change from %d to %d", appStatus, next)
	appStatus = next
}
