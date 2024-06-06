package system

import (
	"fmt"

	"github.com/cockroachdb/errors"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
)

var (
	unrecoverableError error
	debugMessages      []string
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

func AddDebugMessage(format string, a ...interface{}) {
	debugMessages = append(debugMessages, fmt.Sprintf(format, a...))
}

func PopAllDebugMessages() []string {
	res := append([]string{}, debugMessages...)
	debugMessages = []string{}
	return res
}
