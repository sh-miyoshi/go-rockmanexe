package gamehandler

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/newnet/config"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/newnet/netconnpb"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/newnet/object"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/router/gameinfo"
	gameobj "github.com/sh-miyoshi/go-rockmanexe/pkg/router/object"
)

type gameObject struct {
	animObject        objanim.Anim
	actionQueue       *pb.Request_Action
	currentObjectType *int
}

type animInfo struct {
	ownerClientID string
	startCount    int
}

type GameHandler struct {
	info      [2]gameinfo.GameInfo
	objects   map[string]*gameObject // Key: clientID, Value: object情報
	anims     map[string]animInfo    // Key: objectID, Value: anim情報
	gameCount int
}

func NewHandler() *GameHandler {
	return &GameHandler{
		objects:   make(map[string]*gameObject),
		anims:     make(map[string]animInfo),
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
	return nil
}

func (g *GameHandler) AddPlayerObject(clientID string, param object.InitParam) {
	var ginfo *gameinfo.GameInfo
	for i := 0; i < len(g.info); i++ {
		if g.info[i].ClientID == clientID {
			ginfo = &g.info[i]
		}
	}

	// Player Objectを作成
	g.objects[clientID] = &gameObject{}
	plyr := gameobj.NewPlayer(gameinfo.Object{
		ID:            param.ID,
		Type:          gameobj.TypePlayerStand,
		OwnerClientID: clientID,
		HP:            param.HP,
		Pos:           common.Point{X: param.X, Y: param.Y},
		IsReverse:     false,
	}, ginfo, g.objects[clientID].actionQueue)
	g.objects[clientID].animObject = plyr
	id := objanim.New(g.objects[clientID].animObject)
	g.anims[id] = animInfo{
		ownerClientID: clientID,
		startCount:    g.gameCount,
	}
	g.objects[clientID].currentObjectType = plyr.GetCurrentObjectTypePointer()

	g.updateGameInfo()
}

func (g *GameHandler) HandleAction(clientID string, act *pb.Request_Action) error {
	g.objects[clientID].actionQueue = act
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
	if err := anim.MgrProcess(); err != nil {
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
			clientID := g.anims[obj.ObjID].ownerClientID
			if clientID == g.info[i].ClientID {
				// 自分のObject
				objects[i] = append(objects[i], gameinfo.Object{
					ID:            obj.ObjID,
					Type:          *g.objects[clientID].currentObjectType,
					OwnerClientID: clientID,
					HP:            obj.HP,
					Pos:           obj.Pos,
					ActCount:      g.gameCount - g.anims[obj.ObjID].startCount,
					IsReverse:     false,
				})
			} else {
				// 相手のObjectなのでReverseする
				objects[i] = append(objects[i], gameinfo.Object{
					ID:            obj.ObjID,
					Type:          *g.objects[clientID].currentObjectType,
					OwnerClientID: clientID,
					HP:            obj.HP,
					Pos:           common.Point{X: battlecommon.FieldNum.X - obj.Pos.X - 1, Y: obj.Pos.Y},
					ActCount:      g.gameCount - g.anims[obj.ObjID].startCount,
					IsReverse:     true,
				})
			}
		}
	}

	anims := [len(g.info)][]gameinfo.Anim{}
	for _, a := range anim.GetAll() {
		for i := 0; i < len(g.info); i++ {
			pos := a.Pos
			if g.anims[a.ObjID].ownerClientID == g.info[i].ClientID {
				pos.X = battlecommon.FieldNum.X - a.Pos.X - 1
			}

			anims[i] = append(anims[i], gameinfo.Anim{
				ObjectID: a.ObjID,
				Pos:      pos,
				AnimType: a.AnimType,
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
