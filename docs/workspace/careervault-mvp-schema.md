# CareerVault MVP Schema

これは MVP の**最小限の実装スキーマ**です。目的：AI extraction → patch proposal → human review → apply → knowledge accumulate ループを成立させること。

**設計思想**: Flat で人間が読める YAML。AI が会話から抽出でき、人間が git diff で review できる形。

---

## 1. MVP ディレクトリ構成

```
CareerVault/
├── facts/
│   ├── experiences.yaml      # キャリア経験、プロジェクト、職務
│   ├── achievements.yaml     # 成果、改善、実績
│   └── skills.yaml           # 技術スタック、能力
├── profiles/
│   ├── backend-engineer.yaml # ロール別ナラティブ（自動生成）
│   └── platform-engineer.yaml
├── exports/
│   ├── resume-jp.md          # 最終フォーマット出力
│   └── resume-en.md
├── projects.yaml             # プロジェクトメタデータ（会社、チーム、期間）
├── tags.yaml                 # tag 語彙（一元管理）
└── meta.yaml                 # workspace version、メタデータ
```

**MVP に含めない**: inbox/, attachments/, relationships.yaml, domains.yaml, roles.yaml

**この構成の理由**:
- ドメインごと単一ファイル（管理が簡単、責務が明確）
- Flat なディレクトリ（ナビゲーション簡単）
- ネストなし（認知負荷を減らす）
- Patch proposal フロー に適応（1つの semantic change = 1つのファイル修正）

---

## 2. MVP 最小限 Fact スキーマ

### 必須フィールド（常に必要）

```yaml
# facts/experiences.yaml

- fact_id: fact-proj-payment-platform    # 安定的で semantic な ID
  type: experience                        # experience | achievement | skill（このみっつだけ）
  status: confirmed                       # confirmed | proposed
  summary: Payment platform 刷新         # ワンライナー、検索対象
  period: "2022-04 to 2023-01"           # 経験の場合
  company: 株式会社 Example Corp          # 経験の場合
  description: |                          # 自由形式、人間が読める
    Payment service を刷新。
    レイテンシー 30% 削減、コスト 15% 削減。
    技術選定・3チーム調整をリード。
  confidence: high                        # high | medium | low
  source: conversation                    # conversation | manual | import
  created_at: 2026-05-23T10:00:00Z
  tags:
    - backend
    - platform
    - ci-cd
```

**以上。12 フィールド。すべて 1 画面に収まる。すべて人間が読める。**

### オプショナル（情報があれば追加）

```yaml
- fact_id: fact-proj-payment-platform
  # ... 上記必須フィールド ...
  
  # 技術コンテキスト（この fact に関連しているなら）
  tech_stack: [Go, gRPC, PostgreSQL, Kubernetes]
  
  # インパクト指標（定量化可能で確認済みなら）
  impact:
    latency_improvement_pct: 30
    cost_reduction_pct: 15
  
  # Extraction 追跡用（review の手助けになる）
  source_detail: "ユーザーがセッション xyz で言及、latency metrics は prod logs で確認"
```

**ルール**: オプショナルフィールドは、その fact に情報があるときだけ追加。空白フィールドを作らない。

### なぜこれらのフィールドか

| フィールド | 理由 | 対象 |
|-----------|-----|------|
| `fact_id` | patch 全体で安定的に参照される | AI、patch system |
| `type` | experience/achievement/skill を区別 | AI filter、profile 生成 |
| `status` | confirmed vs proposed を追跡 | 人間が信頼度判定、AI が hallucination 回避 |
| `summary` | 素早いスキャン、grep 対象 | 人間の review、検索 |
| `description` | 完全なコンテキスト、ナラティブ | 人間の理解、profile 下書き |
| `confidence` | 不確定性を示す | 人間が review 優先度を判定 |
| `source` | どこから来たのか | 人間が検証 |
| `created_at` | 履歴のための timestamp | workspace audit trail |
| `tags` | 横断的な整理 | profile 生成、検索 |
| `period`, `company` | コンテキスト（experience のみ） | profile 構築、フィルタリング |

それ以外（tech_stack、impact、source_detail）は**状況による**— fact に情報があれば追加。

---

## 3. ファイル形式：1 ファイル、複数 Fact

### experiences.yaml

