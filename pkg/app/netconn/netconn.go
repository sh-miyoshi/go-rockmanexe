package netconn

import (
	"context"
	"crypto/tls"
	"fmt"
	"sync"
	"time"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/effect"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/net/netconnpb"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/object"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/session"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	ConnStateWaiting = iota
	ConnStateOK
	ConnStateError
)

type Config struct {
	StreamAddr     string
	ClientID       string
	ClientKey      string
	ProgramVersion string
	Insecure       bool
}

type ConnectStatus struct {
	Status int
	Error  error
}

type sendObject struct {
	objects       map[string]object.Object
	removeObjects []object.Object
	damages       []damage.Damage
	effects       []effect.Effect
}

type NetConn struct {
	config        Config
	conn          *grpc.ClientConn
	dataStream    pb.NetConn_TransDataClient
	connectStatus ConnectStatus
	sessionID     string

	gameStatus pb.Data_Status
	gameInfo   session.GameInfo
	gameInfoMu sync.Mutex
	sendInfo   sendObject
}

var (
	inst NetConn
)

func Init(conf Config) {
	inst.config = conf
	inst.connectStatus = ConnectStatus{
		Status: ConnStateWaiting,
		Error:  nil,
	}
	inst.sendInfo = sendObject{
		objects: make(map[string]object.Object),
	}
}

func GetInst() *NetConn {
	return &inst
}

func (n *NetConn) GetConnStatus() ConnectStatus {
	return n.connectStatus
}

func (n *NetConn) ConnectRequest() {
	go func() {
		if err := n.connect(); err != nil {
			n.connectStatus = ConnectStatus{
				Status: ConnStateError,
				Error:  err,
			}
			return
		}
		logger.Info("Successfully connected to router")

		// Start data receiver
		go n.dataRecv()

		n.connectStatus = ConnectStatus{
			Status: ConnStateOK,
		}
	}()
}

func (n *NetConn) Disconnect() {
	if n.conn != nil {
		n.conn.Close()
		n.conn = nil
	}
}

func (n *NetConn) GetGameStatus() pb.Data_Status {
	return n.gameStatus
}

func (n *NetConn) GetGameInfo() session.GameInfo {
	return n.gameInfo
}

func (n *NetConn) SendSignal(signal pb.Action_SignalType) error {
	return n.dataStream.Send(&pb.Action{
		SessionID: n.sessionID,
		ClientID:  n.config.ClientID,
		Type:      pb.Action_SENDSIGNAL,
		Data:      &pb.Action_Signal{Signal: signal},
	})
}

func (n *NetConn) SendObject(obj object.Object) {
	obj.ClientID = n.config.ClientID
	n.sendInfo.objects[obj.ID] = obj
}

func (n *NetConn) RemoveObject(objID string) {
	obj := object.Object{
		ID: objID,
	}
	n.sendInfo.removeObjects = append(n.sendInfo.removeObjects, obj)
}

func (n *NetConn) AddDamage(dm damage.Damage) {
	dm.ClientID = n.config.ClientID
	n.sendInfo.damages = append(n.sendInfo.damages, dm)
}

func (n *NetConn) RemoveDamage(id string) {
	n.gameInfoMu.Lock()
	defer n.gameInfoMu.Unlock()

	for i, dm := range n.gameInfo.HitDamages {
		if dm.ID == id {
			n.gameInfo.HitDamages = append(n.gameInfo.HitDamages[:i], n.gameInfo.HitDamages[i+1:]...)
			return
		}
	}
}

func (n *NetConn) SendEffect(eff effect.Effect) {
	eff.ClientID = n.config.ClientID
	n.sendInfo.effects = append(n.sendInfo.effects, eff)
}

func (n *NetConn) RemoveEffect(id string) {
	n.gameInfoMu.Lock()
	defer n.gameInfoMu.Unlock()

	for i, eff := range n.gameInfo.Effects {
		if eff.ID == id {
			n.gameInfo.Effects = append(n.gameInfo.Effects[:i], n.gameInfo.Effects[i+1:]...)
			return
		}
	}
}

