package field

import (
	"bytes"
	"encoding/gob"
	"time"
)

const (
	SizeX = 6
	SizeY = 3
)

type Info struct {
	CurrentTime time.Time
	Objects     []Object

	// TODO PanelInfo
}

func Marshal(fieldInfo *Info) []byte {
	buf := bytes.NewBuffer(nil)
	gob.NewEncoder(buf).Encode(fieldInfo)
	return buf.Bytes()
}

func Unmarshal(fieldInfo *Info, data []byte) {
	buf := bytes.NewBuffer(data)
	_ = gob.NewDecoder(buf).Decode(fieldInfo)
}

func (i *Info) Init() {
}
