package field

import (
	"fmt"
	"math/rand"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/background"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/manager"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/math"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type extendPanelInfo struct {
	info battlecommon.PanelInfo

	objExists       bool
	prevPanelStatus int
	statusCount     int
	typeChangeCount int
	prevPanelType   int
}

var (
	imgPanel       [battlecommon.PanelStatusMax][battlecommon.PanelTypeMax]int
	imgPanelPoison []int
	blackoutCount  = 0
	animCount      = 0
	panels         [][]extendPanelInfo
	animManager    *manager.Manager
)

func Init(animMgr *manager.Manager) error {
	logger.Info("Initialize battle field data")

	animManager = animMgr

	panels = make([][]extendPanelInfo, battlecommon.FieldNum.X)
	for i := 0; i < battlecommon.FieldNum.X; i++ {
		panels[i] = make([]extendPanelInfo, battlecommon.FieldNum.Y)
	}

	// Initialize images
	files := [battlecommon.PanelStatusMax]string{"normal", "crack", "hole", "poison", "empty"}
	for i := 0; i < battlecommon.PanelStatusMax; i++ {
		if i == battlecommon.PanelStatusPoison {
			// 毒沼パネルは別途読み込み
			continue
		}

		fname := fmt.Sprintf("%sbattle/panel_player_%s.png", config.ImagePath, files[i])
		imgPanel[i][battlecommon.PanelTypePlayer] = dxlib.LoadGraph(fname)
		if imgPanel[i][battlecommon.PanelTypePlayer] < 0 {
			return errors.Newf("failed to read player panel image %s", fname)
		}

		fname = fmt.Sprintf("%sbattle/panel_enemy_%s.png", config.ImagePath, files[i])
		imgPanel[i][battlecommon.PanelTypeEnemy] = dxlib.LoadGraph(fname)
		if imgPanel[i][battlecommon.PanelTypeEnemy] < 0 {
			return errors.Newf("failed to read enemy panel image %s", fname)
		}
	}

	imgPanelPoison = make([]int, 6)
	fname := config.ImagePath + "battle/panel_poison.png"
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 64, 34, imgPanelPoison); res == -1 {
		return errors.Newf("failed to read poison panel image %s", fname)
	}

	// Initialize panel info
	for x := 0; x < battlecommon.FieldNum.X; x++ {
		t := battlecommon.PanelTypePlayer
		if x >= battlecommon.FieldNum.X/2 {
			t = battlecommon.PanelTypeEnemy
		}
		for y := 0; y < battlecommon.FieldNum.Y; y++ {
			panels[x][y] = extendPanelInfo{
				info: battlecommon.PanelInfo{
					Status: battlecommon.PanelStatusNormal,
					Type:   t,
				},
				objExists:       false,
				prevPanelStatus: battlecommon.PanelStatusNormal,
				typeChangeCount: 0,
				statusCount:     0,
			}
		}
	}

	// TODO: Map情報から取得する
	mapTypes := []int{
		background.Type秋原町,
		background.Typeアッフリク,
		background.Typeブラックアース,
	}
	if err := background.Set(mapTypes[rand.Intn(len(mapTypes))]); err != nil {
		return errors.Wrap(err, "failed to load background")
	}

	logger.Info("Successfully initialized battle field data")
	return nil
}

func End() {
	logger.Info("Cleanup battle field data")
	for i := 0; i < battlecommon.PanelStatusMax; i++ {
		for j := 0; j < battlecommon.PanelTypeMax; j++ {
			dxlib.DeleteGraph(imgPanel[i][j])
			imgPanel[i][j] = -1
		}
	}
	for i := 0; i < len(imgPanelPoison); i++ {
		dxlib.DeleteGraph(imgPanelPoison[i])
	}

	background.Unset()
	logger.Info("Successfully cleanuped battle field data")
}

