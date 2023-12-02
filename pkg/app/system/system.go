package system

import (
	"errors"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
)

var (
	unrecoverableError error
)

func SetError(msg string) {
	// この関数が呼ばれた場所の呼び出し元情報をセットする
	logger.SetExtraSkipCount(1)
	logger.Error(msg)
	logger.SetExtraSkipCount(0)
	unrecoverableError = errors.New("ゲームプレイ中")
}

func Error() error {
	return unrecoverableError
}
