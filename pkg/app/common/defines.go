package common

const (
	SaveFilePath      = "data/save.dat"
	DefaultLogFile    = "application.log"
	DefaultConfigFile = "data/config.yaml"
	DxlibDLLFilePath  = "data/Dxlib.dll"
	FontFilePath      = "data/font.ttf"
	ChipFilePath      = "data/chipList.yaml"

	MaxUint    = ^uint(0)
	MaxZenny   = 9999999
	MaxChipNum = 99

	MapPlayerHitRange = 10
)

const (
	// DirectUp ...
	DirectUp = 1 << iota
	// DirectLeft ...
	DirectLeft
	// DirectDown ...
	DirectDown
	// DirectRight ...
	DirectRight
)

var (
	ScreenSize = Point{X: 480, Y: 320}

	ImagePath = "data/images/"
	SoundPath = "data/sounds/"

	ProgramVersion = "development"
	EncryptKey     = ""
)
