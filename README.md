# boxfmt

Markdown ファイル内の罫線ボックス (ASCII / Unicode) を検出し、内容の幅に合わせて整形するフォーマッタです。
CJK 文字 (日本語等) の表示幅を正しく計算し、列の揃えやパディングを行います。

## インストール

```bash
go install github.com/HMasataka/claude/boxfmt@latest
```

または、リポジトリをクローンしてビルド:

```bash
git clone https://github.com/HMasataka/boxfmt.git
cd boxfmt
go build
```

## 使い方

```
boxfmt [options] <file>
```

### オプション

| フラグ      | 説明                 |
| ----------- | -------------------- |
| `-w`        | 入力ファイルを上書き |
| `-o <path>` | 指定パスに出力       |

`-w` と `-o` を同時に指定するとエラーになります。
どちらも指定しない場合は標準出力に結果を表示します。

### 例

```bash
# 標準出力に整形結果を表示
boxfmt input.md

# ファイルを上書き
boxfmt -w input.md

# 別ファイルに出力
boxfmt -o output.md input.md
```

## 特徴

- **Unicode 罫線** (`┌ ─ ┐ └ ┘ │ ├ ┤ ┬ ┴ ┼`) と **ASCII 罫線** (`+ - |`) の両方に対応
- **CJK 文字** (日本語・中国語・韓国語) の表示幅を正しく計算してパディング
- **複数列テーブル** の各列を独立して幅揃え
- **インデント保持** -- ボックス全体のインデントを維持
- **タブ展開** -- タブを 4 スペースに変換
- **非ボックス部分はそのまま** -- 通常の Markdown テキストやコードブロック内のボックスには手を加えない

## テスト

```bash
go test ./...
```

ゴールデンファイルの更新:

```bash
go test -run TestGoldenFiles -update
```
