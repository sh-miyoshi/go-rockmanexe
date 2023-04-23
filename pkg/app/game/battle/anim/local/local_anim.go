package localanim

import "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"

var (
	animInst *anim.AnimManager
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
