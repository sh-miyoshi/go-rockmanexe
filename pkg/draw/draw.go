package draw

import (
	"fmt"

	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/common"
)

type NumberOption struct {
	Color        int // defualt is NumberColorWhite
	Centered     bool
	RightAligned bool
	Length       int // Required if RightAligned is tru
}

const (
	// The order of number color depends on the image

	NumberColorWhite int = iota
	NumberColorRed
	NumberColorGreen
	NumberColorWhiteSmall

	numberColorMax
)

const (
	numberSizeX = 15
)

var (
	fontHandle int32 = -1
	imgCode    []int32
	imgNumber  [numberColorMax][]int32
)

func Init() error {
	// Set font
	fontHandle = dxlib.CreateFontToHandle(dxlib.CreateFontToHandleOption{
		FontName: dxlib.StringPtr("k8x12"),
		Size:     dxlib.Int32Ptr(22),
		Thick:    dxlib.Int32Ptr(7),
	})
	if fontHandle == -1 {
		return fmt.Errorf("Failed to create font")
	}

	// Load chip code
	imgCode = make([]int32, 27)
	fname := common.ImagePath + "chipInfo/chip_code.png"
	if res := dxlib.LoadDivGraph(fname, 27, 9, 3, 20, 26, imgCode); res == -1 {
		return fmt.Errorf("Failed to load chip code image %s", fname)
	}

	// Load number data
	tmp := make([]int32, 3*10)
	fname = common.ImagePath + "number.png"
	if res := dxlib.LoadDivGraph(fname, 30, 10, 3, numberSizeX, 26, tmp); res == -1 {
		return fmt.Errorf("Failed to load number image %s", fname)
	}
	// Sort and set to start from 0
	for i := 0; i < 3; i++ {
		imgNumber[i] = make([]int32, 10)
		imgNumber[i][0] = tmp[i*10+9]
		for n := 0; n < 9; n++ {
			imgNumber[i][n+1] = tmp[i*10+n]
		}
	}
	fname = common.ImagePath + "number_small.png"
	if res := dxlib.LoadDivGraph(fname, 10, 10, 1, numberSizeX, 20, tmp); res == -1 {
		return fmt.Errorf("Failed to load small number image %s", fname)
	}
	// Sort and set to start from 0
	imgNumber[NumberColorWhiteSmall] = make([]int32, 10)
	imgNumber[NumberColorWhiteSmall][0] = tmp[9]
	for n := 0; n < 9; n++ {
		imgNumber[NumberColorWhiteSmall][n+1] = tmp[n]
	}

	return nil
}

func String(x int32, y int32, color uint32, format string, a ...interface{}) {
	dxlib.DrawFormatStringToHandle(x, y, color, fontHandle, format, a...)
}

func ChipCode(x int32, y int32, code string, percent int32) {
	index := -1
	if len(code) != 1 {
		panic(fmt.Sprintf("Invalid chip code %s is specified.", code))
	}

	rc := []rune(code)
	if rc[0] >= 'a' && rc[0] <= 'z' {
		index = int(rc[0] - 'a')
	} else if rc[0] >= 'A' && rc[0] <= 'Z' {
		index = int(rc[0] - 'A')
	} else if rc[0] == '*' {
		index = 26
	} else {
		panic(fmt.Sprintf("Invalid chip code %s is specified.", code))
	}

	if percent == 100 {
		dxlib.DrawGraph(x, y, imgCode[index], dxlib.FALSE)
	} else {
		dxlib.DrawExtendGraph(x, y, x+20*percent/100, y+26*percent/100, imgCode[index], dxlib.FALSE)
	}
}

func Number(x int32, y int32, number int32, opts ...NumberOption) {
	nums := []int{}
	for number > 0 {
		nums = append(nums, int(number)%10)
		number /= 10
	}
	for i := 0; i < len(nums)/2; i++ {
		nums[i], nums[len(nums)-i-1] = nums[len(nums)-i-1], nums[i]
	}

	color := NumberColorWhite
	if len(opts) > 0 {
		color = opts[0].Color
		if opts[0].Centered {
			x -= int32(len(nums) * numberSizeX / 2)
		} else if opts[0].RightAligned {
			n := opts[0].Length - len(nums)
			if n < 0 {
				panic(fmt.Sprintf("Failed to show %d with right aligned. requires more %d length", number, -n))
			}
			x += int32(n * numberSizeX)
		}
	}

	for _, n := range nums {
		dxlib.DrawGraph(x, y, imgNumber[color][n], dxlib.TRUE)
		x += numberSizeX
	}
}
