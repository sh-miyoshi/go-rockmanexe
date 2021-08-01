package netconn

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/object"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/net/routerpb"
	"google.golang.org/grpc"
)

type Config struct {
	StreamAddr     string
	ClientID       string
	ClientKey      string
	ProgramVersion string
}

var (
	conn         *grpc.ClientConn
	routerClient pb.RouterClient
	dataStream   pb.Router_PublishDataClient
	sessionID    string
	clientID     string

	status        pb.Data_Status
	fieldInfo     field.Info
	sendObjects   = make(map[string]object.Object)
	removeObjects = []string{}
	sendDamages   = []damage.Damage{}
	sendEffects   = []effect.Effect{}

	exitErr   error
	fieldLock sync.Mutex
)

func Connect(conf Config) error {
	var err error
	conn, err = grpc.Dial(conf.StreamAddr, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("failed to connect server: %w", err)
	}

	routerClient = pb.NewRouterClient(conn)
	authReq := &pb.AuthRequest{
		Id:      conf.ClientID,
		Key:     conf.ClientKey,
		Version: conf.ProgramVersion,
	}
	dataStream, err = routerClient.PublishData(context.TODO(), authReq)
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
	clientID = conf.ClientID

	go dataRecv()

	return nil
}

func Disconnect() {
	if conn != nil {
		conn.Close()
		conn = nil
	}
}

func SendSignal(signal pb.Action_SignalType) error {
	req := &pb.Action{
		SessionID: sessionID,
		ClientID:  clientID,
		Type:      pb.Action_SENDSIGNAL,
		Data:      &pb.Action_Signal{Signal: signal},
	}

	res, err := routerClient.SendAction(context.TODO(), req)
	if err != nil {
		return fmt.Errorf("send signal failed: %w", err)
	}

	if !res.Success {
		return fmt.Errorf("send signal got unexpected response: %s", res.ErrMsg)
	}

	return nil
}

func SendObject(obj object.Object) {
	obj.ClientID = clientID
	sendObjects[obj.ID] = obj
}

func RemoveObject(objID string) {
	removeObjects = append(removeObjects, objID)
}

func SendDamages(damages []damage.Damage) {
	for _, dm := range damages {
		dm.ClientID = clientID
		sendDamages = append(sendDamages, dm)
	}
}

func SendEffect(eff effect.Effect) {
	eff.ClientID = clientID
	sendEffects = append(sendEffects, eff)
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
			logger.Debug("got status update data: %+v", data)
			status = data.GetStatus()
		case pb.Data_DATA:
			b := data.GetRawData()
			var f field.Info
			field.Unmarshal(&f, b)
			fieldLock.Lock()
			fieldInfo = f
			fieldLock.Unlock()
		default:
			exitErr = fmt.Errorf("invalid data type was received: %d", data.Type)
			return
		}
	}
}

// TODO error
func GetStatus() (pb.Data_Status, error) {
	return status, exitErr
}

// TODO error
func GetFieldInfo() (*field.Info, error) {
	return &fieldInfo, exitErr
}

func UpdateObjectsCount() {
	fieldLock.Lock()
	defer fieldLock.Unlock()
	for i, obj := range fieldInfo.Objects {
		if obj.Count == 0 {
			tm := fieldInfo.CurrentTime.Sub(obj.BaseTime)
			fieldInfo.Objects[i].Count = int(tm * 60 / time.Second)
		} else {
			fieldInfo.Objects[i].Count++
		}
	}
}

func RemoveEffects() {
	fieldLock.Lock()
	fieldInfo.Effects = []effect.Effect{}
	fieldLock.Unlock()
}

func RemoveDamage() {
	fieldLock.Lock()
	fieldInfo.HitDamage.ID = ""
	fieldLock.Unlock()
}

func BulkSendFieldInfo() error {
	// Send objects
	for _, obj := range sendObjects {
		req := &pb.Action{
			SessionID: sessionID,
			ClientID:  clientID,
			Type:      pb.Action_UPDATEOBJECT,
			Data: &pb.Action_ObjectInfo{
				ObjectInfo: object.Marshal(obj),
			},
		}

		res, err := routerClient.SendAction(context.TODO(), req)
		if err != nil {
			return fmt.Errorf("send action failed: %w", err)
		}

		if !res.Success {
			return fmt.Errorf("send action got unexpected response: %s", res.ErrMsg)
		}
	}
	// clear sent data
	sendObjects = make(map[string]object.Object)

	for _, objID := range removeObjects {
		req := &pb.Action{
			SessionID: sessionID,
			ClientID:  clientID,
			Type:      pb.Action_REMOVEOBJECT,
			Data:      &pb.Action_ObjectID{ObjectID: objID},
		}

		res, err := routerClient.SendAction(context.TODO(), req)
		if err != nil {
			return fmt.Errorf("remove object failed: %w", err)
		}

		if !res.Success {
			return fmt.Errorf("remove object got unexpected response: %s", res.ErrMsg)
		}
	}
	removeObjects = []string{}

	if len(sendDamages) > 0 {
		req := &pb.Action{
			SessionID: sessionID,
			ClientID:  clientID,
			Type:      pb.Action_NEWDAMAGE,
			Data: &pb.Action_DamageInfo{
				DamageInfo: damage.Marshal(sendDamages),
			},
		}

		res, err := routerClient.SendAction(context.TODO(), req)
		if err != nil {
			return fmt.Errorf("add damages failed: %w", err)
		}

		if !res.Success {
			return fmt.Errorf("add damages got unexpected response: %s", res.ErrMsg)
		}
		sendDamages = []damage.Damage{}
	}

	for _, eff := range sendEffects {
		req := &pb.Action{
			SessionID: sessionID,
			ClientID:  clientID,
			Type:      pb.Action_NEWEFFECT,
			Data: &pb.Action_Effect{
				Effect: effect.Marshal(eff),
			},
		}

		res, err := routerClient.SendAction(context.TODO(), req)
		if err != nil {
			return fmt.Errorf("add effect failed: %w", err)
		}

		if !res.Success {
			return fmt.Errorf("add effect got unexpected response: %s", res.ErrMsg)
		}
	}
	sendEffects = []effect.Effect{}

	return nil

}
