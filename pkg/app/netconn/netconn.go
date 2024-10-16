package netconn

import (
	"context"
	"crypto/tls"
	"strings"
	"sync"

	"github.com/cockroachdb/errors"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/net/netconnpb"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/sysinfo"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/router/gameinfo"
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

type cutinInfo struct {
	IsSet bool
	Info  sysinfo.Cutin
}

type NetConn struct {
	config        Config
	conn          *grpc.ClientConn
	dataStream    pb.NetConn_TransDataClient
	connectStatus ConnectStatus
	sessionID     string
	allUserIDs    []string

	gameStatus    pb.Response_Status
	gameInfo      gameinfo.GameInfo
	gameInfoMu    sync.Mutex
	gameCutinInfo cutinInfo
}

func New(conf Config) *NetConn {
	return &NetConn{
		config:     conf,
		gameStatus: pb.Response_CONNECTWAIT,
		connectStatus: ConnectStatus{
			Status: ConnStateWaiting,
			Error:  nil,
		},
	}
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

func (n *NetConn) GetGameStatus() pb.Response_Status {
	return n.gameStatus
}

func (n *NetConn) GetGameInfo() gameinfo.GameInfo {
	return n.gameInfo
}

func (n *NetConn) PopCutinInfo() (sysinfo.Cutin, bool) {
	if n.gameCutinInfo.IsSet {
		res := n.gameCutinInfo.Info
		n.gameCutinInfo = cutinInfo{
			IsSet: false,
		}
		return res, true
	}
	return sysinfo.Cutin{}, false
}

func (n *NetConn) CleanupSounds() {
	n.gameInfoMu.Lock()
	n.gameInfo.Sounds = []gameinfo.Sound{}
	n.gameInfoMu.Unlock()
}

func (n *NetConn) SendSignal(signalType pb.Request_SignalType, data []byte) error {
	return n.dataStream.Send(&pb.Request{
		SessionID: n.sessionID,
		ClientID:  n.config.ClientID,
		Type:      pb.Request_SENDSIGNAL,
		Data:      &pb.Request_Signal_{Signal: &pb.Request_Signal{Type: signalType, RawData: data}},
	})
}

func (n *NetConn) SendAction(actType pb.Request_ActionType, data []byte) error {
	return n.dataStream.Send(&pb.Request{
		SessionID: n.sessionID,
		ClientID:  n.config.ClientID,
		Type:      pb.Request_ACTION,
		Data: &pb.Request_Act{Act: &pb.Request_Action{
			Type:    actType,
			RawData: data,
		}},
	})
}

func (n *NetConn) GetOpponentUserID() string {
	for _, rawID := range n.allUserIDs {
		t := strings.Split(rawID, ":")
		if len(t) != 2 {
			logger.Error("User ID data maybe broken: %v", n.allUserIDs)
			continue
		}
		cid := t[0]
		uid := t[1]
		if cid == n.config.ClientID {
			return uid
		}
	}

	logger.Error("Failed to get opponent user id in %v", n.allUserIDs)
	return ""
}

func (n *NetConn) connect() error {
	var err error
	n.conn, err = newConn(n.config)
	if err != nil {
		return errors.Wrap(err, "failed to dial router")
	}

	client := pb.NewNetConnClient(n.conn)
	n.dataStream, err = client.TransData(context.Background())
	if err != nil {
		return errors.Wrap(err, "failed to create data stream")
	}

	// Authenticate client
	req := makeAuthReq(n.config.ClientID, n.config.ClientKey, n.config.ProgramVersion)
	if err := n.dataStream.Send(req); err != nil {
		return errors.Wrap(err, "failed to send authn request")
	}

	res, err := n.dataStream.Recv()
	if err != nil {
		return errors.Wrap(err, "failed to recv authn response")
	}
	authRes := res.GetAuthRes()
	if !authRes.Success {
		return errors.Newf("failed to authenticate client: %s", authRes.ErrMsg)
	}
	n.sessionID = authRes.SessionID
	n.allUserIDs = append([]string{}, authRes.AllUserIDs...)

	return nil
}

func (n *NetConn) dataRecv() {
	// Recv data from stream
	for {
		data, err := n.dataStream.Recv()
		if err != nil {
			n.connectStatus = ConnectStatus{
				Status: ConnStateError,
				Error:  errors.Wrap(err, "failed to recv data"),
			}
			return
		}

		switch data.Type {
		case pb.Response_UPDATESTATUS:
			logger.Debug("got status update data: %+v", data)
			n.gameStatus = data.GetStatus()
		case pb.Response_DATA:
			var info gameinfo.GameInfo
			info.Unmarshal(data.GetRawData())
			n.gameInfoMu.Lock()
			n.gameInfo = info
			n.gameInfoMu.Unlock()
		case pb.Response_SYSTEM:
			logger.Debug("got system info: %+v", data)
			sys := data.GetSystem()
			switch sys.GetType() {
			case pb.Response_System_CUTIN:
				var cutin sysinfo.Cutin
				cutin.Unmarshal(sys.GetRawData())
				n.gameCutinInfo = cutinInfo{
					IsSet: true,
					Info:  cutin,
				}
			default:
				logger.Info("invalid data was received, ignore this")
			}
		default:
			n.connectStatus = ConnectStatus{
				Status: ConnStateError,
				Error:  errors.Newf("invalid data type was received: %d", data.Type),
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

func makeAuthReq(id, key, version string) *pb.Request {
	return &pb.Request{
		Type: pb.Request_AUTHENTICATE,
		Data: &pb.Request_Req{
			Req: &pb.Request_AuthRequest{
				Id:      id,
				Key:     key,
				Version: version,
			},
		},
	}
}
