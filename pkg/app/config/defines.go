package config

import "github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"

const (
	SaveFilePath      = "data/save.dat"
	DefaultLogFile    = "application.log"
	DefaultConfigFile = "data/config.yaml"
	DxlibDLLFilePath  = "data/Dxlib.dll"
	FontFilePath      = "data/font.ttf"
	ChipFilePath      = "data/chipList.yaml"
	MapInfoFilePath   = "data/mapInfo.yaml"

	MaxUint    = ^uint(0)
	MaxZenny   = 9999999
	MaxChipNum = 99
)

const (
	DirectUp = 1 << iota
	DirectLeft
	DirectDown
	DirectRight
)

var (
	MaxScreenSize = point.Point{X: 640, Y: 480}
	ScreenSize    = point.Point{X: 480, Y: 320}

	ImagePath = "data/images/"
	SoundPath = "data/sounds/"

	ProgramVersion = "development"
	EncryptKey     = ""
)
