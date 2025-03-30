package skilldraw

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type DrawSnake struct {
}

func (p *DrawSnake) Draw(viewPos point.Point, count int) {
	imageIndex := 2
	if count < processor.SnakeWaitTime-10 {
		imageIndex = 0
	} else if count < processor.SnakeWaitTime {
		imageIndex = 1
	}
	dxlib.DrawRotaGraph(viewPos.X+20, viewPos.Y, 1, 0, images[imageTypeComeOnSnake][imageIndex], true)
}
