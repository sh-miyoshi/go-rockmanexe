package netconn

import (
	"context"
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/net/routerpb"
	"google.golang.org/grpc"
)

var (
	conn       *grpc.ClientConn
	dataStream pb.Router_PublishDataClient
	sessionID  string
)

func Connect() error {
	c := config.Get()
	var err error
	conn, err = grpc.Dial(c.Net.StreamAddr, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("failed to connect server: %w", err)
	}

	client := pb.NewRouterClient(conn)
	authReq := &pb.AuthRequest{
		Id:      c.Net.ClientID,
		Key:     c.Net.ClientKey,
		Version: common.ProgramVersion,
	}
	dataStream, err = client.PublishData(context.TODO(), authReq)
	if err != nil {
		return fmt.Errorf("failed to get data stream: %w", err)
	}

	// At first, receive authenticate response
	authRes, err := dataStream.Recv()
	if err != nil {
		return fmt.Errorf("failed to recv authenticate res: %w", err)
	}
	if authRes.GetType() != pb.Data_AUTHRESPONSE {
		return fmt.Errorf("expect type is auth res, but got: %d", authRes.GetType())
	}
	if res := authRes.GetAuthRes(); !res.Success {
		return fmt.Errorf("failed to auth request: %s", res.ErrMsg)
	}
	sessionID = authRes.GetAuthRes().SessionID

	return nil
}

func Disconnect() {
	if conn != nil {
		conn.Close()
		conn = nil
	}
}

// TODO(add player object, recv data from stream)

// func clientProc(exitErr chan error, clientInfo routerapi.ClientInfo) {
// 	// Add player object
// 	objRes, err := playerActClient.SendAction(context.TODO(), makePlayerObj())
// 	if err != nil {
// 		exitErr <- fmt.Errorf("add player object failed by error: %w", err)
// 		return
// 	}
// 	if !objRes.Success {
// 		exitErr <- fmt.Errorf("add player object failed: %s", objRes.ErrMsg)
// 		return
// 	}

// 	// Recv data from stream
// 	for {
// 		data, err := dataStream.Recv()
// 		if err != nil {
// 			exitErr <- fmt.Errorf("failed to recv data: %w", err)
// 			return
// 		}

// 		switch data.Type {
// 		case pb.Data_UPDATESTATUS:
// 			log.Printf("got status update data: %+v", data)
// 			playerStatusUpdate(data.GetStatus())
// 		case pb.Data_DATA:
// 			// log.Printf("got data: %+v", data)
// 			playerFieldUpdate(data.GetRawData())
// 		default:
// 			exitErr <- fmt.Errorf("invalid data type was received: %d", data.Type)
// 			return
// 		}
// 	}
// }
