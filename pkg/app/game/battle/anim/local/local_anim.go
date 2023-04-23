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

func New(a anim.Anim) string {
	if animInst == nil {
		animInst = anim.NewManager()
	}

	return animInst.New(a)
}

func IsProcessing(animID string) bool {
	if animInst == nil {
		animInst = anim.NewManager()
	}

	return animInst.IsProcessing(animID)
}

func Cleanup() {
	if animInst == nil {
		animInst = anim.NewManager()
	}

	animInst.Cleanup()
}

func Delete(animID string) {
	if animInst == nil {
		animInst = anim.NewManager()
	}

	animInst.Delete(animID)
}

func GetAll() []anim.Param {
	if animInst == nil {
		animInst = anim.NewManager()
	}

	return animInst.GetAll()
}
