package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
)

/*
暗転
自分のいる場所にアクアマンを召喚
	ロックマンを消してアクアマンを表示
土管を表示
2マス前までダメージ
	2マス前は2ヒット
暗転解除
*/

type aquaman struct {
	ID         string
	OwnerID    string
	Power      uint
	TargetType int

	count int
}

func newAquaman(objID string, arg Argument) *aquaman {
	return &aquaman{}
}

func (p *aquaman) Draw() {
}

func (p *aquaman) Process() (bool, error) {
	p.count++

	return false, nil
}

func (p *aquaman) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		AnimType: anim.AnimTypeSkill,
	}
}
