package window

import (
	"fmt"
	"unicode/utf8"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/inputs"
)

const (
	messageSpeed = 2
	lineCharNum  = 19
)

type MessageWindow struct {
	image    int
	messages []string
	msgNum   int
	cursor   int
	count    int
}

func (w *MessageWindow) Init() error {
	fname := common.ImagePath + "msg_frame.png"
	w.image = dxlib.LoadGraph(fname)
	if w.image == -1 {
		return fmt.Errorf("failed to load message frame image %s", fname)
	}
	return nil
}

func (w *MessageWindow) End() {
	dxlib.DeleteGraph(w.image)
}

func (w *MessageWindow) SetMessage(msg string) {
	w.messages = common.SplitJAMsg(msg, lineCharNum)
	w.msgNum = utf8.RuneCount([]byte(msg))
}

func (w *MessageWindow) Draw() {
	dxlib.DrawGraph(40, 205, w.image, true)
	for i, msg := range w.messages {
		last := w.cursor - (i * lineCharNum)
		if last > 0 {
			msg = common.SliceJAMsg(msg, last)
			draw.MessageText(120, 220+i*30, 0x000000, msg)
		}
	}
}

func (w *MessageWindow) Process() bool {
	w.count++
	if w.count%messageSpeed == 0 {
		w.cursor++
	}

	if inputs.CheckKey(inputs.KeyCancel) == 1 || inputs.CheckKey(inputs.KeyEnter) == 1 {
		if w.cursor < w.msgNum {
			w.cursor = w.msgNum
		} else {
			return true
		}
	}
	return false
}
