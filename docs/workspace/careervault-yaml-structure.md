# CareerVault YAML 構造設計

このドキュメントは、CareerVault の YAML 構造を、AI readability、human reviewability、patch diff の見やすさ、merge conflict 最小化、local-first 運用、将来の sync/auth に対応した形で再設計・定義するものである。

**Key Principle**: YAML は単なる保存形式ではなく、AI・Patch・UI・Review をつなぐ中心モデル。

---

## 1. ディレクトリ構成（推奨）

```
CareerVault/
├── facts/
│   ├── experiences/          # 職歴・プロジェクト経験
│   ├── achievements/         # 成果・インパクト
│   ├── skills/               # 技術スタック・領域
│   ├── decisions/            # 判断・意思決定
│   └── relationships.yaml    # Entity relationships (事実間の関連性)
│
├── profiles/
│   ├── backend-engineer/
│   │   ├── profile.yaml
│   │   └── narrative.md      # 人間向け説明（生成後）
│   ├── platform-engineer/
│   │   ├── profile.yaml
│   │   └── narrative.md
│   └── _index.yaml           # 全 profile の metadata
│
├── exports/
│   ├── resumes/
│   │   ├── standard-jp.md
│   │   ├── standard-en.md
│   │   └── _metadata.yaml
│   ├── scout-replies/
│   │   └── _templates.yaml
│   └── _generated-at.yaml
│
├── projects/
│   ├── project-alpha.yaml    # Project metadata (company, period, team)
│   ├── project-beta.yaml
│   └── _index.yaml
│
├── taxonomy/
│   ├── skills.yaml           # Skill taxonomy (categories, levels)
│   ├── roles.yaml            # Role patterns (responsibilities, patterns)
│   ├── tags.yaml             # Custom tags/labels
│   └── domains.yaml          # Industry/domain classification
│
├── inbox/
│   ├── _extraction-session-{date}.yaml    # Extraction session snapshots
│   └── raw-notes.md          # 未整理メモ（temporary）
│
├── attachments/              # Optional: original resume, old materials
│
└── meta/
    ├── workspace.yaml        # Workspace metadata, schema version
    ├── patch-history.yaml    # Patch lifecycle records (not YAML patches)
    └── .schema-version       # Current schema version (e.g., 1.0)
```

### ディレクトリ責務の詳細化

#### `facts/experiences/`

職歴・プロジェクト経験の事実を保存。複数構造案：

**案A: プロジェクト中心**（推奨）
```
facts/experiences/
├── {project-id}-overview.yaml          # Project context
├── {project-id}-role-and-responsibility.yaml   # Role + responsibilities
├── {project-id}-decisions.yaml         # Decisions made
├── {project-id}-impact.yaml            # Outcomes / KPI improvements
└── {project-id}-timeline.yaml          # Key milestones
```

**案B: 時系列中心**
```
facts/experiences/
├── 2026-05-current.yaml                # Current role
├── 2025-01-2026-04.yaml                # Previous role
└── 2023-06-2024-12.yaml                # Earlier role
```

**推奨理由**:
- Patch single semantic change → A 案（プロジェクト粒度で分ける）が自然
- AI extraction も project ごとに proposal を生成しやすい
- Merge conflict が少ない（プロジェクト間が直交）

#### `facts/achievements/`

成果・インパクト・改善の事実。プロジェクト紐付けの有無で分けることもできるが、MVP では単一ファイル。

```
facts/achievements/
├── achievements.yaml         # All achievement facts (structured list)
└── _schema.yaml              # Local schema reference（optional）
```

#### `facts/skills/`

技術スタック・言語・フレームワーク・領域知識の事実。

```
facts/skills/
├── skills.yaml               # Structured list of skills with context
└── _skill-graph.yaml         # Skill relationships（optional, future）
```

#### `facts/decisions/`

技術的・組織的判断の事実。オプショナルだが、「意思決定」を explicit に扱うと AI が context を保持しやすい。

```
facts/decisions/
├── decisions.yaml            # Decision log with rationale, outcomes
```

#### `facts/relationships.yaml`

Entity 間の関連性を明示的に管理（optional、future direction）。例：

```yaml
relationships:
  - source: fact-{project-id}
    target: fact-{achievement-id}
    type: contributed_to    # 貢献している、実現した、など
    confidence: confirmed
  - source: fact-skill-{lang}
    target: fact-{project-id}
    type: used_in
```

