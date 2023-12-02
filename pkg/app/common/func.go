package common

import (
	"errors"
	"unicode/utf8"

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

func SplitJAMsg(msg string, max int) []string {
	if max <= 0 {
		return []string{msg}
	}

	res := []string{}
	for len(msg) > 0 {
		tmp := []byte{}
		for i := 0; i < max; i++ {
			r, size := utf8.DecodeRuneInString(msg)
			msg = msg[size:]
			if string(r) == "\n" {
				break
			}
			tmp = utf8.AppendRune(tmp, r)
			if len(msg) <= 0 {
				break
			}
		}
		res = append(res, string(tmp))
	}

	return res
}

func SliceJAMsg(msg string, end int) string {
	tmp := []byte{}
	for i := 0; i < end; i++ {
		if len(msg) <= 0 {
			break
		}

		r, size := utf8.DecodeRuneInString(msg)
		tmp = utf8.AppendRune(tmp, r)
		msg = msg[size:]
	}
	return string(tmp)
}
