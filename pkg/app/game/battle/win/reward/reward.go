package reward

import (
	"fmt"
	"math/rand"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/chip"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/fps"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/player"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/system"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
)

const (
	TypeMoney int = iota
	TypeChip
)

type EnemyChipInfo struct {
	CharID        int
	ChipID        int
	Code          string
	RequiredLevel int
	IsOnlyOne     bool
}

type EnemyParam struct {
	CharID int
	IsBoss bool
}

type WinArg struct {
	GameTime        int
	DeletedEnemies  []EnemyParam
	PlayerMoveNum   int
	PlayerDamageNum int
}

type Param struct {
	Type          int
	Name          string
	BustingLevel  int
	DeleteTimeSec int

	Zenny int
	Code  string
}

var (
	gotReward     Param
	enemyChipList = []EnemyChipInfo{}
)

func SetEnemyChipList(list []EnemyChipInfo) {
	enemyChipList = list
}

func SetToPlayer(args WinArg, plyr *player.Player) {
	deleteTimeSec := args.GameTime / int(fps.FPS)
	if deleteTimeSec == 0 {
		deleteTimeSec = 1
	}
	bustingLevel := calcBustingLevel(deleteTimeSec, args)

	list := []Param{}
	enemyIDs := map[int]int{}
	for _, e := range args.DeletedEnemies {
		enemyIDs[e.CharID] = e.CharID
	}
	haveOnlyOne := false
	for _, id := range enemyIDs {
		for _, c := range getEnemyChip(id, bustingLevel) {
			if c.IsOnlyOne {
				if plyr.HaveChip(c.ChipID) {
					continue
				} else {
					haveOnlyOne = true
				}
			}

			chipInfo := chip.Get(c.ChipID)
			list = append(list, Param{
				Type: TypeChip,
				Name: chipInfo.Name,
				Code: c.Code,
			})
		}
	}
	if !haveOnlyOne {
		m := getMoney(bustingLevel)
		list = append(list, Param{Type: TypeMoney, Name: fmt.Sprintf("%d ゼニー", m), Zenny: m})
	}

	logger.Debug("Reward list: %+v", list)

	gotReward = getReward(list)
	gotReward.BustingLevel = bustingLevel
	gotReward.DeleteTimeSec = deleteTimeSec
	rewardProc(gotReward, plyr)
	logger.Info("Got reward: %+v", gotReward)
}

func GetParam() Param {
	return gotReward
}

func calcBustingLevel(deleteTimeSec int, args WinArg) int {
	// バスティングレベルの決定
	// ウィルス戦
	//   ～ 5秒:	7point
	//   ～12秒:	6point
	//   ～36秒:	5point
	//   それ以降:	4point
	// ナビ戦
	//   ～30秒:	10point
	//   ～40秒:	 8point
	//   ～50秒:	 6point
	//   それ以降:	 4point
	// 攻撃を受けた回数(のけぞった回数)
	//   0回:		+1point
	//   1回:		+0point
	//   2回:		-1point
	//   3回:		-2point
	//   4回以上:	-3point
	// 移動したマス
	//   ～2マス:	1point
	//   3マス以上:	0point
	// 同時に倒す
	//   2体同時:	2point
	//   3体同時:	4point

	isBoss := false
	for _, e := range args.DeletedEnemies {
		if e.IsBoss {
			isBoss = true
			break
		}
	}

	lv := 4
	if isBoss {
		deadlines := []int{50, 40, 30, -1}
		for i := 0; i < len(deadlines); i++ {
			if deleteTimeSec > deadlines[i] {
				lv += i * 2
				break
			}
		}
	} else {
		deadlines := []int{36, 12, 5, -1}
		for i := 0; i < len(deadlines); i++ {
			if deleteTimeSec > deadlines[i] {
				lv += i
				break
			}
		}
	}

	switch args.PlayerDamageNum {
	case 0:
		lv++
	case 1:
	case 2:
		lv--
	case 3:
		lv -= 2
	default:
		lv -= 3
	}

	if args.PlayerMoveNum < 3 {
		lv++
	}

	// TODO 同時に倒す

	return lv
}

func getMoney(bustingLv int) int {
	table := []int{30, 30, 30, 30, 30, 50, 100, 200, 400, 500, 1000}
	if bustingLv < len(table) {
		return table[bustingLv]
	}
	return 2000
}

func getReward(all []Param) Param {
	if len(all) == 0 {
		system.SetError("no reward data")
		return Param{}
	}

	// TODO 重みづけ
	n := rand.Intn(len(all))
	return all[n]
}

func rewardProc(data Param, plyr *player.Player) {
	switch data.Type {
	case TypeMoney:
		plyr.UpdateMoney(data.Zenny)
	case TypeChip:
		c := chip.GetByName(data.Name)
		plyr.AddChip(c.ID, data.Code)
	}
}

func getEnemyChip(id int, bustingLv int) []EnemyChipInfo {
	res := []EnemyChipInfo{}
	for _, c := range enemyChipList {
		if c.CharID == id && bustingLv >= c.RequiredLevel {
			res = append(res, c)
		}
	}
	return res
}
