package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	routerapi "github.com/sh-miyoshi/go-rockmanexe/pkg/net/api/router"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/net/routerpb"
	"google.golang.org/grpc"
)

func main() {
	// Add client to server
	apiAddr := "http://localhost:8080"
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

	serverAddr := "localhost:80"
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect server: %v", err)
	}
	defer conn.Close()

	client := pb.NewRouterClient(conn)
	req := &pb.AuthRequest{
		Id:      res.ID,
		Key:     res.Key,
		Version: "test",
	}
	dataStream, err := client.PublishData(context.TODO(), req)
	if err != nil {
		log.Fatalf("Failed to get data stream: %v", err)
	}

	log.Printf("data stream: %v", dataStream)
}
