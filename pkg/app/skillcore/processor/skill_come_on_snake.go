package processor

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	comeOnSnakePhaseInit int = iota
	comeOnSnakePhaseCutin
	comeOnSnakePhaseSnakeAppear
	comeOnSnakePhaseSnakeMove
	comeOnSnakePhaseEnd
)

type snake struct {
	x, y    int
	active  bool
	waiting int
}

type ComeOnSnake struct {
	Arg skillcore.Argument

	count     int
	phase     int
	curColumn int
	snakes    []snake

	// アニメーションと移動の制御用
	waitTime     int
	appearDelay  int
	moveSpeed    int
	snakeWaitMax int
}

func (p *ComeOnSnake) Update() (bool, error) {
	if p.count == 0 {
		// 初期化
		p.phase = comeOnSnakePhaseInit
		p.curColumn = 6 // 敵エリアに最も近い列から開始
		p.waitTime = 0
		p.appearDelay = 30  // 蛇の出現アニメーション待ち時間
		p.moveSpeed = 5     // 蛇の移動速度
		p.snakeWaitMax = 60 // 蛇の待機最大時間
		p.snakes = make([]snake, 0)
	}

	switch p.phase {
	case comeOnSnakePhaseInit:
		// CutIn開始
		p.phase = comeOnSnakePhaseCutin
		p.waitTime = 30 // CutInアニメーション時間
		p.Arg.Cutin("カモンスネーク", 500)

	case comeOnSnakePhaseCutin:
		if p.waitTime > 0 {
			p.waitTime--
			break
		}
		p.phase = comeOnSnakePhaseSnakeAppear

	case comeOnSnakePhaseSnakeAppear:
		if p.curColumn < 1 {
			// すべての列の処理が終了
			p.phase = comeOnSnakePhaseSnakeMove
			break
		}

		// 現在の列で穴パネルを探して蛇を生成
		for y := 0; y < common.FieldNum.Y; y++ {
			panelInfo := p.Arg.GetPanelInfo(point.Point{X: p.curColumn, Y: y})
			if panelInfo.Status == common.PanelStatusHole {
				s := snake{
					x:       p.curColumn,
					y:       y,
					active:  true,
					waiting: p.snakeWaitMax - (6-p.curColumn)*10, // 列に応じて待ち時間を調整
				}
				p.snakes = append(p.snakes, s)
			}
		}

		p.curColumn-- // 次の列へ

	case comeOnSnakePhaseSnakeMove:
		allInactive := true
		damage := false

		// すべての蛇を更新
		for i := range p.snakes {
			if !p.snakes[i].active {
				continue
			}

			allInactive = false

			if p.snakes[i].waiting > 0 {
				p.snakes[i].waiting--
				continue
			}

			// 移動処理
			p.snakes[i].x++
			if p.snakes[i].x >= common.FieldNum.X+1 {
				p.snakes[i].active = false
				continue
			}

			// この列での最初の移動時にダメージを与える
			if !damage {
				damage = true
				// 敵全体にダメージを与える
				// p.Arg.GiveDamageToObject("", 10) // ダメージ値を10に固定
			}
		}

		if allInactive {
			p.phase = comeOnSnakePhaseEnd
		}

	case comeOnSnakePhaseEnd:
		return true, nil
	}

	p.count++
	return false, nil
}

func (p *ComeOnSnake) GetCount() int {
	return p.count
}
