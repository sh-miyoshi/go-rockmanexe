# ロックマンエグゼ

## 遊び方

まず、[リリースリンク](https://github.com/sh-miyoshi/go-rockmanexe/releases/download/v0.12/project.zip)からzipファイルをダウンロードします。  
その後、任意の場所にzipファイルを解凍し、rockman.exeをダブルクリックで実行します。  
※dataディレクトリ配下はいじらないでください。

- 各キー配置
  - 決定・チップ使用キー: Z
  - キャンセル・バスターキー: X
  - 移動キー: 矢印(←, →, ↑, ↓)キー
  - ゲーム中のLボタン: A
  - ゲーム中のRボタン: S

## 開発機能をONにする方法

- 1. `data/config.yaml`を編集

  ```config.yaml
  # ↓を追加する
  debug:
    enable_dev_feature: true
  ```

- 2. メニュー画面でLボタン(Aキー)を押す

## ネット対戦のしかた

\[Note\]現在共有サーバはありません。ご自身でサーバーを立ち上げていただいて友人と対戦してください

### GCPを使ってサーバーを起動する方法

- 1. Cloud Runを作成します
  - Container image URL: `docker.io/smiyoshi/rockmanexe-router`
  - Container Port: 16283
  - Environment variables
    - `SESSION_ID`, `CLIENT_1_ID`, `CLIENT_1_KEY`, `CLIENT_2_ID`, `CLIENT_2_KEY`を指定してください
    - 任意の値で大丈夫ですが、そのまま認証情報となるためUUIDのような推測困難な値を推奨します
  - NETWORKINGSのタブから`Use HTTP/2 end-to-end`を有効にする
- 2. 起動ファイルの設定情報を修正します
  - ファイル: `data/config.yaml`
  - 以下のように設定してください

  ```yaml
  net:
    client_id: "" # 環境変数で指定したCLIENT_1_ID(もしくはCLIENT_2_ID)の値
    client_key: "" # 環境変数で指定したCLIENT_1_KEY(もしくはCLIENT_2_KEY)の値
    addr: "rockmanexe-test.a.run.app:443" # Clound RunのURLから「https://」を除いて「:443」を付与したアドレス
  ```

- 3. アプリを起動
  - exeファイルを起動し、「ネット対戦」を選択してください
