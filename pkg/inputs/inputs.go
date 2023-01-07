package inputs

import "fmt"

type KeyType int

const (
	KeyEnter KeyType = iota
	KeyCancel
	KeyLeft
	KeyRight
	KeyUp
	KeyDown
	KeyLButton
	KeyRButton
	KeyDebug

	keyMax
)

const (
	DeviceTypeKeyboard int = iota
	DeviceTypeGamepad
)

type inputDevice interface {
	Init() error
	KeyStateUpdate()
	CheckKey(key KeyType) int
}

var (
	device inputDevice
)

func Init(deviceType int) error {
	switch deviceType {
	case DeviceTypeKeyboard:
		device = &keyboard{}
	case DeviceTypeGamepad:
		device = &pad{}
	default:
		return fmt.Errorf("invalid device type %d specified", deviceType)
	}

	return device.Init()
}

func KeyStateUpdate() {
	if device != nil {
		device.KeyStateUpdate()
	}
}

func CheckKey(key KeyType) int {
	if device != nil {
		return device.CheckKey(key)
	}
	return 0
}
