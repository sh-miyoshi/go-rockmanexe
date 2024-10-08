package player

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"io"
	"os"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/chip"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/ncparts"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
)

const (
	defaultHP            uint = 300
	defaultShotPower     uint = 1
	defaultChargeTime    uint = 180
	defaultChipSelectMax      = 5

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
	ChipSelectMax      int                  `json:"chip_select_max"`
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
		ChipSelectMax:   defaultChipSelectMax,
		Zenny:           0,
		WinNum:          0,
		BackPack:        []ChipInfo{},
		BattleHistories: []History{},
		AllNaviCustomParts: []NaviCustomParts{
			{ID: ncparts.IDAttack1_Pink, IsSet: false},
			{ID: ncparts.IDAttack1_White, IsSet: false},
			{ID: ncparts.IDAttack1_White, IsSet: false},
			{ID: ncparts.IDCharge1_Yellow, IsSet: false},
			{ID: ncparts.IDCharge1_Yellow, IsSet: false},
			{ID: ncparts.IDCharge1_White, IsSet: false},
			{ID: ncparts.IDHP50_White, IsSet: false},
			{ID: ncparts.IDHP50_White, IsSet: false},
			{ID: ncparts.IDHP50_White, IsSet: false},
			{ID: ncparts.IDHP100_Yellow, IsSet: false},
			{ID: ncparts.IDHP100_Yellow, IsSet: false},
			{ID: ncparts.IDCustom1_Blue, IsSet: false},
			{ID: ncparts.IDUnderShirt, IsSet: false},
		},
	}
	res.initChipData()
	res.addPresentChips()
	return res
}

func NewWithSaveData(fname string, key []byte) (*Player, error) {
	var bin []byte

	if key == nil {
		var err error
		bin, err = os.ReadFile(fname)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read save data")
		}
	} else {
		src, err := os.ReadFile(fname)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read save data")
		}
		block, err := aes.NewCipher(key)
		if err != nil {
			return nil, errors.Wrap(err, "failed to init AES")
		}

		iv := src[:aes.BlockSize]
		src = src[aes.BlockSize:]
		if len(bin)%aes.BlockSize != 0 {
			return nil, errors.New("save data is not a multiple of the block size")
		}

		// Decrypt data with AES-CTR mode
		bin = make([]byte, len(src))
		stream := cipher.NewCTR(block, iv)
		stream.XORKeyStream(bin, src)
	}

	var rawData SaveData
	if err := json.Unmarshal(bin, &rawData); err != nil {
		logger.Error("Failed to unmarshal save data: %v", err)
		return nil, errors.New("save data maybe broken or invalid version")
	}

	// 互換性維持
	if rawData.Player.ChargeTime == 0 {
		rawData.Player.ChargeTime = defaultChargeTime
	}
	if rawData.Player.ChipSelectMax == 0 {
		rawData.Player.ChipSelectMax = defaultChipSelectMax
	}

	switch rawData.ProgramVersion {
	case "development":
		logger.Info("Save data is development data")
	case "v0.3", "v0.4", "v0.5", "v0.6", "v0.7", "v0.8", "v0.9", "v0.10", "v0.11", "v0.12":
		logger.Error("Save data version is %s, this is not compatible version.", rawData.ProgramVersion)
		return nil, errors.New("save data is not compatible")
	case "v0.13":
	default:
		logger.Error("Unexpected version %s is in save data", rawData.ProgramVersion)
		return nil, errors.New("invalid save data version")
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
		return errors.Wrap(err, "save data marshal failed")
	}

	if len(key) == 0 {
		logger.Info("Save with no encryption")
	} else {
		logger.Info("Save with encryption")
		src := append([]byte{}, dst...)
		block, err := aes.NewCipher(key)
		if err != nil {
			return errors.Wrap(err, "failed to init AES")
		}

		// The IV needs to be unique, but not secure. Therefore it's common to
		// include it at the beginning of the ciphertext.
		dst = make([]byte, aes.BlockSize+len(src))
		iv := dst[:aes.BlockSize]
		if _, err := io.ReadFull(rand.Reader, iv); err != nil {
			return errors.Wrap(err, "failed to read IV")
		}

		// Encrypt data with AES-CTR mode
		stream := cipher.NewCTR(block, iv)
		stream.XORKeyStream(dst[aes.BlockSize:], src)
	}

	return os.WriteFile(fname, dst, 0644)
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

