package draw

import (
	"fmt"
	"sort"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	netconn "github.com/sh-miyoshi/go-rockmanexe/pkg/app/netconn"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	netconfig "github.com/sh-miyoshi/go-rockmanexe/pkg/net/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/object"
)

type Option struct {
	Reverse        bool
	ViewChip       bool
	ViewHP         int
	ImgUnknownIcon int
}

type DrawManager struct {
	imgObjs        [object.TypeMax][]int
	imgEffects     [effect.TypeMax][]int
	imgUnknownIcon int
	playerObjID    string
}

var (
	inst DrawManager
)

func GetInst() *DrawManager {
	return &inst
}

func Init(playerObjID string) error {
	inst.playerObjID = playerObjID

	if err := inst.loadObjs(); err != nil {
		return fmt.Errorf("load objects failed: %w", err)
	}

	if err := inst.loadEffects(); err != nil {
		return fmt.Errorf("load effects failed: %w", err)
	}

	fname := common.ImagePath + "chipInfo/unknown_icon.png"
	inst.imgUnknownIcon = dxlib.LoadGraph(fname)
	if inst.imgUnknownIcon == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	return nil
}

func (m *DrawManager) End() {
	for _, image := range m.imgObjs {
		for _, img := range image {
			dxlib.DeleteGraph(img)
		}
	}
	for _, image := range m.imgEffects {
		for _, img := range image {
			dxlib.DeleteGraph(img)
		}
	}
	dxlib.DeleteGraph(m.imgUnknownIcon)
}

func (m *DrawManager) DrawObjects() {
	ginfo := netconn.GetInst().GetGameInfo()
	objects := []object.Object{}
	for _, obj := range ginfo.Objects {
		objects = append(objects, obj)
	}

	sort.Slice(objects, func(i, j int) bool {
		ii := objects[i].Y*netconfig.FieldNumX + objects[i].X
		ij := objects[j].Y*netconfig.FieldNumX + objects[j].X
		return ii < ij
	})
	for _, obj := range objects {
		reverse := false

		if obj.ClientID != config.Get().Net.ClientID {
			// enemy object
			reverse = true
		}

		viewHP := 0
		if obj.ID != m.playerObjID {
			viewHP = obj.HP
		}

		drawObject(m.imgObjs, obj, Option{
			Reverse:        reverse,
			ViewHP:         viewHP,
			ViewChip:       obj.ID != m.playerObjID,
			ImgUnknownIcon: m.imgUnknownIcon,
		})
	}
}

func (m *DrawManager) DrawEffects() {
	ginfo := netconn.GetInst().GetGameInfo()
	for _, eff := range ginfo.Effects {
		num, delay := m.GetEffectImageInfo(eff.Type)
		if eff.Count >= num*delay {
			netconn.GetInst().RemoveEffect(eff.ID)
			continue
		}
		imgNo := eff.Count / delay
		drawEffect(m.imgEffects, imgNo, eff)
	}
}

func (m *DrawManager) GetObjectImageInfo(objType int) (imageNum, delay int) {
	return len(m.imgObjs[objType]), object.ImageDelays[objType]
}

func (m *DrawManager) GetEffectImageInfo(effType int) (imageNum, delay int) {
	return len(m.imgEffects[effType]), effect.Delays[effType]
}

