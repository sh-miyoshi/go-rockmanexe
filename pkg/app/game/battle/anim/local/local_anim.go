package localanim

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	effectanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/effect"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	skillanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/skill"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	skillmanager "github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/manager"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

var (
	objanimInst    *objanim.AnimManager
	skillAnimInst  *skillanim.AnimManager
	effectAnimInst *effectanim.AnimManager
	skillMgr       *skillmanager.Manager
)

func newInst() {
	if objanimInst == nil {
		objanimInst = objanim.NewManager()
	}
	if skillAnimInst == nil {
		skillAnimInst = skillanim.NewManager()
	}
	if effectAnimInst == nil {
		effectAnimInst = effectanim.NewManager()
	}
}

func AnimMgrProcess() error {
	newInst()

	if err := skillAnimInst.Update(); err != nil {
		return err
	}

	return effectAnimInst.Update()
}

func AnimMgrDraw() {
	newInst()

	objanimInst.Draw()
	skillAnimInst.Draw()
	effectAnimInst.Draw()
}

func SkillAnimNew(a skillanim.Anim) string {
	newInst()

	return skillAnimInst.New(a)
}

func EffectAnimNew(a effectanim.Anim) string {
	newInst()

	return effectAnimInst.New(a)
}

func AnimIsProcessing(id string) bool {
	newInst()

	if objanimInst.IsProcessing(id) {
		return true
	}

	if skillAnimInst.IsProcessing(id) {
		return true
	}

	if effectAnimInst.IsProcessing(id) {
		return true
	}

	return false
}

func AnimCleanup() {
	newInst()

	skillAnimInst.Cleanup()
	effectAnimInst.Cleanup()
	objanimInst.Cleanup()
}

func AnimDelete(animID string) {
	newInst()

	if skillAnimInst.IsProcessing(animID) {
		skillAnimInst.Delete(animID)
	}
	if effectAnimInst.IsProcessing(animID) {
		effectAnimInst.Delete(animID)
	}
	if objanimInst.IsProcessing(animID) {
		objanimInst.Delete(animID)
	}
}

func AnimGetEffects() []anim.Param {
	newInst()

	return effectAnimInst.GetAll()
}

func ObjAnimMgrProcess(enableDamage bool, blackout bool) error {
	newInst()

	return objanimInst.Process(enableDamage, blackout)
}

func ObjAnimNew(anim objanim.Anim) string {
	newInst()

	return objanimInst.New(anim)
}

func ObjAnimGetObjPos(objID string) point.Point {
	newInst()

	return objanimInst.GetObjPos(objID)
}

func ObjAnimGetObjs(filter objanim.Filter) []objanim.Param {
	newInst()

	return objanimInst.GetObjs(filter)
}

func ObjAnimAddActiveAnim(id string) {
	newInst()

	objanimInst.AddActiveAnim(id)
}

func ObjAnimDeactivateAnim(id string) {
	newInst()

	objanimInst.DeactivateAnim(id)
}

func ObjAnimMakeInvisible(id string, count int) {
	newInst()

	objanimInst.MakeInvisible(id, count)
}

func ObjAnimAddBarrier(id string, hp int) {
	newInst()

	objanimInst.AddBarrier(id, hp)
}

func ObjAnimExistsObject(pos point.Point) string {
	newInst()

	return objanimInst.ExistsObject(pos)
}

func DamageManager() *damage.DamageManager {
	newInst()

	return objanimInst.DamageManager()
}

func SkillManager() *skillmanager.Manager {
	if skillMgr == nil {
		skillMgr = skillmanager.NewManager()
	}
	return skillMgr
}
