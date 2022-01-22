package object

import "github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"

type ObjectParam struct {
	Pos           common.Point
	HP            int
	OnwerCharType int
	AttackNum     int
	Interval      int
	Power         int

	objectID string
	xFlip    bool
}
