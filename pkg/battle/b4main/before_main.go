package b4main

import (
	"fmt"

	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/common"
)

var (
	imgMsg []int32
	delay  = 4
	count  int
)

func Init() error {
	count = 0

	imgMsg = make([]int32, 3)
	fname := common.ImagePath + "battle/msg_start.png"
	if res := dxlib.LoadDivGraph(fname, 3, 1, 3, 274, 32, imgMsg); res == -1 {
		return fmt.Errorf("Failed to load start message image %s", fname)
	}
	return nil
}

func End() {
	for _, img := range imgMsg {
		dxlib.DeleteGraph(img)
	}
	imgMsg = []int32{}
}

func Draw() {
	if len(imgMsg) == 0 {
		// Waiting initialize
		return
	}

	imgNo := count / delay
	if imgNo >= len(imgMsg) {
		imgNo = len(imgMsg) - 1
	}
	dxlib.DrawGraph(105, 125, imgMsg[imgNo], dxlib.TRUE)
}

func Process() bool {
	count++

	if count >= len(imgMsg)*delay+20 {
		return true
	}
	return false
}