func (n *NetConn) BulkSendData() error {
	// TODO 一度の通信で送る

	// Send objects
	for _, obj := range n.sendInfo.objects {
		req := &pb.Action{
			SessionID: n.sessionID,
			ClientID:  n.config.ClientID,
			Type:      pb.Action_UPDATEOBJECT,
			Data: &pb.Action_ObjectInfo{
				ObjectInfo: object.Marshal(obj),
			},
		}

		if err := n.dataStream.Send(req); err != nil {
			return fmt.Errorf("send object failed: %w", err)
		}
	}

	// Send remove objects
	for _, obj := range n.sendInfo.removeObjects {
		req := &pb.Action{
			SessionID: n.sessionID,
			ClientID:  n.config.ClientID,
			Type:      pb.Action_REMOVEOBJECT,
			Data: &pb.Action_ObjectInfo{
				ObjectInfo: object.Marshal(obj),
			},
		}

		if err := n.dataStream.Send(req); err != nil {
			return fmt.Errorf("remove object failed: %w", err)
		}
	}

	// Send damages
	if len(n.sendInfo.damages) > 0 {
		req := &pb.Action{
			SessionID: n.sessionID,
			ClientID:  n.config.ClientID,
			Type:      pb.Action_ADDDAMAGE,
			Data: &pb.Action_DamageInfo{
				DamageInfo: damage.Marshal(n.sendInfo.damages),
			},
		}

		if err := n.dataStream.Send(req); err != nil {
			return fmt.Errorf("send damage failed: %w", err)
		}
	}

	// Send effects
	for _, eff := range n.sendInfo.effects {
		req := &pb.Action{
			SessionID: n.sessionID,
			ClientID:  n.config.ClientID,
			Type:      pb.Action_ADDEFFECT,
			Data: &pb.Action_Effect{
				Effect: effect.Marshal(eff),
			},
		}

		if err := n.dataStream.Send(req); err != nil {
			return fmt.Errorf("add effect failed: %w", err)
		}
	}

	n.sendInfo.Init()
	return nil
}

func (n *NetConn) UpdateDataCount() {
	n.gameInfoMu.Lock()
	defer n.gameInfoMu.Unlock()
	for i, obj := range n.gameInfo.Objects {
		if obj.Count == 0 {
			tm := n.gameInfo.CurrentTime.Sub(obj.BaseTime)
			obj.Count = int(tm * 60 / time.Second)
		} else {
			obj.Count++
		}
		n.gameInfo.Objects[i] = obj
	}
	for i, eff := range n.gameInfo.Effects {
		eff.Count++
		n.gameInfo.Effects[i] = eff
	}
}

func (n *NetConn) connect() error {
	var err error
	n.conn, err = newConn(n.config)
	if err != nil {
		return fmt.Errorf("failed to dial router: %w", err)
	}

	client := pb.NewNetConnClient(n.conn)
	n.dataStream, err = client.TransData(context.Background())
	if err != nil {
		return fmt.Errorf("failed to create data stream: %w", err)
	}

	// Authenticate client
	req := makeAuthReq(n.config.ClientID, n.config.ClientKey, n.config.ProgramVersion)
	if err := n.dataStream.Send(req); err != nil {
		return fmt.Errorf("failed to send authn request: %w", err)
	}

	res, err := n.dataStream.Recv()
	if err != nil {
		return fmt.Errorf("failed to recv authn response: %w", err)
	}
	authRes := res.GetAuthRes()
	if !authRes.Success {
		return fmt.Errorf("failed to authenticate client: %s", authRes.ErrMsg)
	}
	n.sessionID = authRes.SessionID

	return nil
}

func (n *NetConn) dataRecv() {
	// Recv data from stream
	for {
		data, err := n.dataStream.Recv()
		if err != nil {
			n.connectStatus = ConnectStatus{
				Status: ConnStateError,
				Error:  fmt.Errorf("failed to recv data: %w", err),
			}
			return
		}

		switch data.Type {
		case pb.Data_UPDATESTATUS:
			logger.Debug("got status update data: %+v", data)
			n.gameStatus = data.GetStatus()
		case pb.Data_DATA:
			var receivedInfo session.GameInfo
			b := data.GetRawData()
			receivedInfo.Unmarshal(b)
			n.gameInfoMu.Lock()
			n.gameInfo.CurrentTime = receivedInfo.CurrentTime
			n.gameInfo.Objects = receivedInfo.Objects
			n.gameInfo.Panels = receivedInfo.Panels
			n.gameInfo.Effects = append(n.gameInfo.Effects, receivedInfo.Effects...)
			n.gameInfo.HitDamages = append(n.gameInfo.HitDamages, receivedInfo.HitDamages...)
			n.gameInfo.Sounds = append(n.gameInfo.Sounds, receivedInfo.Sounds...)
			n.gameInfoMu.Unlock()
		default:
			n.connectStatus = ConnectStatus{
				Status: ConnStateError,
				Error:  fmt.Errorf("invalid data type was received: %d", data.Type),
			}
			return
		}
	}
}

func newConn(conf Config) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithAuthority(conf.StreamAddr))

	if conf.Insecure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		cred := credentials.NewTLS(&tls.Config{})
		opts = append(opts, grpc.WithTransportCredentials(cred))
	}

	return grpc.Dial(conf.StreamAddr, opts...)
}

func makeAuthReq(id, key, version string) *pb.Action {
	return &pb.Action{
		Type: pb.Action_AUTHENTICATE,
		Data: &pb.Action_Req{
			Req: &pb.Action_AuthRequest{
				Id:      id,
				Key:     key,
				Version: version,
			},
		},
	}
}

func (o *sendObject) Init() {
	o.objects = make(map[string]object.Object)
	o.removeObjects = []object.Object{}
	o.damages = []damage.Damage{}
	o.effects = []effect.Effect{}
}
