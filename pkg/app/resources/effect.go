package resources

import "github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"

const (
	EffectTypeNone int = iota
	EffectTypeHitSmall
	EffectTypeHitBig
	EffectTypeExplode
	EffectTypeCannonHit
	EffectTypeSpreadHit
	EffectTypeVulcanHit1
	EffectTypeVulcanHit2
	EffectTypeWaterBomb
	EffectTypeBlock
	EffectTypeBambooHit
	EffectTypeHeatHit
	EffectTypeExclamation
	EffectTypeFailed
	EffectTypeIceBreak
	EffectTypeExplodeSmall
	EffectTypeSpecialStart
	EffectTypeDeltaRayEdge

	EffectTypeMax
)

type EffectParam struct {
	Type      int
	Pos       point.Point
	RandRange int
}
