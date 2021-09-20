package session

import (
	"time"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/object"
)

func updateObject(objects *[]object.Object, obj object.Object, clientID string, myObject bool) {
	obj.Count = 0
	obj.ClientID = clientID
	if obj.UpdateBaseTime {
		obj.BaseTime = time.Now()
	}

	if !myObject {
		obj.X = config.FieldNumX - obj.X - 1
		obj.PrevX = config.FieldNumX - obj.PrevX - 1
		obj.TargetX = config.FieldNumX - obj.TargetX - 1
	}

	updated := false
	for i, o := range *objects {
		if o.ID == obj.ID {
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

func removeObject(objects *[]object.Object, objID string) {
	newObjs := []object.Object{}
	for _, obj := range *objects {
		if obj.ID != objID {
			newObjs = append(newObjs, obj)
		}
	}
	*objects = newObjs
}
