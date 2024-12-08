package manager

import (
	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	skillmanager "github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/manager"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/sysinfo"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/router/gameinfo"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/router/manager/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/queue"
)

type Manager struct {
	animMgr     *anim.AnimManager
	objAnimMgr  *objanim.AnimManager
	skillMgr    *skillmanager.Manager
	queueIDs    [gameinfo.QueueTypeMax]string
	sysReceiver chan sysinfo.SysInfo
	cutinCount  int
}

func New(sysReceiver chan sysinfo.SysInfo) *Manager {
	res := &Manager{
		animMgr:     anim.NewManager(),
		objAnimMgr:  objanim.NewManager(),
		sysReceiver: sysReceiver,
	}
	res.skillMgr = skillmanager.NewManager()
	for i := 0; i < gameinfo.QueueTypeMax; i++ {
		res.queueIDs[i] = uuid.NewString()
	}

	return res
}

func (m *Manager) Cleanup() {
	if m.animMgr != nil {
		m.animMgr.Cleanup()
	}
	if m.objAnimMgr != nil {
		m.objAnimMgr.Cleanup()
	}
	for _, id := range m.queueIDs {
		queue.Delete(id)
	}
}

func (m *Manager) Update() error {
	if err := m.animMgr.Update(); err != nil {
		return errors.Wrap(err, "anim manage process failed")
	}

	if err := m.objAnimMgr.Process(true, false); err != nil {
		return errors.Wrap(err, "objanim manage process failed")
	}

	if m.cutinCount > 0 {
		m.cutinCount--
		if m.cutinCount == 0 {
			logger.Info("cutin time end")
			m.sysReceiver <- sysinfo.SysInfo{Type: sysinfo.TypeActing}
		}
	}

	return nil
}

func (m *Manager) QueuePush(typ int, info interface{}) {
	queue.Push(m.queueIDs[typ], info)
}

func (m *Manager) QueuePopAll(typ int) []interface{} {
	return queue.PopAll(m.queueIDs[typ])
}

func (m *Manager) AnimNew(a anim.Anim) string {
	return m.animMgr.New(a)
}

func (m *Manager) AnimDelete(animID string) {
	m.animMgr.Delete(animID)
}

func (m *Manager) AnimGetAll() []anim.Param {
	return m.animMgr.GetAll()
}

func (m *Manager) AnimIsProcessing(animID string) bool {
	return m.animMgr.IsProcessing(animID)
}

func (m *Manager) ObjAnimNew(anim objanim.Anim) string {
	return m.objAnimMgr.New(anim)
}

func (m *Manager) ObjAnimGetObjPos(objID string) point.Point {
	return m.objAnimMgr.GetObjPos(objID)
}

func (m *Manager) ObjAnimGetObjs(filter objanim.Filter) []objanim.Param {
	return m.objAnimMgr.GetObjs(filter)
}

func (m *Manager) ObjAnimMakeInvisible(id string, count int) {
	m.objAnimMgr.MakeInvisible(id, count)
}

func (m *Manager) SkillGet(id int, arg skillcore.Argument) skillcore.SkillCore {
	return m.skillMgr.Get(id, arg)
}

func (m *Manager) DamageMgr() *damage.Manager {
	return damage.New(m.objAnimMgr.DamageManager())
}

func (m *Manager) SoundOn(typ resources.SEType) {
	m.QueuePush(gameinfo.QueueTypeSound, &gameinfo.Sound{
		ID:   uuid.New().String(),
		Type: int(typ),
	})
}

func (m *Manager) Cutin(skillName string, count int, clientID string) {
	logger.Info("cutin with %d count", count)
	cutin := sysinfo.Cutin{
		SkillName:     skillName,
		OwnerClientID: clientID,
	}
	m.sysReceiver <- sysinfo.SysInfo{
		Type: sysinfo.TypeCutin,
		Data: cutin.Marshal(),
	}
	m.cutinCount = count
}
