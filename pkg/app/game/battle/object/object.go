package object

import "github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"

type ObjectParam struct {
	Pos           point.Point
	HP            int
	OnwerCharType int
	AttackNum     int
	Interval      int
	Power         int

	objectID string
	xFlip    bool
}
