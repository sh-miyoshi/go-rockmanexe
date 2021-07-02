package session

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/fps"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/db"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/field"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/net/routerpb"
)

const (
	publishInterval = 50 * time.Millisecond // debug
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
	dmMgr      damage.Manager
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
	v.clients[index].fieldInfo.Init()
	v.clients[1-index] = clientInfo{
		clientID:  route.Clients[1-index],
		fieldInfo: &field.Info{},
	}
	v.clients[1-index].fieldInfo.Init()

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
	logAction(action)

	for _, s := range sessionList {
		if s.sessionID == action.SessionID {
			switch action.Type {
			case pb.Action_UPDATEOBJECT:
				var obj field.Object
				field.UnmarshalObject(&obj, action.GetObjectInfo())
				s.fieldLock.Lock()
				s.updateObject(obj, action.ClientID)
				s.fieldLock.Unlock()
				for _, c := range s.clients {
					logger.Debug("Updated objects for Client %s: %+v", c.clientID, c.fieldInfo.Objects)
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
			case pb.Action_REMOVEOBJECT:
				id := action.GetObjectID()
				s.fieldLock.Lock()
				s.removeObject(id)
				s.fieldLock.Unlock()
			case pb.Action_NEWDAMAGE:
				var damages []damage.Damage
				damage.Unmarshal(&damages, action.GetDamageInfo())
				s.dmMgr.Add(damages)
				logger.Debug("Added damges: %+v", damages)
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

	// run process per frame
	go func() {
		fpsMgr := fps.Fps{TargetFPS: 60}
		for {
			// damage process
			// TODO
			// for i, obj := range s.fieldInfo.Objects {
			// 	if dm := s.dmMgr.Hit(obj.X, obj.Y); dm != nil {
			// 		s.fieldLock.Lock()
			// 		s.fieldInfo.Objects[i].DamageChecked = false
			// 		s.fieldInfo.Objects[i].HitDamage = *dm
			// 		s.fieldLock.Unlock()
			// 	}
			// }
			s.dmMgr.Update()

			fpsMgr.Wait()
		}
	}()

	// publish via data stream
	for {
		before := time.Now().UnixNano() / (1000 * 1000)

		// Field data send
		if s.status == statusChipSelectWait || s.status == statusActing {
			s.publishField()
		}

		s.statusUpdate()

		after := time.Now().UnixNano() / (1000 * 1000)
		time.Sleep(publishInterval - time.Duration(after-before))
	}
}

func (s *session) statusUpdate() {
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
			s.changeStatus(statusChipSelectWait)

			// send initial data immediately
			s.publishField()
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
			s.changeStatus(statusActing)
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

			s.changeStatus(s.nextStatus)
			s.nextStatus = -1
		}

		// TODO
		// if s.clients[0 or 1].SendAction(Win or Lose?)
		//   s.clients[0 and? 1].sendQueue <- statusGameEnd
		//   remove(session)?
	}
}

func (s *session) changeStatus(next int) {
	logger.Info("Change state from %d to %d", s.status, next)
	s.status = next
}

func (s *session) publishField() {
	now := time.Now()

	for i := 0; i < len(s.clients); i++ {
		s.clients[i].fieldInfo.CurrentTime = now
		d := &pb.Data{
			Type: pb.Data_DATA,
			Data: &pb.Data_RawData{
				RawData: field.Marshal(s.clients[i].fieldInfo),
			},
		}

		if err := s.clients[i].dataStream.Send(d); err != nil {
			s.exitErr <- fmt.Errorf("field info send failed for client %s: %w", s.clients[i].clientID, err)
		}
	}
}

func (s *session) updateObject(obj field.Object, clientID string) {
	obj.ClientID = clientID
	if obj.UpdateBaseTime {
		obj.BaseTime = time.Now()
	}

	cindex := -1
	for i, c := range s.clients {
		if c.clientID == clientID {
			cindex = i
			break
		}
	}
	if cindex == -1 {
		panic(fmt.Sprintf("no such client %s", clientID))
	}

	updated := false
	for i, o := range s.clients[cindex].fieldInfo.Objects {
		if o.ID == obj.ID {
			if !obj.DamageChecked {
				obj.HitDamage = o.HitDamage
			}
			if !obj.UpdateBaseTime {
				obj.BaseTime = o.BaseTime
			}
			obj.UpdateBaseTime = false
			s.clients[cindex].fieldInfo.Objects[i] = obj
			updated = true
			break
		}
	}

	if !updated {
		s.clients[cindex].fieldInfo.Objects = append(s.clients[cindex].fieldInfo.Objects, obj)
	}

	updated = false
	obj.X = field.SizeX - obj.X - 1
	for i, o := range s.clients[1-cindex].fieldInfo.Objects {
		if o.ID == obj.ID {
			if !obj.DamageChecked {
				obj.HitDamage = o.HitDamage
			}
			if !obj.UpdateBaseTime {
				obj.BaseTime = o.BaseTime
			}
			obj.UpdateBaseTime = false
			s.clients[1-cindex].fieldInfo.Objects[i] = obj
			updated = true
			break
		}
	}

	if !updated {
		s.clients[1-cindex].fieldInfo.Objects = append(s.clients[1-cindex].fieldInfo.Objects, obj)
	}
}

func (s *session) removeObject(objID string) {
	for i, c := range s.clients {
		newObjs := []field.Object{}
		for _, obj := range c.fieldInfo.Objects {
			if obj.ID != objID {
				newObjs = append(newObjs, obj)
			}
		}
		s.clients[i].fieldInfo.Objects = newObjs
	}
}

func logAction(action *pb.Action) {
	msg := fmt.Sprintf(
		"{ SessionID: %s, ClientID: %s, Type: %s, Action: { ",
		action.SessionID,
		action.ClientID,
		action.Type.String(),
	)

	switch action.Type {
	case pb.Action_UPDATEOBJECT:
		var obj field.Object
		field.UnmarshalObject(&obj, action.GetObjectInfo())
		msg += fmt.Sprintf("Objects: %+v", obj)
	case pb.Action_SENDSIGNAL:
		msg += fmt.Sprintf("Signal: %s", action.GetSignal().String())
	case pb.Action_REMOVEOBJECT:
		msg += "TargetObject: " + action.GetObjectID()
	case pb.Action_NEWDAMAGE:
		var dm []damage.Damage
		damage.Unmarshal(&dm, action.GetDamageInfo())
		msg += fmt.Sprintf("Damages: %+v", dm)
	}

	msg += " }}"
	logger.Debug("Got action: %s", msg)
}
