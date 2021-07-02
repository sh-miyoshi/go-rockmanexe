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
	Panels      [SizeX][SizeY]PanelInfo
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

func (i *Info) InitPanel(myClientID, enemyClientID string) {
	for x := 0; x < SizeX; x++ {
		id := myClientID
		if x > 2 {
			id = enemyClientID
		}
		for y := 0; y < SizeY; y++ {
			i.Panels[x][y].OwnerClientID = id
		}
	}
}
