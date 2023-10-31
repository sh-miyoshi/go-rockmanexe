package fps

import (
	"time"
)

var (
	FPS int64 = 60
)

type Fps struct {
	baseTime int64
	count    int64
	current  float32
}

func (f *Fps) Wait() {
	wait := int64(0)
	if f.count == 0 {
		f.baseTime = time.Now().UnixNano() / (1000 * 1000)
	} else {
		c := time.Now().UnixNano() / (1000 * 1000)

		if f.count == FPS-1 {
			// Update current FPS
			f.current = float32(FPS * 1000 / (c - f.baseTime))
		}

		target := f.count*1000/FPS + f.baseTime
		wait = target - c
	}

	if wait > 0 {
		time.Sleep(time.Millisecond * time.Duration(wait))
	}
	f.count = (f.count + 1) % FPS
}

func (f *Fps) Get() float32 {
	return f.current
}
