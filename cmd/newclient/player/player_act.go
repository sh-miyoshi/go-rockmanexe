package player

type act interface {
	Process() bool
	Interval() int
}

type actMove struct {
	// targetX int
	// targetY int
}

func newActMove() *actMove {
	return &actMove{}
}

func (a *actMove) Process() bool {
	return true
}

func (a *actMove) Interval() int {
	return 30
}
