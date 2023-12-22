package system

import (
	"errors"
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
)

var (
	unrecoverableError error
	debugMessages      []string
	font               int = -1
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

func DebugDraw() {
	if config.Get().Debug.ShowDebugData {
		if font == -1 {
			font = dxlib.CreateFontToHandle(dxlib.CreateFontToHandleOption{
				FontName: nil,
				Size:     dxlib.Int32Ptr(22),
				Thick:    dxlib.Int32Ptr(7),
				FontType: dxlib.Int32Ptr(dxlib.DX_FONTTYPE_EDGE),
			})
		}

		for i, msg := range debugMessages {
			dxlib.DrawFormatStringToHandle(0, i*25, 0xffffff, font, msg)
		}
		debugMessages = []string{}
	}
}
