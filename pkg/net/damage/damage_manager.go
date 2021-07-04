package damage

import "github.com/sh-miyoshi/go-rockmanexe/pkg/net/config"

type Manager struct {
	damages []Damage
}

func (m *Manager) Hit(ownerClientID string, objClientID string, objX, objY int) *Damage {
	for _, dm := range m.damages {
		x := objX
		if ownerClientID != dm.ClientID {
			x = config.FieldNumX - objX - 1
		}

		if dm.PosX == x && dm.PosY == objY {
			if dm.TargetType == TargetOwn && dm.ClientID != objClientID {
				continue
			}
			if dm.TargetType == TargetOtherClient && dm.ClientID == objClientID {
				continue
			}

			return &dm
		}
	}
	return nil
}

func (m *Manager) Add(dm []Damage) {
	m.damages = append(m.damages, dm...)
}

func (m *Manager) Update() {
	// TODO remove hit damage data

	newDamages := []Damage{}
	for _, dm := range m.damages {
		dm.TTL--
		if dm.TTL > 0 {
			newDamages = append(newDamages, dm)
		}
	}
	m.damages = newDamages
}
