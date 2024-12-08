package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
)

type chipNameDraw struct {
	name         string
	count        int
	tm           int
	isUserPlayer bool
}

func SetChipNameDraw(name string, isUserPlayer bool) {
	battlecommon.AddSystem(&chipNameDraw{
		name:         name,
		count:        0,
		tm:           10,
		isUserPlayer: isUserPlayer,
	})
}

func (c *chipNameDraw) Draw() {
	r := float64(0)
	if c.count < c.tm {
		r = float64(c.count) / float64(c.tm)
	} else if c.count < c.tm*3 {
		r = 1
	} else if c.count < c.tm*4 {
		r = 1 - float64(c.count-c.tm*3)/float64(c.tm)
	}

	if r <= 0 {
		return
	}

	x := 50
	if !c.isUserPlayer {
		x = 300
	}

	draw.ExtendString(x, 70, r, 0xffffff, c.name)
}

func (c *chipNameDraw) Update() bool {
	c.count++
	return c.count > c.tm*4
}