```yaml
# facts/experiences.yaml
# キャリア経験・プロジェクト一覧

- fact_id: fact-proj-payment-platform
  type: experience
  status: confirmed
  summary: Payment platform 刷新（2022-04〜2023-01）
  period: "2022-04 to 2023-01"
  company: 株式会社 Example
  description: |
    Payment service の刷新をリード。Python service から Go へ移行し、
    latency とコストを最適化。3チーム（決済、platform、infra）を調整。
    
    主な成果：
    - Latency 30% 削減（200ms → 140ms p99）
    - コスト 15% 削減（リソース最適化）
    - チーム開発速度 +20% 向上
  confidence: high
  source: conversation
  created_at: 2026-05-23T10:00:00Z
  tags:
    - backend
    - platform
    - architecture
    - leadership
  tech_stack: [Go, gRPC, PostgreSQL, Kubernetes]
  impact:
    latency_improvement_pct: 30
    cost_reduction_pct: 15
  source_detail: |
    2026-05-23 のユーザー説明から抽出。
    Latency metrics は Datadog レビューで確認。
    コスト削減は infra チームフィードバックで推定。

- fact_id: fact-proj-internal-tools
  type: experience
  status: proposed  # まだ確認されていない
  summary: 内部ツール framework
  period: "2021-06 to 2022-03"
  company: 株式会社 Example
  description: |
    5つのチームで使用される内部 CLI tools framework を開発。
  confidence: medium
  source: conversation
  created_at: 2026-05-24T11:00:00Z
  tags:
    - backend
    - developer-tools
  source_detail: |
    ユーザーが簡潔に言及。
    確認が必要：どの 5 チーム？どのくらい長く使用？定量的な採用度？
```

### achievements.yaml

```yaml
# facts/achievements.yaml
# 成果、改善、学習

- fact_id: fact-ach-latency-improvement
  type: achievement
  status: confirmed
  summary: Payment latency を 30% 削減
  description: |
    Payment service アーキテクチャを Python から Go へ再設計。
    Profiling を通じた hot path 最適化、キャッシング導入、
    データベース connection pooling を実装。
    結果：p99 latency 200ms → 140ms。
  confidence: high
  source: conversation
  created_at: 2026-05-23T10:30:00Z
  tags:
    - performance
    - backend
  impact:
    latency_improvement_pct: 30
    user_experience: 決済確認が高速化、顧客満足度向上
  related_experience: fact-proj-payment-platform

- fact_id: fact-ach-team-coordination
  type: achievement
  status: proposed
  summary: 3チーム統合に向けた統一 API spec 策定を調整
  description: |
    決済、platform、infra チームを調整し、統一された API contract に合意。
    技術的議論をリード。
  confidence: medium
  source: conversation
  created_at: 2026-05-23T10:45:00Z
  tags:
    - leadership
    - cross-team
  related_experience: fact-proj-payment-platform
```

### skills.yaml

```yaml
# facts/skills.yaml
# 技術、スキル、能力

- fact_id: fact-skill-go
  type: skill
  status: confirmed
  summary: Go（システム・バックエンド）
  description: |
    本番環境での 10 年以上の経験。Concurrent systems、
    performance optimization、gRPC の専門知識。
    3人のエンジニアを concurrency patterns で指導した経験あり。
  confidence: high
  source: manual
  created_at: 2026-05-20T00:00:00Z
  tags:
    - language
    - backend
  proficiency: advanced
  used_in_projects:
    - fact-proj-payment-platform
    - fact-proj-internal-tools

- fact_id: fact-skill-leadership
  type: skill
  status: confirmed
  summary: 技術リーダーシップ、チーム調整
  description: |
    アーキテクチャ意思決定、cross-team 技術議論、メンタリング経験。
    3～8 人程度のチーム経験あり。
  confidence: high
  source: manual
  created_at: 2026-05-20T00:00:00Z
  tags:
    - soft-skill
    - leadership
```

---

## 4. Tag 語彙（一元管理、最小限）

```yaml
# tags.yaml
# 厳選された tag リスト（AI は新しい tag を作らない）

tags:
  # 技術ドメイン
  backend: バックエンド、API、データベース
  frontend: フロントエンド、UI、web
  platform: インフラ、platform、DevOps
  
  # 能力
  leadership: チームリーダーシップ、メンタリング
  architecture: システム設計、アーキテクチャ意思決定
  performance: パフォーマンス最適化
  
  # スキル種別
  language: プログラミング言語
  framework: フレームワーク、ライブラリ
  soft-skill: コミュニケーション、リーダーシップなど
  
  # ビジネスインパクト
  cost-reduction: コスト・効率改善
  revenue: 売上・ビジネスインパクト
  developer-experience: ツーリング・生産性向上
```

