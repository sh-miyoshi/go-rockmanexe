package chipsel

const (
	selectMax = 5 // TODO should variable length
)

var (
	selectList []int
	selected   []int
)

// Init ...
func Init(folder []int) {
	selectList = []int{}
	selected = []int{}

	num := len(folder)
	if num > selectMax {
		num = selectMax
	}
	for i := 0; i < num; i++ {
		selectList = append(selectList, folder[i])
	}
}

// Draw ...
func Draw() {

}

// Process ...
func Process() {

}

// GetSelected ...
func GetSelected() []int {
	return selected
}
