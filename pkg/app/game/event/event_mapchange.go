package event

type MapChangeHandler struct{}

func (h *MapChangeHandler) Init(args string) error {
	return nil
}

func (h *MapChangeHandler) Draw() {

}

func (h *MapChangeHandler) Process() error {
	// TODO set event data
	return ErrMapChange
}
