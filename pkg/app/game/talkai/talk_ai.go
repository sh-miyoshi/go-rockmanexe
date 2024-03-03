package talkai

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/background"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/system"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/window"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

const (
	stateInput = iota
	stateOutput
	stateWaiting
)

const (
	inputMaxLen = 80
)

var (
	state             int
	inputHandle       int
	win               window.MessageWindow
	question          string
	serverReqError    error
	serverReqResponse *Response
)

func Init() {
	state = stateInput
	inputHandle = dxlib.MakeKeyInput(inputMaxLen, false, false, false, false, false)
	win.Init()
	background.Set(background.Type秋原町)
	win.SetMessage("", window.FaceTypeRockman)
	dxlib.SetActiveKeyInput(inputHandle)
	b := dxlib.GetColor(0, 0, 0)
	w := dxlib.GetColor(255, 255, 255)
	dxlib.SetKeyInputStringColor(b, b, w, b, b, w, b, b, b, b, b, w, w, b, b, b, b)
}

func End() {
	win.End()
	background.Unset()
	dxlib.DeleteKeyInput(inputHandle)
}

func Draw() {
	background.Draw()
	win.Draw()
	dxlib.DrawBox(45, 75, 430, 140, dxlib.GetColor(232, 184, 56), true)
	dxlib.DrawFormatString(50, 80, 0x000000, "質問を入力してね")
	dxlib.DrawBox(55, 100, 420, 130, 0xffffff, true)

	switch state {
	case stateInput:
		dxlib.DrawKeyInputString(65, 110, inputHandle, true)
	}
}

func Process() bool {
	background.Process()

	switch state {
	case stateInput:
		if dxlib.CheckKeyInput(inputHandle) {
			win.SetMessage("ええと・・・", window.FaceTypeRockman)
			question = inputString(inputHandle)
			serverReqError = nil
			serverReqResponse = nil
			state = stateWaiting
			serverSend()
		}
	case stateOutput:
		return win.Process()
	case stateWaiting:
		win.Process()
		if serverReqError != nil {
			logger.Error("Failed to request server: %v", serverReqError)
			win.SetMessage("送信に失敗しました", window.FaceTypeNone)
			state = stateOutput
			return false
		}
		if serverReqResponse != nil {
			logger.Info("Success to request server: %+v", serverReqResponse)
			msg := ""
			for _, c := range serverReqResponse.Choices {
				msg += strings.ReplaceAll(c.Messages.Content, "\n", "")
			}
			win.SetMessage(msg, window.FaceTypeRockman)
			state = stateOutput
		}
	}
	return false
}

func inputString(handle int) string {
	buf := make([]byte, inputMaxLen)
	dxlib.GetKeyInputString(buf, inputHandle)
	slicedBuf := []byte{}
	for i := 0; i < inputMaxLen; i++ {
		if buf[i] == 0 {
			break
		}
		slicedBuf = append(slicedBuf, buf[i])
	}

	t := japanese.ShiftJIS.NewDecoder()
	str, _, _ := transform.Bytes(t, slicedBuf)
	return string(str)
}

func serverSend() {
	reqBody := Request{
		Model: "gpt-3.5-turbo",
		Messages: []Message{
			{
				Role:    "user",
				Content: "10行以内で簡潔に答えてください",
			},
			{
				Role:    "user",
				Content: question,
			},
		},
	}

	conf := config.Get()

	reqJSON, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", conf.AI.URL, bytes.NewBuffer(reqJSON))
	if err != nil {
		system.SetError(err.Error())
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+conf.AI.APIKey)

	client := &http.Client{}
	go func() {
		resp, err := client.Do(req)
		if err != nil {
			serverReqError = err
			return
		}
		var res Response
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			serverReqError = err
		}
		serverReqResponse = &res
	}()
}
