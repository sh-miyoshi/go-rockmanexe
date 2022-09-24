package vector

import "math"

type Vector struct {
	X, Y float64
}

func New(x, y float64) Vector {
	return Vector{X: x, Y: y}
}

func Dot(v1, v2 Vector) float64 {
	return v1.X*v2.X + v1.Y*v2.Y
}

func Cross(v1, v2 Vector) float64 {
	return v1.X*v2.Y - v2.X*v1.Y
}

func Norm(v Vector) float64 {
	res := v.X*v.X + v.Y*v.Y
	return math.Abs(math.Sqrt(res))
}

func Add(v1, v2 Vector) Vector {
	return Vector{X: v1.X + v2.X, Y: v1.Y + v2.Y}
}

func Sub(v1, v2 Vector) Vector {
	return Vector{X: v1.X - v2.X, Y: v1.Y - v2.Y}
}

func Scale(v Vector, a float64) Vector {
	return Vector{X: v.X * a, Y: v.Y * a}
}

func Normalize(v Vector) Vector {
	n := Norm(v)
	return Vector{X: v.X / n, Y: v.Y / n}
}