func Draw() {
	for x := 0; x < battlecommon.FieldNum.X; x++ {
		for y := 0; y < battlecommon.FieldNum.Y; y++ {
			vx := battlecommon.PanelSize.X * x
			vy := battlecommon.DrawPanelTopY + battlecommon.PanelSize.Y*y
			pnStatus := panels[x][y].info.Status
			// Note:
			//   panelReturnAnimCount以下の場合StatusはNormalになる
			//   HoleとNormalを点滅させるためCountによってイメージを変える
			if panels[x][y].statusCount > 0 {
				if panels[x][y].statusCount < battlecommon.PanelReturnAnimCount && (panels[x][y].statusCount/2)%2 == 0 {
					pnStatus = panels[x][y].prevPanelStatus
				}
			}

			pnType := panels[x][y].info.Type
			if panels[x][y].typeChangeCount > 0 {
				if panels[x][y].typeChangeCount < battlecommon.PanelReturnAnimCount && (panels[x][y].typeChangeCount/2)%2 == 0 {
					pnType = panels[x][y].prevPanelType
				}
			}

			if pnStatus == battlecommon.PanelStatusPoison {
				drawPoisonPanel(vx, vy, pnType)
			} else {
				img := imgPanel[pnStatus][pnType]
				dxlib.DrawGraph(vx, vy, img, true)
			}

			damages := animManager.DamageManager().GetHitDamages(point.Point{X: x, Y: y}, "")
			for _, dm := range damages {
				if dm != nil && dm.ShowHitArea {
					x1 := vx
					y1 := vy
					x2 := vx + battlecommon.PanelSize.X
					y2 := vy + battlecommon.PanelSize.Y
					const s = 5
					dxlib.DrawBox(x1+s, y1+s, x2-s, y2-s, 0xffff00, true)
				}
			}
		}
	}
}

func DrawBlackout() {
	if blackoutCount > 0 {
		dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_ALPHA, 128)
		dxlib.DrawBox(0, 0, config.ScreenSize.X, config.ScreenSize.Y, 0x000000, true)
		dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_NOBLEND, 255)
	}
}

func Update() {
	if blackoutCount > 0 {
		blackoutCount--
		return
	}

	animCount++

	// Cleanup at first
	for x := 0; x < len(panels); x++ {
		for y := 0; y < len(panels[x]); y++ {
			panels[x][y].info.ObjectID = ""
		}
	}

	objs := animManager.ObjAnimGetObjs(objanim.FilterAll)
	for _, obj := range objs {
		panels[obj.Pos.X][obj.Pos.Y].info.ObjectID = obj.ObjID
		if panels[obj.Pos.X][obj.Pos.Y].info.Status == battlecommon.PanelStatusCrack {
			panels[obj.Pos.X][obj.Pos.Y].objExists = true
		}
	}

	// Panel status update
	for x := 0; x < len(panels); x++ {
		for y := 0; y < len(panels[x]); y++ {
			if panels[x][y].statusCount > 0 {
				panels[x][y].statusCount--
			}
			if panels[x][y].statusCount == battlecommon.PanelReturnAnimCount {
				panels[x][y].prevPanelStatus = panels[x][y].info.Status
				panels[x][y].info.Status = battlecommon.PanelStatusNormal
			}

			switch panels[x][y].info.Status {
			case battlecommon.PanelStatusHole:
			case battlecommon.PanelStatusCrack:
				// Objectが乗って離れたらHole状態へ
				if panels[x][y].objExists && panels[x][y].info.ObjectID == "" {
					sound.On(resources.SEPanelBreak)
					panels[x][y].objExists = false
					panels[x][y].info.Status = battlecommon.PanelStatusHole
					panels[x][y].statusCount = battlecommon.DefaultPanelStatusEndCount
				}
			case battlecommon.PanelStatusPoison:
				// 上に載っているオブジェクトのHPを減らす
				if animCount%30 == 0 {
					if panels[x][y].info.ObjectID != "" {
						animManager.DamageManager().New(damage.Damage{
							ID:            uuid.New().String(),
							Power:         1,
							DamageType:    damage.TypeObject,
							TargetObjID:   panels[x][y].info.ObjectID,
							TargetObjType: damage.TargetPlayer | damage.TargetEnemy,
						})
					}
				}
			}

			if panels[x][y].typeChangeCount > 0 {
				panels[x][y].typeChangeCount--
			}
			if panels[x][y].typeChangeCount == battlecommon.PanelReturnAnimCount {
				if panels[x][y].info.ObjectID != "" {
					// オブジェクトが乗っている場合は次に判定持ち越し
					panels[x][y].typeChangeCount = battlecommon.PanelReturnAnimCount + 1
				} else {
					panels[x][y].prevPanelType, panels[x][y].info.Type = panels[x][y].info.Type, panels[x][y].prevPanelType
				}
			}
		}
	}
}

