package sound

import (
	"fmt"

	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
)

type SEType int32

const (
	SENone SEType = iota // required no SE as 0
	SETitleEnter
	SECursorMove
	SEMenuEnter
	SEDenied
	SECancel
	SEGoBattle
	SEEnemyAppear
	SEChipSelectOpen
	SESelect
	SEChipSelectEnd
	SEGaugeMax
	SECannon
	SEBusterCharging
	SEBusterCharged
	SEBusterShot
	SECannonHit
	SEExplode
	SEBusterHit
	SESword
	SERecover
	SEShockWave
	SEGun
	SESpreadHit
	SEBombThrow
	SEPlayerDeleted
	SEDamaged
	SEEnemyDeleted
	SEGotItem
	SEWindowChange
	SEThunderBall
	SEWideShot
	SEBoomerangThrow
	SEWaterLanding
	SEBlock
	SEObjectCreate
	SEWaterpipeAttack
	SEPanelBreak

	seMax
)

var (
	soundEffects = [seMax]int32{}
)

func Init() error {
	basePath := common.SoundPath + "se/"

	soundEffects[SENone] = 0
	soundEffects[SETitleEnter] = dxlib.LoadSoundMem(basePath + "title_enter.mp3")
	soundEffects[SECursorMove] = dxlib.LoadSoundMem(basePath + "cursor_move.mp3")
	soundEffects[SEMenuEnter] = dxlib.LoadSoundMem(basePath + "menu_enter.mp3")
	soundEffects[SEDenied] = dxlib.LoadSoundMem(basePath + "denied.mp3")
	soundEffects[SECancel] = dxlib.LoadSoundMem(basePath + "cancel.mp3")
	soundEffects[SEGoBattle] = dxlib.LoadSoundMem(basePath + "go_battle.mp3")
	soundEffects[SEEnemyAppear] = dxlib.LoadSoundMem(basePath + "enemy_appear.mp3")
	soundEffects[SEChipSelectOpen] = dxlib.LoadSoundMem(basePath + "chip_select_open.mp3")
	soundEffects[SESelect] = dxlib.LoadSoundMem(basePath + "select.mp3")
	soundEffects[SEChipSelectEnd] = dxlib.LoadSoundMem(basePath + "chip_select_end.mp3")
	soundEffects[SEGaugeMax] = dxlib.LoadSoundMem(basePath + "gauge_max.mp3")
	soundEffects[SECannon] = dxlib.LoadSoundMem(basePath + "cannon.mp3")
	soundEffects[SEBusterCharging] = dxlib.LoadSoundMem(basePath + "buster_charging.mp3")
	soundEffects[SEBusterCharged] = dxlib.LoadSoundMem(basePath + "buster_charged.mp3")
	soundEffects[SEBusterShot] = dxlib.LoadSoundMem(basePath + "buster_shot.wav")
	soundEffects[SECannonHit] = dxlib.LoadSoundMem(basePath + "cannon_hit.mp3")
	soundEffects[SEExplode] = dxlib.LoadSoundMem(basePath + "bomb_explode.mp3")
	soundEffects[SEBusterHit] = dxlib.LoadSoundMem(basePath + "shot_hit.wav")
	soundEffects[SESword] = dxlib.LoadSoundMem(basePath + "sword.mp3")
	soundEffects[SERecover] = dxlib.LoadSoundMem(basePath + "recover.mp3")
	soundEffects[SEShockWave] = dxlib.LoadSoundMem(basePath + "shock_wave.mp3")
	soundEffects[SEGun] = dxlib.LoadSoundMem(basePath + "gun.mp3")
	soundEffects[SESpreadHit] = dxlib.LoadSoundMem(basePath + "shot_hit.wav")
	soundEffects[SEBombThrow] = dxlib.LoadSoundMem(basePath + "bomb_throw.mp3")
	soundEffects[SEPlayerDeleted] = dxlib.LoadSoundMem(basePath + "player_deleted.mp3")
	soundEffects[SEDamaged] = dxlib.LoadSoundMem(basePath + "damaged.mp3")
	soundEffects[SEEnemyDeleted] = dxlib.LoadSoundMem(basePath + "enemy_deleted.mp3")
	soundEffects[SEGotItem] = dxlib.LoadSoundMem(basePath + "got_item.mp3")
	soundEffects[SEWindowChange] = dxlib.LoadSoundMem(basePath + "window_change.mp3")
	soundEffects[SEThunderBall] = dxlib.LoadSoundMem(basePath + "thunder_ball.mp3")
	soundEffects[SEWideShot] = dxlib.LoadSoundMem(basePath + "wide_shot.mp3")
	soundEffects[SEBoomerangThrow] = dxlib.LoadSoundMem(basePath + "boomerang_throw.mp3")
	soundEffects[SEWaterLanding] = dxlib.LoadSoundMem(basePath + "water_landing.mp3")
	soundEffects[SEBlock] = dxlib.LoadSoundMem(basePath + "block.mp3")
	soundEffects[SEObjectCreate] = dxlib.LoadSoundMem(basePath + "object_create.wav")
	soundEffects[SEWaterpipeAttack] = dxlib.LoadSoundMem(basePath + "waterpipe_attack.mp3")
	soundEffects[SEPanelBreak] = dxlib.LoadSoundMem(basePath + "panel_break.mp3")

	for i, s := range soundEffects {
		if s == -1 {
			return fmt.Errorf("failed to load %d sound", i)
		}
	}

	dxlib.ChangeVolumeSoundMem(96, soundEffects[SEBusterShot])
	dxlib.ChangeVolumeSoundMem(128, soundEffects[SECannonHit])
	dxlib.ChangeVolumeSoundMem(96, soundEffects[SEBusterHit])
	dxlib.ChangeVolumeSoundMem(128, soundEffects[SEBusterCharging])
	dxlib.ChangeVolumeSoundMem(128, soundEffects[SESpreadHit])
	dxlib.ChangeVolumeSoundMem(96, soundEffects[SEEnemyDeleted])
	dxlib.ChangeVolumeSoundMem(192, soundEffects[SEPlayerDeleted])
	dxlib.ChangeVolumeSoundMem(128, soundEffects[SEPanelBreak])

	return nil
}

func On(typ SEType) {
	if typ == SENone {
		return
	}

	if dxlib.CheckSoundMem(soundEffects[typ]) == 1 {
		if typ != SECannonHit {
			dxlib.StopSoundMem(soundEffects[typ])
		}
	}
	dxlib.PlaySoundMem(soundEffects[typ], dxlib.DX_PLAYTYPE_BACK, dxlib.TRUE)
}
