package main

import (
	"flag"
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/chip"
	chipimage "github.com/sh-miyoshi/go-rockmanexe/pkg/app/chip/image"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/mapinfo"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/system"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/fps"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/inputs"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
)

var (
	version     = ""
	encKey      = ""
	font    int = -1
)

func init() {
	runtime.LockOSThread()
}

func main() {
	fps.FPS = 60
	fpsMgr := fps.Fps{}

	var confFile string
	flag.StringVar(&confFile, "config", config.DefaultConfigFile, "file path of config")
	flag.Parse()

	if err := config.Init(confFile); err != nil {
		msg := fmt.Sprintf("failed to init config: %v", err)
		panic(msg)
	}

	rand.Seed(time.Now().Unix())
	dxlib.Init(config.DxlibDLLFilePath)

	if version != "" {
		config.ProgramVersion = version
	}
	if encKey != "" {
		config.EncryptKey = encKey
	}

	dxlib.SetDoubleStartValidFlag(dxlib.TRUE)
	if config.Get().Debug.RunAlways {
		dxlib.SetAlwaysRunFlag(dxlib.TRUE)
	}

	if config.Get().Debug.UsePrivateResource {
		config.ImagePath = "data/private/images/"
		config.SoundPath = "data/private/sounds/"
	}

	if config.Get().Log.DebugEnabled {
		dxlib.SetOutApplicationLogValidFlag(dxlib.TRUE)
	} else {
		dxlib.SetOutApplicationLogValidFlag(dxlib.FALSE)
	}
	logger.InitLogger(config.Get().Log.DebugEnabled, config.Get().Log.FileName)

	if res := dxlib.AddFontFile(config.FontFilePath); res == nil {
		logger.Error("Failed to load font data %s", config.FontFilePath)
		return
	}

	dxlib.ChangeWindowMode(dxlib.TRUE)
	dxlib.SetWindowSizeChangeEnableFlag(dxlib.FALSE, dxlib.FALSE)
	dxlib.SetGraphMode(config.MaxScreenSize.X, config.MaxScreenSize.Y)
	dxlib.SetWindowSize(int32(config.ScreenSize.X), int32(config.ScreenSize.Y))

	dxlib.DxLib_Init()
	dxlib.SetDrawScreen(dxlib.DX_SCREEN_BACK)

	logger.Info("Program version: %s", config.ProgramVersion)

	count := 0

	if err := appInit(); err != nil {
		logger.Error("Failed to init application: %+v", err)
		system.SetError("ゲーム初期化時")
	}

	logger.Info("Successfully init application.")

	if config.Get().Debug.InitSleepSec > 0 {
		tm := time.Duration(config.Get().Debug.InitSleepSec) * time.Second
		time.Sleep(tm)
	}
MAIN:
	for system.Error() == nil && dxlib.ScreenFlip() == 0 && dxlib.ProcessMessage() == 0 && dxlib.ClearDrawScreen() == 0 {
		inputs.KeyStateUpdate()
		if err := game.Process(); err != nil {
			logger.Error("Failed to play game: %+v", err)
			system.SetError("ゲームプレイ中")
			break MAIN
		}
		game.Draw()

		if dxlib.CheckHitKey(dxlib.KEY_INPUT_ESCAPE) == 1 {
			logger.Info("Game end by escape command")
			break MAIN
		}
		count++

		fpsMgr.Wait()

		// debug情報
		system.AddDebugMessage("[%.1f]", fpsMgr.Get())
		var x, y int
		dxlib.GetMousePoint(&x, &y)
		system.AddDebugMessage("(%d, %d)", x, y)

		debugDraw()
	}

	if err := system.Error(); err != nil {
		sound.BGMStop()
		dxlib.ClearDrawScreen()
		dxlib.DrawFormatString(10, 10, 0xff0000, "%sに回復不可能なエラーが発生しました。", err.Error())
		dxlib.DrawFormatString(10, 40, 0xff0000, "詳細はログを参照してください。")
		dxlib.ScreenFlip()
		dxlib.WaitKey()
	}

	dxlib.DxLib_End()
}

func appInit() error {
	if err := inputs.Init(config.Get().Input.Type); err != nil {
		return fmt.Errorf("inputs init failed: %w", err)
	}
	if err := chip.Init(config.ChipFilePath); err != nil {
		return fmt.Errorf("chip init failed: %w", err)
	}
	if err := chipimage.Init(); err != nil {
		return fmt.Errorf("chip image init failed: %w", err)
	}
	if err := draw.Init(); err != nil {
		return fmt.Errorf("drawing data init failed: %w", err)
	}
	if err := sound.Init(); err != nil {
		return fmt.Errorf("sound init failed: %w", err)
	}
	if err := mapinfo.Init(config.MapInfoFilePath); err != nil {
		return fmt.Errorf("map info init failed: %w", err)
	}

	return nil
}

func debugDraw() {
	if config.Get().Debug.ShowDebugData {
		if font == -1 {
			font = dxlib.CreateFontToHandle(dxlib.CreateFontToHandleOption{
				FontName: nil,
				Size:     dxlib.Int32Ptr(22),
				Thick:    dxlib.Int32Ptr(7),
				FontType: dxlib.Int32Ptr(dxlib.DX_FONTTYPE_EDGE),
			})
		}

		for i, msg := range system.PopAllDebugMessages() {
			dxlib.DrawFormatStringToHandle(0, i*25, 0xffffff, font, msg)
		}
	}
}
