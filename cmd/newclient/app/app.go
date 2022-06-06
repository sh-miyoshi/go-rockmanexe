package app

import (
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
	appStatus = statusWaiting
)

func Process() {
	fpsMgr := fps.Fps{TargetFPS: 60}

	// Main loop
	for {
		switch appStatus {
		case statusWaiting:
			// nothing to do
		case statusChipSelect:
			// Select using chip
			// todo

			statusChange(statusWaitActing)
		case statusWaitActing:
			// 相手がselect完了になるのを待つ
		case statusActing:
		case statusGameEnd:
			// TODO
		}

		fpsMgr.Wait()
	}
}

func statusChange(next int) {
	logger.Info("app status change from %d to %d", appStatus, next)
	appStatus = next
}
