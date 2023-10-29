package collision

import (
	"math"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/mapinfo"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/inputs"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/vector"
)

var (
	walls  = []mapinfo.Wall{}
	events = []mapinfo.Event{}
)

func SetEvents(e []mapinfo.Event) {
	events = append([]mapinfo.Event{}, e...)
}

func SetWalls(w []mapinfo.Wall) {
	walls = append([]mapinfo.Wall{}, w...)
}

func GetEvent(currentX, currentY float64) *mapinfo.Event {
	const r = common.MapPlayerHitRange
	for _, e := range events {
		x2 := (float64(e.X) - currentX) * (float64(e.X) - currentX)
		y2 := (float64(e.Y) - currentY) * (float64(e.Y) - currentY)
		r2 := float64((r + e.R) * (r + e.R))
		if x2+y2 <= r2 {
			if e.RequireAction && inputs.CheckKey(inputs.KeyEnter) == 0 {
				return nil
			}

			return &e
		}
	}
	return nil
}

func NextPos(currentX, currentY float64, goVec vector.Vector) (float64, float64) {
	nextX := currentX + goVec.X
	nextY := currentY + goVec.Y
	hitNum := 0
	for _, w := range walls {
		if isCollision(nextX, nextY, w) {
			v := getWallVec(goVec, w)
			nextX = currentX + v.X
			nextY = currentY + v.Y
			hitNum++
		}
	}

	// 2つ以上にヒットするなら動かさない
	if hitNum >= 2 {
		return currentX, currentY
	}

	return nextX, nextY
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