**MVP では optional**。Fact が standalone で読める設計を優先。

---

## 2. Fact YAML スキーマ（詳細）

### 2.1 Fact Level の構造

すべての fact が持つべき最小スキーマ：

```yaml
# {fact-id}.yaml または facts/{category}/{fact-id}.yaml

fact_id: fact-project-alpha-overview
type: experience_project           # experience_project, achievement, skill, decision
workspace_id: local-careervault    # Workspace identifier（future sync prep）
version: 1                         # Schema version tracking

metadata:
  status: confirmed                # confirmed | proposed | inferred
  confidence: high                 # high | medium | low
  created_at: 2026-05-23T10:00:00+09:00
  updated_at: 2026-05-23T15:30:00+09:00
  source:
    type: conversation             # conversation | import | inference | manual
    session_id: sess_abc123        # Session where created/updated
    reference: "from user statement on 2026-05-23"
  tags:
    - backend
    - platform
    - ci-cd

# Content below is type-specific
```

### 2.2 Experience Project Fact

```yaml
fact_id: fact-project-payment-platform
type: experience_project
metadata:
  status: confirmed
  confidence: high
  created_at: 2026-05-23T10:00:00+09:00
  source:
    type: conversation
    session_id: sess_abc123

# Project Context
project_id: payment-platform-refresh
company: Example Corp
department: Backend Platform
period:
  start: 2022-04
  end: 2023-01
  duration_months: 10
  status: completed                # completed | ongoing | paused
team_size: 8
team_composition: "5 backend, 2 frontend, 1 infra"

# Role & Responsibility
role_hints:                        # Not definitive, hints for profile generation
  - tech-lead
  - senior-engineer
responsibilities:
  - Technical decision making for payment architecture
  - CI/CD workflow optimization
  - Cross-team coordination
  - Mentoring junior engineers
primary_focus: Backend infrastructure and platform improvement

# Key Activities (what actually happened)
activities:
  - action: Selected Go for payment service migration
    rationale: Performance and type safety for financial domain
    outcome: Successfully migrated 80% of payment logic
    impact: Reduced latency by 30%
  - action: Reorganized CI workflow to reduce build time
    outcome: Consolidated 15 separate jobs into 5 parallel stages
    impact: Reduced CI time from 45min to 12min
  - action: Led technical discussions across 3 teams on unified API spec
    outcome: Defined common interface contract
    impact: Enabled smooth integration, reduced incidents

# Technical Stack (what was involved)
tech_stack:
  languages: [Go, Python, TypeScript]
  frameworks: [gRPC, Kafka, React]
  databases: [PostgreSQL, Redis]
  platforms: [AWS, Kubernetes]

# Impact Metrics
impact:
  latency_improvement_pct: 30
  build_time_reduction_pct: 73
  team_velocity_improvement: "estimated +20%"
  incident_reduction: "from ~5 to ~2 per month"
  scope: team + 2 dependent teams

# Relationships to other facts
related_facts:
  - fact-id: fact-project-payment-platform-decisions
    type: decision_log
  - fact-id: fact-achievement-payment-latency-improvement
    type: outcome_of
  - fact-id: fact-skill-go
    type: applied

# Rationale & Confidence Notes
rationale: >
  Extracted from user's detailed account of payment platform refresh.
  All core facts confirmed. Impact metrics partly from retrospective
  review (high confidence on build time, medium on latency estimation).
confidence_notes: >
  High confidence on activities and tech stack.
  Medium confidence on team size (user said "around 8").
  Impact metrics from combined sources: prod logs (latency),
  CI history (build time), team feedback (velocity).

# Uncertainty Tracking
uncertain_fields:
  team_size: "User said 'around 8', exact number unclear"
  exact_timeline: "Start/end in terms of fiscal quarters, exact dates unclear"
```

### 2.3 Achievement Fact

