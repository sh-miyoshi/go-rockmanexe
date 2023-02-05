package app

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/cmd/botclient/player"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/skill"
	netconn "github.com/sh-miyoshi/go-rockmanexe/pkg/app/oldnetconn"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/fps"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/oldnet/netconnpb"
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

func Init(clientID string, conn *netconn.NetConn) {
	dxlib.Disable()
	playerInst = player.New(clientID, conn)
	connInst = conn
	skill.GetInst().Init(playerInst.Object.ID)
	draw.Init(playerInst.Object.ID)
}

func Process() error {
	// set init data to router
	connInst.SendObject(playerInst.Object)
	fpsMgr := fps.Fps{TargetFPS: 60}

	// Main loop
MAIN_LOOP:
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
				continue
			}
		case stateBeforeMain:
			statusChange(stateMain)
		case stateMain:
			if playerInst.Action() {
				statusChange(stateResult)
				continue
			}

			if err := skill.GetInst().Process(); err != nil {
				return fmt.Errorf("skill process failed: %w", err)
			}

			status := connInst.GetGameStatus()
			switch status {
			case pb.Data_CHIPSELECTWAIT:
				statusChange(stateChipSelect)
				continue MAIN_LOOP
			case pb.Data_GAMEEND:
				statusChange(stateResult)
				continue MAIN_LOOP
			}
		case stateResult:
			logger.Info("Reached to state result")
			if playerInst.Object.HP == 0 {
				logger.Info("bot client lose")
			} else {
				logger.Info("bot client win")
			}
			return nil
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
