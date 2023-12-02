package scenario

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/event"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/mapinfo"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

var eno0Args_犬小屋 = event.MapChangeArgs{
	MapID:   mapinfo.ID_秋原町,
	InitPos: point.Point{X: 1400, Y: 500},
}

var Scenario_犬小屋 = [][]event.Scenario{
	{ // EventNo: 0
		{Type: event.TypeChangeMapArea, Values: eno0Args_犬小屋.Marshal()},
	},
}

var eno0Args_秋原町 = event.MapChangeArgs{
	MapID:   mapinfo.ID_犬小屋,
	InitPos: point.Point{X: 300, Y: 200},
}

var Scenario_秋原町 = [][]event.Scenario{
	{ // EventNo: 0
		{Type: event.TypeMessage, Values: []byte("プラグイン！ロックマン．ｅｘｅトランスミッション！")},
		{Type: event.TypeChangeMapArea, Values: eno0Args_秋原町.Marshal()},
	},
}
