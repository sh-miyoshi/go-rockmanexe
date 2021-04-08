package sound

import (
	"fmt"

	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/config"
)

const (
	BGMTitle int = iota

	bgmTypeMax
)

var (
	bgmFiles = [bgmTypeMax]string{
		"bgm_title.mp3",
	}
)

func BGMPlay(typ int) error {
	if config.Get().BGM.Disabled {
		return nil
	}

	BGMStop()

	fname := common.SoundPath + bgmFiles[typ]
	if res := dxlib.PlaySoundFile(fname, dxlib.DX_PLAYTYPE_LOOP); res == -1 {
		return fmt.Errorf("failed to play BGM: %s", fname)
	}

	return nil
}

func BGMStop() {
	if dxlib.CheckSoundFile() == 1 {
		dxlib.StopSoundFile()
	}
}
