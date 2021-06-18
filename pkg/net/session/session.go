package session

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/db"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/field"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/net/routerpb"
)

const (
	publishInterval = 100 * time.Millisecond // debug
)

const (
	statusConnectWait int = iota
	statusChipSelectWait
	statusActing
)

type clientInfo struct {
	active     bool
	chipSent   bool
	clientID   string
	fieldInfo  *field.Info
	dataStream pb.Router_PublishDataServer
}

type session struct {
	sessionID  string
	routeID    string
	clients    [2]clientInfo
	started    bool
	status     int
	nextStatus int
	fieldLock  sync.Mutex
	exitErr    chan error
}

var (
	sessionLock sync.Mutex
	sessionList = []*session{}
)

func Add(clientID string, dataStream pb.Router_PublishDataServer, exitErr chan error) (string, error) {
	route, err := db.GetInst().RouteGetByClient(clientID)
	if err != nil {
		return "", fmt.Errorf("route get failed: %v", err)
	}

	sessionLock.Lock()
	defer sessionLock.Unlock()

	for _, se := range sessionList {
		if se.routeID == route.ID {
			if se.clients[0].clientID == clientID {
				se.clients[0].active = true
				se.clients[0].dataStream = dataStream
			} else if se.clients[1].clientID == clientID {
				se.clients[1].active = true
				se.clients[1].dataStream = dataStream
			}

			return se.sessionID, nil
		}
	}

	// no session in the list
	// so create new session
	sessionID := uuid.New().String()
	v := session{
		sessionID:  sessionID,
		routeID:    route.ID,
		status:     statusConnectWait,
		nextStatus: -1,
		started:    false,
		exitErr:    exitErr,
	}

	index := 0
	if route.Clients[1] == clientID {
		index = 1
	}

	v.clients[index] = clientInfo{
		active:     true,
		clientID:   route.Clients[index],
		dataStream: dataStream,
		fieldInfo:  &field.Info{},
	}
	v.clients[1-index] = clientInfo{
		clientID:  route.Clients[1-index],
		fieldInfo: &field.Info{},
	}
	for _, c := range v.clients {
		c.fieldInfo.Init()
	}

	sessionList = append(sessionList, &v)

	return sessionID, nil
}

func Run(sessionID string) {
	sessionLock.Lock()
	defer sessionLock.Unlock()

	// Run method will be called after session.Add()
	// so the target session is almost in the last.
	for i := len(sessionList) - 1; i >= 0; i-- {
		if sessionList[i].sessionID == sessionID {
			if !sessionList[i].started {
				sessionList[i].started = true
				go sessionList[i].Process()
			}
			return
		}
	}
}

func ActionProc(action *pb.Action) error {
	logger.Debug("Got action: %+v", action)

	for _, s := range sessionList {
		if s.sessionID == action.SessionID {
			switch action.Type {
			case pb.Action_UPDATEOBJECT:
				s.fieldLock.Lock()
				defer s.fieldLock.Unlock()

				var obj field.Object
				field.UnmarshalObject(&obj, action.GetObjectInfo())
				obj.ClientID = action.ClientID
				if obj.UpdateBaseTime {
					obj.BaseTime = time.Now()
					obj.UpdateBaseTime = false
				}

				for i := 0; i < 2; i++ {
					updateObject(&s.clients[i].fieldInfo.Objects, obj)
				}
			case pb.Action_SENDSIGNAL:
				switch action.GetSignal() {
				case pb.Action_CHIPSEND:
					for i, c := range s.clients {
						if c.clientID == action.ClientID {
							s.clients[i].chipSent = true
							break
						}
					}
				case pb.Action_GOCHIPSELECT:
					s.nextStatus = statusChipSelectWait
				}
			default:
				return fmt.Errorf("action %d is not implemented yet", action.Type)
			}
			return nil
		}
	}

	return fmt.Errorf("no such session")
}

