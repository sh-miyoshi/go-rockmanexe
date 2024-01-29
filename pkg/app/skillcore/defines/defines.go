package skilldefines

import "github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"

// TODO: SkillCore側に持たせて、型キャストで取得する

type ShockWaveParam struct {
	InitWait int
	Speed    int
	Direct   int
	ImageNum int
}

func GetShockWaveParam(isPlayer bool) ShockWaveParam {
	if isPlayer {
		return ShockWaveParam{
			InitWait: 9,
			Speed:    3,
			Direct:   config.DirectRight,
			ImageNum: 9,
		}
	} else {
		return ShockWaveParam{
			InitWait: 0,
			Speed:    5,
			Direct:   config.DirectLeft,
			ImageNum: 9,
		}
	}
}
