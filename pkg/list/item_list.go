package list

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/inputs"
)

type ItemList struct {
	current int
	lists   []string
}

func (l *ItemList) SetList(lists []string) {
	l.lists = append([]string{}, lists...)
}

func (l *ItemList) GetList() []string {
	return l.lists
}

func (l *ItemList) GetPointer() int {
	return l.current
}

func (l *ItemList) Process() int {
	if inputs.CheckKey(inputs.KeyEnter) == 1 {
		return l.current
	}
	if inputs.CheckKey(inputs.KeyUp) == 1 {
		if l.current > 0 {
			sound.On(resources.SECursorMove)
			l.current--
		} else {
			sound.On(resources.SEBlock)
		}
	} else if inputs.CheckKey(inputs.KeyDown) == 1 {
		if l.current < len(l.lists)-1 {
			sound.On(resources.SECursorMove)
			l.current++
		} else {
			sound.On(resources.SEBlock)
		}
	}

	return -1
}
