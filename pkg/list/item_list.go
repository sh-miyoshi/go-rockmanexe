package list

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/inputs"
)

type ItemList struct {
	cursor     int
	lists      []string
	scroll     int
	maxShowNum int
}

func (l *ItemList) SetList(lists []string, maxShowNum int) {
	l.lists = append([]string{}, lists...)
	l.maxShowNum = maxShowNum
}

func (l *ItemList) GetList() []string {
	return l.lists
}

func (l *ItemList) GetPointer() int {
	return l.cursor
}

func (l *ItemList) GetScroll() int {
	return l.scroll
}

func (l *ItemList) Process() int {
	if inputs.CheckKey(inputs.KeyEnter) == 1 {
		return l.cursor + l.scroll
	}
	if inputs.CheckKey(inputs.KeyUp)%10 == 1 {
		if l.cursor > 0 {
			sound.On(resources.SECursorMove)
			l.cursor--
		} else if l.maxShowNum > 0 && l.scroll > 0 {
			sound.On(resources.SECursorMove)
			l.scroll--
		}
	} else if inputs.CheckKey(inputs.KeyDown)%10 == 1 {
		n := len(l.lists) - 1
		if l.maxShowNum > 0 && l.maxShowNum < n {
			n = l.maxShowNum - 1
		}

		if l.cursor < n {
			sound.On(resources.SECursorMove)
			l.cursor++
		} else if l.maxShowNum > 0 && l.scroll < len(l.lists)-l.maxShowNum {
			sound.On(resources.SECursorMove)
			l.scroll++
		}
	}

	return -1
}
