package manager

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	effectanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/effect"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	skillanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/skill"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	skillmanager "github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/manager"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type Manager struct {
	objanimInst    *objanim.AnimManager
	skillAnimInst  *skillanim.AnimManager
	effectAnimInst *effectanim.AnimManager
	skillMgr       *skillmanager.Manager
}

func NewManager() *Manager {
	return &Manager{
		objanimInst:    objanim.NewManager(),
		skillAnimInst:  skillanim.NewManager(),
		effectAnimInst: effectanim.NewManager(),
		skillMgr:       skillmanager.NewManager(),
	}
}

func (m *Manager) Update(isActive bool) error {
	if err := m.skillAnimInst.Update(); err != nil {
		return err
	}

	if err := m.effectAnimInst.Update(); err != nil {
		return err
	}

	// WIP: object update
	return m.objanimInst.Update(isActive)
}

func (m *Manager) Draw() {
	m.objanimInst.Draw()
	m.skillAnimInst.Draw()
	m.effectAnimInst.Draw()
}

func (m *Manager) SkillAnimNew(a skillanim.Anim) string {
	return m.skillAnimInst.New(a)
}

func (m *Manager) EffectAnimNew(a effectanim.Anim) string {
	return m.effectAnimInst.New(a)
}

func (m *Manager) ObjAnimNew(anim objanim.Anim) string {
	return m.objanimInst.New(anim)
}

func (m *Manager) IsAnimProcessing(id string) bool {
	if m.objanimInst.IsProcessing(id) {
		return true
	}

	if m.skillAnimInst.IsProcessing(id) {
		return true
	}

	if m.effectAnimInst.IsProcessing(id) {
		return true
	}

	return false
}

func (m *Manager) SetActiveAnim(id string) {
	if m.objanimInst.IsProcessing(id) {
		m.objanimInst.SetActiveAnim(id)
		return
	}

	if m.skillAnimInst.IsProcessing(id) {
		// WIP
		return
	}

	if m.effectAnimInst.IsProcessing(id) {
		// WIP
		return
	}
}

func (m *Manager) DeactivateAnim(id string) {
	if m.objanimInst.IsProcessing(id) {
		m.objanimInst.DeactivateAnim(id)
		return
	}

	if m.skillAnimInst.IsProcessing(id) {
		// WIP
		return
	}

	if m.effectAnimInst.IsProcessing(id) {
		// WIP
		return
	}
}

func (m *Manager) Cleanup() {
	m.skillAnimInst.Cleanup()
	m.effectAnimInst.Cleanup()
	m.objanimInst.Cleanup()
}

func (m *Manager) AnimDelete(animID string) {
	if m.skillAnimInst.IsProcessing(animID) {
		m.skillAnimInst.Delete(animID)
	}
	if m.effectAnimInst.IsProcessing(animID) {
		m.effectAnimInst.Delete(animID)
	}
	if m.objanimInst.IsProcessing(animID) {
		m.objanimInst.Delete(animID)
	}
}

func (m *Manager) AnimGetSkills() []anim.Param {
	return m.skillAnimInst.GetAll()
}

func (m *Manager) ObjAnimGetObjPos(objID string) point.Point {
	return m.objanimInst.GetObjPos(objID)
}

func (m *Manager) ObjAnimGetObjs(filter objanim.Filter) []objanim.Param {
	return m.objanimInst.GetObjs(filter)
}

func (m *Manager) ObjAnimMakeInvisible(id string, count int) {
	m.objanimInst.MakeInvisible(id, count)
}

func (m *Manager) ObjAnimAddBarrier(id string, hp int) {
	m.objanimInst.AddBarrier(id, hp)
}

func (m *Manager) ObjAnimExistsObject(pos point.Point) string {
	return m.objanimInst.ExistsObject(pos)
}

// WIP: managerを直接見せなくてもいいようにしたい
func (m *Manager) DamageManager() *damage.DamageManager {
	return m.objanimInst.DamageManager()
}

// WIP: managerを直接見せなくてもいいようにしたい
func (m *Manager) SkillManager() *skillmanager.Manager {
	return m.skillMgr
}
