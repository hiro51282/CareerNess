# Example CareerVault Workspace

このディレクトリは、CareerNess MVP Schema に準拠した**完全な example workspace**です。

架空の senior backend engineer（「Example User」）のキャリア情報を構造化した、実装可能な minimum viable example。

---

## ファイル構成

```
CareerVault/
├── facts/
│   ├── experiences.yaml    # 3 つの experience fact
│   ├── achievements.yaml   # 4 つの achievement fact
│   └── skills.yaml         # 5 つの skill fact
├── profiles/
│   └── backend-engineer.yaml  # AI 生成 profile（confirmed fact から）
├── exports/
│   └── resume-jp.md        # 最終フォーマット出力（markdown）
├── projects.yaml           # プロジェクト metadata（オプショナル）
├── tags.yaml               # Tag 語彙（クローズドリスト）
├── meta.yaml               # Workspace metadata・状態
└── README.md               # このファイル
```

**合計**: 約 600 行 YAML + 200 行 markdown

---

## 構造的な特徴

### 1. Fact Status Lifecycle

このディレクトリには、`status` の異なる fact が混在しています：

- **confirmed** (7個): ユーザーが明示的に確認した事実
  - Payment platform experience
  - Latency・CI optimization achievements
  - Go・Microservices・Performance optimization skills
  - Leadership・Mentoring achievements

- **proposed** (3個): AI が抽出したが、ユーザーがまだ確認していない
  - Internal tools experience
  - Team coordination achievement
  - Infra optimization experience

このため、`meta.yaml` では以下のように状態を追跡しています：

```yaml
current_state:
  total_facts: 10
  confirmed_facts: 7
  proposed_facts: 3
```

**何を学べるか**: Patch proposal → user review → status update ループの準備。

### 2. Confidence Tracking

各 fact に confidence レベルがあります：

- `high`: 複数の information source で確認済み
- `medium`: 主に user の言葉だけ、詳細未確認
- `low`: (このファイルには含まれていない)

Example: 
```yaml
# Confirmed, high confidence
- fact_id: fact-ach-latency-improvement
  confidence: high
  source_detail: "Datadog logs + team retrospective で確認"

# Proposed, medium confidence
- fact_id: fact-proj-internal-tools
  confidence: medium
  source_detail: "User が概要を述べたが、詳細は未確認"
```

**何を学べるか**: AI extraction では、自動で全て `confirmed` にしない。必ず confidence を明示し、不確定性を追跡する。

### 3. Source Tracking

すべての fact が `source` と `source_detail` を持ちます：

```yaml
source: conversation  # conversation | manual | import
source_detail: |
  2026-05-23 の user との会話から抽出。
  Latency metrics は Datadog で確認。
```

**何を学べるか**: Fact が「どこから来たのか」を追跡可能に。Review・修正時に根拠が見える。

### 4. Tags は Closed Vocabulary

`tags.yaml` で 20 個の tag を定義し、fact はその中から選びます。

```yaml
# facts/skills.yaml
- fact_id: fact-skill-go
  tags:
    - language
    - backend
    - systems-programming
```

**何を学べるか**: AI が任意の tag を作らない。Profile 生成が consistent。

### 5. Profile は Generated, Not Authored

`profiles/backend-engineer.yaml` は AI が生成したもの（人間が手書きしていない）。

```yaml
# Generated profile
generated_by: AI extraction agent
source_facts:
  - fact-proj-payment-platform (confirmed)
  - fact-skill-go (confirmed)
  # ... など
```

Profile は source facts への参照を持つので、「何が使われたのか」が見える。

**何を学べるか**: Profile は facts から生成されるもの。Fact が正本。

---

## 使用方法

### 1. このファイルを参考にしながら、自分の CareerVault を構築

```bash
$ cp -r implementation/examples/CareerVault ~/my-careervault
$ cd ~/my-careervault

# facts/*.yaml を編集
# 自分の experience、achievement、skills を追加
```

### 2. Schema validation（今後）

```bash
# スキーマに準拠しているか check
$ careervault validate ~/my-careervault
```

### 3. AI extraction テスト（次フェーズ）

```bash
# AI に自分の career を話す
# AI が patch-001（fact upsert）を生成
# User が review → approve → apply
```

### 4. Profile 生成（次フェーズ）

```bash
# confirmed facts から profile を再生成
$ careervault generate-profile backend-engineer
```

---

## 学ぶべきポイント

このファイルから以下を学べます：

1. **Flat structure**: ネストを避け、yaml が短くなるように設計
2. **Human readability**: Grep、diff、cat で natural に読める
3. **AI transparency**: source_detail、confidence で uncertainty を明示
4. **Semantic granularity**: 1 fact = 1 autonomous information
5. **Status tracking**: proposed → confirmed のライフサイクル
6. **Closed tag vocabulary**: consistency を優先
7. **Traceability**: source を全て記録、修正可能

---

## 実装のチェックリスト

このファイルを実装ガイドにして、以下を確認しながら進める：

- [ ] facts/experiences.yaml を period・company・description で構造化
- [ ] facts/achievements.yaml で impact を定量化（%、metrics）
- [ ] facts/skills.yaml で proficiency レベルを明示
- [ ] すべての fact が status（confirmed | proposed）を持つ
- [ ] すべての fact が source_detail（根拠）を持つ
- [ ] Tags が tags.yaml に列挙されたもののみ
- [ ] Profiles は confirmed facts のみから生成
- [ ] meta.yaml で workspace 全体の状態を追跡可能

---

## 次のステップ

1. **AI extraction agent** を構築：会話 → patch-001 生成
2. **Patch reviewer UI** を実装：before/after YAML 表示
3. **Patch apply logic** を実装：fact_upsert operation
4. **Profile generator** を実装：template + facts → narrative

---

## 注記

- このファイルは example です。実装時はファイルをコピーして、自分のデータで置き換えてください。
- YAML の形式・フィールド名は MVP schema に厳密に従ってください。
- 疑問がある場合は、`docs/workspace/careervault-mvp-schema.md` を参照。