```yaml
fact_id: fact-achievement-payment-latency-improvement
type: achievement
metadata:
  status: confirmed
  confidence: high
  created_at: 2026-05-23T11:00:00+09:00
  source:
    type: conversation
    session_id: sess_abc123

# What was achieved
title: Payment latency improvement (30% reduction)
description: >
  Led architecture redesign and implementation optimization in payment
  service, reducing 99p latency from 200ms to 140ms.

# Impact Measurement
impact:
  primary_metric: latency
  metric_type: performance
  before: "200ms (p99)"
  after: "140ms (p99)"
  improvement_pct: 30
  duration: 4 months
  timeframe_confirmed: true

# Context & Attribution
context:
  related_project: fact-project-payment-platform
  related_role_hints: [backend-lead, platform-engineer]
  team_contribution: "Led design and coordination, 60% implementation"
  dependencies: "Required buy-in from 2 dependent services"

# Technical Details
technical_approach:
  - Implemented connection pooling for database
  - Optimized hot path with profiling data
  - Introduced caching layer for idempotent operations
  - Migrated to gRPC for inter-service communication

# Business/Product Impact
business_impact:
  user_experience: "Faster payment confirmation, reduced customer frustration"
  operational_cost: "Reduced infrastructure cost by ~15% (fewer instances needed)"
  reliability: "Fewer timeout-related incidents"

# Verification
verification:
  method: production_metrics
  data_source: Datadog monitoring
  confirmed_by: team lead review
  external_validation: customer report of improved experience

related_facts:
  - fact-id: fact-project-payment-platform
    type: part_of_project
  - fact-id: fact-achievement-team-coordination
    type: related
```

### 2.4 Skill Fact

```yaml
fact_id: fact-skill-go-systems
type: skill
metadata:
  status: confirmed
  confidence: high
  source:
    type: conversation
    session_id: sess_abc123

# Skill Identification
skill_name: Go (systems/backend programming)
category: Programming Language
proficiency: advanced
level_evidence: 10+ years production experience, designed concurrent systems

# Context (where/how used)
used_in_projects:
  - fact-project-payment-platform
  - fact-project-infra-optimization
  - fact-project-internal-tools

depth_breadth:
  breadth: "Full stack (from CLI tools to distributed systems)"
  depth: "Concurrency patterns, memory model, performance profiling"
  specialization: "High-performance backend systems, gRPC, Kubernetes"

# Sub-skills / Related Concepts
sub_skills:
  - Goroutines and channel patterns
  - Memory profiling and optimization
  - gRPC protocol and implementations
  - Testing (unit, integration, benchmarking)

# Key Learnings / Specific Knowledge
key_learnings:
  - "Goroutine scheduling and avoiding common pitfalls (context leaks)"
  - "Database connection pooling patterns in Go"
  - "Zero-copy techniques and safe concurrent data structures"

# Teaching / Mentoring
has_taught: true
teaching_context: "Mentored 3 junior engineers on Go concurrency patterns"
```

### 2.5 Decision Fact

```yaml
fact_id: fact-decision-go-migration
type: decision
metadata:
  status: confirmed
  confidence: high
  source:
    type: conversation
    session_id: sess_abc123

# Decision Details
decision_title: "Migrate payment service to Go"
context:
  problem: Python service was bottleneck for payment processing (high CPU, memory)
  constraints:
    - Must maintain backward compatibility
    - 2 month timeline
    - Team had no prior Go experience
  stakeholders: [payments-team, platform-team, CTO]

# Decision Process
rationale: >
  Go selected for: type safety (vs Python), performance (vs Node),
  easier deployment (single binary, vs Java), good ecosystem for
  networking libraries (gRPC, Kafka clients).
reasoning_type: technical_decision
involved_parties:
  - yourself: "Led technical evaluation"
  - tech_committee: "Reviewed and approved"
  - team: "Provided input on learning curve"

# Outcome
outcome: Successfully migrated 80% of payment logic in 3 months
result_quality: exceeded_expectations
learnings:
  - Go learning curve was faster than expected
  - gRPC integration reduced service communication latency
  - Concurrency model required new debugging patterns
```

---

## 3. 構造設計の原則

### 3.1 ID Strategy

**Stable Identifier の重要性**: Fact は長く参照される。ID が変わると、patch / profile / relationship が破損する。

推奨スキーム：

```
fact-{category}-{semantic-id}

例：
  fact-project-payment-platform
  fact-achievement-latency-improvement
  fact-skill-go-systems
  fact-decision-go-migration
```

**NOT**:
- Timestamp-based ID（semantic meaning が失われる）
- Sequential ID（移動・削除時に破損する）
- Auto-generated UUID（human review 困難）

### 3.2 Status Lifecycle

```
proposed → inferred → confirmed → archived
```

