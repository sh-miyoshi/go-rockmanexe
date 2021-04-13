package player

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/chip"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/common"
)

const (
	defaultHP        uint = 100
	defaultShotPower uint = 1
	separater             = "#"

	// FolderSize ...
	FolderSize = 30
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

	WinNum    int
	PlayCount uint
}

// New returns player data with default values
func New() *Player {
	return &Player{
		HP:        defaultHP,
		ShotPower: defaultShotPower,
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
			{ID: chip.IDVulcan1, Code: "b"},
			{ID: chip.IDVulcan1, Code: "b"},
			{ID: chip.IDVulcan1, Code: "d"},
			{ID: chip.IDVulcan1, Code: "d"},
			{ID: chip.IDSpreadGun, Code: "n"},
			{ID: chip.IDSpreadGun, Code: "n"},
			{ID: chip.IDSpreadGun, Code: "m"},
			{ID: chip.IDSpreadGun, Code: "m"},
		},
		WinNum: 0,
	}

}

// TODO NewWithSaveData(fname string) (*Player, error)

func (p *Player) Save(fname string, key []byte) error {
	// Convert player info to string
	var buf bytes.Buffer
	buf.WriteString(common.ProgramVersion)
	buf.WriteString(separater)
	buf.WriteString(strconv.FormatUint(uint64(p.PlayCount), 10))
	buf.WriteString(separater)
	buf.WriteString(strconv.FormatUint(uint64(p.HP), 10))
	buf.WriteString(separater)
	buf.WriteString(strconv.FormatUint(uint64(p.ShotPower), 10))
	buf.WriteString(separater)
	buf.WriteString(strconv.FormatInt(int64(p.WinNum), 10))
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
