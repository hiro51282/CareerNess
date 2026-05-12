# Export Pipeline

export は CareerNess の提出物層だが、正本ではない。pipeline を設計するときも、この「完全に派生データである」という前提を崩さないことが重要である。

## Principles

- export は disposable である
- export は regenerate 可能である
- export は media-specific である
- export は facts を上書きしない

## Source Order

```text
Confirmed / proposed facts
  ↓
Selected profile view
  ↓
Export template
  ↓
Generated export
```

必要なら proposed facts を含めて下書きを作ってよいが、その場合は confirmed only ではないことを表示する。

## Pipeline Stages

### 1. Input Selection

- 対象 profile
- 参照 fact 範囲
- audience
- format

### 2. Constraint Assembly

- 文字数
- 見出し構造
- 媒体固有の制約
- 強調観点

### 3. Generation

AI または deterministic formatter が export draft を生成する。

### 4. Review

ユーザーは出力を確認し、必要なら再生成や profile 修正に戻る。

### 5. Save

保存してよいが、canonical truth 扱いしない。再生成できる前提で version を軽く持つ。

## Export Types

- resume markdown
- plain text for form input
- scout reply draft
- short interview summary

## Fact Safety Rules

- export から facts を自動逆流させない
- 魅せ方のための誇張を fact 側に戻さない
- confirmed でない内容は明示的に扱う

## Regeneration Policy

- template 変更時は再生成してよい
- profile 更新時は再生成してよい
- export の manual edit は許容してよい

manual edit した export は、それ自体を正本に格上げしない。

## Prohibited Patterns

- export だけが残り fact source が失われる運用
- 提出先向け wording を facts に書き戻す
- 媒体制約に引きずられて profile を壊す

## Open Questions

- export に source fact reference を埋め込むかは未確定
- PDF 生成責務をどこまで持つかは未確定
