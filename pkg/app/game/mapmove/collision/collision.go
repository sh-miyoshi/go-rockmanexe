package collision

import (
	"math"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/mapinfo"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/vector"
)

var (
	walls = []mapinfo.Wall{}
)

func SetWalls(w []mapinfo.Wall) {
	walls = append([]mapinfo.Wall{}, w...)
}

func NextPos(current common.Point, goVec vector.Vector) common.Point {
	nextX := float64(current.X) + goVec.X
	nextY := float64(current.Y) + goVec.Y
	for _, w := range walls {
		if isCollision(nextX, nextY, w) {
			v := getWallVec(goVec, w)
			nextX = float64(current.X) + v.X
			nextY = float64(current.X) + v.Y
		}
	}

	return common.Point{X: int(nextX), Y: int(nextY)}
}

func getWallVec(goVec vector.Vector, wall mapinfo.Wall) vector.Vector {
	n := vector.New(-float64(wall.Y2-wall.Y1), float64(wall.X2-wall.X1))
	n = vector.Normalize(n)
	return vector.Sub(goVec, vector.Scale(n, vector.Dot(goVec, n)))
}

func isCollision(x, y float64, wall mapinfo.Wall) bool {
	s := vector.New(float64(wall.X2-wall.X1), float64(wall.Y2-wall.Y1))
	a := vector.New(x-float64(wall.X1), y-float64(wall.Y1))
	b := vector.Sub(a, s)
	sa := vector.Cross(s, a)
	d := math.Abs(sa) / vector.Norm(s)
	const r = common.MapPlayerHitRange

	if d > r {
		return false
	}

	if vector.Dot(a, s)*vector.Dot(b, s) <= 0 {
		return true
	}

	return r > vector.Norm(a) || r > vector.Norm(b)
}