- `proposed`: AI が提案、ユーザーが未確認
- `inferred`: AI が推論、ユーザーが「多分そう」まで承認
- `confirmed`: ユーザーが明示的に確認済み
- `archived`: 古い、不要になったが削除しない（history追跡用）

**Rule**: `confirmed` に遷移するには明示的な ユーザー承認が必須。AI だけでは遷移しない。

### 3.3 Confidence Tracking

```yaml
confidence: high | medium | low
confidence_notes: string
uncertain_fields:
  field_name: explanation
```

**Why**: Fact が完全でない場合も多い。「どこが不確定か」を明示することで：
- AI が推理しすぎない
- Profile 生成時に慎重に扱える
- ユーザーが追加情報を提供しやすい

### 3.4 Source Tracking

```yaml
source:
  type: conversation | import | inference | manual
  session_id: sess_abc123      # Which conversation
  reference: string             # "from user statement...", "extracted from resume..."
  user_statement: |             # Optional: actual quote from user
    "I led payment platform refresh from Apr 2022 to Jan 2023..."
```

**Why**:
- Fact の根拠を追跡可能に
- Patch review 時に「どこから来たのか」が見える
- 誤解・変更があった場合に遡れる

### 3.5 Human Readability

```yaml
# Good: semantic, greppable, reviewable
fact_id: fact-project-payment-platform
description: >
  Redesigned and optimized payment service,
  reducing latency by 30% and cost by 15%.
activities:
  - action: Selected Go for service migration
    outcome: 80% of logic migrated successfully
```

**NOT**:
```yaml
# Bad: nested, implicit, AI-only
data:
  p: [{a: "migration", c: 80, d: "success"}, ...]
```

### 3.6 Patch Diff Friendliness

**原則**: YAML をそのまま diff で見た時に、何が変わったかが自明であること。

Good:
```diff
- confidence: medium
+ confidence: high
+ confidence_notes: "Confirmed by team retrospective review"
  updated_at: 2026-05-23T10:00:00+09:00
- updated_at: 2026-05-23T09:00:00+09:00
+ updated_at: 2026-05-23T15:30:00+09:00
```

**NOT**:
```diff
- data: [{...1000 lines...}]
+ data: [{...1000 lines, 1 field changed...}]
```

→ 巨大ネストは避ける。flat 中心。

### 3.7 Tag Strategy

```yaml
tags:
  - backend           # 技術領域
  - platform          # プロダクト・チーム領域
  - ci-cd             # 機能領域
  - system-design     # Skill/capability 領域
  - leadership        # Role 領域
  - 3p-coordination   # Soft skill 領域
```

**Taxonomy は `taxonomy/tags.yaml` で集中管理**:

```yaml
# taxonomy/tags.yaml
tags:
  backend:
    description: Backend development, systems
    aliases: [backend-engineering]
  platform:
    description: Platform / Infrastructure work
  ci-cd:
    description: CI/CD pipeline, build systems
  # ... more
```

**Why**:
- Profile 生成時に consistent に扱える
- AI が新規 tag を勝手に発明しない
- Search/filter が uniform に

---

## 4. Patch との相性設計

### 4.1 Semantic Change の粒度

**1 patch = 1 semantic change** に対応しやすい構造：

```yaml
# Patch 1: New fact proposed
patch_id: patch-001
operations:
  - type: upsert_fact
    target: facts/experiences/payment-platform-overview.yaml
    entity_id: fact-project-payment-platform
    change:
      before: null
      after: {full new fact}
    status: proposed

# Patch 2: Status upgraded to confirmed
patch_id: patch-002
operations:
  - type: mark_fact_status
    target: facts/experiences/payment-platform-overview.yaml
    entity_id: fact-project-payment-platform
    change:
      before: proposed
      after: confirmed
    status: approved

# Patch 3: Achievement fact added (separate semantic change)
patch_id: patch-003
operations:
  - type: upsert_fact
    target: facts/achievements/achievements.yaml
    entity_id: fact-achievement-latency-improvement
    change:
      before: null
      after: {new fact}
    status: proposed
```

**Why**:
- 1 approval = 1 clear decision
- Rollback が細粒度
- Conflict が少ない

### 4.2 File-per-fact vs. Multi-fact Files

**2 つの戦略：**

#### 戦略 A: File-per-fact（推奨）

