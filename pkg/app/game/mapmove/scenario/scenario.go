package scenario

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/event"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/mapinfo"
)

var Scenario_犬小屋 = [][]event.Scenario{
	{ // EventNo: 0
		{Type: event.TypeChangeMapArea},
	},
}

var eno0Args = event.MapChangeArgs{
	MapID:   mapinfo.ID_犬小屋,
	InitPos: common.Point{X: 300, Y: 200},
}

var Scenario_秋原町 = [][]event.Scenario{
	{ // EventNo: 0
		{Type: event.TypeChangeMapArea, Values: eno0Args.Marshal()},
	},
}
