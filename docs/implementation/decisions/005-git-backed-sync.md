# 005 Git-backed Sync はスコープ外とする

## Status

Accepted

## Context

`docs/agent-plan.md` の Phase A-3 では「CareerVault を Git リポジトリとして管理し、GitHub/GitLab で同期を提供する」ことをマルチデバイス対応のオプションとして提案した。

Reviewer からの指摘（agent-review ブランチ / e58571e）により、この設計が Local-first 原則と緊張関係にあることが明らかになった。

具体的な問題：

- push/pull タイミングと patch apply のタイミングが競合し、YAML の整合性が保証できない
- GitHub 上のリモートリポジトリが「second truth」になりやすい。特に conflict 解消時にどちらが正本かが曖昧になる
- Git に慣れていないユーザーに「workspace を Git リポジトリとして管理する」という概念を説明するコストが大きい
- ADR-003（Local-First Storage）の「ユーザーが正本を所有する」精神と、外部サービス（GitHub/GitLab）へのデータ push が緊張関係にある

## Decision

**Git-backed sync を CareerNESS の提供機能としてスコープに含めない。**

- CareerNESS は Git リポジトリ管理機能を提供しない
- ユーザーが自身の判断で CareerVault を手動で Git 管理・GitHub push することは妨げない。CareerVault は YAML / markdown ベースであり Git と相性がよい
- CareerNESS 側での「sync 機能」としての設計・実装・サポートはしない
- マルチデバイス同期はユーザーの自己責任（既存の Dropbox / iCloud / rsync 等でも代替可能）

この判断を非目標として `docs/product/non-goals.md` に追記する。

## Consequences

### Positive

- マルチデバイス同期の設計複雑性を完全に排除できる
- GitHub が「second truth」になるリスクを回避できる
- Git を知らないユーザーへの説明・サポートコストがなくなる
- ADR-003（Local-First Storage）との整合性を保てる

### Negative

- マルチデバイスでの CareerVault 共有はユーザー自身の責任になる
- デバイス喪失時のデータ保護はユーザー任せになる（この点は既存の README でも述べられている）
- sync を期待していたユーザーには機能不足に見える可能性がある

## Rejected Alternatives

- **Git-backed sync をオプション機能として提供する**: Local-first と緊張関係にあり、"second truth" 問題を根本解決できない
- **独自 cloud sync を提供する**: ADR-003 に反し、運営側がキャリアデータの保管責任を持つことになる

## 関連ドキュメント

- `docs/implementation/decisions/003-local-first-storage.md`（正本はローカルに置く）
- `docs/product/non-goals.md`（huge cloud storage を持たない）
- `docs/product/vision.md`（未確定事項: 複数 device 間の同期を公式にどう扱うかは未確定 → 本 ADR で確定）
