package main

import (
	"bytes"
	"context"
	"encoding/json"
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

	conn, err := grpc.Dial(streamAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect server: %v", err)
	}
	defer conn.Close()

	client := pb.NewRouterClient(conn)
	req := &pb.AuthRequest{
		Id:      client1.ID,
		Key:     client1.Key,
		Version: "test",
	}
	dataStream, err := client.PublishData(context.TODO(), req)
	if err != nil {
		log.Fatalf("Failed to get data stream: %v", err)
	}

	log.Printf("data stream: %v", dataStream)
	// TODO recv data from stream
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
