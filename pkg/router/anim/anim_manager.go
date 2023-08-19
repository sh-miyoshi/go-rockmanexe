package anim

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
)

var (
	animManagers     = make(map[string]*anim.AnimManager)
	objanimManagers  = make(map[string]*objanim.AnimManager)
	clientAnimMgrMap = make(map[string]string) // Key: clientID, Value: animManagerID
)

func NewManager(clientIDs [2]string) string {
	id := uuid.New().String()
	animManagers[id] = anim.NewManager()
	objanimManagers[id] = objanim.NewManager()
	for i := 0; i < len(clientIDs); i++ {
		clientAnimMgrMap[clientIDs[i]] = id
	}
	return id
}

func Cleanup(mgrID string) {
	animManagers[mgrID].Cleanup()
	objanimManagers[mgrID].Cleanup()
	delete(animManagers, mgrID)
	delete(objanimManagers, mgrID)
}

func MgrProcess(mgrID string) error {
	if mgrID == "" {
		return nil
	}

	if err := animManagers[mgrID].Process(); err != nil {
		return fmt.Errorf("anim manage process failed: %w", err)
	}

	if err := objanimManagers[mgrID].Process(true, false); err != nil {
		return fmt.Errorf("objanim manage process failed: %w", err)
	}

	return nil
}

func AnimNew(clientID string, a anim.Anim) string {
	mgrID := clientAnimMgrMap[clientID]
	return animManagers[mgrID].New(a)
}

func AnimDelete(clientID string, animID string) {
	mgrID := clientAnimMgrMap[clientID]
	animManagers[mgrID].Delete(animID)
}

func AnimGetAll(mgrID string) []anim.Param {
	return animManagers[mgrID].GetAll()
}

func AnimIsProcessing(clientID string, animID string) bool {
	mgrID := clientAnimMgrMap[clientID]
	return animManagers[mgrID].IsProcessing(animID)
}

func ObjAnimNew(clientID string, a objanim.Anim) string {
	mgrID := clientAnimMgrMap[clientID]
	return objanimManagers[mgrID].New(a)
}

func ObjAnimGetObjs(clientID string, filter objanim.Filter) []objanim.Param {
	mgrID := clientAnimMgrMap[clientID]
	return objanimManagers[mgrID].GetObjs(filter)
}

func ObjAnimGetObjPos(clientID string, objID string) common.Point {
	mgrID := clientAnimMgrMap[clientID]
	return objanimManagers[mgrID].GetObjPos(objID)
}

func DamageManager(clientID string) *damage.DamageManager {
	mgrID := clientAnimMgrMap[clientID]
	return objanimManagers[mgrID].DamageManager()
}

func DamageNew(clientID string, dm damage.Damage) string {
	if dm.TargetObjType == damage.TargetEnemy {
		// ダメージでは反転させる
		dm.Pos.X = battlecommon.FieldNum.X - dm.Pos.X - 1
	}
	return DamageManager(clientID).New(dm)
}
