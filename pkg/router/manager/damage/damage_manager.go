package damage

import (
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
)

type Manager struct {
	dmMgr *damage.DamageManager
}

func New(dmMgr *damage.DamageManager) *Manager {
	return &Manager{
		dmMgr: dmMgr,
	}
}

func (m *Manager) New(dm damage.Damage) string {
	if dm.TargetObjType == damage.TargetEnemy {
		// ダメージでは反転させる
		dm.Pos.X = battlecommon.FieldNum.X - dm.Pos.X - 1
	}
	if dm.OwnerClientID == "" {
		// 本来はプログラムバグなのでエラーだがエラーハンドリングがめんどいので一旦ログのみにする
		logger.Error("network damage requires owner client id")
	}

	return m.dmMgr.New(dm)
}

func (m *Manager) Exists(id string) bool {
	return m.dmMgr.Exists(id)
}
