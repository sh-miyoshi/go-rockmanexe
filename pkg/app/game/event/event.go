package event

const (
	TypeChangeMapArea int = iota
)

type Handler interface {
	Init(args string) error
	Process()
}

func Draw() {

}

func Process() {

}