```
facts/experiences/
├── payment-platform-overview.yaml     # fact-project-payment-platform
├── payment-platform-decisions.yaml    # fact-project-payment-platform-decisions
└── ...
```

**利点**:
- Patch が単一 file
- Merge conflict が少ない
- Diff が見やすい
- History が granular

#### 戦略 B: Category-per-file

```
facts/experiences/
├── all-experiences.yaml
  - - fact_id: fact-project-payment-platform
        ...
    - fact_id: fact-project-alpha
        ...
```

**利点**:
- Directory が flatten（nav が簡単）
- Category view が単一 file

**欠点**:
- Patch が大きくなりやすい
- Merge conflict が多い
- 単一ファイルへの複数パッチは risk

**推奨**: 戦略 A。MVP ではシンプルさを優先し、file-per-fact が吉。

### 4.3 Version Tracking

**Option 1: Per-fact versioning**

```yaml
fact_id: fact-project-payment-platform
version: 3                    # Incremented on meaningful change
status: confirmed
```

**Option 2: Workspace-level revision**

```
meta/workspace.yaml
current_revision: 42
last_patch_id: patch-2026-05-23-015
```

**推奨**: Option 2。Workspace-level を source of truth にして、patch ごとに increment。

---

## 5. AI Extraction 向けベストプラクティス

### 5.1 Clarity Rules

```yaml
# Clear: Each claim is separate, status is explicit
- action: Selected Go for migration
  rationale: Type safety, performance
  outcome: 80% migrated
  impact: 30% latency reduction
  status: confirmed           # ← Explicit

# NOT: Vague, combined claims
- "Selected Go which helped with migration and latency"
  status: proposed            # Unclear what's confirmed
```

### 5.2 Fact Extraction Query Template

When AI extracts from conversation:

```
Extracted from user input:
  Timeline: 2022-04 to 2023-01 ✓ confirmed
  Project: payment platform refresh ✓ confirmed
  Team size: ~8 people ◐ medium confidence
  Impact: latency improved 30% ◐ estimated from conversation
  
Proposed fact status: proposed (user hasn't reviewed yet)
AI confidence: high on activities, medium on metrics
```

### 5.3 Uncertain Field Handling

**Never do**:
```yaml
# BAD: AI filled in gaps
company_size: 500              # Not mentioned by user
role_title: Backend Lead       # User didn't say this
```

**Do**:
```yaml
# GOOD: Mark as uncertain
role_hints:
  - tech-lead                  # Inferred from responsibilities
uncertain_fields:
  company_size: "Not mentioned, inferred as ~500 from domain knowledge"
  exact_title: "User described as 'technical lead role', exact title unclear"
```

### 5.4 Conflicting Information Handling

```yaml
# When AI detects inconsistency:
impact:
  latency_improvement: "30% (from prod logs)"
  latency_improvement_alternative: "25% (from user estimate)"
  source_of_discrepancy: "Prod logs show consistent improvement; user estimate was rough"
  confirmed_value: 30
```

---

## 6. Local-First Knowledge Vault 設計

### 6.1 Human Greppability

```bash
$ grep -r "latency" facts/
facts/experiences/payment-platform-overview.yaml:    impact: Reduced latency by 30%
facts/achievements/achievements.yaml:title: Payment latency improvement (30% reduction)

$ grep -r "Go" taxonomy/skills.yaml
skills.yaml:- Go (systems/backend programming)

$ grep -r "2022-04" facts/
facts/experiences/payment-platform-overview.yaml:  start: 2022-04
```

**Why**: Backup, search, revision, context — すべて shell commands で可能に。

### 6.2 Plain Text Recovery

WorkspaceDB が broken な場合でも、YAML を手で読んで復旧可能。

```yaml
# 人間が読んで理解できる粒度
fact_id: fact-achievement-latency-improvement
title: Payment latency improvement (30% reduction)
description: >
  Reduced 99p latency from 200ms to 140ms
  in payment service through architecture
  redesign and optimization.
```

### 6.3 No Implicit Dependencies

**NOT**:
```yaml
# reference.yaml
- id: 123
  text: "payment platform"

# fact.yaml
- name: ref: 123           # Implicit reference
```

**DO**:
```yaml
# fact.yaml
related_facts:
  - fact_id: fact-project-payment-platform
    type: part_of_project
```

Explicit が local-first を支える。

