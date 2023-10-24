package event

type EndHandler struct {
}

func (h *EndHandler) Draw() {
}

func (h *EndHandler) Process() (int, error) {
	return ResultEnd, nil
}
