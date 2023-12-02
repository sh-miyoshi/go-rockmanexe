package point

import "fmt"

type Point struct {
	X int
	Y int
}

func (p Point) String() string {
	return fmt.Sprintf("(%d, %d)", p.X, p.Y)
}

func (p Point) Equal(a Point) bool {
	return p.X == a.X && p.Y == a.Y
}
