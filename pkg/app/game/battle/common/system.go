package common

type System interface {
	Draw()
	Update() bool
}

var (
	systemTable = []System{}
)

func AddSystem(sys System) {
	systemTable = append(systemTable, sys)
}

func SystemDraw() {
	for _, s := range systemTable {
		s.Draw()
	}
}

func SystemProcess() {
	newTable := []System{}
	for _, s := range systemTable {
		if !s.Update() {
			newTable = append(newTable, s)
		}
	}
	systemTable = newTable
}
