package player

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/chip"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/ncparts"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/stretchr/stew/slice"
)

const (
	defaultHP         uint = 200
	defaultShotPower  uint = 1
	defaultChargeTime uint = 180

	FolderSize          = 30
	SameChipNumInFolder = 4
)

type ChipInfo struct {
	ID   int    `json:"id"`
	Code string `json:"code"`
}

type History struct {
	OpponentID string    `json:"opponent_id"`
	Date       time.Time `json:"date"`
	IsWin      bool      `json:"is_win"`
}

type NaviCustomParts struct {
	ID    int  `json:"id"`
	IsSet bool `json:"is_set"`
	X     int  `json:"x"`
	Y     int  `json:"y"`
}

type Player struct {
	HP                 uint                 `json:"hp"`
	ShotPower          uint                 `json:"shot_power"`
	ChargeTime         uint                 `json:"charge_time"`
	Zenny              uint                 `json:"zenny"`
	ChipFolder         [FolderSize]ChipInfo `json:"chip_folder"`
	WinNum             int                  `json:"win_num"`
	PlayCount          uint                 `json:"play_count"`
	BackPack           []ChipInfo           `json:"back_pack"`
	BattleHistories    []History            `json:"battle_histories"`
	AllNaviCustomParts []NaviCustomParts    `json:"navi_custom_parts"`
}

type SaveData struct {
	Player         Player `json:"player"`
	ProgramVersion string `json:"program_version"`
}

// New returns player data with default values
func New() *Player {
	res := &Player{
		HP:              defaultHP,
		ShotPower:       defaultShotPower,
		ChargeTime:      defaultChargeTime,
		Zenny:           0,
		WinNum:          0,
		BackPack:        []ChipInfo{},
		BattleHistories: []History{},
		AllNaviCustomParts: []NaviCustomParts{
			{ID: ncparts.IDAttack1, IsSet: false},
			{ID: ncparts.IDCharge1, IsSet: false},
			{ID: ncparts.IDHP50, IsSet: false},
			{ID: ncparts.IDHP100, IsSet: false},
			{ID: ncparts.IDUnderShirt, IsSet: false},
		},
	}
	res.setChipFolder()
	res.addPresentChips()
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
		src, err := ioutil.ReadFile(fname)
		if err != nil {
			return nil, fmt.Errorf("failed to read save data: %w", err)
		}
		block, err := aes.NewCipher(key)
		if err != nil {
			return nil, fmt.Errorf("failed to init AES: %w", err)
		}

		iv := src[:aes.BlockSize]
		src = src[aes.BlockSize:]
		if len(bin)%aes.BlockSize != 0 {
			return nil, fmt.Errorf("save data is not a multiple of the block size")
		}

		// Decrypt data with AES-CTR mode
		bin = make([]byte, len(src))
		stream := cipher.NewCTR(block, iv)
		stream.XORKeyStream(bin, src)
	}

	var rawData SaveData
	if err := json.Unmarshal(bin, &rawData); err != nil {
		logger.Error("Failed to unmarshal save data: %v", err)
		return nil, fmt.Errorf("save data maybe broken or invalid version")
	}

	// 互換性維持
	if rawData.Player.ChargeTime == 0 {
		rawData.Player.ChargeTime = defaultChargeTime
	}

	switch rawData.ProgramVersion {
	case "development":
		logger.Info("Save data is development data")
	case "v0.3", "v0.4", "v0.5", "v0.6", "v0.7", "v0.8", "v0.9", "v0.10":
		logger.Info("Save data version is %s, but it is compatible with the current version.", rawData.ProgramVersion)
	case "v0.11":
	default:
		logger.Error("Unexpected version %s is in save data", rawData.ProgramVersion)
		return nil, fmt.Errorf("invalid save data version")
	}

	rawData.Player.addPresentChips()
	return &rawData.Player, nil
}

