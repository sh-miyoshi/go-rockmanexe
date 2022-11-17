package win

import (
	"fmt"
	"strings"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/chip"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/enemy"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/titlemsg"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/inputs"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/player"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
)

const (
	stateMsg int = iota
	stateFrameIn
	stateResult

	stateMax
)

var (
	imgFrame      int
	imgZenny      int
	count         int
	state         int
	deleteTimeSec int
	bustingLevel  int
	reward        rewardInfo
	winMsgInst    *titlemsg.TitleMsg
)

type WinArg struct {
	GameTime        int
	DeletedEnemies  []enemy.EnemyParam
	PlayerMoveNum   int
	PlayerDamageNum int
}

func Init(args WinArg, plyr *player.Player) error {
	state = stateMsg
	deleteTimeSec = args.GameTime / 60
	if deleteTimeSec == 0 {
		deleteTimeSec = 1
	}
	count = 0

	fname := common.ImagePath + "battle/result_frame.png"
	imgFrame = dxlib.LoadGraph(fname)
	if imgFrame == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = common.ImagePath + "battle/zenny.png"
	imgZenny = dxlib.LoadGraph(fname)
	if imgZenny == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = common.ImagePath + "battle/msg_win.png"
	var err error
	winMsgInst, err = titlemsg.New(fname, 0)

	if err := sound.BGMPlay(sound.BGMWin); err != nil {
		return fmt.Errorf("failed to play bgm: %v", err)
	}

	bustingLevel = calcBustingLevel(args)

	m := getMoney(bustingLevel)
	list := []rewardInfo{
		{Type: rewardTypeMoney, Name: fmt.Sprintf("%d ゼニー", m), Value: m, Image: imgZenny},
	}
	enemyIDs := map[int]int{}
	for _, e := range args.DeletedEnemies {
		enemyIDs[e.CharID] = e.CharID
	}
	for _, id := range enemyIDs {
		for _, c := range enemy.GetEnemyChip(id, bustingLevel) {
			chipInfo := chip.Get(c.ChipID)
			list = append(list, rewardInfo{
				Type:  rewardTypeChip,
				Name:  chipInfo.Name,
				Value: c.Code,
				Image: chipInfo.Image,
			})
		}
	}
	logger.Debug("Reward list: %+v", list)

	reward = getReward(list)
	rewardProc(reward, plyr)
	logger.Info("Got reward: %+v", reward)

	return err
}

func End() {
	dxlib.DeleteGraph(imgFrame)
	dxlib.DeleteGraph(imgZenny)
	if winMsgInst != nil {
		winMsgInst.End()
		winMsgInst = nil
	}
	state = stateMsg
}

func Process() bool {
	count++

	switch state {
	case stateMsg:
		if winMsgInst != nil && winMsgInst.Process() {
			stateChange(stateFrameIn)
			return false
		}
	case stateFrameIn:
		if count > 60 {
			stateChange(stateResult)
			sound.On(sound.SEGotItem)
			return false
		}
	case stateResult:
		if inputs.CheckKey(inputs.KeyEnter) == 1 {
			sound.On(sound.SESelect)
			return true
		}
	}

	return false
}

func Draw() {
	baseX := common.ScreenSize.X/2 - 195
	baseY := common.ScreenSize.Y/2 - 130

	switch state {
	case stateMsg:
		if winMsgInst != nil {
			winMsgInst.Draw()
		}
	case stateFrameIn:
		x := count * baseX / 60
		if x > baseX {
			x = baseX
		}
		dxlib.DrawGraph(x, baseY, imgFrame, true)
	case stateResult:
		dxlib.DrawGraph(baseX, baseY, imgFrame, true)
		dxlib.DrawGraph(baseX+227, baseY+144, reward.Image, true)
		draw.String(baseX+60, baseY+200, 0xffffff, reward.Name)
		if reward.Type == rewardTypeChip {
			// Show chip code
			c := strings.ToUpper(reward.Value.(string))
			draw.String(baseX+195, baseY+200, 0xffffff, c)
		}
		showDeleteTime(baseX, baseY)
		draw.Number(baseX+315, baseY+95, bustingLevel)
	}
}

func stateChange(nextState int) {
	logger.Info("Change battle result win state from %d to %d", state, nextState)
	if nextState < 0 || nextState >= stateMax {
		panic(fmt.Sprintf("Invalid next battle result win state: %d", nextState))
	}
	state = nextState
	count = 0
}

func showDeleteTime(baseX, baseY int) {
	tm := deleteTimeSec

	min := tm / 60
	sec := tm % 60
	if min > 99 {
		min = 99
	}
	zero := 0
	draw.Number(baseX+255, baseY+47, min, draw.NumberOption{Padding: &zero, Length: 2})
	draw.String(baseX+288, baseY+50, 0xffffff, "：")
	draw.Number(baseX+305, baseY+47, sec, draw.NumberOption{Padding: &zero, Length: 2})
}

func calcBustingLevel(args WinArg) int {
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
		if enemy.IsBoss(e.CharID) {
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
