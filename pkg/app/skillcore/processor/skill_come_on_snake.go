package processor

import (
	"github.com/cockroachdb/errors"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	SnakeWaitTime = 60
)

type Snake struct {
	Count   int
	ViewPos point.Point
}

type ComeOnSnake struct {
	Arg skillcore.Argument

	count  int
	snakes []Snake
}

func (p *ComeOnSnake) Update() (bool, error) {
	p.count++
	if p.count == 1 {
		p.Arg.Cutin("カモンスネーク", 500)
		if p.Arg.TargetType == damage.TargetEnemy {
			for x := battlecommon.FieldNum.X - 1; x >= 0; x-- {
				for y := 0; y < battlecommon.FieldNum.Y; y++ {
					pos := point.Point{X: x, Y: y}
					pn := p.Arg.GetPanelInfo(pos)
					if pn.Type == battlecommon.PanelTypePlayer && pn.Status == battlecommon.PanelStatusHole {
						p.snakes = append(p.snakes, newSnake(pos))
					}
				}
			}
		} else {
			return true, errors.New("enemy cannot use ComeOnSnake")
		}
	}
	// スキル名表示中は待つ
	if p.count < 90 {
		return false, nil
	}

	if len(p.snakes) == 0 {
		return true, nil
	}

	// スネークの更新と削除処理
	newSnakes := make([]Snake, 0, len(p.snakes))
	for i := range p.snakes {
		finished, err := p.snakes[i].Update()
		if err != nil {
			return false, err
		}
		if !finished {
			newSnakes = append(newSnakes, p.snakes[i])
		}
	}
	p.snakes = newSnakes

	return false, nil
}

func (p *ComeOnSnake) GetCount() int {
	return p.count
}

func (p *ComeOnSnake) GetSnakes() []Snake {
	return p.snakes
}

func newSnake(initPos point.Point) Snake {
	return Snake{
		Count:   0,
		ViewPos: battlecommon.ViewPos(initPos),
	}
}

func (p *Snake) Update() (bool, error) {
	p.Count++
	if p.Count < SnakeWaitTime {
		return false, nil
	}

	const spd = 2
	p.ViewPos.X += spd
	if p.ViewPos.X > config.ScreenSize.X {
		return true, nil
	}

	return false, nil
}
