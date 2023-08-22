package game

import (
	"errors"
	"fmt"
	"time"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
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
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/title"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/player"
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

	stateMax
)

var (
	state           = stateTitle
	count      uint = 0
	playerInfo *player.Player
)

// Process ...
func Process() error {
	background.Process()
	fade.Process()

	if playerInfo != nil {
		playerInfo.PlayCount++
		// Countermeasures against buffer overflow
		if playerInfo.PlayCount >= common.MaxUint {
			playerInfo.PlayCount = common.MaxUint - 1
		}
	}

	switch state {
	case stateTitle:
		if count == 0 {
			if err := title.Init(); err != nil {
				return fmt.Errorf("game process in state title failed: %w", err)
			}
		}
		if err := title.Process(); err != nil {
			if errors.Is(err, title.ErrStartInit) {
				playerInfo = player.New()
			} else if errors.Is(err, title.ErrStartContinue) {
				var key []byte
				if common.EncryptKey == "" {
					key = nil
				} else {
					key = []byte(common.EncryptKey)
				}

				var err error
				playerInfo, err = player.NewWithSaveData(common.SaveFilePath, key)
				if err != nil {
					return fmt.Errorf("failed to continue: %w", err)
				}
			} else {
				return fmt.Errorf("failed to process title: %w", err)
			}
			title.End()
			stateChange(stateMenu)
			return nil
		}
	case stateMenu:
		if count == 0 {
			if err := menu.Init(playerInfo); err != nil {
				return fmt.Errorf("game process in state menu init failed: %w", err)
			}
		}
		result, err := menu.Process()
		if err != nil {
			menu.End()
			return fmt.Errorf("game process in state menu failed: %w", err)
		}
		if result != menu.ResultContinue {
			menu.End()
		}
		switch result {
		case menu.ResultGoBattle:
			if err := battle.Init(playerInfo, menu.GetBattleEnemies()); err != nil {
				return fmt.Errorf("battle init failed at menu: %w", err)
			}
			stateChange(stateBattle)
			return nil
		case menu.ResultGoNetBattle:
			stateChange(stateNetBattle)
			return nil
		case menu.ResultGoMap:
			stateChange(stateMap)
			return nil
		case menu.ResultGoScratch:
			stateChange(stateScratch)
			return nil
		case menu.ResultGoNaviCustom:
			stateChange(stateNaviCustom)
			return nil
		}
	case stateBattle:
		if err := battle.Process(); err != nil {
			battle.End()
			if errors.Is(err, battle.ErrWin) {
				playerInfo.WinNum++
				key := []byte(common.EncryptKey)
				if err := playerInfo.Save(common.SaveFilePath, key); err != nil {
					return fmt.Errorf("save failed: %w", err)
				}
				stateChange(stateMenu)
				return nil
			} else if errors.Is(err, battle.ErrLose) {
				stateChange(stateMenu)
				return nil
			}

			return fmt.Errorf("battle process failed: % w", err)
		}
	case stateNetBattle:
		if count == 0 {
			if err := netbattle.Init(playerInfo); err != nil {
				return fmt.Errorf("game process in state net battle failed: %w", err)
			}
		}

		if err := netbattle.Process(); err != nil {
			netbattle.End()

			history := player.History{
				Date:       time.Now(),
				OpponentID: net.GetInst().GetOpponentUserID(),
			}

			if errors.Is(err, battle.ErrWin) {
				history.IsWin = true
				playerInfo.BattleHistories = append(playerInfo.BattleHistories, history)
				key := []byte(common.EncryptKey)
				if err := playerInfo.Save(common.SaveFilePath, key); err != nil {
					return fmt.Errorf("save failed: %w", err)
				}
				stateChange(stateMenu)
				return nil
			} else if errors.Is(err, battle.ErrLose) {
				history.IsWin = false
				playerInfo.BattleHistories = append(playerInfo.BattleHistories, history)
				key := []byte(common.EncryptKey)
				if err := playerInfo.Save(common.SaveFilePath, key); err != nil {
					return fmt.Errorf("save failed: %w", err)
				}
				stateChange(stateMenu)
				return nil
			}

			return fmt.Errorf("battle process failed: %w", err)
		}
	case stateMap:
		if count == 0 {
			if err := mapmove.Init(); err != nil {
				return fmt.Errorf("game process in state map move failed: %w", err)
			}
		}
		if err := mapmove.Process(); err != nil {
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
			return fmt.Errorf("map move process failed: %w", err)
		}
	case stateMapChange:
		mapmove.Init()
		stateChange(stateEvent)
		return nil
	case stateScratch:
		if count == 0 {
			scratch.Init()
		}
		scratch.Process()
	case stateEvent:
		res, err := event.Process()
		if err != nil {
			return fmt.Errorf("event process failed: %w", err)
		}
		switch res {
		case event.ResultMapChange:
			stateChange(stateMapChange)
		case event.ResultEnd:
			stateChange(stateMap)
		}
		return nil
	case stateNaviCustom:
		if count == 0 {
			navicustom.Init()
		}
		navicustom.Process()
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
	}

	fade.Draw()
}

func stateChange(nextState int) {
	if nextState < 0 || nextState >= stateMax {
		common.SetError(fmt.Sprintf("Invalid next game state: %d", nextState))
		return
	}
	state = nextState
	count = 0
}
