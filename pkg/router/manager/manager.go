package manager

import (
	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	effectanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/effect"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	skillanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/skill"
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
	skillAnimMgr  *skillanim.AnimManager
	effectAnimMgr *effectanim.AnimManager
	objAnimMgr    *objanim.AnimManager
	skillMgr      *skillmanager.Manager
	queueIDs      [gameinfo.QueueTypeMax]string
	sysReceiver   chan sysinfo.SysInfo
	cutinCount    int
}

func New(sysReceiver chan sysinfo.SysInfo) *Manager {
	res := &Manager{
		skillAnimMgr:  skillanim.NewManager(),
		effectAnimMgr: effectanim.NewManager(),
		objAnimMgr:    objanim.NewManager(),
		sysReceiver:   sysReceiver,
	}
	res.skillMgr = skillmanager.NewManager()
	for i := 0; i < gameinfo.QueueTypeMax; i++ {
		res.queueIDs[i] = uuid.NewString()
	}

	return res
}

func (m *Manager) Cleanup() {
	if m.objAnimMgr != nil {
		m.objAnimMgr.Cleanup()
	}
	if m.skillAnimMgr != nil {
		m.skillAnimMgr.Cleanup()
	}
	if m.effectAnimMgr != nil {
		m.effectAnimMgr.Cleanup()
	}
	for _, id := range m.queueIDs {
		queue.Delete(id)
	}
}

func (m *Manager) Update() error {
	if err := m.objAnimMgr.Update(true); err != nil {
		return errors.Wrap(err, "objanim manage process failed")
	}
	if err := m.skillAnimMgr.Update(true); err != nil {
		return errors.Wrap(err, "skillanim manage process failed")
	}
	if err := m.effectAnimMgr.Update(); err != nil {
		return errors.Wrap(err, "effectanim manage process failed")
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

func (m *Manager) SkillAnimNew(a skillanim.Anim) string {
	return m.skillAnimMgr.New(a)
}

func (m *Manager) EffectAnimNew(a effectanim.Anim) string {
	return m.effectAnimMgr.New(a)
}

func (m *Manager) AnimDelete(animID string) {
	if m.skillAnimMgr.IsProcessing(animID) {
		m.skillAnimMgr.Delete(animID)
	}
	if m.effectAnimMgr.IsProcessing(animID) {
		m.effectAnimMgr.Delete(animID)
	}
	if m.objAnimMgr.IsProcessing(animID) {
		m.objAnimMgr.Delete(animID)
	}
}

func (m *Manager) AnimGetAll() []anim.Param {
	return append(m.skillAnimMgr.GetAll(), m.effectAnimMgr.GetAll()...)
}

func (m *Manager) AnimIsProcessing(animID string) bool {
	if m.skillAnimMgr.IsProcessing(animID) {
		return true
	}
	if m.effectAnimMgr.IsProcessing(animID) {
		return true
	}
	if m.objAnimMgr.IsProcessing(animID) {
		return true
	}
	return false
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

func (m *Manager) ObjAnimAddBarrier(id string, hp int) {
	m.objAnimMgr.AddBarrier(id, hp)
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