func (p *Player) HaveChip(chipID int) bool {
	for _, c := range p.ChipFolder {
		if c.ID == chipID {
			return true
		}
	}
	for _, c := range p.BackPack {
		if c.ID == chipID {
			return true
		}
	}
	return false
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
		return errors.New("reached to max chip num")
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

func (p *Player) initChipData() {
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

		for _, c := range chip.GetIDList() {
			p.BackPack = append(p.BackPack, ChipInfo{
				ID:   c,
				Code: "*",
			})
		}
	} else {
		// For production
		p.ChipFolder = [FolderSize]ChipInfo{
			{ID: chip.IDCannon, Code: "b"},
			{ID: chip.IDCannon, Code: "b"},
			{ID: chip.IDCannon, Code: "c"},
			{ID: chip.IDCannon, Code: "c"},
			{ID: chip.IDRecover50, Code: "*"},
			{ID: chip.IDRecover50, Code: "*"},
			{ID: chip.IDRecover50, Code: "*"},
			{ID: chip.IDRecover50, Code: "*"},
			{ID: chip.IDShockWave, Code: "s"},
			{ID: chip.IDShockWave, Code: "s"},
			{ID: chip.IDShockWave, Code: "b"},
			{ID: chip.IDShockWave, Code: "b"},
			{ID: chip.IDSpreadGun, Code: "n"},
			{ID: chip.IDSpreadGun, Code: "n"},
			{ID: chip.IDSword, Code: "s"},
			{ID: chip.IDSword, Code: "s"},
			{ID: chip.IDWideSword, Code: "s"},
			{ID: chip.IDWideSword, Code: "s"},
			{ID: chip.IDVulcan1, Code: "b"},
			{ID: chip.IDVulcan1, Code: "b"},
			{ID: chip.IDWideShot1, Code: "m"},
			{ID: chip.IDHeatShot, Code: "f"},
			{ID: chip.IDFlameLine1, Code: "f"},
			{ID: chip.IDFlameLine1, Code: "f"},
			{ID: chip.IDBoomerang1, Code: "m"},
			{ID: chip.IDBoomerang1, Code: "m"},
			{ID: chip.IDBambooLance, Code: "n"},
			{ID: chip.IDBambooLance, Code: "n"},
			{ID: chip.IDCrackout, Code: "*"},
			{ID: chip.IDCrackout, Code: "*"},
		}

		p.BackPack = []ChipInfo{
			{ID: chip.IDHighCannon, Code: "d"},
			{ID: chip.IDHighCannon, Code: "d"},
			{ID: chip.IDMegaCannon, Code: "g"},
			{ID: chip.IDMiniBomb, Code: "*"},
			{ID: chip.IDMiniBomb, Code: "*"},
			{ID: chip.IDMiniBomb, Code: "*"},
			{ID: chip.IDMiniBomb, Code: "*"},
			{ID: chip.IDRecover80, Code: "g"},
			{ID: chip.IDRecover150, Code: "h"},
			{ID: chip.IDSword, Code: "s"},
			{ID: chip.IDSword, Code: "s"},
			{ID: chip.IDWideSword, Code: "s"},
			{ID: chip.IDWideSword, Code: "s"},
			{ID: chip.IDDoubleCrack, Code: "*"},
			{ID: chip.IDDoubleCrack, Code: "*"},
			{ID: chip.IDTripleCrack, Code: "*"},
			{ID: chip.IDTripleCrack, Code: "*"},
		}
	}
}

func (p *Player) addPresentChips() {
	// v0.12時点では何もなし
}

func (p *Player) updatePlayerStatus() {
	p.HP = defaultHP
	p.ShotPower = defaultShotPower
	p.ChargeTime = defaultChargeTime
	p.ChipSelectMax = defaultChipSelectMax

	// ナビカスによるステータス上昇
	for _, parts := range p.AllNaviCustomParts {
		if parts.IsSet {
			info := ncparts.Get(parts.ID)
			switch info.ID {
			case ncparts.IDAttack1_Pink:
				p.ShotPower++
			case ncparts.IDCharge1_Yellow:
				p.ChargeTime -= 20
			case ncparts.IDHP50_White:
				p.HP += 50
			case ncparts.IDHP100_Yellow:
				p.HP += 100
			case ncparts.IDCustom1_Blue:
				p.ChipSelectMax++
			}
		}
	}
}
