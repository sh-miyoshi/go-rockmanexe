package damage

import (
	"bytes"
	"encoding/gob"
	"errors"
)

const (
	TargetOwn int = iota
	TargetOtherClient
)

type Damage struct {
	ID            string
	ClientID      string
	PosX          int
	PosY          int
	Power         int
	TTL           int
	TargetType    int
	HitEffectType int
	BigDamage     bool
	ViewOfsX      int32
	ViewOfsY      int32
	ShowHitArea   bool
}

func Marshal(dm []Damage) []byte {
	buf := bytes.NewBuffer(nil)
	gob.NewEncoder(buf).Encode(&dm)
	return buf.Bytes()
}

func Unmarshal(dm *[]Damage, data []byte) {
	buf := bytes.NewBuffer(data)
	_ = gob.NewDecoder(buf).Decode(dm)
}

func (d *Damage) Validate() error {
	if d.ID == "" {
		return errors.New("id is empty")
	}

	if d.ClientID == "" {
		return errors.New("client id is empty")
	}

	if d.TTL <= 0 {
		return errors.New("TTL must be positive number")
	}

	return nil
}
