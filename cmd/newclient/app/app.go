package app

import (
	netconn "github.com/sh-miyoshi/go-rockmanexe/pkg/app/newnetconn"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/fps"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
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

func Init(clientID string) {
	playerInst = newPlayer(clientID)
}

func Process() error {
	// set init data to router
	netconn.GetInst().SendObject(playerInst.Object)
	fpsMgr := fps.Fps{TargetFPS: 60}

	// Main loop
	for {
		switch appStatus {
		case statusWaiting:
			// nothing to do
		case statusChipSelect:
			// Select using chip
			if err := playerInst.ChipSelect(); err != nil {
				return err
			}

			statusChange(statusWaitActing)
		case statusWaitActing:
			// 相手がselect完了になるのを待つ
		case statusActing:
		case statusGameEnd:
			// TODO
		}

		if err := netconn.GetInst().BulkSendData(); err != nil {
			return err
		}

		fpsMgr.Wait()
	}
}

func statusChange(next int) {
	logger.Info("app status change from %d to %d", appStatus, next)
	appStatus = next
}
