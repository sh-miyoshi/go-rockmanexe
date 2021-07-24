package damage

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/config"
)

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

func (m *Manager) Add(dm []Damage) error {
	for _, d := range dm {
		if err := d.Validate(); err != nil {
			return fmt.Errorf("damage validation failed: %w", err)
		}
		m.damages = append(m.damages, d)
	}
	return nil
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

func (m *Manager) GetHitAreas(clientID string) [][2]int {
	res := [][2]int{}

	for _, dm := range m.damages {
		if dm.ShowHitArea {
			if dm.ClientID == clientID {
				res = append(res, [2]int{dm.PosX, dm.PosY})
			} else {
				res = append(res, [2]int{config.FieldNumX - dm.PosX, dm.PosY})
			}
		}
	}

	return res
}
