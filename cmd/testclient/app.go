package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/field"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/net/routerpb"
	"google.golang.org/grpc"
)

const (
	statusWaiting int = iota
	statusChipSelect
	statusWaitActing
	statusActing
)

var (
	playerStatus    = statusWaiting
	playerActClient pb.RouterClient
	fieldInfo       field.Info
	playerObject    field.Object
)

func playerInit() error {
	conn, err := grpc.Dial(streamAddr, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("grpc dial failed: %w", err)
	}
	playerActClient = pb.NewRouterClient(conn)

	playerObject = field.Object{
		ID:   uuid.New().String(),
		Type: field.ObjectTypeRockmanStand,
		HP:   100,
		X:    1,
		Y:    1,
	}

	if _, err := playerActClient.SendAction(context.TODO(), makePlayerObj()); err != nil {
		return fmt.Errorf("failed to get data stream: %w", err)
	}

	return nil
}

func playerProc(exitErr chan error) {
	for {
		switch playerStatus {
		case statusWaiting:
			// nothing to do
		case statusChipSelect:
			// Select using chip
			n := rand.Intn(2) + 1
			time.Sleep(time.Duration(n) * time.Second)
			playerObject.Chips = []int{1, 3, -1} // debug

			// Finished chip select, so send action
			if _, err := playerActClient.SendAction(context.TODO(), makePlayerObj()); err != nil {
				exitErr <- fmt.Errorf("failed to get data stream: %w", err)
				return
			}

			if _, err := playerActClient.SendAction(context.TODO(), &pb.Action{
				SessionID: sessionID,
				ClientID:  clientID,
				Type:      pb.Action_SENDSIGNAL,
				Data:      &pb.Action_Signal{Signal: pb.Action_CHIPSEND},
			}); err != nil {
				exitErr <- fmt.Errorf("failed to get data stream: %w", err)
				return
			}

			statusChange(statusWaitActing)
		case statusWaitActing:
			// 相手がselect完了になるのを待つ
		case statusActing:
			// TODO
		}
		time.Sleep(16 * time.Millisecond)
	}
}

func playerStatusUpdate(status pb.Data_Status) {
	switch status {
	case pb.Data_CHIPSELECTWAIT:
		if playerStatus == statusWaiting {
			statusChange(statusChipSelect)
			return
		}
	case pb.Data_ACTING:
		if playerStatus == statusActing {
			// nothing to do
			return
		}

		if playerStatus == statusWaitActing {
			statusChange(statusActing)
			return
		}
	}

	msg := fmt.Sprintf("unexpected status: player status %d, got status %d", playerStatus, status)
	panic(msg)
}

func playerFieldUpdate(data []byte) {
	field.Unmarshal(&fieldInfo, data)
	// log.Printf("Update field data to %+v", fieldInfo)
}

func statusChange(next int) {
	log.Printf("player status change from %d to %d", playerStatus, next)
	playerStatus = next
}

func makePlayerObj() *pb.Action {
	return &pb.Action{
		SessionID: sessionID,
		ClientID:  clientID,
		Type:      pb.Action_UPDATEOBJECT,
		Data: &pb.Action_ObjectInfo{
			ObjectInfo: field.MarshalObject(playerObject),
		},
	}
}
