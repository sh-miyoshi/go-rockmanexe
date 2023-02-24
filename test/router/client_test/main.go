package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	netconn "github.com/sh-miyoshi/go-rockmanexe/pkg/app/newnetconn"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/newnet/action"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/newnet/netconnpb"
	netobj "github.com/sh-miyoshi/go-rockmanexe/pkg/newnet/object"
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

const (
	streamAddr = "localhost:16283"
)

func main() {
	logger.InitLogger(true, "")

	conn := netconn.New(netconn.Config{
		StreamAddr:     streamAddr,
		ClientID:       "tester1",
		ClientKey:      "testtest",
		ProgramVersion: "testclient",
		Insecure:       true,
	})
	conn.ConnectRequest()

	// Waiting connection
	for n := 0; ; n++ {
		if n > 100 {
			exitByError(errors.New("failed to connect to router"))
		}

		st := conn.GetConnStatus()
		if st.Status == netconn.ConnStateOK {
			break
		}
		if st.Status == netconn.ConnStateError {
			exitByError(st.Error)
		}

		time.Sleep(100 * time.Millisecond)
	}

	logger.Info("Success to connect router")

	obj := netobj.InitParam{
		ID: uuid.New().String(),
		HP: 10,
		X:  1,
		Y:  1,
	}
	if err := conn.SendSignal(pb.Request_INITPARAMS, obj.Marshal()); err != nil {
		exitByError(err)
	}

	go runClient2()

	// Client2の起動を待つ
	time.Sleep(300 * time.Millisecond)

	startTime := time.Now()
	appStatus := stateWaiting
MAIN_LOOP:
	for {
		current := time.Now()
		const timeout = 5 * time.Second
		if current.After(startTime.Add(timeout)) {
			exitByError(errors.New("timeout main loop"))
		}

		switch appStatus {
		case stateWaiting:
			status := conn.GetGameStatus()
			if status == pb.Response_CHIPSELECTWAIT {
				stateChange(&appStatus, stateChipSelect)
			}
		case stateChipSelect:
			// obj.Chips = []object.ChipInfo{
			// 	{ID: 1, Code: "*"},
			// }

			// conn.SendObject(obj)

			conn.SendSignal(pb.Request_CHIPSELECT, nil)
			stateChange(&appStatus, stateWaitSelect)
		case stateWaitSelect:
			status := conn.GetGameStatus()
			if status == pb.Response_ACTING {
				stateChange(&appStatus, stateMain)
			}
		case stateMain:
			status := conn.GetGameStatus()
			if status == pb.Response_GAMEEND {
				stateChange(&appStatus, stateResult)
			}

			// Check action
			// 1. Move
			ok := false
			move := action.Move{
				ObjectID: obj.ID,
				Type:     action.MoveTypeAbs,
				AbsPosX:  2,
				AbsPosY:  1,
			}
			conn.SendAction(pb.Request_MOVE, common.Point{X: 1, Y: 1}, move.Marshal())
			info := conn.GetGameInfo()
			for i := 0; i < 10; i++ {
				info = conn.GetGameInfo()
				myObj := info.Objects[obj.ID]
				if myObj.Pos.X == 2 && myObj.Pos.Y == 1 {
					ok = true
					logger.Info("Success to move")
					break
				}
				time.Sleep(30 * time.Millisecond)
			}
			if !ok {
				exitByError(fmt.Errorf("failed to move: %+v", info))
			}

			// 2. Buster
			buster := action.Buster{
				ObjectID: obj.ID,
				Power:    1,
			}
			conn.SendAction(pb.Request_BUSTER, common.Point{X: 2, Y: 1}, buster.Marshal())
			ok = false
		BUSTER_CHECK_LOOP:
			for i := 0; i < 10; i++ {
				info := conn.GetGameInfo()
				for _, obj := range info.Objects {
					if obj.OwnerClientID != "tester1" && obj.HP == 0 {
						ok = true
						logger.Info("Successfully damaged by buster")
						break BUSTER_CHECK_LOOP
					}
				}
				time.Sleep(30 * time.Millisecond)
			}
			if !ok {
				exitByError(fmt.Errorf("failed to add buster: %+v", info))
			}

			// TODO
			break MAIN_LOOP

			// 3. Use Chip
			// TODO
		case stateResult:
			logger.Info("Successfully state change to result")
			break MAIN_LOOP
		}
	}

	// 終了処理をできるように若干待つ
	time.Sleep(1 * time.Second)
	logger.Info("Successfully closed app")
}

func stateChange(state *int, next int) {
	logger.Info("state change from %d to %d", *state, next)
	*state = next
}

func exitByError(err error) {
	// この関数が呼ばれた場所の呼び出し元情報をセットする
	logger.SetExtraSkipCount(1)
	logger.Error("Failed to run test: %+v", err)
	os.Exit(1)
}

func runClient2() {
	conn := netconn.New(netconn.Config{
		StreamAddr:     streamAddr,
		ClientID:       "tester2",
		ClientKey:      "testtest",
		ProgramVersion: "testclient",
		Insecure:       true,
	})
	conn.ConnectRequest()

	for n := 0; ; n++ {
		if n > 100 {
			exitByError(errors.New("failed to connect to router"))
		}

		st := conn.GetConnStatus()
		if st.Status == netconn.ConnStateOK {
			break
		}
		if st.Status == netconn.ConnStateError {
			exitByError(st.Error)
		}

		time.Sleep(100 * time.Millisecond)
	}

	obj := netobj.InitParam{
		ID: uuid.New().String(),
		HP: 1,
		X:  1,
		Y:  1,
	}
	conn.SendSignal(pb.Request_INITPARAMS, obj.Marshal())

	appStatus := stateWaiting
	for {
		switch appStatus {
		case stateWaiting:
			status := conn.GetGameStatus()
			if status == pb.Response_CHIPSELECTWAIT {
				appStatus = stateChipSelect
			}
		case stateChipSelect:
			// obj.Chips = []object.ChipInfo{
			// 	{ID: 1, Code: "*"},
			// }

			// conn.SendObject(obj)
			// conn.BulkSendData()
			conn.SendSignal(pb.Request_CHIPSELECT, nil)
			appStatus = stateWaitSelect
		case stateWaitSelect:
			status := conn.GetGameStatus()
			if status == pb.Response_ACTING {
				appStatus = stateMain
				continue
			}
		case stateMain:
			time.Sleep(300 * time.Millisecond)

			// 負けたことにする
			// obj.HP = 0
			// conn.SendObject(obj)
			// conn.BulkSendData()
			return
		}
	}
}
