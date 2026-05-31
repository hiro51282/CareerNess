# Fact Schema

CareerNess の fact は「肩書」ではなく「行動・判断・成果」を中心に保存する。role や title は保持してよいが、truth の中心には置かない。

## Design Goal

- profile と export を再生成できる粒度で保持する
- AI inference と user-confirmed truth を区別する
- human-readable な YAML を正本にする

## Core Shape

1 fact 1 claim を基本とする。複数主張を巨大な文章塊に押し込まない。

```yaml
id: fact-project-alpha-ci-speedup
type: achievement
status: proposed
project_id: project-alpha
period:
  start: 2023-04
  end: 2023-09
action: CI workflow を整理した
decision: 並列化より job 分割とキャッシュ整理を優先した
impact:
  summary: CI 実行時間を短縮した
  metrics:
    before: "25m"
    after: "12m"
context:
  constraints:
    - monorepo
    - 3 team shared pipeline
  collaborators:
    - backend-team
    - platform-team
evidence:
  sources:
    - conversation:msg-102
confidence: medium
updated_at: 2026-05-12
```

## Recommended Fields

- `id`
- `type`
- `status`
- `project_id`
- `action`
- `decision`
- `impact`
- `context`
- `evidence`
- `confidence`

## Status Model

- `proposed`
- `inferred`
- `confirmed`
- `rejected`

`confirmed` は user approval 後のみ。

## Field Semantics

### `action`

何をしたか。最小単位の行動。

### `decision`

何を判断したか。技術選定、優先順位、方針変更などを含む。

### `impact`

何が変わったか。定量・定性のどちらでもよい。

### `context`

どんな制約、チーム、状況だったか。

## Optional Fields

- `role_hints`
- `title_hints`
- `tags`
- `related_fact_ids`
- `employer_id`
- `location`

これらは補助であり、fact の中心ではない。

## Title And Role Handling

- title は保存してよい
- role は保存してよい
- ただし title / role は derived metadata に近く、truth の中心ではない
- action / decision / impact / context が先で、title / role は後から導出・整理されうる
- title / role だけで fact を成立させない

悪い例:

```yaml
id: fact-001
title: Tech Lead
status: confirmed
```

これだけでは後から profile を作る情報として弱い。

## Granularity Rule

- 1 fact はできるだけ一つの action/decision/impact に寄せる
- 同一 project に複数 facts が並ぶことは正常
- narrative ではなく queryable unit を目指す

## Evidence

evidence を厳格な監査証跡にしすぎる必要はないが、少なくとも「なぜこの fact があるのか」は追えるようにする。

- conversation message
- imported resume note
- user manual edit

## Prohibited Patterns

- marketing copy を fact として保存する
- 根拠不明の数値を impact に入れる
- confirmed fact に AI 推測を混ぜる
- export wording を facts に逆流させる

## Phase 1 での Minimal Fix 方針（2026-05-31 確定）

Reviewer レビュー（e58571e）と Planner 再評価（4d80096）を経て以下を確定した。

**Phase 1 完了時点で `action / decision / impact` の3フィールドを空で追加する。**

理由: Phase 2 での fact schema 拡張時に、confirmed facts のマイグレーションが必要になることを回避するため。空フィールドの場合は `description` にフォールバック表示する。

```yaml
# Phase 1 で追加する最小フィールド（値は空でよい）
action: ""       # 何をしたか（空の場合は description を表示）
decision: ""     # 何を判断したか
impact:
  summary: ""    # 何が変わったか
```

Phase 2 以降でユーザーがこれらのフィールドを埋めていくことで、structured fact への移行が段階的に進む。Phase 1 で蓄積した `confirmed` facts に対してマイグレーションスクリプトを走らせる必要がなくなる。

## Open Questions

- `decision` を複数値配列にするか単数フィールドにするか → Phase 2 で確定
- fact file を project 単位にまとめるか 1 fact 1 file にするか → Phase 2 で確定
