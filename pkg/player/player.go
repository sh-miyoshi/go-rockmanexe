package player

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/chip"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
)

const (
	defaultHP        uint = 100
	defaultShotPower uint = 1
	separater             = "#"

	FolderSize          = 30
	SameChipNumInFolder = 4
)

// ChipInfo ...
type ChipInfo struct {
	ID   int
	Code string
}

// Player ...
type Player struct {
	HP         uint
	ShotPower  uint
	Zenny      uint
	ChipFolder [FolderSize]ChipInfo
	WinNum     int
	PlayCount  uint
	BackPack   []ChipInfo
}

// New returns player data with default values
func New() *Player {
	res := &Player{
		HP:        defaultHP,
		ShotPower: defaultShotPower,
		Zenny:     0,
		WinNum:    0,
		BackPack:  []ChipInfo{},
	}
	res.setChipFolder()
	return res
}

func NewWithSaveData(fname string, key []byte) (*Player, error) {
	var bin []byte

	if key == nil {
		var err error
		bin, err = ioutil.ReadFile(fname)
		if err != nil {
			return nil, fmt.Errorf("failed to read save data: %w", err)
		}
	} else {
		return nil, fmt.Errorf("not implement yet")
	}

	data := strings.Split(string(bin), separater)
	if len(data) != 6+FolderSize+1 {
		logger.Error("save data requires %d data, but got %d", 6+FolderSize+1, len(data))
		return nil, fmt.Errorf("save data maybe broken or invalid version")
	}

	version := data[0]
	if version != common.ProgramVersion {
		logger.Error("Invalid version save data. expect %s, but got %s", common.ProgramVersion, version)
		return nil, fmt.Errorf("version miss matched")
	}

	playCnt, err := strconv.ParseUint(data[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse play count: %w", err)
	}
	hp, err := strconv.ParseUint(data[2], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse hp: %w", err)
	}
	shot, err := strconv.ParseUint(data[3], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse shot power: %w", err)
	}
	zenny, err := strconv.ParseUint(data[4], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse zenny: %w", err)
	}
	win, err := strconv.ParseInt(data[5], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse win num: %w", err)
	}

	// TODO BackPack

	res := &Player{
		PlayCount: uint(playCnt),
		HP:        uint(hp),
		ShotPower: uint(shot),
		Zenny:     uint(zenny),
		WinNum:    int(win),
	}

	for i := 0; i < FolderSize; i++ {
		var id int
		var code string
		if _, err := fmt.Sscanf(data[6+i], "%d%s", &id, &code); err != nil {
			return nil, fmt.Errorf("failed to parse chip %d: %w", i, err)
		}
		res.ChipFolder[i].ID = id
		res.ChipFolder[i].Code = code
	}

	return res, nil
}

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
	buf.WriteString(strconv.FormatUint(uint64(p.Zenny), 10))
	buf.WriteString(separater)
	buf.WriteString(strconv.FormatInt(int64(p.WinNum), 10))
	buf.WriteString(separater)

	// TODO BackPack

	for _, c := range p.ChipFolder {
		buf.WriteString(fmt.Sprintf("%d%s%s", c.ID, c.Code, separater))
	}

	var dst []byte

	if key == nil {
		dst = buf.Bytes()
	} else {
		src := buf.Bytes()
		block, err := aes.NewCipher(key)
		if err != nil {
			return fmt.Errorf("failed to init AES: %w", err)
		}

		// The IV needs to be unique, but not secure. Therefore it's common to
		// include it at the beginning of the ciphertext.
		dst = make([]byte, aes.BlockSize+len(src))
		iv := dst[:aes.BlockSize]
		if _, err := io.ReadFull(rand.Reader, iv); err != nil {
			return fmt.Errorf("failed to read IV: %w", err)
		}

		// Encrypt data with AES-CTR mode
		stream := cipher.NewCTR(block, iv)
		stream.XORKeyStream(dst[aes.BlockSize:], src)
	}

	return ioutil.WriteFile(fname, dst, 0644)
}

func (p *Player) setChipFolder() {
	// For debug
	// p.ChipFolder = [FolderSize]ChipInfo{
	// 	{ID: chip.IDShockWave, Code: "*"},
	// 	{ID: chip.IDShockWave, Code: "*"},
	// 	{ID: chip.IDShockWave, Code: "*"},
	// 	{ID: chip.IDShockWave, Code: "*"},
	// 	{ID: chip.IDShockWave, Code: "*"},
	// 	{ID: chip.IDShockWave, Code: "*"},
	// 	{ID: chip.IDShockWave, Code: "*"},
	// 	{ID: chip.IDShockWave, Code: "*"},
	// 	{ID: chip.IDShockWave, Code: "*"},
	// 	{ID: chip.IDShockWave, Code: "*"},
	// 	{ID: chip.IDShockWave, Code: "*"},
	// 	{ID: chip.IDShockWave, Code: "*"},
	// 	{ID: chip.IDShockWave, Code: "*"},
	// 	{ID: chip.IDShockWave, Code: "*"},
	// 	{ID: chip.IDShockWave, Code: "*"},
	// 	{ID: chip.IDShockWave, Code: "*"},
	// 	{ID: chip.IDShockWave, Code: "*"},
	// 	{ID: chip.IDShockWave, Code: "*"},
	// 	{ID: chip.IDShockWave, Code: "*"},
	// 	{ID: chip.IDShockWave, Code: "*"},
	// 	{ID: chip.IDShockWave, Code: "*"},
	// 	{ID: chip.IDShockWave, Code: "*"},
	// 	{ID: chip.IDShockWave, Code: "*"},
	// 	{ID: chip.IDShockWave, Code: "*"},
	// 	{ID: chip.IDShockWave, Code: "*"},
	// 	{ID: chip.IDShockWave, Code: "*"},
	// 	{ID: chip.IDShockWave, Code: "*"},
	// 	{ID: chip.IDShockWave, Code: "*"},
	// 	{ID: chip.IDShockWave, Code: "*"},
	// 	{ID: chip.IDShockWave, Code: "*"},
	// }

	// For production
	p.ChipFolder = [FolderSize]ChipInfo{
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
	}
}
