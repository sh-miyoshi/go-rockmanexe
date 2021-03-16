package enemy

import (
	"errors"
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/anim"
)

// EnemyParam ...
type EnemyParam struct {
	ID   string
	PosX int
	PosY int
	HP   int
}

type enemyObject interface {
	anim.Anim
	Get() *EnemyParam
	Init() error
	End()
}

var (
	ErrGameEnd = errors.New("game end")

	enemies []enemyObject
)

func Init() error {
	// Decide enemies
	// debug(set debug param)
	enemies = append(enemies, getObject(idMetall, EnemyParam{
		PosX: 4,
		PosY: 1,
		HP:   40,
	}))

	// Init enemy data
	for _, e := range enemies {
		if err := e.Init(); err != nil {
			return err
		}
	}
	return nil
}

func End() {
	// Cleanup existsting enemy data
	for _, e := range enemies {
		e.End()
	}
}

func MgrProcess() error {
	newEnemies := []enemyObject{}
	for _, e := range enemies {
		end, err := e.Process()
		if err != nil {
			return fmt.Errorf("Anim process failed: %w", err)
		}

		if end {
			e.End()
			continue
		}
		newEnemies = append(newEnemies, e)
	}

	if len(newEnemies) == 0 {
		return ErrGameEnd
	}
	enemies = newEnemies

	return nil
}

func MgrDraw() {
	for _, e := range enemies {
		e.Draw()
	}
}

func GetEnemies() []*EnemyParam {
	res := []*EnemyParam{}
	for _, e := range enemies {
		res = append(res, e.Get())
	}
	return res
}
