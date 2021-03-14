# TODO List

## 全体

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

## バトル

- stateOpening
  - decide enemies
  - show enemies
  - force state change to chip select
- stateChipSelect
- stateBeforeMain
- stateMain
  - player action
    - move
    - rock buster
    - chip use
    - damaged
  - enemy action
    - unique action
      - move
      - attack
    - chip use(?)
    - damaged
- stateResult