### 6.4 Minimal Structured Metadata

```yaml
metadata:
  status: confirmed
  created_at: 2026-05-23T10:00:00+09:00
  updated_at: 2026-05-23T15:30:00+09:00
  source:
    type: conversation
    session_id: sess_abc123
```

**NOT**:
```yaml
metadata:
  version: 42
  hash: abc123def456
  last_ai_model: claude-opus-4-7
  edit_depth: 3
  # ... excessive tracking
```

---

## 7. 将来の Sync / Auth を考慮した設計

### 7.1 Workspace ID Stability

```yaml
workspace_id: local-careervault        # Fixed per workspace
fact_id: fact-project-payment-platform # Stable across syncs
```

**Why**: Cloud sync を導入した時、device A ↔ B で ID が変わらない必要。

### 7.2 Session ID Tracking

```yaml
source:
  session_id: sess_abc123    # Which UI session created/modified
  auth_user: hiro51282@gmail.com  # Future: track who made the change
```

**Why**: Multi-device や shared workspace でも「誰が、いつ」わかる。

### 7.3 Conflict Markers（Future）

```yaml
# 将来的に merge conflict が発生した時
impact:
  latency_improvement_pct:
    - "30"         # Device A
    - "25"         # Device B
  _conflict_source: "Concurrent modification"
```

MVP では不要だが、schema に余裕を持たせておく。

### 7.4 Auth Scoping（Future）

```yaml
metadata:
  access_level: private                # Future: private | shared | public
  shared_with: []                      # Future: list of user IDs
```

MVP では単一ユーザーなので空でいい。

### 7.5 Schema Versioning

```yaml
# meta/workspace.yaml
current_schema_version: "1.0"
schema_migration_history:
  - from: "0.9"
    to: "1.0"
    date: 2026-05-23
    migration_script: "scripts/migrate-0-9-to-1-0.py"
```

**Why**: Schema が進化した時、古い workspace の migration path が必要。

---

## 8. Merge Conflict 最小化戦略

### 8.1 File Granularity

**避ける**:
```
facts/all-facts.yaml          # 1000+ lines, frequent changes
```

**推奨**:
```
facts/experiences/
├── project-1.yaml
├── project-2.yaml
├── project-3.yaml
```

同時編集が同じ file に hit しにくい。

### 8.2 Tag は集中管理

```
# Avoid scattered tags in each fact
facts/project-1.yaml:
  tags: [backend, ci-cd]
facts/project-2.yaml:
  tags: [backend, ci-cd]

# Instead: centralized
taxonomy/tags.yaml
  backend: { ... }
  ci-cd: { ... }
```

参照だけなら conflict なし。

### 8.3 Metadata フィールドの位置

```yaml
# GOOD: metadata first（多くのパッチが content に集中するので）
fact_id: ...
type: ...
metadata:
  ...

content:
  # 実際の内容
```

### 8.4 List Append Strategy

```yaml
# NOT: ordered list with indices
activities:
  - index: 1
    action: ...
  - index: 2
    action: ...

# DO: unordered, ID-based
activities:
  - activity_id: activity-001
    action: ...
  - activity_id: activity-002
    action: ...
```

Merge が自動化しやすい。

---

## 9. アンチパターン

### 9.1 Deeply Nested Structures

```yaml
# NOT
data:
  experiences:
    projects:
      - details:
          context:
            metadata:
              - values

# DO
fact_id: fact-project-payment-platform
experience_type: project
details:
  ...
metadata:
  ...
```

Flat + moderate depth が読みやすく、diff がクリア。

### 9.2 Implicit Status

```yaml
# NOT: status が implied
activities:
  - "Selected Go for migration"      # Confirmed? Proposed?

# DO: explicit
activities:
  - action: Selected Go for migration
    status: confirmed
    confidence: high
```

### 9.3 Mutable Derived Fields

```yaml
# NOT: derived field が原本のように扱われる
fact_id: fact-project-alpha
title_auto_generated: "Alpha Project (2022-04 to 2023-01)"
# → 後で period が変わると title が stale に

# DO: derived は computed only
fact_id: fact-project-alpha
period:
  start: 2022-04
  end: 2023-01
# → UI/AI が computed display を作る
```

### 9.4 Comment-as-Metadata

