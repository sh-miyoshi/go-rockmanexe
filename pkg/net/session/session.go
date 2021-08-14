package session

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/fps"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/db"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/object"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/net/routerpb"
)

const (
	publishInterval = 100 * time.Millisecond // debug
)

const (
	statusConnectWait int = iota
	statusChipSelectWait
	statusActing
	statusGameEnd
)

type sessionError struct {
	generatorClientID string
	reason            error
}

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
	exitErr    chan sessionError
}

var (
	errSendFailed = errors.New("send failed")

	sessionLock sync.Mutex
	sessionList = make(map[string]*session)
)

func Add(clientID string, dataStream pb.Router_PublishDataServer) (string, error) {
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
		exitErr:    make(chan sessionError),
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

	sessionList[sessionID] = &v

	return sessionID, nil
}

func Run(sessionID string) {
	s := sessionList[sessionID]
	logger.Info("start new session for route %s", s.routeID)
	logger.Debug("client info: %+v", s.clients)
	sessionLock.Lock()
	s.started = true
	sessionLock.Unlock()
	s.mainProcess()
}

func ActionProc(action *pb.Action) error {
	logAction(action)

	s, ok := sessionList[action.SessionID]
	if !ok {
		return fmt.Errorf("no such session")
	}

	switch action.Type {
	case pb.Action_UPDATEOBJECT:
		var obj object.Object
		object.Unmarshal(&obj, action.GetObjectInfo())
		s.fieldLock.Lock()
		for i := 0; i < len(s.clients); i++ {
			myObj := s.clients[i].clientID == action.ClientID
			updateObject(&s.clients[i].fieldInfo.Objects, obj, action.ClientID, myObj)
		}
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
		case pb.Action_PLAYERDEAD:
			s.nextStatus = statusGameEnd
		}
	case pb.Action_REMOVEOBJECT:
		id := action.GetObjectID()
		s.fieldLock.Lock()
		for i := 0; i < len(s.clients); i++ {
			removeObject(&s.clients[i].fieldInfo.Objects, id)
		}
		s.fieldLock.Unlock()
	case pb.Action_NEWDAMAGE:
		var damages []damage.Damage
		damage.Unmarshal(&damages, action.GetDamageInfo())
		if err := s.dmMgr.Add(damages); err != nil {
			return fmt.Errorf("failed to add damages: %w", err)
		}
		logger.Debug("Added damages: %+v", damages)
	case pb.Action_NEWEFFECT:
		var eff effect.Effect
		effect.Unmarshal(&eff, action.GetEffect())
		s.fieldLock.Lock()
		for i := 0; i < len(s.clients); i++ {
			if s.clients[i].clientID != eff.ClientID {
				eff.X = config.FieldNumX - eff.X - 1
			}
			s.clients[i].fieldInfo.Effects = append(s.clients[i].fieldInfo.Effects, eff)
		}
		s.fieldLock.Unlock()
		logger.Debug("Added effect: %+v", eff)
	case pb.Action_ADDSOUND:
		s.fieldLock.Lock()
		for i := 0; i < len(s.clients); i++ {
			s.clients[i].fieldInfo.Sounds = append(s.clients[i].fieldInfo.Sounds, action.GetSeType())
		}
		s.fieldLock.Unlock()
	default:
		return fmt.Errorf("action %d is not implemented yet", action.Type)
	}

	return nil
}

func (s *session) mainProcess() {
	// init info
	s.clients[0].fieldInfo.InitPanel(s.clients[0].clientID, s.clients[1].clientID)
	s.clients[1].fieldInfo.InitPanel(s.clients[1].clientID, s.clients[0].clientID)

	cancel := make(chan struct{})
	go s.frameProc(cancel)
	go s.dataSend(cancel)

	err := <-s.exitErr
	if errors.Is(err.reason, errSendFailed) {
		s.publishGameEnd(err.generatorClientID)
	}
	close(cancel)
	sessionLock.Lock()
	delete(sessionList, s.sessionID)
	sessionLock.Unlock()

	if err.reason != nil && !errors.Is(err.reason, errSendFailed) {
		logger.Error("Run failed: %v", err)
	}
}

func (s *session) frameProc(cancel chan struct{}) {
	fpsMgr := fps.Fps{TargetFPS: 60}
	for {
		select {
		case <-cancel:
			return
		default:
			if s.status == statusActing {
				// damage process
				for i, c := range s.clients {
					for _, obj := range c.fieldInfo.Objects {
						if !obj.Hittable {
							continue
						}

						if dm := s.dmMgr.Hit(c.clientID, obj.ClientID, obj.X, obj.Y); dm != nil {
							s.fieldLock.Lock()
							s.clients[i].fieldInfo.HitDamages = append(s.clients[i].fieldInfo.HitDamages, *dm)
							s.fieldLock.Unlock()
							logger.Debug("Hit damage for %s: %+v", c.clientID, dm)
						}
					}
				}
				s.dmMgr.Update()
			}

			fpsMgr.Wait()
		}
	}
}

