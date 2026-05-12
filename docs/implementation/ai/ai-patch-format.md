# AI Patch Format

CareerNess の patch proposal は、AI が workspace を直接編集しないための中核フォーマットである。ここで重要なのは、差分の適用方法そのものよりも、何をどう変えるつもりなのかを人間が読めることにある。

## Goal

- AI の変更責務を proposal に限定する
- ユーザーが「このセッションで何が変わるか」を確認できるようにする
- facts / profiles / exports で危険度の異なる変更を区別する
- apply 前に validation と approval を挟めるようにする
- rollback と patch history の単位をそろえる

## Non-Goal

- Git patch 互換を最優先にすること
- 任意バイナリ更新を表現すること
- AI が自由形式テキストだけで変更を押し通せるようにすること

## Format Principles

- patch は machine-readable かつ human-readable である
- patch 1 つは 1 semantic change を表す
- patch 本体と patch summary を分ける
- file overwrite ではなく operation の列として表現する
- facts 変更では rationale と confidence を必須に近く扱う
- AI inference は confirmed fact と同じ形で確定しない

## Patch Envelope

推奨は YAML である。CareerVault の正本が YAML ベースであるため、patch も同じ読解負荷で扱える方がよい。

```yaml
patch_id: patch-2026-05-12-001
workspace_id: local-careervault
session_id: sess_abc123
created_at: 2026-05-12T10:15:00+09:00
created_by: ai
kind: workspace_patch
summary: >
  Added two proposed facts to project alpha and regenerated platform profile draft.
status: proposed
operations:
  - op_id: op-001
    type: upsert_fact
    target: facts/projects/project-alpha.yaml
    entity_id: fact-project-alpha-ci-speedup
    change:
      before: null
      after:
        status: proposed
        action: CI workflow を整理した
        impact: CI 実行時間を短縮した
    rationale: >
      Derived from the user statement about shortening slow CI pipelines.
    confidence: medium
  - op_id: op-002
    type: replace_generated_profile
    target: profiles/platform/profile.yaml
    change:
      before_ref: profile-version-3
      after_ref: profile-version-4
    rationale: >
      Regenerated from confirmed and proposed platform-related facts.
    confidence: high
```

## Required Fields

### Patch Level

- `patch_id`
- `workspace_id`
- `session_id`
- `created_at`
- `created_by`
- `kind`
- `summary`
- `status`
- `operations`

### Operation Level

- `op_id`
- `type`
- `target`
- `change`
- `rationale`
- `confidence`

fact を触る operation では次も重要である。

- `entity_id`
- `fact_status_after`
- `review_required`

## Operation Types

MVP では operation 種別を絞る。

- `create_file`
- `update_file`
- `delete_file`
- `upsert_fact`
- `mark_fact_status`
- `replace_generated_profile`
- `replace_generated_export`
- `append_history_record`

`delete_file` は高リスクなので、MVP では UI 上で明示確認を必須にする。AI が自動適用してよい前提にはしない。

## Semantic Change Rule

1 patch = 1 semantic change を原則とする。

- fact 追加
- fact status 更新
- profile draft 再生成
- export draft 再生成

は別 patch として review できる方がよい。1 patch に複数 semantic change を詰め込むと、approval の意味が薄くなる。

## Fact Patch Rules

fact patch は特に保守的に扱う。

- `confirmed` を AI 単独で新規作成しない
- AI 由来の追加は `proposed` か `inferred` を基本とする
- 既存 confirmed fact の意味変更は差分を細かく見せる
- action / decision / impact / context のうち、欠けている項目は空欄のまま提案してよい
- 数値や期間の補完は根拠が明示できる場合に限る

## Summary View

UI に見せるのは raw patch だけでは不十分である。最低限、次の要約を生成する。

- 追加される facts
- 変更される facts
- 再生成される profiles
- 再生成される exports
- 要確認の inference
- 削除や上書きなどの高リスク操作

## Validation Before Approval

approval 画面に出す前に、AppServer または local workspace 側で最低限の validation を行う。

- target path が workspace 配下か
- operation type が許可リスト内か
- YAML schema に収まるか
- fact status の遷移が許可されているか
- confirmed fact を提案のみで上書きしていないか

## Prohibited Patterns

- patch summary だけで raw operation を持たない
- reasoning を持たない fact update
- export 更新を根拠に fact を暗黙更新する
- patch 外での silent file mutation
- workspace 外パスを target に含める
- fact 追加と export 再生成を一体の不可分 patch にする

## Open Questions

- JSON Patch 互換レイヤを後で持つかは未確定
- patch signing を入れるかは未確定
- operation ごとの conflict marker をどこまで持つかは未確定