**MVP の rule**: Tag は**クローズドリスト**— AI は任意の tag を追加できない。Profile 生成はこの tag リストを使ってフィルタリング・グループ化する。

---

## 5. プロジェクトメタデータ

```yaml
# projects.yaml
# オプショナル：プロジェクトレベルのコンテキスト（会社、チーム、期間）
# 複数プロジェクトを整理したいときだけ使う

projects:
  - id: proj-payment-platform
    company: 株式会社 Example
    title: Payment Platform 刷新
    period: "2022-04 to 2023-01"
    team_size: 8
    linked_experiences:
      - fact-proj-payment-platform
    
  - id: proj-internal-tools
    company: 株式会社 Example
    title: 内部ツール Framework
    period: "2021-06 to 2022-03"
    team_size: 1
    linked_experiences:
      - fact-proj-internal-tools
```

**注記**: MVP では projects.yaml は**オプショナル**。単一ユーザーの MVP では、facts が period・company をインラインで記述できれば十分。複数関連 fact をクロスリファレンスしたいときだけ projects.yaml を追加。

---

## 6. Profile（生成的、作成ではない）

```yaml
# profiles/backend-engineer.yaml
# AI が確認済み fact + ユーザー入力から生成

id: backend-engineer
target_role: バックエンドエンジニア
created_at: 2026-05-24T15:00:00Z
source_facts:
  - fact-proj-payment-platform (confirmed)
  - fact-proj-internal-tools (proposed ※ review 対象)
  - fact-skill-go (confirmed)
narrative: |
  本番環境での 10 年以上の経験を持つバックエンドエンジニア。
  高パフォーマンスシステム構築が専門。Payment インフラ、
  マイクロサービス、パフォーマンス最適化に特化。
  
  Payment platform 刷新（2022-23）ではリードを担当し、
  latency 30% 削減を実現。
  
  強み：
  - Go systems programming、gRPC、分散システム
  - パフォーマンス profiling・最適化
  - cross-team 技術リーダーシップ
```

**重要**: Profile は**派生物**。source facts への参照を含む。人間が「何が使われたか」（そして「何が proposed だが未確認か」）を見えるようにする。

---

## 7. Export（フォーマット出力）

```markdown
# resume-jp.md
# 提出用の最終フォーマット出力

## 職務経歴書

### Payment Platform 刷新 (2022年4月〜2023年1月)
株式会社Example  

**職務**: バックエンド技術リード

Payment service の刷新に従事。アーキテクチャ設計から実装まで一貫して担当。
Python から Go への移行をリード。Latency 30% 削減、コスト 15% 削減を実現。

**技術スタック**: Go, gRPC, PostgreSQL, Kubernetes

---

**出典**: fact-proj-payment-platform、fact-ach-latency-improvement
**生成日**: 2026-05-24
**注記**: 提出前に必ず review してください。proposed な fact は含まれていません。
```

---

## 8. メタデータ

```yaml
# meta.yaml

workspace:
  id: local-careervault
  created_at: 2026-05-20T00:00:00Z
  schema_version: "1.0-mvp"
  last_updated: 2026-05-24T15:00:00Z

current_state:
  confirmed_facts: 4
  proposed_facts: 2
  total_facts: 6
  generated_profiles: 2
  generated_exports: 2

version_notes:
  1.0-mvp: "Extraction + patch review ループ向けの最小限スキーマ"
```

---

## 9. Patch モデル（変更がどう流れるか）

### Patch 例 1：新しい experience を追加（proposed）

```yaml
# AI が生成して、user に review 用に送信する patch

patch_id: patch-2026-05-24-001
kind: fact_upsert
status: proposed
created_at: 2026-05-24T11:00:00Z

summary: |
  新しい experience fact を追加：内部ツール framework。
  Confidence：medium。
  User は team size と adoption metrics を確認すべき。

operations:
  - op_id: op-001
    type: upsert_fact
    target: facts/experiences.yaml
    fact_id: fact-proj-internal-tools
    new_fact:
      fact_id: fact-proj-internal-tools
      type: experience
      status: proposed
      summary: 内部ツール framework
      period: "2021-06 to 2022-03"
      company: 株式会社 Example
      description: |
        5つのチームで使用される内部 CLI tools framework を開発。
      confidence: medium
      source: conversation
      created_at: 2026-05-24T11:00:00Z
      tags:
        - backend
        - developer-tools
      source_detail: |
        User が簡潔に言及。確認が必要：
        - 正確には どの 5 チーム？
        - どのくらい長く使用されている？
        - 採用指標はあるか？
```

