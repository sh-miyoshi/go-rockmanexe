package main

import (
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/cockroachdb/errors"

	"github.com/google/uuid"
	netconn "github.com/sh-miyoshi/go-rockmanexe/pkg/app/netconn"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/action"
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
)

const (
	streamAddr = "localhost:16283"
)

func main() {
	logger.InitLogger(true, "")

	// 1. 同時通信
	var wg sync.WaitGroup
	wg.Add(4)
	go runClient("tester1", &wg, true)
	go runClient("tester2", &wg, false)

	go runClient("tester3", &wg, true)
	go runClient("tester4", &wg, false)

	wg.Wait()
	logger.Info("Successfully finished multiple sessions")

	//   2. 再接続
	time.Sleep(1 * time.Second)
	wg.Add(2)
	go runClient("tester1", &wg, true)
	go runClient("tester2", &wg, false)
	wg.Wait()
	logger.Info("Successfully finished second time session")

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

func runClient(clientID string, wg *sync.WaitGroup, isWinner bool) {
	defer wg.Done()

	n := rand.Intn(30)
	time.Sleep(time.Duration(n) * time.Millisecond)

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

	startTime := time.Now()
	appStatus := stateWaiting
	atkFlag := false
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
				obj := netobj.InitParam{
					ID: uuid.New().String(),
					HP: 10,
					X:  1,
					Y:  1,
				}
				conn.SendSignal(pb.Request_INITPARAMS, obj.Marshal())
				stateChange(&appStatus, stateChipSelect)
			}
		case stateChipSelect:
			conn.SendSignal(pb.Request_CHIPSELECT, nil)
			stateChange(&appStatus, stateWaitSelect)
		case stateWaitSelect:
			status := conn.GetGameStatus()
			if status == pb.Response_ACTING {
				stateChange(&appStatus, stateMain)
				continue
			}
		case stateMain:
			if isWinner && !atkFlag {
				atkFlag = true
				buster := action.Buster{
					Power: 10,
				}
				conn.SendAction(pb.Request_BUSTER, buster.Marshal())
			}

			status := conn.GetGameStatus()
			if status == pb.Response_GAMEEND {
				logger.Info("got game end at %s", clientID)
				return
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func stateChange(state *int, next int) {
	logger.Info("state change from %d to %d", *state, next)
	*state = next
}