func GetPanelInfo(pos point.Point) battlecommon.PanelInfo {
	if pos.X < 0 || pos.X >= battlecommon.FieldNum.X || pos.Y < 0 || pos.Y >= battlecommon.FieldNum.Y {
		return battlecommon.PanelInfo{}
	}

	// Update objectID to latest
	panels[pos.X][pos.Y].info.ObjectID = animManager.ObjAnimExistsObject(pos)

	return panels[pos.X][pos.Y].info
}

func SetBlackoutCount(cnt int) {
	blackoutCount = cnt
}

func IsBlackout() bool {
	return blackoutCount > 0
}

func ChangePanelType(pos point.Point, pnType int, endCount int) {
	if pos.X < 0 || pos.X >= battlecommon.FieldNum.X || pos.Y < 0 || pos.Y >= battlecommon.FieldNum.Y {
		return
	}

	panels[pos.X][pos.Y].prevPanelType = panels[pos.X][pos.Y].info.Type
	panels[pos.X][pos.Y].info.Type = pnType
	panels[pos.X][pos.Y].typeChangeCount = endCount
}

func ChangePanelStatus(pos point.Point, pnStatus int, endCount int) {
	if pos.X < 0 || pos.X >= battlecommon.FieldNum.X || pos.Y < 0 || pos.Y >= battlecommon.FieldNum.Y {
		return
	}

	if panels[pos.X][pos.Y].info.Status == pnStatus {
		return
	}

	if pnStatus == battlecommon.PanelStatusHole {
		if panels[pos.X][pos.Y].info.ObjectID != "" {
			panels[pos.X][pos.Y].info.Status = battlecommon.PanelStatusCrack
		} else {
			panels[pos.X][pos.Y].info.Status = battlecommon.PanelStatusHole
		}
	} else {
		panels[pos.X][pos.Y].info.Status = pnStatus
	}

	panels[pos.X][pos.Y].statusCount = endCount
}

func Set4x4Area() {
	battlecommon.FieldNum = point.Point{X: 8, Y: 4}
	config.ScreenSize = point.Point{X: 640, Y: 480}
	battlecommon.DrawPanelTopY = config.ScreenSize.Y - (battlecommon.PanelSize.Y * battlecommon.FieldNum.Y) - 30
	dxlib.SetWindowSize(640, 480)
}

func ResetSet4x4Area() {
	if Is4x4Area() {
		battlecommon.FieldNum = point.Point{X: 6, Y: 3}
		config.ScreenSize = point.Point{X: 480, Y: 320}
		battlecommon.DrawPanelTopY = config.ScreenSize.Y - (battlecommon.PanelSize.Y * battlecommon.FieldNum.Y) - 30
		dxlib.SetWindowSize(480, 320)
	}
}

func Is4x4Area() bool {
	return battlecommon.FieldNum.Equal(point.Point{X: 8, Y: 4})
}

func drawPoisonPanel(vx, vy int, pnType int) {
	n := (animCount / 15) % (len(imgPanelPoison) * 2)
	img := imgPanelPoison[math.MountainIndex(n, len(imgPanelPoison)*2)]
	dxlib.DrawBox(vx, vy, vx+battlecommon.PanelSize.X, vy+battlecommon.PanelSize.Y, 0x000000, true)
	dxlib.DrawGraph(vx, vy, imgPanel[battlecommon.PanelStatusEmpty][pnType], true)
	dxlib.DrawGraph(vx+8, vy+8, img, true)
}
