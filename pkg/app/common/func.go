package common

import (
	"errors"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
)

func Abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func MountainIndex(i, max int) int {
	if i >= max/2 {
		return max - i - 1
	} else {
		return i
	}
}

func SetError(msg string) {
	// この関数が呼ばれた場所の呼び出し元情報をセットする
	logger.SetExtraSkipCount(1)
	logger.Error(msg)
	logger.SetExtraSkipCount(0)
	IrreversibleError = errors.New("ゲームプレイ中")
}
