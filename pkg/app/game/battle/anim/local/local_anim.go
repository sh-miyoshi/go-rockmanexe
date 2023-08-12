package localanim

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
)

var (
	animInst    *anim.AnimManager
	objanimInst *objanim.AnimManager
)

func AnimMgrProcess() error {
	if animInst == nil {
		animInst = anim.NewManager()
	}

	return animInst.Process()
}

func AnimMgrDraw() {
	if animInst == nil {
		animInst = anim.NewManager()
	}

	animInst.MgrDraw()
}

func AnimNew(a anim.Anim) string {
	if animInst == nil {
		animInst = anim.NewManager()
	}

	return animInst.New(a)
}

func AnimIsProcessing(animID string) bool {
	if animInst == nil {
		animInst = anim.NewManager()
	}

	return animInst.IsProcessing(animID)
}

func AnimCleanup() {
	if animInst == nil {
		animInst = anim.NewManager()
	}

	animInst.Cleanup()
}

func AnimDelete(animID string) {
	if animInst == nil {
		animInst = anim.NewManager()
	}

	animInst.Delete(animID)
}

func AnimGetAll() []anim.Param {
	if animInst == nil {
		animInst = anim.NewManager()
	}

	return animInst.GetAll()
}

func ObjAnimMgrProcess(enableDamage bool, blackout bool) error {
	if objanimInst == nil {
		objanimInst = objanim.NewManager()
	}

	return objanimInst.Process(enableDamage, blackout)
}

func ObjAnimMgrDraw() {
	if objanimInst == nil {
		objanimInst = objanim.NewManager()
	}

	objanimInst.Draw()
}

func ObjAnimNew(anim objanim.Anim) string {
	if objanimInst == nil {
		objanimInst = objanim.NewManager()
	}

	return objanimInst.New(anim)
}

func ObjAnimIsProcessing(animID string) bool {
	if objanimInst == nil {
		objanimInst = objanim.NewManager()
	}

	return objanimInst.IsProcessing(animID)
}

func ObjAnimCleanup() {
	if objanimInst == nil {
		objanimInst = objanim.NewManager()
	}

	objanimInst.Cleanup()
}

func ObjAnimDelete(animID string) {
	if objanimInst == nil {
		objanimInst = objanim.NewManager()
	}

	objanimInst.Delete(animID)
}

func ObjAnimGetObjPos(objID string) common.Point {
	if objanimInst == nil {
		objanimInst = objanim.NewManager()
	}

	return objanimInst.GetObjPos(objID)
}

func ObjAnimGetObjs(filter objanim.Filter) []objanim.Param {
	if objanimInst == nil {
		objanimInst = objanim.NewManager()
	}

	return objanimInst.GetObjs(filter)
}

func ObjAnimAddActiveAnim(id string) {
	if objanimInst == nil {
		objanimInst = objanim.NewManager()
	}

	objanimInst.AddActiveAnim(id)
}

func ObjAnimDeactivateAnim(id string) {
	if objanimInst == nil {
		objanimInst = objanim.NewManager()
	}

	objanimInst.DeactivateAnim(id)
}

func ObjAnimMakeInvisible(id string, count int) {
	if objanimInst == nil {
		objanimInst = objanim.NewManager()
	}

	objanimInst.MakeInvisible(id, count)
}

func ObjAnimExistsObject(pos common.Point) string {
	if objanimInst == nil {
		objanimInst = objanim.NewManager()
	}

	return objanimInst.ExistsObject(pos)
}

func DamageManager() *damage.DamageManager {
	if objanimInst == nil {
		objanimInst = objanim.NewManager()
	}

	return objanimInst.DamageManager()
}
