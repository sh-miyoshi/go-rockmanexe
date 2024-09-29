package math

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

func ReverseIf(val int, isReverse bool) int {
	if isReverse {
		return -val
	}
	return val
}

func Sign(x int) int {
	if x > 0 {
		return 1
	} else if x < 0 {
		return -1
	}
	return 0
}
