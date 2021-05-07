package game

import (
	"errors"
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/menu"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/player"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/title"
)

const (
	stateTitle int = iota
	stateBattle
	stateMenu

	stateMax
)

var (
	state           = stateTitle
	count      uint = 0
	playerInfo *player.Player
)

// Process ...
func Process() error {
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
				// TODO implement
				return fmt.Errorf("start with continue is not implemented yet")
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
				if common.EncryptKey == "" {
					// Save without encryption(debug mode)
					playerInfo.Save(common.SaveFilePath, nil)
				} else {
					key := []byte(common.EncryptKey)
					playerInfo.Save(common.SaveFilePath, key)
				}
				stateChange(stateMenu)
				return nil
			} else if errors.Is(err, battle.ErrLose) {
				stateChange(stateMenu)
				return nil
			}

			return fmt.Errorf("battle process failed: % w", err)
		}
	}
	count++
	return nil
}

// Draw ...
func Draw() {
	if count == 0 {
		// skip if initialize phase
		return
	}

	switch state {
	case stateTitle:
		title.Draw()
	case stateMenu:
		menu.Draw()
	case stateBattle:
		battle.Draw()
	}
}

func stateChange(nextState int) {
	if nextState < 0 || nextState >= stateMax {
		panic(fmt.Sprintf("Invalid next game state: %d", nextState))
	}
	state = nextState
	count = 0
}
