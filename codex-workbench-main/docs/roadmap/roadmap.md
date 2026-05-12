# Roadmap

CareerNess の roadmap は、MVP を中心に段階的に考える。最初から理想形を作るのではなく、設計思想を壊さずに最小限の流れを成立させることを優先する。

## Phase 0: Concept Fixing

目的:

- Local-first の意味を明確にする
- workspace boundary を文書化する
- facts / profiles / exports の責務を固定する
- AI capability boundary を定義する

この phase では、実装よりも設計固定を優先する。

## Phase 1: MVP Workspace Flow

目的:

- login
- workspace attach
- chat
- fact extraction
- profile generation
- export

この phase で成立すべきこと:

- attach された workspace だけを対象に AI が動く
- facts を review しながら保存できる
- profile を facts から生成できる
- export を派生物として扱える

## Phase 2: Better Structure and Review

目的:

- fact extraction 精度の改善
- patch proposal の見やすさ改善
- taxonomy / role / tag の整理
- 未確定情報の扱い改善

ここでは AI 品質よりも、「間違いをどう止めるか」を改善する。

## Phase 3: Reuse and Iteration

目的:

- role 別 profile の再利用
- 既存 facts からの再生成
- export variation の拡充
- workspace の運用しやすさ向上

## Phase 4: Optional Cloud Enhancements

目的:

- 必要最小限の同期や補助機能の検討
- 運用監視の整備
- deploy の安定化

前提:

- 正本が cloud に移るわけではない
- Local-first と user-owned data を壊さない

## 明確に後回しにするもの

- enterprise HR platform 化
- local LLM mandatory 化
- GPU hosting
- huge cloud storage
- 大規模 infra 最適化
