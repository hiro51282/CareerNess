# AI Tooling Model

CareerNess の AI は「何でもできる agent」ではなく、workspace-scoped な支援系ツール群の上で動く。tooling model の目的は、能力を増やすことではなく、能力境界を固定することにある。

## Design Principle

- AI は workspace system の一部として動く
- tool は capability whitelist である
- direct filesystem mutation は tool から分離する
- user-owned truth は tool 経由でも越権更新しない

## Tool Categories

### Read Tools

- workspace tree を読む
- facts / profiles / exports を読む
- schema / taxonomy を読む
- patch history を読む

read tool は attach された workspace root を必須引数に持つ。

### Analysis Tools

- fact extraction
- inconsistency detection
- profile synthesis
- export drafting
- duplicate fact suggestion

analysis tool は proposal を返してよいが、apply はしない。

### Patch Tools

- workspace patch generation
- patch validation
- preview diff generation

patch tool は file write ではなく patch object を返す。

### Apply Tools

- approved patch apply
- history append
- rollback preparation

apply tool は user approval token か approval state を前提とする。AI が任意に叩く前提にはしない。

## Minimum Safe Tool Set

MVP では次があれば十分である。

- `read_workspace`
- `read_file`
- `list_facts`
- `generate_patch_proposal`
- `validate_patch_proposal`
- `apply_approved_patch`
- `list_patch_history`
- `prepare_rollback_patch`

## Explicit Boundaries

tool には次の制約を組み込む。

- workspace root 外パスは拒否
- hidden retention 用の外部保存をしない
- attach されていない別 workspace を参照しない
- unrestricted shell 相当の escape hatch を用意しない

## Fact Safety Rules

fact に触る tool は追加制約を持つ。

- AI-only write を禁止
- `confirmed` への直行遷移を禁止
- source evidence を持たない数値補完を禁止
- export を source of truth として reverse write しない

## Human-Readable Output

tool の返り値は内部都合の opaque blob より、レビューしやすい構造を優先する。

- why this change
- what files are affected
- what facts are touched
- which fields are uncertain

## AppServer Responsibility

AppServer は tool を束ねる orchestration layer である。

- conversation context を渡す
- workspace-scoped capability を付与する
- proposal と approval を接続する
- truth を自分で保持しない

AppServer 自体が canonical state store になってはいけない。

## Prohibited Patterns

- AI に raw filesystem write 権限を渡す
- profile generation tool が fact update を内包する
- export tool が暗黙に history 以外の file mutation を行う
- tool response 内で workspace 外参照を許す

## Open Questions

- local client tool と AppServer-side tool の分担粒度は未確定
- attach 時にどこまで directory-level capability を絞るかは未確定
