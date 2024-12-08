package gamehandler

import (
	"fmt"

	"github.com/cockroachdb/errors"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/system"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/net/netconnpb"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/object"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/session"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/sysinfo"
	routeranim "github.com/sh-miyoshi/go-rockmanexe/pkg/router/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/router/gameinfo"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/router/manager"
	gameobj "github.com/sh-miyoshi/go-rockmanexe/pkg/router/object"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	clientNum = 2
)

type gameObject struct {
	info         gameinfo.GameInfo
	playerObject *gameobj.Player
}

type GameHandler struct {
	manager   *manager.Manager
	objects   [clientNum]*gameObject // Key: clientID, Value: object情報
	gameCount int
}

func NewHandler() session.GameLogic {
	return &GameHandler{
		gameCount: 0,
	}
}

func (g *GameHandler) Init(clientIDs [clientNum]string, sysReceiver chan sysinfo.SysInfo) error {
	g.objects[0] = newGameObject(clientIDs[0], clientIDs[1])
	g.objects[1] = newGameObject(clientIDs[1], clientIDs[0])
	g.manager = manager.New(sysReceiver)

	logger.Info("Successfully initalized game handler by clients %+v", clientIDs)
	return nil
}

func (g *GameHandler) Cleanup() {
	for i := 0; i < clientNum; i++ {
		if g.objects[i].playerObject != nil {
			g.objects[i].playerObject.End()
		}
	}

	if g.manager != nil {
		g.manager.Cleanup()
	}
}

func (g *GameHandler) AddPlayerObject(clientID string, param object.InitParam) error {
	index := g.indexForClient(clientID)
	if index == -1 {
		logger.Error("cannot find game object for client %s", clientID)
		return errors.New("failed to find game object, it maybe called this point before Init())")
	}

	g.objects[index].playerObject = gameobj.NewPlayer(gameinfo.Object{
		ID:            param.ID,
		Type:          gameobj.TypePlayerStand,
		OwnerClientID: clientID,
		HP:            param.HP,
		Pos:           point.Point{X: param.X, Y: param.Y},
		IsReverse:     false,
	}, g.manager, gameinfo.FieldFuncs{
		GetPanelInfo:      g.objects[index].info.GetPanelInfo,
		ChangePanelStatus: g.changePanelStatus,
		ChangePanelType:   g.changePanelType,
	})
	g.manager.ObjAnimNew(g.objects[index].playerObject)

	g.updateGameInfo()
	logger.Info("Successfully add client %s with %+v", clientID, param)
	return nil
}

func (g *GameHandler) HandleAction(clientID string, act *pb.Request_Action) error {
	index := g.indexForClient(clientID)
	logger.Info("Got action %d from %s", act.GetType(), clientID)
	g.objects[index].playerObject.HandleAction(act)
	return nil
}

func (g *GameHandler) GetInfo(clientID string) []byte {
	index := g.indexForClient(clientID)
	return g.objects[index].info.Marshal()
}

func (g *GameHandler) UpdateGameStatus() {
	if err := g.manager.Update(); err != nil {
		system.SetError(fmt.Sprintf("Failed to manage animation: %+v", err))
		return
	}

	g.updateGameInfo()
}

func (g *GameHandler) IsGameEnd() bool {
	for _, obj := range g.objects {
		// Note: 1対1以外の場合は追加の考慮が必要
		if obj.playerObject.GetParam().HP <= 0 {
			return true
		}
	}

	return false
}