**User review**: Patch を読んで、判断：
- そのまま accept → status は `confirmed` に
- Reject → apply しない
- 修正要望 → user が YAML を直接編集、AI に改提案を依頼

### Patch 例 2：Fact status を更新（confirmed にマーク）

```yaml
patch_id: patch-2026-05-24-002
kind: fact_status_update
status: approved
created_at: 2026-05-24T12:00:00Z

summary: |
  User が内部ツール experience fact を確認。
  Status を proposed → confirmed に更新。

operations:
  - op_id: op-002
    type: update_fact_status
    target: facts/experiences.yaml
    fact_id: fact-proj-internal-tools
    status_change:
      from: proposed
      to: confirmed
    rationale: "User が team size = 5、18ヶ月使用を確認"
```

### Patch 例 3：Profile を再生成

```yaml
patch_id: patch-2026-05-24-003
kind: profile_regenerate
status: approved
created_at: 2026-05-24T12:30:00Z

summary: |
  Backend-engineer profile を確認済み fact から再生成。
  内部ツール experience を新たに含む。

operations:
  - op_id: op-003
    type: regenerate_profile
    target: profiles/backend-engineer.yaml
    profile_id: backend-engineer
    facts_used:
      - fact-proj-payment-platform (confirmed)
      - fact-proj-internal-tools (confirmed)
      - fact-skill-go (confirmed)
```

---

## 10. AI Extraction テンプレート

User が AI に career について話したとき：

```
[User が payment platform work について AI に話す]

→ AI が patch-001 を生成（proposed fact ）：

fact_id: fact-proj-payment-platform
type: experience
status: proposed
summary: "Payment platform 刷新"
description: |
  [User の説明を構造化抽出]
confidence: high
source: conversation
source_detail: |
  User の発言：
  - 「2022 年 4月〜2023 年 1月、payment platform 刷新をリード」
  - 「Python から Go への移行」
  - 「Latency を 200ms から 140ms に削減」
  - 「チーム約 8 人」
  
  抽出した事実（high confidence）:
  ✓ Timeline：2022-04 to 2023-01
  ✓ Action：Go migration
  ✓ Impact：latency 30% 削減
  
  抽出した事実（medium confidence）:
  ◐ Team size：「約 8 人」（user の言葉、正確ではない）
  ◐ Cost impact：アーキテクチャ変更から推定、直接述べられていない

→ User が patch を review：
  - 正しければ：approve → status は confirmed に
  - 調整が必要なら：YAML を直接編集、AI にフィードバック
  - 間違っていれば：reject → AI が clarification

→ 確認済み fact が accumulate
→ AI が確認済み fact から profile/export を生成
```

**重要**: Proposed fact はシステムに留まる、見える、review 可能。AI が uncertainty を隠さない。

---

## 11. Patch による修正パターン

### パターン 1：experiences.yaml に新 fact を append

```yaml
# 修正前
- fact_id: fact-proj-payment-platform
  type: experience
  ...

# Patch operation：append
- fact_id: fact-proj-internal-tools  # 新しい fact を追加
  type: experience
  ...
```

単純な list append。Git diff は clean。

### パターン 2：status フィールドのみ更新

```yaml
# 修正前
- fact_id: fact-proj-internal-tools
  status: proposed

# 修正後
- fact_id: fact-proj-internal-tools
  status: confirmed
```

最小限の diff。一目で review できる。

### パターン 3：Fact を拡張（impact metrics を追加）

```yaml
# 修正前
- fact_id: fact-proj-payment-platform
  type: experience
  description: ...
  # (impact セクションなし)

# 修正後
- fact_id: fact-proj-payment-platform
  type: experience
  description: ...
  impact:
    latency_improvement_pct: 30
    cost_reduction_pct: 15
```

加算的。既存フィールドへの破壊的変更なし。

---

## 12. 推奨 Workflow（MVP）

