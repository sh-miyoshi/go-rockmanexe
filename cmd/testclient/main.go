package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	routerapi "github.com/sh-miyoshi/go-rockmanexe/pkg/net/api/router"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/net/routerpb"
	"google.golang.org/grpc"
)

var (
	apiAddr    = "http://localhost:8080"
	streamAddr = "localhost:80"
)

func main() {
	// Add clients and route to server
	client1 := clientAdd()
	client2 := clientAdd()
	routeAdd(client1.ID, client2.ID)

	exitErr := make(chan error)
	go clientProc(exitErr, client1)
	go clientProc(exitErr, client2)

	err := <-exitErr
	log.Fatalf("Run failed: %v", err)
}

func clientAdd() routerapi.ClientInfo {
	httpRes, err := http.Post(apiAddr+"/api/v1/client", "text/plain", nil)
	if err != nil {
		log.Fatalf("Failed to add client: %v", err)
	}
	defer httpRes.Body.Close()

	var res routerapi.ClientInfo
	if httpRes.StatusCode == http.StatusOK {
		if err := json.NewDecoder(httpRes.Body).Decode(&res); err != nil {
			log.Fatalf("Failed to decode client add response: %v", err)
		}
		log.Printf("Client Info: %+v", res)
	} else {
		log.Fatalf("Client add request returns unexpected status %s", httpRes.Status)
	}
	return res
}

func routeAdd(id1, id2 string) routerapi.RouteInfo {
	req := routerapi.RouteAddRequest{
		Clients: [2]string{id1, id2},
	}

	body, _ := json.Marshal(req)
	httpRes, err := http.Post(apiAddr+"/api/v1/route", "application/json", bytes.NewReader(body))
	if err != nil {
		log.Fatalf("Failed to add route: %v", err)
	}
	defer httpRes.Body.Close()

	var res routerapi.RouteInfo
	if httpRes.StatusCode == http.StatusOK {
		if err := json.NewDecoder(httpRes.Body).Decode(&res); err != nil {
			log.Fatalf("Failed to decode route add response: %v", err)
		}
		log.Printf("Route Info: %+v", res)
	} else {
		log.Fatalf("Route add request returns unexpected status %s", httpRes.Status)
	}
	return res
}

func clientProc(exitErr chan error, clientInfo routerapi.ClientInfo) {
	conn, err := grpc.Dial(streamAddr, grpc.WithInsecure())
	if err != nil {
		exitErr <- fmt.Errorf("failed to connect server: %w", err)
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
	}

	// At first, receive authenticate response
	authRes, err := dataStream.Recv()
	if err != nil {
		exitErr <- fmt.Errorf("failed to recv authenticate res: %w", err)
	}
	if authRes.GetType() != pb.Data_AUTHRESPONSE {
		exitErr <- fmt.Errorf("expect type is auth res, but got: %d", authRes.GetType())
	}
	if res := authRes.GetAuthRes(); !res.Success {
		exitErr <- fmt.Errorf("failed to auth request: %s", res.ErrMsg)
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
			// TODO
		case pb.Data_DATA:
			log.Printf("got data: %+v", data)
		default:
			exitErr <- fmt.Errorf("invalid data type was received: %d", data.Type)
			return
		}
	}
}
