package draw

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/system"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

type NumberOption struct {
	Color        int // defualt is NumberColorWhite
	Centered     bool
	RightAligned bool
	Padding      *int
	Length       int // Required if RightAligned is true or Padding is set
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
	defaultFont int = -1
	msgFont     int = -1
	paTitleFont int = -1
	imgCode     []int
	imgNumber   [numberColorMax][]int
)

func Init() error {
	// Set font
	defaultFont = dxlib.CreateFontToHandle(dxlib.CreateFontToHandleOption{
		FontName: dxlib.StringPtr("k8x12"),
		Size:     dxlib.Int32Ptr(22),
		Thick:    dxlib.Int32Ptr(7),
	})
	if defaultFont == -1 {
		return fmt.Errorf("failed to create default font")
	}

	msgFont = dxlib.CreateFontToHandle(dxlib.CreateFontToHandleOption{
		FontName: dxlib.StringPtr("k8x12"),
		Size:     dxlib.Int32Ptr(24),
		Thick:    dxlib.Int32Ptr(4),
	})
	if msgFont == -1 {
		return fmt.Errorf("failed to create message font")
	}

	paTitleFont = dxlib.CreateFontToHandle(dxlib.CreateFontToHandleOption{
		FontName: dxlib.StringPtr("k8x12"),
		Size:     dxlib.Int32Ptr(28),
		Thick:    dxlib.Int32Ptr(9),
		EdgeSize: dxlib.Int32Ptr(4),
		Italic:   dxlib.Int32Ptr(dxlib.TRUE),
		FontType: dxlib.Int32Ptr(dxlib.DX_FONTTYPE_EDGE),
	})
	if paTitleFont == -1 {
		return fmt.Errorf("failed to create program advance title font")
	}

	// Load chip code
	imgCode = make([]int, 27)
	fname := config.ImagePath + "chipInfo/chip_code.png"
	if res := dxlib.LoadDivGraph(fname, 27, 9, 3, 20, 26, imgCode); res == -1 {
		return fmt.Errorf("failed to load chip code image %s", fname)
	}

	// Load number data
	tmp := make([]int, 3*10)
	fname = config.ImagePath + "number.png"
	if res := dxlib.LoadDivGraph(fname, 30, 10, 3, numberSizeX, 26, tmp); res == -1 {
		return fmt.Errorf("failed to load number image %s", fname)
	}
	// Sort and set to start from 0
	for i := 0; i < 3; i++ {
		imgNumber[i] = make([]int, 10)
		imgNumber[i][0] = tmp[i*10+9]
		for n := 0; n < 9; n++ {
			imgNumber[i][n+1] = tmp[i*10+n]
		}
	}
	fname = config.ImagePath + "number_small.png"
	if res := dxlib.LoadDivGraph(fname, 10, 10, 1, numberSizeX, 20, tmp); res == -1 {
		return fmt.Errorf("failed to load small number image %s", fname)
	}
	// Sort and set to start from 0
	imgNumber[NumberColorWhiteSmall] = make([]int, 10)
	imgNumber[NumberColorWhiteSmall][0] = tmp[9]
	for n := 0; n < 9; n++ {
		imgNumber[NumberColorWhiteSmall][n+1] = tmp[n]
	}

	return nil
}

func String(x int, y int, color uint, format string, a ...interface{}) {
	dxlib.DrawFormatStringToHandle(x, y, color, defaultFont, format, a...)
}

func MessageText(x int, y int, color uint, format string, a ...interface{}) {
	dxlib.DrawFormatStringToHandle(x, y, color, msgFont, format, a...)
}

func PAText(x int, y int) {
	dxlib.DrawStringToHandle(x, y, 0xffffff, paTitleFont, "プログラムアドバンス")
}

func ExtendString(x int, y int, exRateY float64, color uint, format string, a ...interface{}) {
	dxlib.DrawExtendFormatStringToHandle(x, y, 1, exRateY, color, defaultFont, format, a...)
}

func ChipCode(x int, y int, code string, percent int) {
	index := -1
	if len(code) != 1 {
		system.SetError(fmt.Sprintf("Invalid chip code %s is specified.", code))
		return
	}

	rc := []rune(code)
	if rc[0] >= 'a' && rc[0] <= 'z' {
		index = int(rc[0] - 'a')
	} else if rc[0] >= 'A' && rc[0] <= 'Z' {
		index = int(rc[0] - 'A')
	} else if rc[0] == '*' {
		index = 26
	} else {
		system.SetError(fmt.Sprintf("Invalid chip code %s is specified.", code))
		return
	}

	if percent == 100 {
		dxlib.DrawGraph(x, y, imgCode[index], false)
	} else {
		dxlib.DrawExtendGraph(x, y, x+20*percent/100, y+26*percent/100, imgCode[index], false)
	}
}

func Number(x int, y int, number int, opts ...NumberOption) {
	nums := []int{}
	for number > 0 {
		nums = append(nums, number%10)
		number /= 10
	}
	for i := 0; i < len(nums)/2; i++ {
		nums[i], nums[len(nums)-i-1] = nums[len(nums)-i-1], nums[i]
	}

	color := NumberColorWhite
	if len(opts) > 0 {
		color = opts[0].Color
		if opts[0].Centered {
			x -= len(nums) * numberSizeX / 2
		} else if opts[0].RightAligned {
			n := opts[0].Length - len(nums)
			if n < 0 {
				system.SetError(fmt.Sprintf("Failed to show %d with right aligned. requires more %d length", number, -n))
				return
			}
			x += n * numberSizeX
		} else if opts[0].Padding != nil {
			v := *opts[0].Padding
			n := opts[0].Length - len(nums)
			for i := 0; i < n; i++ {
				nums = append([]int{v}, nums...)
			}
		}
	}

	for _, n := range nums {
		dxlib.DrawGraph(x, y, imgNumber[color][n], true)
		x += numberSizeX
	}
}
