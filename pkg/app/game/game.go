package game

import (
	"fmt"
	"time"

	"github.com/cockroachdb/errors"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/fade"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/background"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/event"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/mapmove"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/menu"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/navicustom"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/net"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/scratch"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/talkai"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/title"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/mapinfo"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/player"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/system"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	stateTitle int = iota
	stateBattle
	stateNetBattle
	stateMenu
	stateMap
	stateMapChange
	stateScratch
	stateEvent
	stateNaviCustom
	stateTalkAI

	stateMax
)

var (
	state           = stateTitle
	count      uint = 0
	playerInfo *player.Player
)

func Update() error {
	background.Update()
	fade.Update()

	if playerInfo != nil {
		playerInfo.PlayCount++
		// Countermeasures against buffer overflow
		if playerInfo.PlayCount >= config.MaxUint {
			playerInfo.PlayCount = config.MaxUint - 1
		}
	}

	switch state {
	case stateTitle:
		if count == 0 {
			if err := title.Init(); err != nil {
				return errors.Wrap(err, "game process in state title failed")
			}
		}
		if err := title.Update(); err != nil {
			if errors.Is(err, title.ErrStartInit) {
				playerInfo = player.New()
			} else if errors.Is(err, title.ErrStartContinue) {
				var key []byte
				if config.EncryptKey == "" {
					key = nil
				} else {
					key = []byte(config.EncryptKey)
				}

				var err error
				playerInfo, err = player.NewWithSaveData(config.SaveFilePath, key)
				if err != nil {
					return errors.Wrap(err, "failed to continue")
				}
			} else {
				return errors.Wrap(err, "failed to process title")
			}
			title.End()
			stateChange(stateMenu)
			return nil
		}
	case stateMenu:
		if count == 0 {
			if err := menu.Init(playerInfo); err != nil {
				return errors.Wrap(err, "game process in state menu init failed")
			}
		}
		result, err := menu.Update()
		if err != nil {
			menu.End()
			return errors.Wrap(err, "game process in state menu failed")
		}
		if result != menu.ResultContinue {
			menu.End()
		}
		switch result {
		case menu.ResultGoBattle:
			if err := battle.Init(playerInfo, menu.GetBattleEnemies()); err != nil {
				return errors.Wrap(err, "battle init failed at menu")
			}
			stateChange(stateBattle)
			return nil
		case menu.ResultGoNetBattle:
			stateChange(stateNetBattle)
			return nil
		case menu.ResultGoMap:
			// debug: 初期イベントをセット
			args := event.MapChangeArgs{MapID: mapinfo.ID_秋原町, InitPos: point.Point{X: 1400, Y: 500}}
			event.SetScenarios([]event.Scenario{
				{Type: event.TypeChangeMapArea, Values: args.Marshal()},
			})
			stateChange(stateEvent)
			return nil
		case menu.ResultGoScratch:
			stateChange(stateScratch)
			return nil
		case menu.ResultGoNaviCustom:
			stateChange(stateNaviCustom)
			return nil
		case menu.ResultGoTalkAI:
			stateChange(stateTalkAI)
			return nil
		}
	case stateBattle:
		if err := battle.Update(); err != nil {
			battle.End()
			if errors.Is(err, battle.ErrWin) {
				playerInfo.WinNum++
				key := []byte(config.EncryptKey)
				if err := playerInfo.Save(config.SaveFilePath, key); err != nil {
					return errors.Wrap(err, "save failed")
				}
				stateChange(stateMenu)
				return nil
			} else if errors.Is(err, battle.ErrLose) {
				stateChange(stateMenu)
				return nil
			}

			return errors.Wrap(err, "battle process failed")
		}
	case stateNetBattle:
		if count == 0 {
			if err := netbattle.Init(playerInfo); err != nil {
				return errors.Wrap(err, "game process in state net battle failed")
			}
		}

		if err := netbattle.Update(); err != nil {
			netbattle.End()

			history := player.History{
				Date:       time.Now(),
				OpponentID: net.GetInst().GetOpponentUserID(),
			}

			if errors.Is(err, battle.ErrWin) {
				history.IsWin = true
				playerInfo.BattleHistories = append(playerInfo.BattleHistories, history)
				key := []byte(config.EncryptKey)
				if err := playerInfo.Save(config.SaveFilePath, key); err != nil {
					return errors.Wrap(err, "save failed")
				}
				stateChange(stateMenu)
				return nil
			} else if errors.Is(err, battle.ErrLose) {
				history.IsWin = false
				playerInfo.BattleHistories = append(playerInfo.BattleHistories, history)
				key := []byte(config.EncryptKey)
				if err := playerInfo.Save(config.SaveFilePath, key); err != nil {
					return errors.Wrap(err, "save failed")
				}
				stateChange(stateMenu)
				return nil
			}

			return errors.Wrap(err, "battle process failed")
		}
	case stateMap:
		if count == 0 {
			if err := mapmove.Init(); err != nil {
				return errors.Wrap(err, "game process in state map move failed")
			}
		}
		if err := mapmove.Update(); err != nil {
			if errors.Is(err, mapmove.ErrGoBattle) {
				mapmove.End()
				stateChange(stateBattle)
				return nil
			} else if errors.Is(err, mapmove.ErrGoMenu) {
				mapmove.End()
				stateChange(stateMenu)
				return nil
			} else if errors.Is(err, mapmove.ErrGoEvent) {
				stateChange(stateEvent)
				return nil
			}
			return errors.Wrap(err, "map move process failed")
		}
	case stateMapChange:
		var args event.MapChangeArgs
		args.Unmarshal(event.GetStoredValues())
		mapmove.MapChange(args.MapID, args.InitPos)
		stateChange(stateEvent)
		return nil
	case stateScratch:
		if count == 0 {
			scratch.Init()
		}
		scratch.Update()
	case stateEvent:
		res, err := event.Update()
		if err != nil {
			return errors.Wrap(err, "event process failed")
		}
		switch res {
		case event.ResultMapChange:
			stateChange(stateMapChange)
			return nil
		case event.ResultEnd:
			stateChange(stateMap)
			return nil
		}
	case stateNaviCustom:
		if count == 0 {
			navicustom.Init(playerInfo)
		}
		if navicustom.Update() {
			navicustom.End()
			stateChange(stateMenu)
			return nil
		}
	case stateTalkAI:
		if count == 0 {
			talkai.Init()
		}
		if talkai.Update() {
			talkai.End()
			stateChange(stateMenu)
			return nil
		}
	}
	count++
	return nil
}

func Draw() {
	if count == 0 {
		// skip if initialize phase
		return
	}

	background.Draw()

	switch state {
	case stateTitle:
		title.Draw()
	case stateMenu:
		menu.Draw()
	case stateBattle:
		battle.Draw()
	case stateNetBattle:
		netbattle.Draw()
	case stateMap:
		mapmove.Draw()
	case stateMapChange:
		mapmove.Draw()
	case stateScratch:
		scratch.Draw()
	case stateEvent:
		mapmove.Draw()
		event.Draw()
	case stateNaviCustom:
		navicustom.Draw()
	case stateTalkAI:
		talkai.Draw()
	}

	fade.Draw()
}

func stateChange(nextState int) {
	if nextState < 0 || nextState >= stateMax {
		system.SetError(fmt.Sprintf("Invalid next game state: %d", nextState))
		return
	}
	logger.Info("game state change from %d to %d", state, nextState)
	state = nextState
	count = 0
}
