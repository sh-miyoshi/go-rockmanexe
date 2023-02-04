package main

import (
	"errors"
	"os"
	"time"

	"github.com/google/uuid"
	netconn "github.com/sh-miyoshi/go-rockmanexe/pkg/app/netconn"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/net/netconnpb"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/object"
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

	// TODO test
	obj := object.Object{
		ID:             uuid.New().String(),
		ClientID:       "tester1",
		Type:           object.TypeRockmanStand,
		HP:             10,
		X:              1,
		Y:              1,
		Hittable:       true,
		UpdateBaseTime: true,
	}
	conn.SendObject(obj)
	if err := conn.BulkSendData(); err != nil {
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
			if status == pb.Data_CHIPSELECTWAIT {
				stateChange(&appStatus, stateChipSelect)
			}
		case stateChipSelect:
			obj.Chips = []object.ChipInfo{
				{ID: 1, Code: "*"},
			}

			conn.SendObject(obj)
			conn.SendSignal(pb.Action_CHIPSEND)
			stateChange(&appStatus, stateWaitSelect)
		case stateWaitSelect:
			status := conn.GetGameStatus()
			if status == pb.Data_ACTING {
				stateChange(&appStatus, stateMain)
			}
		case stateMain:
			status := conn.GetGameStatus()
			if status == pb.Data_GAMEEND {
				stateChange(&appStatus, stateResult)
			}
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

	obj := object.Object{
		ID:             uuid.New().String(),
		ClientID:       "tester2",
		Type:           object.TypeRockmanStand,
		HP:             10,
		X:              1,
		Y:              1,
		Hittable:       true,
		UpdateBaseTime: true,
	}
	conn.SendObject(obj)
	conn.BulkSendData()

	appStatus := stateWaiting
	for {
		switch appStatus {
		case stateWaiting:
			status := conn.GetGameStatus()
			if status == pb.Data_CHIPSELECTWAIT {
				appStatus = stateChipSelect
			}
		case stateChipSelect:
			obj.Chips = []object.ChipInfo{
				{ID: 1, Code: "*"},
			}

			conn.SendObject(obj)
			conn.BulkSendData()
			conn.SendSignal(pb.Action_CHIPSEND)
			appStatus = stateWaitSelect
		case stateWaitSelect:
			status := conn.GetGameStatus()
			if status == pb.Data_ACTING {
				appStatus = stateMain
				continue
			}
		case stateMain:
			time.Sleep(300 * time.Millisecond)

			// 負けたことにする
			obj.HP = 0
			conn.SendObject(obj)
			conn.BulkSendData()
			conn.SendSignal(pb.Action_PLAYERDEAD)
			return
		}
	}
}