func (m *DrawManager) loadObjs() error {
	fname := common.ImagePath + "battle/character/player_move.png"
	m.imgObjs[object.TypeRockmanMove] = make([]int, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 100, 100, m.imgObjs[object.TypeRockmanMove]); res == -1 {
		return fmt.Errorf("failed to load player move image: %s", fname)
	}

	fname = common.ImagePath + "battle/character/player_damaged.png"
	m.imgObjs[object.TypeRockmanDamage] = make([]int, 6)
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 100, 100, m.imgObjs[object.TypeRockmanDamage]); res == -1 {
		return fmt.Errorf("failed to load player damage image: %s", fname)
	}
	// 1 -> 2,3  2-4 3-5
	m.imgObjs[object.TypeRockmanDamage][4] = m.imgObjs[object.TypeRockmanDamage][2]
	m.imgObjs[object.TypeRockmanDamage][5] = m.imgObjs[object.TypeRockmanDamage][3]
	m.imgObjs[object.TypeRockmanDamage][2] = m.imgObjs[object.TypeRockmanDamage][1]
	m.imgObjs[object.TypeRockmanDamage][3] = m.imgObjs[object.TypeRockmanDamage][1]

	fname = common.ImagePath + "battle/character/player_shot.png"
	m.imgObjs[object.TypeRockmanShot] = make([]int, 6)
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 180, 100, m.imgObjs[object.TypeRockmanShot]); res == -1 {
		return fmt.Errorf("failed to load player shot image: %s", fname)
	}

	fname = common.ImagePath + "battle/character/player_cannon.png"
	m.imgObjs[object.TypeRockmanCannon] = make([]int, 6)
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 100, 100, m.imgObjs[object.TypeRockmanCannon]); res == -1 {
		return fmt.Errorf("failed to load player cannon image: %s", fname)
	}

	fname = common.ImagePath + "battle/character/player_sword.png"
	m.imgObjs[object.TypeRockmanSword] = make([]int, 7)
	if res := dxlib.LoadDivGraph(fname, 7, 7, 1, 128, 128, m.imgObjs[object.TypeRockmanSword]); res == -1 {
		return fmt.Errorf("failed to load player sword image: %s", fname)
	}

	fname = common.ImagePath + "battle/character/player_bomb.png"
	m.imgObjs[object.TypeRockmanBomb] = make([]int, 7)
	if res := dxlib.LoadDivGraph(fname, 5, 5, 1, 100, 114, m.imgObjs[object.TypeRockmanBomb]); res == -1 {
		return fmt.Errorf("failed to load player bomb image: %s", fname)
	}
	m.imgObjs[object.TypeRockmanBomb][5] = m.imgObjs[object.TypeRockmanBomb][4]
	m.imgObjs[object.TypeRockmanBomb][6] = m.imgObjs[object.TypeRockmanBomb][4]

	fname = common.ImagePath + "battle/character/player_buster.png"
	m.imgObjs[object.TypeRockmanBuster] = make([]int, 6)
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 180, 100, m.imgObjs[object.TypeRockmanBuster]); res == -1 {
		return fmt.Errorf("failed to load player buster image: %s", fname)
	}

	fname = common.ImagePath + "battle/character/player_pick.png"
	m.imgObjs[object.TypeRockmanPick] = make([]int, 6)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 96, 124, m.imgObjs[object.TypeRockmanPick]); res == -1 {
		return fmt.Errorf("failed to load player pick image: %s", fname)
	}
	m.imgObjs[object.TypeRockmanPick][4] = m.imgObjs[object.TypeRockmanPick][3]
	m.imgObjs[object.TypeRockmanPick][5] = m.imgObjs[object.TypeRockmanPick][3]

	m.imgObjs[object.TypeRockmanStand] = make([]int, 1)
	m.imgObjs[object.TypeRockmanStand][0] = m.imgObjs[object.TypeRockmanMove][0]

	skillPath := common.ImagePath + "battle/skill/"
	fname = skillPath + "キャノン_atk.png"
	tmp := make([]int, 24)
	if res := dxlib.LoadDivGraph(fname, 24, 8, 3, 120, 140, tmp); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	m.imgObjs[object.TypeNormalCannonAtk] = make([]int, 8)
	m.imgObjs[object.TypeHighCannonAtk] = make([]int, 8)
	m.imgObjs[object.TypeMegaCannonAtk] = make([]int, 8)
	for i := 0; i < 8; i++ {
		m.imgObjs[object.TypeNormalCannonAtk][i] = tmp[i]
		m.imgObjs[object.TypeHighCannonAtk][i] = tmp[i+8]
		m.imgObjs[object.TypeMegaCannonAtk][i] = tmp[i+16]
	}

	fname = skillPath + "キャノン_body.png"
	tmp = make([]int, 15)
	if res := dxlib.LoadDivGraph(fname, 15, 5, 3, 46, 40, tmp); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	m.imgObjs[object.TypeNormalCannonBody] = make([]int, 5)
	m.imgObjs[object.TypeHighCannonBody] = make([]int, 5)
	m.imgObjs[object.TypeMegaCannonBody] = make([]int, 5)
	for i := 0; i < 5; i++ {
		m.imgObjs[object.TypeNormalCannonBody][i] = tmp[i]
		m.imgObjs[object.TypeHighCannonBody][i] = tmp[i+5]
		m.imgObjs[object.TypeMegaCannonBody][i] = tmp[i+10]
	}

	fname = skillPath + "ミニボム.png"
	m.imgObjs[object.TypeMiniBomb] = make([]int, 5)
	if res := dxlib.LoadDivGraph(fname, 5, 5, 1, 40, 30, m.imgObjs[object.TypeMiniBomb]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = skillPath + "ソード.png"
	tmp = make([]int, 12)
	if res := dxlib.LoadDivGraph(fname, 12, 4, 3, 160, 150, tmp); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	m.imgObjs[object.TypeSword] = make([]int, 4)
	m.imgObjs[object.TypeWideSword] = make([]int, 4)
	m.imgObjs[object.TypeLongSword] = make([]int, 4)
	for i := 0; i < 4; i++ {
		m.imgObjs[object.TypeSword][i] = tmp[i]
		m.imgObjs[object.TypeWideSword][i] = tmp[i+8]
		m.imgObjs[object.TypeLongSword][i] = tmp[i+4]
	}

	fname = skillPath + "リカバリー.png"
	m.imgObjs[object.TypeRecover] = make([]int, 8)
	if res := dxlib.LoadDivGraph(fname, 8, 8, 1, 84, 144, m.imgObjs[object.TypeRecover]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = skillPath + "スプレッドガン_atk.png"
	m.imgObjs[object.TypeSpreadGunAtk] = make([]int, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 75, 76, m.imgObjs[object.TypeSpreadGunAtk]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = skillPath + "スプレッドガン_body.png"
	m.imgObjs[object.TypeSpreadGunBody] = make([]int, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 56, 76, m.imgObjs[object.TypeSpreadGunBody]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = skillPath + "バルカン.png"
	m.imgObjs[object.TypeVulcan] = make([]int, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 66, 50, m.imgObjs[object.TypeVulcan]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = skillPath + "ウェーブ_body.png"
	m.imgObjs[object.TypePick] = make([]int, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 128, 136, m.imgObjs[object.TypePick]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = skillPath + "サンダーボール.png"
	m.imgObjs[object.TypeThunderBall] = make([]int, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 64, 80, m.imgObjs[object.TypeThunderBall]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = skillPath + "ワイドショット_body.png"
	m.imgObjs[object.TypeWideShotBody] = make([]int, 3)
	if res := dxlib.LoadDivGraph(fname, 3, 3, 1, 56, 66, m.imgObjs[object.TypeWideShotBody]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = skillPath + "ワイドショット_begin.png"
	m.imgObjs[object.TypeWideShotBegin] = make([]int, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 90, 147, m.imgObjs[object.TypeWideShotBegin]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = skillPath + "ワイドショット_move.png"
	m.imgObjs[object.TypeWideShotMove] = make([]int, 3)
	if res := dxlib.LoadDivGraph(fname, 3, 3, 1, 90, 148, m.imgObjs[object.TypeWideShotMove]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = skillPath + "ショックウェーブ.png"
	m.imgObjs[object.TypeShockWave] = make([]int, 7)
	if res := dxlib.LoadDivGraph(fname, 7, 7, 1, 100, 140, m.imgObjs[object.TypeShockWave]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	return nil
}

func (m *DrawManager) loadEffects() error {
	fname := common.ImagePath + "battle/effect/hit_small.png"
	m.imgEffects[effect.TypeHitSmallEffect] = make([]int, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 40, 44, m.imgEffects[effect.TypeHitSmallEffect]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = common.ImagePath + "battle/effect/hit_big.png"
	m.imgEffects[effect.TypeHitBigEffect] = make([]int, 6)
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 90, 76, m.imgEffects[effect.TypeHitBigEffect]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = common.ImagePath + "battle/effect/explode.png"
	m.imgEffects[effect.TypeExplodeEffect] = make([]int, 16)
	if res := dxlib.LoadDivGraph(fname, 16, 8, 2, 110, 124, m.imgEffects[effect.TypeExplodeEffect]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = common.ImagePath + "battle/effect/cannon_hit.png"
	m.imgEffects[effect.TypeCannonHitEffect] = make([]int, 7)
	if res := dxlib.LoadDivGraph(fname, 7, 7, 1, 110, 136, m.imgEffects[effect.TypeCannonHitEffect]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = common.ImagePath + "battle/effect/spread_and_bamboo_hit.png"
	m.imgEffects[effect.TypeSpreadHitEffect] = make([]int, 6)
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 92, 88, m.imgEffects[effect.TypeSpreadHitEffect]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	tmp := make([]int, 8)
	fname = common.ImagePath + "battle/effect/vulcan_hit.png"
	if res := dxlib.LoadDivGraph(fname, 8, 8, 1, 50, 58, tmp); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	m.imgEffects[effect.TypeVulcanHit1Effect] = []int{}
	m.imgEffects[effect.TypeVulcanHit2Effect] = []int{}
	for i := 0; i < 4; i++ {
		m.imgEffects[effect.TypeVulcanHit1Effect] = append(m.imgEffects[effect.TypeVulcanHit1Effect], tmp[i])
		m.imgEffects[effect.TypeVulcanHit2Effect] = append(m.imgEffects[effect.TypeVulcanHit2Effect], tmp[i+4])
	}

	return nil
}
