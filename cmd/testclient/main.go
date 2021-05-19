package main

import (
	"context"
	"log"

	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/net/routerpb"
	"google.golang.org/grpc"
)

func main() {
	serverAddr := "localhost:80"
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect server: %v", err)
	}
	defer conn.Close()

	client := pb.NewRouterClient(conn)
	req := &pb.AuthRequest{
		Id:      "tester",
		Key:     "testclient",
		Version: "test",
	}
	dataStream, err := client.PublishData(context.TODO(), req)
	if err != nil {
		log.Fatalf("Failed to get data stream: %v", err)
	}

	log.Printf("data stream: %v", dataStream)
}
