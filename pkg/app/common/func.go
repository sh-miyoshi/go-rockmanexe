package common

func Abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func AbsInt32(a int32) int32 {
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
