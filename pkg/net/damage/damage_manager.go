package damage

type Manager struct {
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
	m.damages = append(m.damages, dm...)
}

func (m *Manager) Update() {
	newDamages := []Damage{}
	for _, dm := range m.damages {
		dm.TTL--
		if dm.TTL > 0 {
			newDamages = append(newDamages, dm)
		}
	}
	m.damages = newDamages
}
