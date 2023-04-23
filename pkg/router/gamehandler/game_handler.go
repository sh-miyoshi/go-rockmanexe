package gamehandler

import (
	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/config"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/net/netconnpb"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/object"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/session"
	routeranim "github.com/sh-miyoshi/go-rockmanexe/pkg/router/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/router/gameinfo"
	gameobj "github.com/sh-miyoshi/go-rockmanexe/pkg/router/object"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/router/queue"
)

type gameObject struct {
	animObject        objanim.Anim
	actionQueueID     string
	currentObjectType *int
}

type GameHandler struct {
	info      [2]gameinfo.GameInfo
	objects   map[string]*gameObject // Key: clientID, Value: object情報
	gameCount int
	animMgrID string
}

func NewHandler() session.GameLogic {
	return &GameHandler{
		objects:   make(map[string]*gameObject),
		gameCount: 0,
	}
}

func (g *GameHandler) Init(clientIDs [2]string) error {
	for y := 0; y < config.FieldNumY; y++ {
		hx := config.FieldNumX / 2
		for x := 0; x < hx; x++ {
			g.info[0].ClientID = clientIDs[0]
			g.info[0].Panels[x][y] = gameinfo.PanelInfo{OwnerClientID: clientIDs[0]}
			g.info[0].Panels[x+hx][y] = gameinfo.PanelInfo{OwnerClientID: clientIDs[1]}
			g.info[1].ClientID = clientIDs[1]
			g.info[1].Panels[x][y] = gameinfo.PanelInfo{OwnerClientID: clientIDs[1]}
			g.info[1].Panels[x+hx][y] = gameinfo.PanelInfo{OwnerClientID: clientIDs[0]}
		}
	}
	g.animMgrID = routeranim.NewManager(clientIDs)

	logger.Info("Successfully initalized game handler by clients %+v", clientIDs)
	return nil
}

func (g *GameHandler) Cleanup() {
	for clientID, obj := range g.objects {
		damage.RemoveForClient(clientID)
		queue.Delete(obj.actionQueueID)
	}
	// TODO: remove data from anim, objanim
}

func (g *GameHandler) AddPlayerObject(clientID string, param object.InitParam) {
	var ginfo *gameinfo.GameInfo
	for i := 0; i < len(g.info); i++ {
		if g.info[i].ClientID == clientID {
			ginfo = &g.info[i]
		}
	}
	if ginfo == nil {
		logger.Error("cannot find game info for client %s", clientID)
		// TODO: return error to client
		return
	}

	// Player Objectを作成
	g.objects[clientID] = &gameObject{
		actionQueueID: uuid.New().String(),
	}
	plyr := gameobj.NewPlayer(gameinfo.Object{
		ID:            param.ID,
		Type:          gameobj.TypePlayerStand,
		OwnerClientID: clientID,
		HP:            param.HP,
		Pos:           common.Point{X: param.X, Y: param.Y},
		IsReverse:     false,
	}, ginfo, g.objects[clientID].actionQueueID)
	g.objects[clientID].animObject = plyr
	objanim.New(g.objects[clientID].animObject)
	g.objects[clientID].currentObjectType = plyr.GetCurrentObjectTypePointer()

	g.updateGameInfo()
	logger.Info("Successfully add client %s with %+v", clientID, param)
}

func (g *GameHandler) HandleAction(clientID string, act *pb.Request_Action) error {
	logger.Info("Got action %d from %s", act.GetType(), clientID)
	queue.Push(g.objects[clientID].actionQueueID, act)
	return nil
}

func (g *GameHandler) GetInfo(clientID string) []byte {
	for i := 0; i < len(g.info); i++ {
		if g.info[i].ClientID == clientID {
			return g.info[i].Marshal()
		}
	}
	return nil
}

func (g *GameHandler) UpdateGameStatus() {
	if err := routeranim.MgrProcess(g.animMgrID); err != nil {
		logger.Error("Failed to manage animation: %+v", err)
		// TODO: 処理を終了する
	}
	if err := objanim.MgrProcess(true, false); err != nil {
		logger.Error("Failed to manage object animation: %+v", err)
		// TODO: 処理を終了する
	}
	damage.MgrProcess()

	g.updateGameInfo()
}

func (g *GameHandler) IsGameEnd() bool {
	for _, obj := range g.objects {
		// Note: 1対1以外の場合は追加の考慮が必要
		if obj.animObject.GetParam().HP <= 0 {
			return true
		}
	}

	return false
}

func (g *GameHandler) updatePanelObject() {
	for i := 0; i < len(g.info); i++ {
		// Cleanup at first
		for y := 0; y < battlecommon.FieldNum.Y; y++ {
			for x := 0; x < battlecommon.FieldNum.X; x++ {
				g.info[i].Panels[x][y].ObjectID = ""
			}
		}
		for _, obj := range g.info[i].Objects {
			g.info[i].Panels[obj.Pos.X][obj.Pos.Y].ObjectID = obj.ID
		}
	}
}

func (g *GameHandler) updateGameInfo() {
	objects := [len(g.info)][]gameinfo.Object{}
	for _, obj := range objanim.GetObjs(objanim.FilterAll) {
		for i := 0; i < len(g.info); i++ {
			var info gameobj.NetInfo
			info.Unmarshal(obj.ExtraInfo)

			if info.OwnerClientID == g.info[i].ClientID {
				logger.Debug("g.objects: %+v", g.objects)
				// 自分のObject
				objects[i] = append(objects[i], gameinfo.Object{
					ID:            obj.ObjID,
					Type:          *g.objects[info.OwnerClientID].currentObjectType,
					OwnerClientID: info.OwnerClientID,
					HP:            obj.HP,
					Pos:           obj.Pos,
					ActCount:      info.ActCount,
					IsReverse:     false,
				})
			} else if _, ok := g.objects[info.OwnerClientID]; ok {
				// 相手のObjectなのでReverseする
				objects[i] = append(objects[i], gameinfo.Object{
					ID:            obj.ObjID,
					Type:          *g.objects[info.OwnerClientID].currentObjectType,
					OwnerClientID: info.OwnerClientID,
					HP:            obj.HP,
					Pos:           common.Point{X: battlecommon.FieldNum.X - obj.Pos.X - 1, Y: obj.Pos.Y},
					ActCount:      info.ActCount,
					IsReverse:     true,
				})
			}
		}
	}

	anims := [len(g.info)][]gameinfo.Anim{}
	for _, a := range routeranim.GetAll(g.animMgrID) {
		for i := 0; i < len(g.info); i++ {
			var info routeranim.NetInfo
			info.Unmarshal(a.ExtraInfo)

			pos := a.Pos
			if info.OwnerClientID != g.info[i].ClientID {
				pos.X = battlecommon.FieldNum.X - a.Pos.X - 1
			}

			anims[i] = append(anims[i], gameinfo.Anim{
				ObjectID: a.ObjID,
				Pos:      pos,
				DrawType: a.DrawType,
				AnimType: info.AnimType,
				ActCount: info.ActCount,
			})
		}
	}

	// TODO: lock
	for i := 0; i < len(g.info); i++ {
		g.info[i].Objects = objects[i]
		g.info[i].Anims = anims[i]
	}
	g.updatePanelObject()
	g.gameCount++
}
