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
	publishInterval = 500 * time.Millisecond // debug
)

const (
	statusConnectWait int = iota
	statusChipSelectWait
	statusActing
)

type clientInfo struct {
	active    bool
	chipSent  bool
	clientID  string
	sendQueue chan *pb.Data
}

type session struct {
	sessionID string
	routeID   string
	clients   [2]clientInfo

	started   bool
	status    int
	fieldLock sync.Mutex
	fieldInfo map[string]*field.Info
}

var (
	sessionLock sync.Mutex
	sessionList = []*session{}
)

func Add(clientID string, sendQueue chan *pb.Data) (string, error) {
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
				se.clients[0].sendQueue = sendQueue
			} else if se.clients[1].clientID == clientID {
				se.clients[1].active = true
				se.clients[1].sendQueue = sendQueue
			}

			return se.sessionID, nil
		}
	}

	// no session in the list
	// so create new session
	sessionID := uuid.New().String()
	v := session{
		sessionID: sessionID,
		routeID:   route.ID,
		status:    statusConnectWait,
		started:   false,
		fieldInfo: make(map[string]*field.Info),
	}
	v.fieldInfo[route.Clients[0]] = &field.Info{}
	v.fieldInfo[route.Clients[0]].Init()
	v.fieldInfo[route.Clients[1]] = &field.Info{}
	v.fieldInfo[route.Clients[1]].Init()

	index := 0
	if route.Clients[1] == clientID {
		index = 1
	}

	v.clients[index] = clientInfo{
		active:    true,
		clientID:  route.Clients[index],
		sendQueue: sendQueue,
	}
	v.clients[1-index] = clientInfo{
		clientID: route.Clients[1-index],
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
	for _, s := range sessionList {
		if s.sessionID == action.SessionID {
			switch action.Type {
			case pb.Action_UPDATEOBJECT:
				s.fieldLock.Lock()
				defer s.fieldLock.Unlock()

				var obj field.Object
				field.UnmarshalObject(&obj, action.GetObjectInfo())
				for _, c := range s.clients {
					if c.clientID == action.ClientID {
						s.fieldInfo[c.clientID].MyArea[obj.X][obj.Y] = obj
					} else {
						s.fieldInfo[c.clientID].EnemyArea[obj.X][obj.Y] = obj
					}
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
			d := &pb.Data{
				Type: pb.Data_DATA,
				Data: &pb.Data_RawData{
					RawData: field.Marshal(s.fieldInfo[s.clients[0].clientID]),
				},
			}
			s.clients[0].sendQueue <- d

			d.Data = &pb.Data_RawData{
				RawData: field.Marshal(s.fieldInfo[s.clients[1].clientID]),
			}
			s.clients[1].sendQueue <- d
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

				s.clients[0].sendQueue <- d
				s.clients[1].sendQueue <- d
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

				s.clients[0].sendQueue <- d
				s.clients[1].sendQueue <- d
				s.clients[0].chipSent = false
				s.clients[1].chipSent = false
				s.status = statusActing
			}

		case statusActing:
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
