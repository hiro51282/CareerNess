# Deploy Strategy

CareerNess の deploy は、MVP 段階では小さく始める。思想に対して過剰な infra を先に持ち込まないことが重要である。

## 基本方針

- small scale first
- avoid infra complexity
- avoid heavy cloud dependency
- Local-first を壊さない

## MVP で優先すること

- Browser と AppServer が安定して動くこと
- workspace attach と AI orchestration が成立すること
- facts / profiles / exports の責務分離を壊さないこと
- 運用者が過剰な infra 負債を抱えないこと

## MVP で避けること

- Kubernetes 前提の deploy
- 複雑な multi-region 構成
- 大規模 queue / event bus 前提の設計
- user data を抱える巨大 cloud storage

## なぜ小さく始めるのか

CareerNess の難しさは infra 規模ではなく、workspace boundary と AI capability boundary をどう自然に成立させるかにある。したがって、MVP では infra の一般論よりも、設計思想を壊さない最小構成を選ぶ方が合理的である。

## 現実的な初期構成

- 単一 AppServer
- シンプルな認証連携
- 最小限の運用ログ
- stateless に近い backend

この構成で十分なのは、キャリア正本が cloud DB に集中しない前提だからである。

## 将来の拡張可能性

将来的に利用増加や運用要件が見えた場合、AWS や Terraform による構成管理へ進む余地はある。ただし、それは先に固定すべき前提ではない。

将来ありうるもの:

- AWS 上での標準的な app hosting
- Terraform による環境定義
- 監視や secret 管理の整備

現時点で固定しないもの:

- 本格的な distributed systems 化
- 複雑な multi-tenant data plane
- GPU workload hosting

## 未確定事項

- 最初の hosting 先は未確定
- secret 管理の具体手段は未確定
