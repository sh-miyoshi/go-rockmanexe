package common

import (
	"unicode/utf8"
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
		if string(r) == "\n" {
			break
		}
		tmp = utf8.AppendRune(tmp, r)
		msg = msg[size:]
	}
	return string(tmp)
}
