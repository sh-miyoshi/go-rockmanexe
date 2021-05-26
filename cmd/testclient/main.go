package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	routerapi "github.com/sh-miyoshi/go-rockmanexe/pkg/net/api/router"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/net/routerpb"
	"google.golang.org/grpc"
)

var (
	streamAddr = "localhost:80"
	sessionID  = ""
	clientID   = ""
)

func main() {
	flag.StringVar(&clientID, "client", "", "client id")
	flag.StringVar(&clientID, "c", "", "client id")
	flag.Parse()

	if clientID == "" {
		fmt.Println("Please set client ID")
		return
	}

	if err := playerInit(); err != nil {
		log.Fatalf("Failed to init player info: %v", err)
		return
	}

	// run with debug client
	clientKey := "testtest"

	exitErr := make(chan error)
	go clientProc(exitErr, routerapi.ClientInfo{
		ID:  clientID,
		Key: clientKey,
	})
	go playerProc(exitErr)

	err := <-exitErr
	log.Fatalf("Run failed: %v", err)
}

func clientProc(exitErr chan error, clientInfo routerapi.ClientInfo) {
	conn, err := grpc.Dial(streamAddr, grpc.WithInsecure())
	if err != nil {
		exitErr <- fmt.Errorf("failed to connect server: %w", err)
		return
	}
	defer conn.Close()

	client := pb.NewRouterClient(conn)
	req := &pb.AuthRequest{
		Id:      clientInfo.ID,
		Key:     clientInfo.Key,
		Version: "test",
	}
	dataStream, err := client.PublishData(context.TODO(), req)
	if err != nil {
		exitErr <- fmt.Errorf("failed to get data stream: %w", err)
		return
	}

	// At first, receive authenticate response
	authRes, err := dataStream.Recv()
	if err != nil {
		exitErr <- fmt.Errorf("failed to recv authenticate res: %w", err)
		return
	}
	if authRes.GetType() != pb.Data_AUTHRESPONSE {
		exitErr <- fmt.Errorf("expect type is auth res, but got: %d", authRes.GetType())
		return
	}
	if res := authRes.GetAuthRes(); !res.Success {
		exitErr <- fmt.Errorf("failed to auth request: %s", res.ErrMsg)
		return
	}
	sessionID = authRes.GetAuthRes().SessionID

	// Add player object
	objRes, err := playerActClient.SendAction(context.TODO(), makePlayerObj())
	if err != nil {
		exitErr <- fmt.Errorf("add player object failed by error: %w", err)
		return
	}
	if !objRes.Success {
		exitErr <- fmt.Errorf("add player object failed: %s", objRes.ErrMsg)
		return
	}

	// Recv data from stream
	for {
		data, err := dataStream.Recv()
		if err != nil {
			exitErr <- fmt.Errorf("failed to recv data: %w", err)
			return
		}

		switch data.Type {
		case pb.Data_UPDATESTATUS:
			log.Printf("got status update data: %+v", data)
			playerStatusUpdate(data.GetStatus())
		case pb.Data_DATA:
			// log.Printf("got data: %+v", data)
			playerFieldUpdate(data.GetRawData())
		default:
			exitErr <- fmt.Errorf("invalid data type was received: %d", data.Type)
			return
		}
	}
}
