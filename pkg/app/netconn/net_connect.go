package netconn

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
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

type sendInfo struct {
	objects       map[string]object.Object
	removeObjects []string
	damages       []damage.Damage
	effects       []effect.Effect
	sounds        []sound.SEType
}

var (
	// variables for router connection
	conn         *grpc.ClientConn
	routerClient pb.RouterClient
	dataStream   pb.Router_PublishDataClient
	sessionID    string
	clientID     string

	// variables for application data
	status     pb.Data_Status
	fieldInfo  field.Info
	sendData   sendInfo
	allUserIDs []string

	// variables for system control
	exitErr   error
	fieldLock sync.Mutex
)

func Connect(conf Config) error {
	initVars()

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
	res := authRes.GetAuthRes()
	if !res.Success {
		return fmt.Errorf("failed to auth request: %s", res.ErrMsg)
	}

	sessionID = res.SessionID
	clientID = conf.ClientID
	allUserIDs = append([]string{}, res.AllUserIDs...)
	sendData.Init()

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
	sendData.objects[obj.ID] = obj
}

func RemoveObject(objID string) {
	sendData.removeObjects = append(sendData.removeObjects, objID)
}

func SendDamages(damages []damage.Damage) {
	for _, dm := range damages {
		dm.ClientID = clientID
		sendData.damages = append(sendData.damages, dm)
	}
}

func SendEffect(eff effect.Effect) {
	eff.ClientID = clientID
	sendData.effects = append(sendData.effects, eff)
}

func AddSound(se sound.SEType) {
	sendData.sounds = append(sendData.sounds, se)
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
			fieldInfo.CurrentTime = f.CurrentTime
			fieldInfo.Objects = f.Objects
			fieldInfo.Panels = f.Panels
			fieldInfo.Effects = append(fieldInfo.Effects, f.Effects...)
			fieldInfo.HitDamages = append(fieldInfo.HitDamages, f.HitDamages...)
			fieldInfo.Sounds = append(fieldInfo.Sounds, f.Sounds...)
			fieldLock.Unlock()
		default:
			exitErr = fmt.Errorf("invalid data type was received: %d", data.Type)
			return
		}
	}
}

func GetOpponentUserID() string {
	for _, rawID := range allUserIDs {
		t := strings.Split(rawID, ":")
		if len(t) != 2 {
			logger.Error("User ID data maybe broken: %v", allUserIDs)
			continue
		}
		cid := t[0]
		uid := t[1]
		if cid == clientID {
			return uid
		}
	}

	logger.Error("Failed to get opponent user id in %v", allUserIDs)
	return ""
}

func GetStatus() pb.Data_Status {
	return status
}

func GetFieldInfo() field.Info {
	return fieldInfo
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

func RemoveDamage(id string) {
	fieldLock.Lock()
	defer fieldLock.Unlock()
	for i, dm := range fieldInfo.HitDamages {
		if dm.ID == id {
			fieldInfo.HitDamages = append(fieldInfo.HitDamages[:i], fieldInfo.HitDamages[i+1:]...)
			return
		}
	}
}

func RemoveSounds() {
	fieldLock.Lock()
	fieldInfo.Sounds = []int32{}
	fieldLock.Unlock()
}

func BulkSendFieldInfo() error {
	if exitErr != nil {
		return fmt.Errorf("already exit in recv: %w", exitErr)
	}

	// TODO refactoring
	// Send objects
	for _, obj := range sendData.objects {
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

	for _, objID := range sendData.removeObjects {
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

	if len(sendData.damages) > 0 {
		req := &pb.Action{
			SessionID: sessionID,
			ClientID:  clientID,
			Type:      pb.Action_NEWDAMAGE,
			Data: &pb.Action_DamageInfo{
				DamageInfo: damage.Marshal(sendData.damages),
			},
		}

		res, err := routerClient.SendAction(context.TODO(), req)
		if err != nil {
			return fmt.Errorf("add damages failed: %w", err)
		}

		if !res.Success {
			return fmt.Errorf("add damages got unexpected response: %s", res.ErrMsg)
		}
	}

	for _, eff := range sendData.effects {
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

	for _, se := range sendData.sounds {
		req := &pb.Action{
			SessionID: sessionID,
			ClientID:  clientID,
			Type:      pb.Action_ADDSOUND,
			Data: &pb.Action_SeType{
				SeType: int32(se),
			},
		}

		res, err := routerClient.SendAction(context.TODO(), req)
		if err != nil {
			return fmt.Errorf("add sound failed: %w", err)
		}

		if !res.Success {
			return fmt.Errorf("add sound got unexpected response: %s", res.ErrMsg)
		}
	}

	sendData.Init()
	return nil
}

func (i *sendInfo) Init() {
	i.objects = make(map[string]object.Object)
	i.removeObjects = []string{}
	i.damages = []damage.Damage{}
	i.effects = []effect.Effect{}
	i.sounds = []sound.SEType{}
}

func initVars() {
	Disconnect()
	sessionID = ""
	clientID = ""
	status = pb.Data_CONNECTWAIT
	fieldInfo = field.Info{}
	sendData = sendInfo{}
	allUserIDs = []string{}
	exitErr = nil
}
