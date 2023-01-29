package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	netconn "github.com/sh-miyoshi/go-rockmanexe/pkg/app/netconn"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/object"
)

const (
	streamAddr = "localhost:16283"
)

func main() {
	var clientID string
	flag.StringVar(&clientID, "c", "", "client id")
	flag.Parse()

	if clientID == "" {
		fmt.Println("Please set client ID")
		return
	}

	logger.InitLogger(true, "")

	clientKey := "testtest"
	connInst := netconn.New(netconn.Config{
		StreamAddr:     streamAddr,
		ClientID:       clientID,
		ClientKey:      clientKey,
		ProgramVersion: "testclient",
		Insecure:       true,
	})
	connInst.ConnectRequest()

	// Waiting connection
	for n := 0; ; n++ {
		if n > 100 {
			exitByError(errors.New("failed to connect to router"))
		}

		st := connInst.GetConnStatus()
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
		ClientID:       clientID,
		Type:           object.TypeRockmanStand,
		HP:             10,
		X:              1,
		Y:              1,
		Hittable:       true,
		UpdateBaseTime: true,
	}
	connInst.SendObject(obj)
	if err := connInst.BulkSendData(); err != nil {
		exitByError(err)
	}

	// TODO
	// 	switch appStatus {
	// case stateWaiting:
	// 	status := connInst.GetGameStatus()
	// 	if status == pb.Data_CHIPSELECTWAIT {
	// 		statusChange(stateOpening)
	// 	}
	// case stateOpening:
	// 	statusChange(stateChipSelect)
	// case stateChipSelect:
	// 	// Select using chip
	// 	if err := playerInst.ChipSelect(); err != nil {
	// 		return err
	// 	}

	// 	statusChange(stateWaitSelect)
	// case stateWaitSelect:
	// 	status := connInst.GetGameStatus()
	// 	if status == pb.Data_ACTING {
	// 		statusChange(stateBeforeMain)
	// 		continue
	// 	}
	// case stateBeforeMain:
	// 	statusChange(stateMain)
	// case stateMain:
	// 	if playerInst.Action() {
	// 		statusChange(stateResult)
	// 		continue
	// 	}

	// 	if err := skill.GetInst().Process(); err != nil {
	// 		return fmt.Errorf("skill process failed: %w", err)
	// 	}

	// 	status := connInst.GetGameStatus()
	// 	switch status {
	// 	case pb.Data_CHIPSELECTWAIT:
	// 		statusChange(stateChipSelect)
	// 		continue MAIN_LOOP
	// 	case pb.Data_GAMEEND:
	// 		statusChange(stateResult)
	// 		continue MAIN_LOOP
	// 	}
	// case stateResult:
	// 	logger.Info("Reached to state result")
	// 	if playerInst.Object.HP == 0 {
	// 		logger.Info("bot client lose")
	// 	} else {
	// 		logger.Info("bot client win")
	// 	}
	// 	return nil
	// }

	// if err := connInst.BulkSendData(); err != nil {
	// 	return err
	// }

	// 終了処理をできるように若干待つ
	time.Sleep(1 * time.Second)
	logger.Info("Successfully closed app")
}

func exitByError(err error) {
	// この関数が呼ばれた場所の呼び出し元情報をセットする
	logger.SetExtraSkipCount(1)
	logger.Error("Failed to run test: %+v", err)
	os.Exit(1)
}
