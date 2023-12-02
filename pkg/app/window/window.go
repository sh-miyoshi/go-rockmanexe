package window

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/inputs"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/locale/ja"
)

const (
	FaceTypeNone = iota
	FaceTypeRockman

	faceTypeMax
)

const (
	messageSpeed = 2
	lineCharNum  = 19
	maxLineNum   = 3
)

type MessageWindow struct {
	imgFrame int
	imgFaces [faceTypeMax][]int
	messages [][]string
	msgNum   int
	cursor   int
	count    int
	page     int
	faceType int
}

func (w *MessageWindow) Init() error {
	fname := config.ImagePath + "msg_frame.png"
	w.imgFrame = dxlib.LoadGraph(fname)
	if w.imgFrame == -1 {
		return fmt.Errorf("failed to load message frame image %s", fname)
	}

	fname = config.ImagePath + "face/ロックマン.png"
	w.imgFaces[FaceTypeRockman] = make([]int, 5)
	if res := dxlib.LoadDivGraph(fname, 5, 5, 1, 60, 72, w.imgFaces[FaceTypeRockman]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	w.faceType = FaceTypeNone
	return nil
}

func (w *MessageWindow) End() {
	dxlib.DeleteGraph(w.imgFrame)
	for i := 0; i < faceTypeMax; i++ {
		for _, img := range w.imgFaces[i] {
			dxlib.DeleteGraph(img)
		}
		w.imgFaces[i] = []int{}
	}
}

func (w *MessageWindow) SetMessage(msg string, faceType int) {
	messages := ja.SplitMsg(msg, lineCharNum)

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
	w.faceType = faceType
}

func (w *MessageWindow) Draw() {
	dxlib.DrawGraph(40, 205, w.imgFrame, true)
	if w.faceType != FaceTypeNone {
		n := (w.count / 10) % len(w.imgFaces[w.faceType])
		dxlib.DrawGraph(50, 225, w.imgFaces[w.faceType][n], false)
	}

	readNum := 0
	for i, msg := range w.messages[w.page] {
		last := w.cursor - readNum
		if last > 0 {
			msg = ja.SliceMsg(msg, last)
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
