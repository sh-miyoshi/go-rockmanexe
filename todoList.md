# TODO List

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

## 未実装項目

### 全体

- panic時のエラーの表示の仕方
- つづきから
- キー設定の変更と保存

### バトル時

- stateOpening
  - decide enemies
  - show enemies
  - force state change to chip select
- stateBeforeMain
- stateMain
  - player action
    - damaged
      - Animationや無敵処理など
    - image delay
  - その他
    - チップ情報を左下に表示
    - ダメージ情報にのけぞるかのフラグ
    - 敵名表示
    - ヒット後貫通しない攻撃ならスキルを止める
    - sort anim by type
- stateResult
  - lose
  - win

### その他の状態

- すべて

## チップ情報

### 優先的実装予定

- リカバリー10

### 後回し

- クラックアウト
- エリアスチール
- アタック+10

### 悩み中

- エアシュート
- バルカン1
- ・・・
