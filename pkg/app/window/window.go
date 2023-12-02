package window

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/inputs"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
)

const (
	messageSpeed = 2
	lineCharNum  = 19
	maxLineNum   = 3
)

type MessageWindow struct {
	image    int
	messages [][]string
	msgNum   int
	cursor   int
	count    int
	page     int
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
	messages := common.SplitJAMsg(msg, lineCharNum)

	// 複数行のMessageをmaxLineNumごとのMessage配列に分割する
	tmp := []string{}
	w.messages = [][]string{}
	for _, m := range messages {
		tmp = append(tmp, m)
		if len(tmp) >= maxLineNum {
			w.messages = append(w.messages, tmp)
			tmp = []string{}
		}
	}
	if len(tmp) > 0 {
		w.messages = append(w.messages, tmp)
		w.msgNum = messageCount(w.messages[0])
	} else {
		w.msgNum = 0
	}

	w.cursor = 0
	w.page = 0
}

func (w *MessageWindow) Draw() {
	dxlib.DrawGraph(40, 205, w.image, true)
	readNum := 0
	for i, msg := range w.messages[w.page] {
		last := w.cursor - readNum
		if last > 0 {
			msg = common.SliceJAMsg(msg, last)
			draw.MessageText(120, 220+i*30, 0x000000, msg)
			readNum += utf8.RuneCount([]byte(msg))
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
		} else if w.page >= len(w.messages)-1 {
			// 最終ページなら終了
			return true
		} else {
			w.cursor = 0
			w.page++
			w.msgNum = messageCount(w.messages[w.page])
		}
	}
	return false
}

func messageCount(messages []string) int {
	res := 0
	for _, msg := range messages {
		res += utf8.RuneCount([]byte(msg))
		res -= strings.Count(msg, "\n")
	}
	logger.Debug("message count: %d for message %+v", res, messages)
	return res
}
