package player

import (
	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/netconn"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/action"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/net/netconnpb"
)

type Act interface {
	Update() bool
	Interval() int
}

/*

type ActTemplate struct {
}

func NewActTemplate() *ActTemplate {
	return &ActTemplate{}
}

func (a *ActTemplate) Update() bool {
	return false
}

func (a *ActTemplate) Interval() int {
	return 60
}

*/

type ActWait struct {
	waitFrame int
}

func NewActWait(waitFrame int) *ActWait {
	return &ActWait{
		waitFrame: waitFrame,
	}
}

func (a *ActWait) Update() bool {
	return true
}

func (a *ActWait) Interval() int {
	return a.waitFrame
}

type ActMove struct {
	targetX int
	targetY int
	conn    *netconn.NetConn
}

func NewActMove(targetX, targetY int, conn *netconn.NetConn) *ActMove {
	return &ActMove{
		targetX: targetX,
		targetY: targetY,
		conn:    conn,
	}
}

func (a *ActMove) Update() bool {
	logger.Debug("Move to (%d, %d)", a.targetX, a.targetY)
	move := action.Move{
		Type:    action.MoveTypeAbs,
		AbsPosX: a.targetX,
		AbsPosY: a.targetY,
	}
	a.conn.SendAction(pb.Request_MOVE, move.Marshal())
	return true
}

func (a *ActMove) Interval() int {
	return 30
}

type ActBuster struct {
	count     int
	conn      *netconn.NetConn
	shotPower uint
	charged   bool
}

func NewActBuster(conn *netconn.NetConn) *ActBuster {
	return &ActBuster{
		count: 0,
		conn:  conn,

		// debug
		shotPower: 1,
		charged:   false,
	}
}

func (a *ActBuster) Update() bool {
	buster := action.Buster{
		Power: 1,
	}
	a.conn.SendAction(pb.Request_BUSTER, buster.Marshal())
	return true
}

func (a *ActBuster) Interval() int {
	return 30
}

type ActSkill struct {
	clientID string
	chipID   int
	conn     *netconn.NetConn
	count    int
	id       string
}

func NewActSkill(chipID int, clientID string, conn *netconn.NetConn) *ActSkill {
	return &ActSkill{
		chipID:   chipID,
		clientID: clientID,
		conn:     conn,
		count:    0,
		id:       uuid.New().String(),
	}
}

func (a *ActSkill) Update() bool {
	if a.count == 0 {
		chipInfo := action.UseChip{
			AnimID:           a.id,
			ChipUserClientID: a.clientID,
			ChipID:           a.chipID,
		}
		a.conn.SendAction(pb.Request_CHIPUSE, chipInfo.Marshal())
	}
	a.count++

	info := a.conn.GetGameInfo()
	for _, anim := range info.Anims {
		if anim.ObjectID == a.id {
			return false // まだ処理中
		}
	}

	return true
}

func (a *ActSkill) Interval() int {
	return 60
}
