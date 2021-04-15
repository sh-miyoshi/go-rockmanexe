# TODO List

## バグ

- ミニボムがplayer animとあってない

## マイルストーン

### v1.0

- 音楽関係
  - BGM
  - SE
- 仕上げ
  - logの埋め込み
  - release方法確立

### v1.1

- セーブ・続きから
- menu
  - Record
- 勝利時にチップをゲット
- フォルダ編集
- 敵追加
  - ビリー
- チップ追加
- プログラムアドバンス(PA)

## 機能一覧

- タイトル
  - はじめから
  - つづきから
- フィールド
  - 移動
  - 当たり判定
  - 敵出現
  - NPC
- バトル
  - 敵登場
  - チップ選択
  - ソウルユニゾン
  - PA
  - 戦闘開始
  - 戦闘
  - 結果
- メニュー
  - フォルダ編集
  - サブチップ(?)
  - データライブラリ
  - ロックマン
    - ナビカス
  - Eメール
  - キーアイテム
  - セーブ
- 通信対戦
- ゲームオーバー
- エンディング
- シナリオ

## Sound Effects

### シーンごと

- title
  - [x] 決定
- menu
  - top
    - [x] 選択
    - [x] 決定
  - folder
    - [x] スクロール
    - [x] キャンセル
    - [x] Denied(のちに決定)
  - gobattle
    - [x] 選択
    - [x] キャンセル
    - [x] 決定(敵登場)
  - record
    - [x] キャンセル
- battle
  - opening
    - [x] 敵登場
  - chip select
    - [x] フレームイン
    - [x] カーソル移動
    - [x] チップ選択
    - [x] キャンセル
    - [x] 決定
  - before main
  - main
    - [ ] キャラアニメーション
    - [ ] スキル
    - [x] カスタムゲージマックス
    - [ ] フォルダオープン
  - result win
    - [ ] 敵デリート
    - [ ] アイテム登場
    - [ ] 決定
  - result lose
    - [ ] player dead

### 項目ごと

- skill
  - [x] Cannon
  - [x] HighCannon
  - [x] MegaCannon
  - [ ] MiniBomb
  - [ ] Sword
  - [ ] WideSword
  - [ ] LongSword
  - [ ] ShockWave
  - [ ] Recover
  - [ ] SpreadGun
  - [ ] Vulcan1
- animation
  - player
    - [ ] damaged
    - [x] buster
    - [x] charge
    - [x] charged
- effect
  - [ ] TypeHitSmall
  - [ ] TypeHitBig
  - [ ] TypeExplode
  - [ ] TypeCannonHit
  - [ ] TypeSpreadHit
  - [ ] TypeVulcanHit1
  - [ ] TypeVulcanHit2

## 未実装項目

### 全体

- panic時のエラーの表示の仕方
- つづきから
- キー設定の変更と保存

### バトル時

- stateMain
  - その他
    - チップ情報を左下に表示
    - ダメージ情報にのけぞるかのフラグ
      - のけぞらない処理の追加
    - 敵名表示
    - ヒット後貫通しない攻撃ならスキルを止める
- stateResult
  - lose
  - win
    - バスティングレベル
    - チップゲット

### メニュー

- Top
  - 左下にアニメーション
- Chip Folder
  - Chip Folder ListのCodeのフォント
  - フォルダ編集
  - チップの説明文
- Go Battle
- Record
- (Settings)

### その他の状態

- すべて
