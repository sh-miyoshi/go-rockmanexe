package damage

import (
	"sync"
)

type Manager struct {
	dmLock  sync.Mutex
	damages []Damage
}

func (m *Manager) Hit(objX, objY int) *Damage {
	for _, dm := range m.damages {
		if dm.PosX == objX && dm.PosY == objY {
			return &dm
		}
	}
	return nil
}

func (m *Manager) Add(dm []Damage) {
	m.dmLock.Lock()
	defer m.dmLock.Unlock()
	m.damages = append(m.damages, dm...)
}

/*
	TODO
	MgrProcは16msごとに処理をしたい
	Sessionごとに必要
	field.ObjectsはSessionに持っている
*/
