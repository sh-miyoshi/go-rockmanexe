package sound

import (
	"github.com/cockroachdb/errors"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

var (
	soundEffects = [resources.SEMax]int{}
)

func Init() error {
	if config.Get().Sound.SE.Disabled {
		return nil
	}

	basePath := config.SoundPath + "se/"

	soundEffects[resources.SENone] = 0
	soundEffects[resources.SETitleEnter] = dxlib.LoadSoundMem(basePath + "title_enter.mp3")
	soundEffects[resources.SECursorMove] = dxlib.LoadSoundMem(basePath + "cursor_move.mp3")
	soundEffects[resources.SEMenuEnter] = dxlib.LoadSoundMem(basePath + "menu_enter.mp3")
	soundEffects[resources.SEDenied] = dxlib.LoadSoundMem(basePath + "denied.mp3")
	soundEffects[resources.SECancel] = dxlib.LoadSoundMem(basePath + "cancel.mp3")
	soundEffects[resources.SEGoBattle] = dxlib.LoadSoundMem(basePath + "go_battle.mp3")
	soundEffects[resources.SEEnemyAppear] = dxlib.LoadSoundMem(basePath + "enemy_appear.mp3")
	soundEffects[resources.SEChipSelectOpen] = dxlib.LoadSoundMem(basePath + "chip_select_open.mp3")
	soundEffects[resources.SESelect] = dxlib.LoadSoundMem(basePath + "select.mp3")
	soundEffects[resources.SEChipSelectEnd] = dxlib.LoadSoundMem(basePath + "chip_select_end.mp3")
	soundEffects[resources.SEGaugeMax] = dxlib.LoadSoundMem(basePath + "gauge_max.mp3")
	soundEffects[resources.SECannon] = dxlib.LoadSoundMem(basePath + "cannon.mp3")
	soundEffects[resources.SEBusterCharging] = dxlib.LoadSoundMem(basePath + "buster_charging.mp3")
	soundEffects[resources.SEBusterCharged] = dxlib.LoadSoundMem(basePath + "buster_charged.mp3")
	soundEffects[resources.SEBusterShot] = dxlib.LoadSoundMem(basePath + "buster_shot.wav")
	soundEffects[resources.SECannonHit] = dxlib.LoadSoundMem(basePath + "cannon_hit.mp3")
	soundEffects[resources.SEExplode] = dxlib.LoadSoundMem(basePath + "bomb_explode.mp3")
	soundEffects[resources.SEBusterHit] = dxlib.LoadSoundMem(basePath + "shot_hit.wav")
	soundEffects[resources.SESword] = dxlib.LoadSoundMem(basePath + "sword.mp3")
	soundEffects[resources.SERecover] = dxlib.LoadSoundMem(basePath + "recover.mp3")
	soundEffects[resources.SEShockWave] = dxlib.LoadSoundMem(basePath + "shock_wave.mp3")
	soundEffects[resources.SEGun] = dxlib.LoadSoundMem(basePath + "gun.mp3")
	soundEffects[resources.SESpreadHit] = dxlib.LoadSoundMem(basePath + "shot_hit.wav")
	soundEffects[resources.SEBombThrow] = dxlib.LoadSoundMem(basePath + "bomb_throw.mp3")
	soundEffects[resources.SEPlayerDeleted] = dxlib.LoadSoundMem(basePath + "player_deleted.mp3")
	soundEffects[resources.SEDamaged] = dxlib.LoadSoundMem(basePath + "damaged.mp3")
	soundEffects[resources.SEEnemyDeleted] = dxlib.LoadSoundMem(basePath + "enemy_deleted.mp3")
	soundEffects[resources.SEGotItem] = dxlib.LoadSoundMem(basePath + "got_item.mp3")
	soundEffects[resources.SEWindowChange] = dxlib.LoadSoundMem(basePath + "window_change.mp3")
	soundEffects[resources.SEThunderBall] = dxlib.LoadSoundMem(basePath + "thunder_ball.mp3")
	soundEffects[resources.SEWideShot] = dxlib.LoadSoundMem(basePath + "wide_shot.mp3")
	soundEffects[resources.SEBoomerangThrow] = dxlib.LoadSoundMem(basePath + "boomerang_throw.mp3")
	soundEffects[resources.SEWaterLanding] = dxlib.LoadSoundMem(basePath + "water_landing.mp3")
	soundEffects[resources.SEBlock] = dxlib.LoadSoundMem(basePath + "block.mp3")
	soundEffects[resources.SEObjectCreate] = dxlib.LoadSoundMem(basePath + "object_create.wav")
	soundEffects[resources.SEWaterpipeAttack] = dxlib.LoadSoundMem(basePath + "waterpipe_attack.mp3")
	soundEffects[resources.SEPanelBreak] = dxlib.LoadSoundMem(basePath + "panel_break.mp3")
	soundEffects[resources.SEPanelBreakShort] = dxlib.LoadSoundMem(basePath + "panel_break_short.mp3")
	soundEffects[resources.SEPAPrepare] = dxlib.LoadSoundMem(basePath + "pa_prepare.mp3")
	soundEffects[resources.SEPACreated] = dxlib.LoadSoundMem(basePath + "pa_created.mp3")
	soundEffects[resources.SEDreamSword] = dxlib.LoadSoundMem(basePath + "dream_sword.mp3")
	soundEffects[resources.SEFlameAttack] = dxlib.LoadSoundMem(basePath + "flame_attack.wav")
	soundEffects[resources.SEAreaSteal] = dxlib.LoadSoundMem(basePath + "area_steal.mp3")
	soundEffects[resources.SEAreaStealHit] = dxlib.LoadSoundMem(basePath + "area_steal_hit.mp3")
	soundEffects[resources.SERunOK] = dxlib.LoadSoundMem(basePath + "run_ok.mp3")
	soundEffects[resources.SERunFailed] = dxlib.LoadSoundMem(basePath + "run_failed.mp3")
	soundEffects[resources.SECountBombCountdown] = dxlib.LoadSoundMem(basePath + "count_bomb_countdown.mp3")
	soundEffects[resources.SECountBombEnd] = dxlib.LoadSoundMem(basePath + "count_bomb_end.mp3")
	soundEffects[resources.SETornado] = dxlib.LoadSoundMem(basePath + "tornado.mp3")
	soundEffects[resources.SEFailed] = dxlib.LoadSoundMem(basePath + "failed.mp3")
	soundEffects[resources.SEBubbleShot] = dxlib.LoadSoundMem(basePath + "bubble_shot.mp3")
	soundEffects[resources.SEChing] = dxlib.LoadSoundMem(basePath + "ching.mp3")
	soundEffects[resources.SEDeltaRayEdgeEnd] = dxlib.LoadSoundMem(basePath + "delta_ray_edge_end.mp3")
	soundEffects[resources.SEPanelReturn] = dxlib.LoadSoundMem(basePath + "panel_return.mp3")
	soundEffects[resources.SEMakePoison] = dxlib.LoadSoundMem(basePath + "make_poison.mp3")
	soundEffects[resources.SEAirShoot] = dxlib.LoadSoundMem(basePath + "air_shoot.mp3")

	for i, s := range soundEffects {
		if s == -1 {
			return errors.Newf("failed to load %d sound", i)
		}
	}

	dxlib.ChangeVolumeSoundMem(96, soundEffects[resources.SEBusterShot])
	dxlib.ChangeVolumeSoundMem(128, soundEffects[resources.SECannonHit])
	dxlib.ChangeVolumeSoundMem(96, soundEffects[resources.SEBusterHit])
	dxlib.ChangeVolumeSoundMem(128, soundEffects[resources.SEBusterCharging])
	dxlib.ChangeVolumeSoundMem(128, soundEffects[resources.SESpreadHit])
	dxlib.ChangeVolumeSoundMem(96, soundEffects[resources.SEEnemyDeleted])
	dxlib.ChangeVolumeSoundMem(192, soundEffects[resources.SEPlayerDeleted])
	dxlib.ChangeVolumeSoundMem(128, soundEffects[resources.SEPanelBreak])
	dxlib.ChangeVolumeSoundMem(128, soundEffects[resources.SEDeltaRayEdgeEnd])
	dxlib.ChangeVolumeSoundMem(128, soundEffects[resources.SEPanelReturn])
	dxlib.ChangeVolumeSoundMem(128, soundEffects[resources.SEMakePoison])
	dxlib.ChangeVolumeSoundMem(192, soundEffects[resources.SEAirShoot])

	return nil
}

func On(typ resources.SEType) {
	if config.Get().Sound.SE.Disabled {
		return
	}

	if typ == resources.SENone {
		return
	}

	if dxlib.CheckSoundMem(soundEffects[typ]) == 1 {
		if typ != resources.SECannonHit {
			dxlib.StopSoundMem(soundEffects[typ])
		}
	}
	dxlib.PlaySoundMem(soundEffects[typ], dxlib.DX_PLAYTYPE_BACK, true)
}

func SEClear() {
	if config.Get().Sound.SE.Disabled {
		return
	}

	for _, s := range soundEffects {
		dxlib.StopSoundMem(s)
	}
}
