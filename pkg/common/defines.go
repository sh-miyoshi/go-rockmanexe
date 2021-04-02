package common

const (
	// ScreenX is x size of screen
	ScreenX = 480
	// ScreenY is y size of screen
	ScreenY = 320
	// SaveFilePath ...
	SaveFilePath = "data/save.dat"
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
	// ImagePath ...
	ImagePath = "data/images/"

	ProgramVersion = "development"
)
