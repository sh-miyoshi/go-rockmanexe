package netconn

import (
	"context"
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
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

func dataRecv(exitErr chan error) {
	// Recv data from stream
	for {
		data, err := dataStream.Recv()
		if err != nil {
			exitErr <- fmt.Errorf("failed to recv data: %w", err)
			return
		}

		switch data.Type {
		case pb.Data_UPDATESTATUS:
			logger.Debug("got status update data: %+v", data)
			// playerStatusUpdate(data.GetStatus())
		case pb.Data_DATA:
			// playerFieldUpdate(data.GetRawData())
		default:
			exitErr <- fmt.Errorf("invalid data type was received: %d", data.Type)
			return
		}
	}
}