func (g *GameHandler) updateGameInfo() {
	objects := [clientNum][]gameinfo.Object{}
	for _, obj := range g.manager.ObjAnimGetObjs(objanim.FilterAll) {
		for i := 0; i < clientNum; i++ {
			var info gameobj.NetInfo
			info.Unmarshal(obj.ExtraInfo)

			if info.OwnerClientID == g.objects[i].info.ClientID {
				// 自分のObject
				objects[i] = append(objects[i], gameinfo.Object{
					ID:            obj.ObjID,
					Type:          g.objects[i].playerObject.GetAnimObjectType(),
					OwnerClientID: info.OwnerClientID,
					HP:            obj.HP,
					Pos:           obj.Pos,
					ActCount:      info.ActCount,
					IsReverse:     false,
					IsInvincible:  info.IsInvincible,
				})
			} else if index := g.indexForClient(info.OwnerClientID); index >= 0 && g.objects[index].playerObject != nil {
				// 相手のObjectなのでReverseする
				objects[i] = append(objects[i], gameinfo.Object{
					ID:            obj.ObjID,
					Type:          g.objects[index].playerObject.GetAnimObjectType(),
					OwnerClientID: info.OwnerClientID,
					HP:            obj.HP,
					Pos:           point.Point{X: battlecommon.FieldNum.X - obj.Pos.X - 1, Y: obj.Pos.Y},
					ActCount:      info.ActCount,
					IsReverse:     true,
					IsInvincible:  info.IsInvincible,
				})
			}
		}
	}

	anims := [clientNum][]gameinfo.Anim{}
	for _, a := range g.manager.AnimGetAll() {
		for i := 0; i < clientNum; i++ {
			var info routeranim.NetInfo
			info.Unmarshal(a.ExtraInfo)

			pos := a.Pos
			if info.OwnerClientID != g.objects[i].info.ClientID {
				pos.X = battlecommon.FieldNum.X - a.Pos.X - 1
			}

			anims[i] = append(anims[i], gameinfo.Anim{
				ObjectID:      a.ObjID,
				OwnerClientID: info.OwnerClientID,
				Pos:           pos,
				DrawType:      a.DrawType,
				AnimType:      info.AnimType,
				ActCount:      info.ActCount,
				DrawParam:     info.DrawParam[:],
			})
		}
	}

	effects := []gameinfo.Effect{}
	for _, e := range g.manager.QueuePopAll(gameinfo.QueueTypeEffect) {
		effects = append(effects, *e.(*gameinfo.Effect))
	}

	sounds := []gameinfo.Sound{}
	for _, s := range g.manager.QueuePopAll(gameinfo.QueueTypeSound) {
		sounds = append(sounds, *s.(*gameinfo.Sound))
	}

	// TODO: lock
	for i := 0; i < clientNum; i++ {
		g.objects[i].info.Update(objects[i], anims[i], effects, sounds)
	}
	g.gameCount++
}

func (g *GameHandler) indexForClient(clientID string) int {
	if clientID == "" {
		logger.Debug("index target clientID is nil")
		return -1
	}

	for i := 0; i < clientNum; i++ {
		if g.objects[i].info.ClientID == clientID {
			return i
		}
	}
	return -1
}

func (g *GameHandler) changePanelStatus(clientID string, pos point.Point, crackType int, endCount int) {
	// WIP: endCount

	index := g.indexForClient(clientID)
	for i := 0; i < clientNum; i++ {
		if i == index {
			g.objects[i].info.PanelChange(pos, crackType)
		} else {
			// 敵によるPanelBreakの場合場所を反転させる
			bpos := point.Point{X: battlecommon.FieldNum.X - pos.X - 1, Y: pos.Y}
			g.objects[i].info.PanelChange(bpos, crackType)
		}
	}
}

func (g *GameHandler) changePanelType(clientID string, pos point.Point, pnType int, endCount int) {
	// WIP: endCount

	index := g.indexForClient(clientID)
	for i := 0; i < clientNum; i++ {
		if i == index {
			g.objects[i].info.ChangePanelType(pos, clientID)
		} else {
			// 敵によるPanelBreakの場合場所を反転させる
			bpos := point.Point{X: battlecommon.FieldNum.X - pos.X - 1, Y: pos.Y}
			g.objects[i].info.ChangePanelType(bpos, clientID)
		}
	}
}

func newGameObject(myClientID string, opponentClientID string) *gameObject {
	res := &gameObject{}
	res.info.Init(myClientID, opponentClientID)
	return res
}
