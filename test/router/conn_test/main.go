package main

import (
	"errors"
	"os"
	"sync"
	"time"

	netconn "github.com/sh-miyoshi/go-rockmanexe/pkg/app/netconn"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
)

const (
	streamAddr = "localhost:16283"
)

// テスト内容
//   1. 同時通信
//     tester1, tester2
//     tester3, tester4
//   2. 再接続

func main() {
	logger.InitLogger(true, "")

	var wg sync.WaitGroup
	wg.Add(2)
	go runClient("tester1", &wg)
	go runClient("tester2", &wg)

	wg.Wait()
	logger.Info("Finished first battle for tester1, tester2")

	// TODO

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

func runClient(clientID string, wg *sync.WaitGroup) {
	defer wg.Done()

	conn := netconn.New(netconn.Config{
		StreamAddr:     streamAddr,
		ClientID:       clientID,
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
	logger.Info("Successfully connected client %s to router", clientID)

	// TODO
}
