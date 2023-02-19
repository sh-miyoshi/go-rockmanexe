package netconn

import (
	"context"
	"crypto/tls"
	"fmt"
	"strings"
	"sync"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/newnet/netconnpb"
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

type NetConn struct {
	config        Config
	conn          *grpc.ClientConn
	dataStream    pb.NetConn_TransDataClient
	connectStatus ConnectStatus
	sessionID     string
	allUserIDs    []string

	gameStatus pb.Response_Status
	gameInfoMu sync.Mutex
}

func New(conf Config) *NetConn {
	res := &NetConn{}

	res.config = conf
	res.connectStatus = ConnectStatus{
		Status: ConnStateWaiting,
		Error:  nil,
	}

	return res
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

func (n *NetConn) SendSignal(signalType pb.Request_SignalType, data []byte) error {
	return n.dataStream.Send(&pb.Request{
		SessionID: n.sessionID,
		ClientID:  n.config.ClientID,
		Type:      pb.Request_SENDSIGNAL,
		Data:      &pb.Request_Signal_{Signal: &pb.Request_Signal{Type: signalType, RawData: data}},
	})
}

func (n *NetConn) SendAction(actType pb.Request_ActionType, currentPos common.Point, data []byte) error {
	return n.dataStream.Send(&pb.Request{
		SessionID: n.sessionID,
		ClientID:  n.config.ClientID,
		Type:      pb.Request_ACTION,
		Data: &pb.Request_Act{Act: &pb.Request_Action{
			Type:    actType,
			X:       int64(currentPos.X),
			Y:       int64(currentPos.Y),
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
				Error:  fmt.Errorf("failed to recv data: %w", err),
			}
			return
		}

		switch data.Type {
		case pb.Response_UPDATESTATUS:
			logger.Debug("got status update data: %+v", data)
			n.gameStatus = data.GetStatus()
		case pb.Response_DATA:
			n.gameInfoMu.Lock()
			// TODO
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
