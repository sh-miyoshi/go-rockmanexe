package app

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/cmd/botclient/player"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/netconn"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/fps"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/net/netconnpb"
	netobj "github.com/sh-miyoshi/go-rockmanexe/pkg/net/object"
)

const (
	stateWaiting int = iota
	stateOpening
	stateChipSelect
	stateWaitSelect
	stateBeforeMain
	stateMain
	stateResult
	stateCutin
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
}

func Process() error {
	fps.FPS = 60
	fpsMgr := fps.Fps{}

	// Main loop
MAIN_LOOP:
	for {
		switch appStatus {
		case stateWaiting:
			status := connInst.GetGameStatus()
			if status == pb.Response_CHIPSELECTWAIT {
				obj := netobj.InitParam{
					ID: playerInst.ID,
					HP: playerInst.HP,
					X:  playerInst.Pos.X,
					Y:  playerInst.Pos.Y,
				}
				if err := connInst.SendSignal(pb.Request_INITPARAMS, obj.Marshal()); err != nil {
					return fmt.Errorf("failed to send init object param: %w", err)
				}

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
			if status == pb.Response_ACTING {
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

			status := connInst.GetGameStatus()
			switch status {
			case pb.Response_CHIPSELECTWAIT:
				statusChange(stateChipSelect)
				continue MAIN_LOOP
			case pb.Response_GAMEEND:
				statusChange(stateResult)
				continue MAIN_LOOP
			case pb.Response_CUTIN:
				statusChange(stateCutin)
				continue MAIN_LOOP
			}
		case stateResult:
			logger.Info("Reached to state result")
			if playerInst.HP == 0 {
				logger.Info("bot client lose")
			} else {
				logger.Info("bot client win")
			}
			return nil
		case stateCutin:
			status := connInst.GetGameStatus()
			if status == pb.Response_ACTING {
				statusChange(stateMain)
			}
		}

		fpsMgr.Wait()
	}
}

func statusChange(next int) {
	logger.Info("app status change from %d to %d", appStatus, next)
	appStatus = next
}
