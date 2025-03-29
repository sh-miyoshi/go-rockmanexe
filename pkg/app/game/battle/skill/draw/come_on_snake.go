package skilldraw

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type DrawSnake struct {
}

func (p *DrawSnake) Draw(viewPos point.Point, count int) {
	dxlib.DrawRotaGraph(viewPos.X, viewPos.Y, 1, 0, images[imageTypeComeOnSnake][0], true)
}