```yaml
# NOT
fact_id: fact-project-alpha
# TODO: Confirm team size with manager
# This fact needs review before profile generation
activities: ...

# DO
fact_id: fact-project-alpha
metadata:
  review_needed: true
  review_reason: "Team size needs confirmation"
  review_scheduled_for: "2026-05-30"
activities: ...
```

YAML comment は version control で落ちやすい。

### 9.5 Single Large File for Categories

```yaml
# NOT: facts/all-experiences.yaml (1000+ lines)

# DO: facts/experiences/project-{id}.yaml (50-150 lines each)
```

### 9.6 Unreferenced Relationships

```yaml
# NOT: implicit relationship
fact_id: fact-project-alpha
related_project_name: "payment platform"  # String reference, fragile

# DO: explicit ID reference
fact_id: fact-project-alpha
related_facts:
  - fact_id: fact-project-payment-platform
    type: part_of_initiative
```

ID reference なら AI が validate できる。

---

## 10. MVP で不要な過剰設計

### 10.1 スキップしてよい：

- **Embeddings**: `embeddings/` は future。YAML text search で充分。
- **Full graph database model**: `relationships.yaml` もMVP では optional。
- **Approval workflow per field**: Fact 単位の approval でいい。
- **Multi-language support in schema**: 日本語で start。Localization は later。
- **Encryption at rest in YAML**: Local filesystem trust で start。Cloud sync 時に add。
- **Complex schema validation**: JSON Schema や GraphQL まで不要。Runtime Python validation で start。
- **Change reconciliation logic**: Merge conflict automatic resolve は out of scope。
- **Patch signing**: Trust user's own workspace for MVP。Auth 導入後に add。
- **Workspace federation**: Single workspace per user。Multi-workspace は v2。

### 10.2 MVP で必須：

- **Stable IDs**: `fact-{category}-{name}`
- **Status tracking**: `confirmed | proposed | inferred`
- **Source tracking**: どこから来たのか
- **Timestamp**: created_at, updated_at
- **Patch proposal model**: 承認フロー
- **File-per-fact granularity**: Diff が見やすい
- **Tag taxonomy centralization**: 散らばらない
- **Relationship references**: ID-based, explicit

---

## 11. Migration Strategy（既存 workspace）

もし既存 workspace があれば：

1. **Audit current structure**: `facts/`, `profiles/`, `exports/` の既存ファイルを確認
2. **Create migration script**: Python script で YAML を新スキーマに normalize
3. **Preserve history**: Old facts に `archived: true` を加えて保持（削除しない）
4. **Incremental adoption**: 新しい fact は新スキーマ、古い fact は gradual migrate
5. **Test patches**: New schema で patch が正しく diff / apply できるか確認

---

## 12. スキーマバージョニング

```yaml
# meta/workspace.yaml
workspace:
  id: local-careervault
  schema_version: "1.0"
  created_at: 2026-05-23
  last_updated: 2026-05-23

schema_versions:
  "1.0":
    date: 2026-05-23
    breaking_changes: null  # First version
    features:
      - stable fact IDs
      - status lifecycle (proposed/inferred/confirmed)
      - source tracking
      - confidence levels
      - relationship references
  # "1.1": future...
```

---

## まとめ: Design Checklist

Your YAML structure should:

- [ ] ID は stable かつ semantic
- [ ] Status は explicit かつ lifecycle が clear
- [ ] Source / confidence が fact ごと
- [ ] File granularity は small（50-150 lines / file）
- [ ] Tag は taxonomy 集中管理
- [ ] Relationship は ID reference
- [ ] Metadata は structural, comment not
- [ ] Nested depth は 2-3 level max
- [ ] Patch diff は見て何が変わったか明自
- [ ] grep / find で human search 可能
- [ ] Future sync/auth を考慮（workspace_id, session_id）
- [ ] アンチパターンを避けている

---

## 参考: 実装チェックリスト

1. `facts/experiences/`, `facts/achievements/`, `facts/skills/`, `facts/decisions/` を create
2. Fact schema の TypeScript/Go interface を定義
3. YAML → Object のデserialize logic を実装
4. Patch apply logic で schema validation を implement
5. AI extraction で proposed/inferred status を respect
6. Profile / export 生成で fact reference を追跡
7. UI で fact diff を見やすく表示
8. Test で merge conflict が出にくいことを confirm
9. Migration script（既存 workspace 対応）
10. Documentation を workspace に置く（README）