func (s *session) dataSend(cancel chan struct{}) {
	for {
		select {
		case <-cancel:
			return
		default:
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
				logger.Info("Update status send failed for client %s: %v", s.clients[0].clientID, err)
				s.exitErr <- sessionError{
					generatorClientID: s.clients[0].clientID,
					reason:            errSendFailed,
				}
				return
			}
			if err := s.clients[1].dataStream.Send(d); err != nil {
				logger.Info("Update status send failed for client %s: %v", s.clients[1].clientID, err)
				s.exitErr <- sessionError{
					generatorClientID: s.clients[1].clientID,
					reason:            errSendFailed,
				}
				return
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
				logger.Info("Update status send failed for client %s: %v", s.clients[0].clientID, err)
				s.exitErr <- sessionError{
					generatorClientID: s.clients[0].clientID,
					reason:            errSendFailed,
				}
				return
			}
			if err := s.clients[1].dataStream.Send(d); err != nil {
				logger.Info("Update status send failed for client %s: %v", s.clients[1].clientID, err)
				s.exitErr <- sessionError{
					generatorClientID: s.clients[1].clientID,
					reason:            errSendFailed,
				}
				return
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
			case statusGameEnd:
				sendSt = pb.Data_GAMEEND
			default:
				s.exitErr <- sessionError{
					reason: fmt.Errorf("invalid next status: %d", s.nextStatus),
				}
				return
			}

			d := &pb.Data{
				Type: pb.Data_UPDATESTATUS,
				Data: &pb.Data_Status_{
					Status: sendSt,
				},
			}
			for i := 0; i < len(s.clients); i++ {
				if err := s.clients[i].dataStream.Send(d); err != nil {
					logger.Info("Update status send failed for client %s: %v", s.clients[i].clientID, err)
					s.exitErr <- sessionError{
						generatorClientID: s.clients[i].clientID,
						reason:            errSendFailed,
					}
					return
				}
			}

			s.changeStatus(s.nextStatus)
			s.nextStatus = -1
		}
	case statusGameEnd:
		s.exitErr <- sessionError{
			reason: nil,
		}
		logger.Info("Finished session %s by game end", s.sessionID)
	}
}

func (s *session) changeStatus(next int) {
	logger.Info("Change state from %d to %d", s.status, next)
	s.status = next
}

func (s *session) publishField() {
	now := time.Now()

	for i := 0; i < len(s.clients); i++ {
		s.fieldLock.Lock()
		s.clients[i].fieldInfo.CurrentTime = now

		for x := 0; x < config.FieldNumX; x++ {
			for y := 0; y < config.FieldNumY; y++ {
				s.clients[i].fieldInfo.Panels[x][y].ShowHitArea = false
			}
		}
		for _, pos := range s.dmMgr.GetHitAreas(s.clients[i].clientID) {
			s.clients[i].fieldInfo.Panels[pos[0]][pos[1]].ShowHitArea = true
		}

		d := &pb.Data{
			Type: pb.Data_DATA,
			Data: &pb.Data_RawData{
				RawData: field.Marshal(s.clients[i].fieldInfo),
			},
		}
		s.fieldLock.Unlock()

		if err := s.clients[i].dataStream.Send(d); err != nil {
			logger.Info("Field info send failed for client %s: %v", s.clients[i].clientID, err)
			s.exitErr <- sessionError{
				generatorClientID: s.clients[i].clientID,
				reason:            errSendFailed,
			}
		}

		s.fieldLock.Lock()
		s.clients[i].fieldInfo.Effects = []effect.Effect{}
		s.clients[i].fieldInfo.HitDamages = []damage.Damage{}
		s.clients[i].fieldInfo.Sounds = []int32{}
		s.fieldLock.Unlock()
	}
}

func (s *session) publishGameEnd(sendClientID string) {
	d := &pb.Data{
		Type: pb.Data_UPDATESTATUS,
		Data: &pb.Data_Status_{
			Status: pb.Data_GAMEEND,
		},
	}

	for i := 0; i < len(s.clients); i++ {
		if s.clients[i].clientID == sendClientID {
			continue
		}
		// do not require error handling, because this is final send in session
		s.clients[i].dataStream.Send(d)
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
		var obj object.Object
		object.Unmarshal(&obj, action.GetObjectInfo())
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
