package dxlib

import "github.com/sh-miyoshi/dxlib"

type CreateFontToHandleOption struct {
	FontName *string
	Size     *int32
	Thick    *int32
	FontType *int32
	CharSet  *int32
	EdgeSize *int32
	Italic   *int32
	Handle   *int32
}

type DrawRotaGraphOption struct {
	ReverseXFlag *int32
	ReverseYFlag *int32
}

const (
	DX_BLENDMODE_INVSRC  = dxlib.DX_BLENDMODE_INVSRC
	DX_BLENDMODE_ADD     = dxlib.DX_BLENDMODE_ADD
	DX_BLENDMODE_NOBLEND = dxlib.DX_BLENDMODE_NOBLEND
	DX_BLENDMODE_ALPHA   = dxlib.DX_BLENDMODE_ALPHA

	DX_PLAYTYPE_LOOP = dxlib.DX_PLAYTYPE_LOOP
	DX_PLAYTYPE_BACK = dxlib.DX_PLAYTYPE_BACK

	DX_FONTTYPE_EDGE = dxlib.DX_FONTTYPE_EDGE

	DX_SCREEN_BACK = dxlib.DX_SCREEN_BACK
)

const (
	TRUE  = 1
	FALSE = 0
)

const (
	KEY_INPUT_Z      = dxlib.KEY_INPUT_Z
	KEY_INPUT_X      = dxlib.KEY_INPUT_X
	KEY_INPUT_LEFT   = dxlib.KEY_INPUT_LEFT
	KEY_INPUT_RIGHT  = dxlib.KEY_INPUT_RIGHT
	KEY_INPUT_UP     = dxlib.KEY_INPUT_UP
	KEY_INPUT_DOWN   = dxlib.KEY_INPUT_DOWN
	KEY_INPUT_A      = dxlib.KEY_INPUT_A
	KEY_INPUT_S      = dxlib.KEY_INPUT_S
	KEY_INPUT_D      = dxlib.KEY_INPUT_D
	KEY_INPUT_ESCAPE = dxlib.KEY_INPUT_ESCAPE
)

const (
	DX_INPUT_PAD1 = dxlib.DX_INPUT_PAD1
	DX_INPUT_KEY  = dxlib.DX_INPUT_KEY
)
