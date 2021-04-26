package sound

import (
	"fmt"

	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/config"
)

const (
	BGMTitle int = iota
	BGMMenu
	BGMBattle
	BGMWin
	BGMLose

	bgmTypeMax
)

var (
	bgmFiles = [bgmTypeMax]string{
		"title.mp3",
		"menu.mp3",
		"battle.mp3",
		"win.mp3",
		"lose.mp3",
	}
	bgmHandle int32 = -1
)

func BGMPlay(typ int) error {
	if config.Get().BGM.Disabled {
		return nil
	}

	BGMStop()

	fname := common.SoundPath + "bgm/" + bgmFiles[typ]
	if bgmHandle = dxlib.LoadSoundMem(fname); bgmHandle == -1 {
		return fmt.Errorf("failed to load BGM: %s", fname)
	}

	dxlib.PlaySoundMem(bgmHandle, dxlib.DX_PLAYTYPE_LOOP, dxlib.TRUE)

	return nil
}

func BGMStop() {
	if bgmHandle != -1 && dxlib.CheckSoundMem(bgmHandle) == 1 {
		dxlib.StopSoundMem(bgmHandle)
	}
}
