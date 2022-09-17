package app

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/cmd/botclient/player"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/skill"
	netconn "github.com/sh-miyoshi/go-rockmanexe/pkg/app/netconn"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/fps"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/net/netconnpb"
)

const (
	stateWaiting int = iota
	stateOpening
	stateChipSelect
	stateWaitSelect
	stateBeforeMain
	stateMain
	stateResult
)

var (
	appStatus  = stateWaiting
	playerInst *player.Player
	connInst   *netconn.NetConn
)

func Init(clientID string) {
	dxlib.Disable()
	skill.GetInst().Init()
	playerInst = player.New(clientID)
	connInst = netconn.GetInst()
}

func Process() error {
	// set init data to router
	connInst.SendObject(playerInst.Object)
	fpsMgr := fps.Fps{TargetFPS: 60}

	// Main loop
	for {
		switch appStatus {
		case stateWaiting:
			status := connInst.GetGameStatus()
			if status == pb.Data_CHIPSELECTWAIT {
				statusChange(stateOpening)
			}
		case stateOpening:
			statusChange(stateChipSelect)
		case stateChipSelect:
			// Select using chip
			if err := playerInst.ChipSelect(); err != nil {
				return err
			}

			statusChange(stateWaitSelect)
		case stateWaitSelect:
			status := connInst.GetGameStatus()
			if status == pb.Data_ACTING {
				statusChange(stateBeforeMain)
			}
		case stateBeforeMain:
			statusChange(stateMain)
		case stateMain:
			if playerInst.Action() {
				statusChange(stateResult)
			}

			if err := skill.GetInst().Process(); err != nil {
				return fmt.Errorf("skill process failed: %w", err)
			}
		case stateResult:
			panic("not implemented yet")
		}

		if err := connInst.BulkSendData(); err != nil {
			return err
		}

		fpsMgr.Wait()
	}
}

func statusChange(next int) {
	logger.Info("app status change from %d to %d", appStatus, next)
	appStatus = next
}
