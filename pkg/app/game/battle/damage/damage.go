package damage

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
)

const (
	TargetPlayer int = 1 << iota
	TargetEnemy
)

const (
	ElementNone int = iota
	ElementFire
	ElementWater
	ElementElec
	ElementWood
)

const (
	TypePosition int = iota
	TypeObject
)

type Damage struct {
	OwnerClientID string // ネット対戦時のDamageを発生させたOwner
	ID            string // Damage ID
	Power         int
	PushRight     int  // ヒット時に右に押されるカウント
	PushLeft      int  // ヒット時に左に押されるカウント
	HitEffectType int  // ヒット時に表示されるEffect
	BigDamage     bool // trueならのけぞる
	Element       int

	DamageType int // ダメージの種類

	// DamageTypeがTypePositionの時使うパラメータ
	Pos           common.Point // (TypePosition)発生箇所
	TTL           int          // (TypePosition)ダメージが残り続ける時間
	ShowHitArea   bool         // (TypePosition)足元にダメージ箇所を表示するか
	TargetObjType int          // (TypePosition)ダメージを受けるObjectのタイプ

	// DamageTypeがTypeObjectの時使うパラメータ
	TargetObjID string // (TypeObject)ダメージを受けるObjectのID

	// TODO: インビジ貫通
}

type DamageManager struct {
	damages map[string]*Damage
}

func IsWeakness(charType int, dm Damage) bool {
	switch charType {
	case ElementFire:
		return dm.Element == ElementWater
	case ElementWater:
		return dm.Element == ElementElec
	case ElementElec:
		return dm.Element == ElementWood
	case ElementWood:
		return dm.Element == ElementFire
	}
	return false
}

func NewManager() *DamageManager {
	return &DamageManager{
		damages: make(map[string]*Damage),
	}
}

func (m *DamageManager) New(dm Damage) string {
	dm.ID = uuid.New().String()
	m.damages[dm.ID] = &dm
	logger.Debug("Add damage: %+v to damage manager", dm)
	return dm.ID
}

func (m *DamageManager) Process() {
	for id, d := range m.damages {
		if d.DamageType == TypeObject {
			delete(m.damages, id)
		}

		if d.DamageType == TypePosition {
			d.TTL--
			if d.TTL <= 0 {
				delete(m.damages, id)
			}
		}
	}
}

func (m *DamageManager) GetHitDamage(pos common.Point, objID string) *Damage {
	for _, d := range m.damages {
		switch d.DamageType {
		case TypeObject:
			if d.TargetObjID == objID {
				return d
			}
		case TypePosition:
			if d.Pos.Equal(pos) {
				return d
			}
		}
	}
	return nil
}

func (m *DamageManager) Exists(id string) bool {
	_, ok := m.damages[id]
	return ok
}

func (m *DamageManager) Remove(id string) {
	delete(m.damages, id)
}

func (m *DamageManager) RemoveAll() {
	m.damages = make(map[string]*Damage)
}

func (m *DamageManager) PrintAllData() {
	msg := "current damages: { "
	for id, d := range m.damages {
		msg += fmt.Sprintf("%s: %+v, ", id, *d)
	}
	msg += " }"
	logger.Debug(msg)
}
