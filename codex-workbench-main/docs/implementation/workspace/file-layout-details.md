# File Layout Details

workspace layout は「どこに何を置くか」だけでなく、「AI がどこをどう扱うか」の説明でもある。ここでは directory レベルより一段実装寄りの推奨を定める。

## Recommended Layout

```text
CareerVault/
├── facts/
│   ├── projects/
│   ├── people/
│   └── timeline/
├── profiles/
│   ├── platform/
│   ├── backend/
│   └── leadership/
├── exports/
│   ├── resume/
│   └── replies/
├── projects/
├── taxonomy/
├── inbox/
├── history/
└── meta/
```

## Directory Notes

### `facts/projects/`

project に紐づく facts を置く。MVP では project 単位ファイルでもよい。

### `profiles/`

role-oriented view を置く。正本ではなく regenerate 可能な派生表現である。

### `exports/`

媒体依存の最終出力を置く。dispose/recreate 可能な層である。

### `history/`

patch history と rollback 補助情報を置く。Git 非前提でも戻せるための最小単位。

### `meta/`

schema version, workspace manifest, settings を置く。

## File Format Bias

- facts: YAML
- profiles: YAML + markdown text blocks でもよい
- exports: markdown / text
- history: YAML or JSON lines

CareerNess は JSON より YAML を優先する。理由は human-readable design と diff readability である。

## Naming

- 人間が見て意味が分かる slug を使う
- random uuid only naming は避ける
- path 自体が意味を持ちすぎないよう、内容の id も持つ

## Layout Rules

- facts / profiles / exports を同一ファイルに混在させない
- `inbox/` を正本の長期置き場にしない
- attach 対象外の private cache を workspace 内に紛れ込ませない

## Open Questions

- `attachments/` を常設するか
- `profiles/` の単位を role か audience かで切るか
