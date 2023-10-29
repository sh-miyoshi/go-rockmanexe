package event

type EndHandler struct {
}

func (h *EndHandler) Draw() {
}

func (h *EndHandler) Process() (bool, error) {
	resultCode = ResultEnd
	return true, nil
}
