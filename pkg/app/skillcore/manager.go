package skillcore

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type Manager struct {
	damageMgr    *damage.DamageManager
	GetObjectPos func(objID string) point.Point
}

func NewManager(damageMgr *damage.DamageManager, GetObjectPos func(objID string) point.Point) *Manager {
	return &Manager{
		damageMgr:    damageMgr,
		GetObjectPos: GetObjectPos,
	}
}

func (m *Manager) Get(id int, arg Argument) SkillCore {
	// switch id {
	// }

	// TODO: 不正なIDの場合はエラーをセットしたいが、現状実装途中なので呼び出し元で参照しないようにする
	return nil
}
