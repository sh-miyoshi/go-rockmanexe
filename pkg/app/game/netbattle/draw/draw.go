package draw

import "github.com/cockroachdb/errors"

var (
	objDrawInst  objectDraw
	animDrawInst animDraw
)

func Init(playerObjectID string) error {
	if err := objDrawInst.Init(playerObjectID); err != nil {
		return errors.Wrap(err, "object draw init failed")
	}

	if err := animDrawInst.Init(); err != nil {
		return errors.Wrap(err, "anim draw init failed")
	}

	return nil
}

func End() {
	objDrawInst.End()
	animDrawInst.End()
}

func Draw() {
	objDrawInst.Draw()
	animDrawInst.Draw()
}
