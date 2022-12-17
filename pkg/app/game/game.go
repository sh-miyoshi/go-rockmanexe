package game

import (
	"errors"
	"fmt"
	"time"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/background"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/event"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/mapmove"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/menu"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/scratch"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/title"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/netconn"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/player"
)

const (
	stateTitle int = iota
	stateBattle
	stateNetBattle
	stateMenu
	stateMap
	stateScratch
	stateEvent

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
		if err := menu.Process(); err != nil {
			menu.End()
			if errors.Is(err, menu.ErrGoBattle) {
				stateChange(stateBattle)
				return nil
			} else if errors.Is(err, menu.ErrGoNetBattle) {
				stateChange(stateNetBattle)
				return nil
			} else if errors.Is(err, menu.ErrGoMap) {
				stateChange(stateMap)
				return nil
			} else if errors.Is(err, menu.ErrGoScratch) {
				stateChange(stateScratch)
				return nil
			}
			return fmt.Errorf("game process in state menu failed: %w", err)
		}
	case stateBattle:
		if count == 0 {
			if err := battle.Init(playerInfo, menu.GetBattleEnemies()); err != nil {
				return fmt.Errorf("game process in state battle failed: %w", err)
			}
		}
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
				OpponentID: netconn.GetInst().GetOpponentUserID(),
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
			mapmove.End()
			if errors.Is(err, mapmove.ErrGoBattle) {
				stateChange(stateBattle)
				return nil
			} else if errors.Is(err, mapmove.ErrGoMenu) {
				stateChange(stateMenu)
				return nil
			} else if errors.Is(err, mapmove.ErrGoEvent) {
				stateChange(stateEvent)
				return nil
			}
			return fmt.Errorf("map move process failed: %w", err)
		}
	case stateScratch:
		if count == 0 {
			scratch.Init()
		}
		scratch.Process()
	case stateEvent:
		event.Process()
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
	case stateScratch:
		scratch.Draw()
	case stateEvent:
		event.Draw()
	}
}

func stateChange(nextState int) {
	if nextState < 0 || nextState >= stateMax {
		panic(fmt.Sprintf("Invalid next game state: %d", nextState))
	}
	state = nextState
	count = 0
}