```
1. User と AI が会話
   ↓
2. AI が fact を抽出、patch を生成（operations: [fact_upsert]）
   Patch status：proposed
   すべての新 fact status：proposed
   ↓
3. User が patch を diff viewer で review
   - Before/after YAML が見える
   - AI の extraction confidence が見える
   - Approve、reject、clarification 依頼が可能
   ↓
4a. User が approve → AI が patch を apply
    → fact が workspace に write される
    → fact status は「proposed」のまま（user がまだ確認していない）
    ↓
4b. User が実際の YAML を review
    → エディタで直接編集可能
    → AI にフィードバック or 手動で確認
    ↓
5. User が fact を確認
   → patch-002 を生成（kind: fact_status_update）
   → status を proposed → confirmed に
   → User が patch-002 を approve
   → Patch が apply される
   ↓
6. AI が profile/export を生成可能に
   → Profile は確認済み fact だけを使用
   → Profile は source fact への参照を含む
   ↓
7. Loop：fact が accumulate するにつれ、profile が改善される
   User が fact を追加 → 確認 → AI が出力を改善
```

---

## 13. MVP vs. v2+ Roadmap

### MVP（今）

**必須**:
- ドメインごと単一ファイル（experiences、achievements、skills）
- Fact あたり 12 コアフィールド
- クローズドな tag 語彙
- Patch proposal + human review + apply ループ
- Flat YAML、ネストなし

**MVP に含めない**:
- relationships.yaml（tag や related_experience から推定可）
- decisions.yaml（achievements に統合）
- Fact あたりの複雑メタデータ
- Embeddings、fulltext search
- Sync、auth、multi-device
- Schema versioning（時期尚早）

### v2（MVP がループを検証した後）

**追加**:
- 明示的 decision log（fact type：decision）
- Relationship graph（relationships.yaml）
- Fact あたりのリッチなメタデータ
- Workspace versioning、patch history
- Multi-user、auth scoping
- Cloud sync 対応
- Semantic search 向け embeddings

---

## 14. Checklist：MVP Schema は準備できましたか？

- [ ] すべての fact が安定的な `fact_id` を持っている（形式：`fact-{type}-{name}`）
- [ ] すべての fact が `status` を持っている（confirmed | proposed、空白なし）
- [ ] すべての fact が `description` を持っている（人間が読める、1～5 文）
- [ ] Tag は `tags.yaml` からのみ（新規 tag の作成なし）
- [ ] Empty/optional フィールドがない（fact に情報がない場合、フィールドを作らない）
- [ ] 各ファイルが <300 行（読みやすい）
- [ ] Patch diff が <10 行 diff/operation（何が変わったかが見える）
- [ ] YAML が正しく parse される（ツールで validation）
- [ ] `grep "fact_id"` ですべての fact が返される
- [ ] 新規 fact は常に `status: proposed`（確認されるまで）
- [ ] AI extraction に `source_detail` が含まれている（説明）

---

## 15. 例：完全な最小 MVP Workspace

```
CareerVault/
├── facts/
│   ├── experiences.yaml    # 4 個の experience fact、120 行
│   ├── achievements.yaml   # 2 個の achievement fact、60 行
│   └── skills.yaml         # 3 個の skill fact、80 行
├── profiles/
│   └── backend-engineer.yaml  # 確認済み fact から生成
├── exports/
│   └── resume-jp.md        # 最終フォーマット出力
├── tags.yaml               # 15 個の tag、20 行
├── projects.yaml           # オプショナル、30 行
├── meta.yaml               # 10 行
└── README.md               # ガイド（オプショナル）
```

**合計**: 約 400 行 YAML。完全に人間が読める。Git friendly。AI friendly。

---

## 主要原則（改めて強調）

1. **最小限優先**: Fact に情報がなければフィールドを追加しない
2. **Flat 構造**: 最大 2 レベルのネスト（fact_id + フィールド）
3. **人間が読める**: grep、diff、cat がすべて自然に機能
4. **AI transparent**: source_detail が extraction を説明
5. **Status は explicit**: confirmed vs proposed に曖昧さなし
6. **1 semantic change = 1 patch**: Review が明確
7. **クローズドな tag**: 任意の語彙なし
8. **Local-first**: Cloud 前提なし、plain YAML fallback

---

## 次のステップ

1. このスキーマで example workspace を作成
2. AI extraction を patch 生成に（直接書き込みではなく）
3. Patch review UI を構築（before/after YAML 表示）
4. 完全なループをテスト：extract → review → apply → profile 生成
5. Merge conflict が少ないことを検証
6. フィードバック収集、反復

このスキーマは**2～3 週間で実装可能**で、**3～6 ヶ月の利用に十分**なように設計されています。
