package gamehandler

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/chip"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/newnet/action"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/newnet/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/newnet/object"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/router/gameinfo"
	gameobj "github.com/sh-miyoshi/go-rockmanexe/pkg/router/object"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/router/skill"
)

type objectInfo struct {
	OwnerClientID string
	Type          int
	StartCount    int
}

type GameHandler struct {
	info          gameinfo.GameInfo
	objInfo       map[string]*objectInfo
	playerObjects map[string]*gameobj.Player
	gameCount     int
}

func NewHandler() *GameHandler {
	return &GameHandler{
		objInfo:       make(map[string]*objectInfo),
		playerObjects: make(map[string]*gameobj.Player),
		gameCount:     0,
	}
}

func (g *GameHandler) Init(clientIDs [2]string) error {
	for y := 0; y < config.FieldNumY; y++ {
		hx := config.FieldNumX / 2
		for x := 0; x < hx; x++ {
			g.info.Panels[x][y] = gameinfo.PanelInfo{OwnerClientID: clientIDs[0]}
			g.info.Panels[x+hx][y] = gameinfo.PanelInfo{OwnerClientID: clientIDs[1]}
		}
	}
	g.info.ReverseClientID = clientIDs[1]
	return nil
}

func (g *GameHandler) AddPlayerObject(clientID string, param object.InitParam) {
	x := param.X
	if g.info.ReverseClientID == clientID {
		x = battlecommon.FieldNum.X - x - 1
	}

	g.objInfo[param.ID] = &objectInfo{
		OwnerClientID: clientID,
		Type:          gameobj.TypePlayerStand,
		StartCount:    g.gameCount,
	}
	g.playerObjects[clientID] = gameobj.NewPlayer(gameinfo.Object{
		ID:            param.ID,
		Type:          g.objInfo[param.ID].Type,
		OwnerClientID: clientID,
		HP:            param.HP,
		Pos:           common.Point{X: x, Y: param.Y},
	}, &g.info)
	objanim.New(g.playerObjects[clientID])
	g.updateGameInfo()
}

func (g *GameHandler) MoveObject(clientID string, moveInfo action.Move) {
	g.playerObjects[clientID].AddMove(moveInfo)
}

func (g *GameHandler) AddBuster(clientID string, busterInfo action.Buster) {
	g.playerObjects[clientID].AddBuster(busterInfo)
}

func (g *GameHandler) UseChip(clientID string, chipInfo action.UseChip) {
	c := chip.Get(chipInfo.ChipID)
	logger.Debug("Use chip: %+v", c)

	var targetType int
	if g.info.ReverseClientID == clientID {
		if c.ForMe {
			targetType = damage.TargetEnemy
		} else {
			targetType = damage.TargetPlayer
		}
	} else {
		if c.ForMe {
			targetType = damage.TargetPlayer
		} else {
			targetType = damage.TargetEnemy
		}
	}

	s := skill.GetByChip(chipInfo.ChipID, skill.Argument{
		AnimObjID:  chipInfo.AnimID,
		OwnerID:    chipInfo.ChipUserClientID,
		Power:      c.Power,
		TargetType: targetType,

		GameInfo: &g.info,
	})
	anim.New(s)

	// TODO: player_act
}

func (g *GameHandler) GetInfo() []byte {
	return g.info.Marshal()
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
	// Cleanup at first
	for y := 0; y < battlecommon.FieldNum.Y; y++ {
		for x := 0; x < battlecommon.FieldNum.X; x++ {
			g.info.Panels[x][y].ObjectID = ""
		}
	}
	for _, obj := range g.info.Objects {
		g.info.Panels[obj.Pos.X][obj.Pos.Y].ObjectID = obj.ID
	}
}

func (g *GameHandler) updateGameInfo() {
	// Cleanup at first
	g.info.Objects = []gameinfo.Object{}
	g.info.Anims = []gameinfo.Anim{}

	for _, obj := range objanim.GetObjs(objanim.FilterAll) {
		g.info.Objects = append(g.info.Objects, gameinfo.Object{
			ID:            obj.ObjID,
			Type:          g.objInfo[obj.ObjID].Type,
			OwnerClientID: g.objInfo[obj.ObjID].OwnerClientID,
			HP:            obj.HP,
			Pos:           obj.Pos,
			ActCount:      g.gameCount - g.objInfo[obj.ObjID].StartCount,
		})
	}
	for _, a := range anim.GetAll() {
		g.info.Anims = append(g.info.Anims, gameinfo.Anim{
			ObjectID: a.ObjID,
			Pos:      a.Pos,
			AnimType: a.AnimType,
		})
	}

	g.updatePanelObject()
	g.gameCount++
}
