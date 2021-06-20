package netconn

import (
	"context"
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/field"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/net/routerpb"
	"google.golang.org/grpc"
)

var (
	conn       *grpc.ClientConn
	client     pb.RouterClient
	dataStream pb.Router_PublishDataClient
	sessionID  string
	exitErr    error
	status     pb.Data_Status
	fieldInfo  field.Info
	clientID   string
)

func Connect(streamAddr string, cID, cKey string) error {
	clientID = cID

	var err error
	conn, err = grpc.Dial(streamAddr, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("failed to connect server: %w", err)
	}

	client = pb.NewRouterClient(conn)
	authReq := &pb.AuthRequest{
		Id:      cID,
		Key:     cKey,
		Version: "test",
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

	go dataRecv()

	return nil
}

func Disconnect() {
	if conn != nil {
		conn.Close()
		conn = nil
	}
}

func SendObject(obj field.Object) error {
	req := &pb.Action{
		SessionID: sessionID,
		ClientID:  clientID,
		Type:      pb.Action_UPDATEOBJECT,
		Data: &pb.Action_ObjectInfo{
			ObjectInfo: field.MarshalObject(obj),
		},
	}

	res, err := client.SendAction(context.TODO(), req)
	if err != nil {
		return fmt.Errorf("send action failed: %w", err)
	}

	if !res.Success {
		return fmt.Errorf("send action got unexpected response: %s", res.ErrMsg)
	}

	return nil
}

func SendSignal(signal pb.Action_SignalType) error {
	req := &pb.Action{
		SessionID: sessionID,
		ClientID:  clientID,
		Type:      pb.Action_SENDSIGNAL,
		Data:      &pb.Action_Signal{Signal: signal},
	}

	res, err := client.SendAction(context.TODO(), req)
	if err != nil {
		return fmt.Errorf("send signal failed: %w", err)
	}

	if !res.Success {
		return fmt.Errorf("send signal got unexpected response: %s", res.ErrMsg)
	}

	return nil
}

func RemoveObject(objID string) error {
	req := &pb.Action{
		SessionID: sessionID,
		ClientID:  clientID,
		Type:      pb.Action_REMOVEOBJECT,
		Data:      &pb.Action_ObjectID{ObjectID: objID},
	}

	res, err := client.SendAction(context.TODO(), req)
	if err != nil {
		return fmt.Errorf("remove object failed: %w", err)
	}

	if !res.Success {
		return fmt.Errorf("remove object got unexpected response: %s", res.ErrMsg)
	}

	return nil
}

func dataRecv() {
	// Recv data from stream
	for {
		data, err := dataStream.Recv()
		if err != nil {
			exitErr = fmt.Errorf("failed to recv data: %w", err)
			return
		}

		switch data.Type {
		case pb.Data_UPDATESTATUS:
			status = data.GetStatus()
		case pb.Data_DATA:
			b := data.GetRawData()
			var f field.Info
			field.Unmarshal(&f, b)
			fieldInfo = f
		default:
			exitErr = fmt.Errorf("invalid data type was received: %d", data.Type)
			return
		}
	}
}

func GetStatus() (pb.Data_Status, error) {
	return status, exitErr
}

func GetFieldInfo() (*field.Info, error) {
	return &fieldInfo, exitErr
}
