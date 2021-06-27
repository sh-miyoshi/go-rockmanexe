package fps

import (
	"time"
)

type Fps struct {
	TargetFPS int64

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

		if f.count == f.TargetFPS-1 {
			// Update current FPS
			f.current = float32(f.TargetFPS * 1000 / (c - f.baseTime))
		}

		target := f.count*1000/f.TargetFPS + f.baseTime
		wait = target - c
	}

	if wait > 0 {
		time.Sleep(time.Millisecond * time.Duration(wait))
	}
	f.count = (f.count + 1) % f.TargetFPS
}

func (f *Fps) Get() float32 {
	return f.current
}
