package player

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/chip"
)

const (
	defaultHP        uint = 100
	defaultShotPower uint = 1
	separater             = "#"

	// FolderSize ...
	FolderSize = 10 // debug
)

// ChipInfo ...
type ChipInfo struct {
	ID   int
	Code string
}

// Player ...
type Player struct {
	HP        uint
	ShotPower uint
	// Zenny      uint
	ChipFolder [FolderSize]ChipInfo

	WinNum  int
	LoseNum int
	// PlayTime
}

// New returns player data with default values
func New() *Player {
	return &Player{
		HP:        defaultHP,
		ShotPower: defaultShotPower,
		ChipFolder: [FolderSize]ChipInfo{
			{ID: chip.IDRecover10, Code: "*"},
			{ID: chip.IDMiniBomb, Code: "l"},
			{ID: chip.IDSword, Code: "a"},
			{ID: chip.IDWideSword, Code: "a"},
			{ID: chip.IDLongSword, Code: "a"},
			{ID: chip.IDCannon, Code: "b"},
			{ID: chip.IDCannon, Code: "b"},
			{ID: chip.IDCannon, Code: "c"},
			{ID: chip.IDCannon, Code: "c"},
			{ID: chip.IDCannon, Code: "*"},
		},
		WinNum:  0,
		LoseNum: 0,
	}

	/* Correct Init Chip Folder
	ChipFolder: [FolderSize]ChipInfo{
		{ID: chip.IDCannon, Code: "b"},
		{ID: chip.IDCannon, Code: "b"},
		{ID: chip.IDCannon, Code: "c"},
		{ID: chip.IDCannon, Code: "c"},
		{ID: chip.IDHighCannon, Code: "d"},
		{ID: chip.IDHighCannon, Code: "d"},
		{ID: chip.IDMiniBomb, Code: "l"},
		{ID: chip.IDMiniBomb, Code: "l"},
		{ID: chip.IDMiniBomb, Code: "*"},
		{ID: chip.IDMiniBomb, Code: "*"},
		{ID: chip.IDSword, Code: "s"},
		{ID: chip.IDSword, Code: "s"},
		{ID: chip.IDSword, Code: "s"},
		{ID: chip.IDSword, Code: "s"},
		{ID: chip.IDWideSword, Code: "s"},
		{ID: chip.IDWideSword, Code: "s"},
		{ID: chip.IDRecover10, Code: "l"},
		{ID: chip.IDRecover10, Code: "l"},
		{ID: chip.IDRecover10, Code: "*"},
		{ID: chip.IDRecover10, Code: "*"},

		{ID: chip.IDRecover30, Code: "l"},
		{ID: chip.IDRecover30, Code: "l"},
		{ID: chip.IDバルカン1, Code: "b"},
		{ID: chip.IDバルカン1, Code: "b"},
		{ID: chip.IDバルカン1, Code: "d"},
		{ID: chip.IDバルカン1, Code: "d"},
		{ID: chip.IDスプレッドガン, Code: "n"},
		{ID: chip.IDスプレッドガン, Code: "n"},
		{ID: chip.IDスプレッドガン, Code: "m"},
		{ID: chip.IDスプレッドガン, Code: "m"},
	}
	*/
}

// TODO NewWithSaveData(fname string) (*Player, error)

func (p *Player) Save(fname string, key []byte) error {
	// Convert player info to string
	var buf bytes.Buffer
	buf.WriteString(strconv.FormatUint(uint64(p.HP), 10))
	buf.WriteString(separater)
	buf.WriteString(strconv.FormatUint(uint64(p.ShotPower), 10))
	buf.WriteString(separater)
	buf.WriteString(strconv.FormatInt(int64(p.WinNum), 10))
	buf.WriteString(separater)
	buf.WriteString(strconv.FormatInt(int64(p.LoseNum), 10))
	buf.WriteString(separater)
	for _, c := range p.ChipFolder {
		buf.WriteString(fmt.Sprintf("%d%s#", c.ID, c.Code))
	}

	var dst []byte

	if key == nil {
		dst = buf.Bytes()
	} else {
		// TODO Encryption
	}

	return ioutil.WriteFile(fname, dst, 0644)
}
