//go:build mac
// +build mac

package dxlib

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

// TODO: 正しい値を入れる

const (
	DX_BLENDMODE_INVSRC  = 0
	DX_BLENDMODE_ADD     = 0
	DX_BLENDMODE_NOBLEND = 0
	DX_BLENDMODE_ALPHA   = 0

	DX_PLAYTYPE_LOOP = 0
	DX_PLAYTYPE_BACK = 0

	DX_FONTTYPE_EDGE = 0

	DX_SCREEN_BACK = 0
)

const (
	TRUE  = 1
	FALSE = 0
)

const (
	KEY_INPUT_Z      = 0
	KEY_INPUT_X      = 0
	KEY_INPUT_LEFT   = 0
	KEY_INPUT_RIGHT  = 0
	KEY_INPUT_UP     = 0
	KEY_INPUT_DOWN   = 0
	KEY_INPUT_A      = 0
	KEY_INPUT_S      = 0
	KEY_INPUT_D      = 0
	KEY_INPUT_ESCAPE = 0
)

const (
	DX_INPUT_PAD1 = 0
	DX_INPUT_KEY  = 0
)
