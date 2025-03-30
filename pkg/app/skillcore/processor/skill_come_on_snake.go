package processor

import (
	"github.com/cockroachdb/errors"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	SnakeWaitTime = 60
)

type Snake struct {
	Count   int
	ViewPos point.Point
	Arg     skillcore.Argument
}

type ComeOnSnake struct {
	Arg skillcore.Argument

	count         int
	snakes        []Snake
	candidatePos  []point.Point
	nextSnakeTime int
	currentX      int
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
						p.candidatePos = append(p.candidatePos, pos)
					}
				}
			}
			p.nextSnakeTime = 90 // スキル名表示後に最初の蛇を出現
			p.currentX = p.candidatePos[0].X
		} else {
			return true, errors.New("enemy cannot use ComeOnSnake")
		}
	}

	// 一定間隔で新しい蛇を生成
	if len(p.candidatePos) > 0 && p.count >= p.nextSnakeTime {
		if p.currentX != p.candidatePos[0].X {
			p.nextSnakeTime = p.count + 30 // 次の蛇は30フレーム後
			p.currentX = p.candidatePos[0].X
		} else {
			p.nextSnakeTime = p.count + 10 // 次の蛇は10フレーム後
		}
		p.snakes = append(p.snakes, newSnake(p.candidatePos[0], p.Arg))
		p.candidatePos = p.candidatePos[1:]
	}

	if len(p.snakes) == 0 && len(p.candidatePos) == 0 {
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

func newSnake(initPos point.Point, arg skillcore.Argument) Snake {
	return Snake{
		Count:   0,
		ViewPos: battlecommon.ViewPos(initPos),
		Arg:     arg,
	}
}

func (p *Snake) Update() (bool, error) {
	p.Count++
	if p.Count == 1 {
		p.Arg.SoundOn(resources.SEEnemyAppear)
	} else if p.Count == SnakeWaitTime {
		p.Arg.SoundOn(resources.SEBoomerangThrow)
	}

	if p.Count < SnakeWaitTime {
		return false, nil
	}

	const spd = 16
	p.ViewPos.X += spd
	if p.ViewPos.X > config.ScreenSize.X {
		return true, nil
	}

	return false, nil
}
