package sound

import (
	"bytes"
	"encoding/gob"
)

type Sound struct {
	ClientID string
	SEType   int
}

func Marshal(se Sound) []byte {
	buf := bytes.NewBuffer(nil)
	gob.NewEncoder(buf).Encode(&se)
	return buf.Bytes()
}

func Unmarshal(se *Sound, data []byte) {
	buf := bytes.NewBuffer(data)
	_ = gob.NewDecoder(buf).Decode(se)
}
