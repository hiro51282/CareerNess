# Approval Flow

CareerNess では、AI の提案と workspace の更新を同一視しない。approval flow は UX の追加機能ではなく、Local-first と user-owned data を守るための必須境界である。

## Goal

- AI proposal と user-owned truth を分離する
- 変更内容を user review 可能にする
- hallucinated fact の混入を減らす
- apply 後に何が変わったか追跡できるようにする

## Core Flow

```text
User input
  ↓
AI extraction / generation
  ↓
Patch proposal
  ↓
Validation
  ↓
User review
  ↓
Approve / reject / edit
  ↓
Apply to workspace
  ↓
History record
```

## Review Stages

### 1. Proposal Creation

AppServer は会話と workspace context から patch proposal を作る。ここではまだ workspace 正本は変更しない。

### 2. Validation

proposal をそのまま見せる前に機械的 validation を行う。

- schema violation がないか
- target path が許可領域か
- confirmed fact を unsafe に更新していないか
- delete / overwrite が高リスク扱いになっているか

### 3. User Review

ユーザーに見せる画面では、diff と意味の両方が必要である。

- どのファイルが変わるか
- どの fact が追加・変更されるか
- proposed / inferred / confirmed の状態
- AI がそう提案した理由
- どこが未確定か

### 4. Decision

最低限の選択肢は次の 3 つでよい。

- approve
- reject
- revise before apply

MVP では partial approval がなくてもよいが、operation 単位 reject は将来的に価値が高い。

### 5. Apply

apply は AI ではなく workspace update layer の責務である。承認済み patch だけが実ファイル変更になる。

## Risk-Based Approval

すべての変更を同じ危険度で扱わない。

### High Risk

- confirmed fact の意味変更
- file delete
- project の統合や分割
- profile ではなく facts 側の大量更新

### Medium Risk

- proposed fact の追加
- inferred fact の status 更新
- profile regenerate

### Low Risk

- export regenerate
- metadata update
- history append

低リスクでも silent apply は基本にしない。CareerNess は「便利だから隠す」より「見えるから信頼できる」を優先する。

## Editing During Review

review 中にユーザーが文面を直せることは重要だが、proposal と manual edit は区別して保持する。

- AI proposal
- user-edited patch
- applied result

この 3 つを混同すると、後で「AI が言ったのか、ユーザーが直したのか」が消える。

## Fact Confirmation Rule

- `proposed` は AI 提案または未確認抽出
- `inferred` は文脈推定を含む提案
- `confirmed` はユーザー確認後
- `rejected` は不採用だが履歴上参照可能

AI は `confirmed` を直接発行しない。承認操作を経て初めて `confirmed` になる。

## Failure Handling

approval 後でも apply に失敗しうる。

- lock conflict
- schema drift
- target file moved
- patch validation mismatch

この場合は「承認済みだが未適用」の状態を返し、再計算または再承認の要否を明示する。

## Prohibited Patterns

- 会話送信と同時に facts を保存する
- 承認 UI を経ずに profile / export を既定保存する
- 確認されていない AI inference を confirmed fact として apply する
- review で表示した diff と実際の apply 内容がずれる

## Open Questions

- partial approval の UX を MVP で持つかは未確定
- user-edited patch をどこまで structured に再吸収するかは未確定
