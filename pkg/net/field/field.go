package field

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/object"
)

type Info struct {
	CurrentTime time.Time
	Objects     []object.Object
	Panels      [config.FieldNumX][config.FieldNumY]PanelInfo
	Effects     []effect.Effect
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
	for x := 0; x < config.FieldNumX; x++ {
		id := myClientID
		if x > 2 {
			id = enemyClientID
		}
		for y := 0; y < config.FieldNumY; y++ {
			i.Panels[x][y].OwnerClientID = id
		}
	}
}
