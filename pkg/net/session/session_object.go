package session

import (
	"time"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/field"
)

func updateObject(objects *[]field.Object, obj field.Object, clientID string, myObject bool) {
	obj.ClientID = clientID
	if obj.UpdateBaseTime {
		obj.BaseTime = time.Now()
	}

	if !myObject {
		obj.X = config.FieldNumX - obj.X - 1
	}

	updated := false
	for i, o := range *objects {
		if o.ID == obj.ID {
			if !obj.DamageChecked {
				obj.HitDamage = o.HitDamage
			}
			if !obj.UpdateBaseTime {
				obj.BaseTime = o.BaseTime
			}
			obj.UpdateBaseTime = false
			(*objects)[i] = obj
			updated = true
			break
		}
	}

	if !updated {
		*objects = append(*objects, obj)
	}
}

func removeObject(objects *[]field.Object, objID string) {
	newObjs := []field.Object{}
	for _, obj := range *objects {
		if obj.ID != objID {
			newObjs = append(newObjs, obj)
		}
	}
	*objects = newObjs
}
