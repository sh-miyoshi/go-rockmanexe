package sound

import (
	"fmt"

	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/common"
)

type SEType int32

const (
	SETitleEnter SEType = iota
	SECursorMove
	SEMenuEnter
	SEDenied
	SECancel
	SEGoBattle
	SEEnemyAppear
	SEChipSelectOpen
	SESelect

	seMax
)

var (
	soundEffects = [seMax]int32{}
)

func Init() error {
	basePath := common.SoundPath + "se/"

	soundEffects[SETitleEnter] = dxlib.LoadSoundMem(basePath + "title_enter.mp3")
	soundEffects[SECursorMove] = dxlib.LoadSoundMem(basePath + "cursor_move.mp3")
	soundEffects[SEMenuEnter] = dxlib.LoadSoundMem(basePath + "menu_enter.mp3")
	soundEffects[SEDenied] = dxlib.LoadSoundMem(basePath + "denied.mp3")
	soundEffects[SECancel] = dxlib.LoadSoundMem(basePath + "cancel.mp3")
	soundEffects[SEGoBattle] = dxlib.LoadSoundMem(basePath + "go_battle.mp3")
	soundEffects[SEEnemyAppear] = dxlib.LoadSoundMem(basePath + "enemy_appear.mp3")
	soundEffects[SEChipSelectOpen] = dxlib.LoadSoundMem(basePath + "chip_select_open.mp3")
	soundEffects[SESelect] = dxlib.LoadSoundMem(basePath + "select.mp3")

	for i, s := range soundEffects {
		if s == -1 {
			return fmt.Errorf("failed to load %d sound", i)
		}
	}

	return nil
}

func On(typ SEType) {
	if dxlib.CheckSoundMem(soundEffects[typ]) == 1 {
		dxlib.StopSoundMem(soundEffects[typ])
	}
	dxlib.PlaySoundMem(soundEffects[typ], dxlib.DX_PLAYTYPE_BACK, dxlib.TRUE)
}
