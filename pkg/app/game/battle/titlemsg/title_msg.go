package titlemsg

import (
	"fmt"

	"github.com/sh-miyoshi/dxlib"
)

const (
	delay = 4
)

type TitleMsg struct {
	imgMsg    []int32
	count     int
	waitCount int
}

func New(fname string, waitCount int) (*TitleMsg, error) {
	res := TitleMsg{
		count:     0,
		imgMsg:    make([]int32, 3),
		waitCount: waitCount,
	}

	if loadRes := dxlib.LoadDivGraph(fname, 3, 1, 3, 274, 32, res.imgMsg); loadRes == -1 {
		return nil, fmt.Errorf("failed to load image %s", fname)
	}

	return &res, nil
}

func (m *TitleMsg) End() {
	for _, img := range m.imgMsg {
		dxlib.DeleteGraph(img)
	}
}

func (m *TitleMsg) Draw() {
	imgNo := m.count / delay
	if imgNo >= len(m.imgMsg) {
		imgNo = len(m.imgMsg) - 1
	}
	dxlib.DrawGraph(105, 125, m.imgMsg[imgNo], dxlib.TRUE)
}

func (m *TitleMsg) Process() bool {
	m.count++
	return m.count >= len(m.imgMsg)*delay+20+m.waitCount
}