func (s *session) Process() {
	logger.Info("start new session for route %s", s.routeID)
	logger.Debug("client info: %+v", s.clients)

	// publish via data stream
	for {
		time.Sleep(publishInterval)

		// Field data send
		if s.status == statusChipSelectWait || s.status == statusActing {
			now := time.Now()

			s.clients[0].fieldInfo.CurrentTime = now
			d := &pb.Data{
				Type: pb.Data_DATA,
				Data: &pb.Data_RawData{
					RawData: field.Marshal(s.clients[0].fieldInfo),
				},
			}
			if err := s.clients[0].dataStream.Send(d); err != nil {
				s.exitErr <- fmt.Errorf("field info send failed for client %s: %w", s.clients[0].clientID, err)
			}

			s.clients[1].fieldInfo.CurrentTime = now
			d.Data = &pb.Data_RawData{
				RawData: field.Marshal(s.clients[1].fieldInfo),
			}
			if err := s.clients[1].dataStream.Send(d); err != nil {
				s.exitErr <- fmt.Errorf("field info send failed for client %s: %w", s.clients[1].clientID, err)
			}
		}

		switch s.status {
		case statusConnectWait:
			// check ready
			if s.clients[0].active && s.clients[1].active {
				d := &pb.Data{
					Type: pb.Data_UPDATESTATUS,
					Data: &pb.Data_Status_{
						Status: pb.Data_CHIPSELECTWAIT,
					},
				}

				if err := s.clients[0].dataStream.Send(d); err != nil {
					s.exitErr <- fmt.Errorf("update status send failed for client %s: %w", s.clients[0].clientID, err)
				}
				if err := s.clients[1].dataStream.Send(d); err != nil {
					s.exitErr <- fmt.Errorf("update status send failed for client %s: %w", s.clients[1].clientID, err)
				}
				s.status = statusChipSelectWait
			}
		case statusChipSelectWait:
			if s.clients[0].chipSent && s.clients[1].chipSent {
				d := &pb.Data{
					Type: pb.Data_UPDATESTATUS,
					Data: &pb.Data_Status_{
						Status: pb.Data_ACTING,
					},
				}

				if err := s.clients[0].dataStream.Send(d); err != nil {
					s.exitErr <- fmt.Errorf("update status send failed for client %s: %w", s.clients[0].clientID, err)
				}
				if err := s.clients[1].dataStream.Send(d); err != nil {
					s.exitErr <- fmt.Errorf("update status send failed for client %s: %w", s.clients[1].clientID, err)
				}
				s.clients[0].chipSent = false
				s.clients[1].chipSent = false
				s.status = statusActing
			}

		case statusActing:
			if s.nextStatus != -1 {
				sendSt := pb.Data_CONNECTWAIT
				switch s.nextStatus {
				case statusChipSelectWait:
					sendSt = pb.Data_CHIPSELECTWAIT
				default:
					s.exitErr <- fmt.Errorf("invalid next status: %d", s.nextStatus)
					return
				}

				d := &pb.Data{
					Type: pb.Data_UPDATESTATUS,
					Data: &pb.Data_Status_{
						Status: sendSt,
					},
				}
				if err := s.clients[0].dataStream.Send(d); err != nil {
					s.exitErr <- fmt.Errorf("update status send failed for client %s: %w", s.clients[0].clientID, err)
				}
				if err := s.clients[1].dataStream.Send(d); err != nil {
					s.exitErr <- fmt.Errorf("update status send failed for client %s: %w", s.clients[1].clientID, err)
				}

				s.status = s.nextStatus
				s.nextStatus = -1
			}

			// TODO
			// if s.clients[0 or 1].SendAction(MoveToChipSel)
			//   s.clients[0 and 1].sendQueue <- statusChipSelectWait
			//   s.status = statusChipSelectWait
			// else if s.clients[0 or 1].SendAction(Win or Lose?)
			//   s.clients[0 and? 1].sendQueue <- statusGameEnd
			//   remove(session)?
		}
	}
}

func updateObject(objs *[]field.Object, obj field.Object) {
	for i, o := range *objs {
		if o.ID == obj.ID {
			(*objs)[i] = obj
			return
		}
	}

	*objs = append(*objs, obj)
}