func (p *Player) Save(fname string, key []byte) error {
	dst, err := json.Marshal(SaveData{
		Player:         *p,
		ProgramVersion: config.ProgramVersion,
	})
	if err != nil {
		return fmt.Errorf("save data marshal failed: %w", err)
	}

	if len(key) == 0 {
		logger.Info("Save with no encryption")
	} else {
		logger.Info("Save with encryption")
		src := append([]byte{}, dst...)
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

func (p *Player) UpdateMoney(diff int) {
	tmp := int(p.Zenny) + diff
	if tmp > config.MaxZenny {
		tmp = config.MaxZenny
	} else if tmp < 0 {
		tmp = 0
	}

	p.Zenny = uint(tmp)
}

func (p *Player) AddChip(id int, code string) error {
	n := 0
	for _, c := range p.ChipFolder {
		if c.ID == id && c.Code == code {
			n++
		}
	}
	for _, c := range p.BackPack {
		if c.ID == id && c.Code == code {
			n++
		}
	}

	if n >= config.MaxChipNum {
		return fmt.Errorf("reached to max chip num")
	}

	p.BackPack = append(p.BackPack, ChipInfo{
		ID:   id,
		Code: code,
	})
	return nil
}

func (p *Player) SetNaviCustomParts(parts []NaviCustomParts) {
	p.AllNaviCustomParts = append([]NaviCustomParts{}, parts...)
	p.updatePlayerStatus()
}

func (p *Player) IsUnderShirt() bool {
	for _, parts := range p.AllNaviCustomParts {
		if parts.ID == ncparts.IDUnderShirt {
			return parts.IsSet
		}
	}

	return false
}

func (p *Player) setChipFolder() {
	if config.Get().Debug.UseDebugFolder {
		// For debug
		p.ChipFolder = [FolderSize]ChipInfo{
			{ID: chip.IDCannon, Code: "*"},
			{ID: chip.IDCannon, Code: "*"},
			{ID: chip.IDCannon, Code: "*"},
			{ID: chip.IDCannon, Code: "*"},
			{ID: chip.IDCannon, Code: "*"},
			{ID: chip.IDCannon, Code: "*"},
			{ID: chip.IDCannon, Code: "*"},
			{ID: chip.IDCannon, Code: "*"},
			{ID: chip.IDCannon, Code: "*"},
			{ID: chip.IDCannon, Code: "*"},
			{ID: chip.IDCannon, Code: "*"},
			{ID: chip.IDCannon, Code: "*"},
			{ID: chip.IDCannon, Code: "*"},
			{ID: chip.IDCannon, Code: "*"},
			{ID: chip.IDCannon, Code: "*"},
			{ID: chip.IDCannon, Code: "*"},
			{ID: chip.IDCannon, Code: "*"},
			{ID: chip.IDCannon, Code: "*"},
			{ID: chip.IDCannon, Code: "*"},
			{ID: chip.IDCannon, Code: "*"},
			{ID: chip.IDCannon, Code: "*"},
			{ID: chip.IDCannon, Code: "*"},
			{ID: chip.IDCannon, Code: "*"},
			{ID: chip.IDCannon, Code: "*"},
			{ID: chip.IDCannon, Code: "*"},
			{ID: chip.IDCannon, Code: "*"},
			{ID: chip.IDCannon, Code: "*"},
			{ID: chip.IDCannon, Code: "*"},
			{ID: chip.IDCannon, Code: "*"},
			{ID: chip.IDCannon, Code: "*"},
		}
	} else {
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
}

func (p *Player) addPresentChips() {
	presentChips := []ChipInfo{
		{ID: chip.IDCrackout, Code: "*"},
		{ID: chip.IDDoubleCrack, Code: "*"},
		{ID: chip.IDTripleCrack, Code: "*"},
		{ID: chip.IDBambooLance, Code: "g"},
		{ID: chip.IDCountBomb, Code: "k"},
		{ID: chip.IDAttack10, Code: "*"},
		{ID: chip.IDQuickGauge, Code: "q"},
	}

	for _, c := range presentChips {
		if slice.Contains(p.ChipFolder, c) {
			continue
		}
		if slice.Contains(p.BackPack, c) {
			continue
		}

		p.BackPack = append(p.BackPack, c)
	}
}

func (p *Player) updatePlayerStatus() {
	p.HP = defaultHP
	p.ShotPower = defaultShotPower
	p.ChargeTime = defaultChargeTime

	// ナビカスによるステータス上昇
	for _, parts := range p.AllNaviCustomParts {
		if parts.IsSet {
			info := ncparts.Get(parts.ID)
			switch info.ID {
			case ncparts.IDAttack1:
				p.ShotPower++
			case ncparts.IDCharge1:
				p.ChargeTime -= 20
			case ncparts.IDHP50:
				p.HP += 50
			case ncparts.IDHP100:
				p.HP += 100
			}
		}
	}
}
