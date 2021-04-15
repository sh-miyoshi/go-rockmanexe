package common

const (
	// ScreenX is x size of screen
	ScreenX = 480
	// ScreenY is y size of screen
	ScreenY = 320

	SaveFilePath      = "data/save.dat"
	DefaultLogFile    = "application.log"
	DefaultConfigFile = "data/config.yaml"
	DxlibDLLFilePath  = "data/Dxlib.dll"
	FontFilePath      = "data/font.ttf"
	ChipFilePath      = "data/chipList.yaml"

	MaxUint = ^uint(0)
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
	ImagePath = "data/images/"
	SoundPath = "data/sounds/"

	ProgramVersion = "development"
)
