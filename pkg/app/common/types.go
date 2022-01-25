package common

type Point struct {
	X int
	Y int
}

func (p *Point) Add(a Point) Point {
	return Point{X: p.X + a.X, Y: p.Y + a.Y}
}
