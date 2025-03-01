# airules

`airules` は、AI アシスタントを搭載したエディタ（Windsurf や Cursor など）の設定ファイルを適切な場所にインストールするための CLI ツールです。

## 機能

- Windsurf の設定ファイルをインストールする
- Cursor の設定ファイルをインストールする
- ローカル設定ファイルとグローバル設定ファイルの選択的インストール

## インストール方法

### Go を使用してインストール

```bash
go install github.com/hashiiiii/airules@latest
```

### ソースからビルド

```bash
git clone https://github.com/hashiiiii/airules.git
cd airules
go build -o bin/airules ./cmd/airules
```

## 使用方法

### 基本的なコマンド

```bash
# Windsurf の設定ファイルをインストール（ローカルとグローバル両方）
airules windsurf

# Cursor の設定ファイルをインストール（ローカルとグローバル両方）
airules cursor

# Windsurf のローカル設定ファイルのみをインストール
airules windsurf -l
# または
airules windsurf --local

# Cursor のグローバル設定ファイルのみをインストール
airules cursor -g
# または
airules cursor --global

# バージョン情報を表示
airules version

# ヘルプを表示
airules -h
```

## 設定ファイルの場所

### Windsurf

- ローカル設定ファイル: カレントディレクトリの `cascade.local.json`
- グローバル設定ファイル:
  - macOS/Linux: `~/.config/windsurf/cascade.global.json`
  - Windows: `%APPDATA%\Windsurf\cascade.global.json`

### Cursor

- ローカル設定ファイル: カレントディレクトリの `prompt_library.local.json`
- グローバル設定ファイル:
  - macOS: `~/Library/Application Support/Cursor/prompt_library.global.json`
  - Linux: `~/.config/cursor/prompt_library.global.json`
  - Windows: `%APPDATA%\Cursor\prompt_library.global.json`

## テンプレートのカスタマイズ

テンプレートファイルは `vendor/rules-for-ai` ディレクトリ内にあります。これらのファイルを編集することで、インストールされる設定をカスタマイズできます。

## ライセンス

MIT
