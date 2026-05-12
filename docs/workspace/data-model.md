# Data Model

CareerNess の中心は「構造化されたキャリア情報」である。この文書では、主要概念の責務を固定する。もっとも重要なのは、Fact、Profile、Export を混同しないことである。

## Model Principles

- Fact は正本に近い情報である
- Profile は見せ方であり、正本ではない
- Export は外部提出用の派生物である
- Project、Tag、Role は Fact と Profile を支える補助概念である
- AI は Fact を発明してはならない
- role/title は truth の中心ではなく derived metadata に近い

## Fact

### Responsibility

Fact は、ユーザーのキャリアに関する確認可能な事実を保持する。

対象例:

- いつ、どこで、何をしていたか
- どんな責務や判断を担っていたか
- どの技術や業務に関与したか
- どの成果や改善を出したか
- 何が未確認で、何が確認済みか

Fact は、最終文章として美しくある必要はない。重要なのは、後から Profile や Export を再構成できるだけの意味が落ちていないことだ。

### Must Not

- 応募先に合わせた誇張表現を入れない
- 断定できない推測を混ぜない
- 一つの fact に複数の異なる主張を雑に詰め込まない
- profile 用の言い回しを正本扱いしない

### Example

```md
project: payment-platform-refresh
period: 2022-04 to 2023-01
role_hints:
  - tech-lead
facts:
  - Go への一部移行方針を技術選定した
  - CI 実行時間を短縮するために workflow を整理した
  - 3 team 間の調整を主導した
status: confirmed
```

## Profile

### Responsibility

Profile は、facts を特定の目的や audience に合わせて再構成した見せ方である。

たとえば:

- backend engineer 向けプロフィール
- platform / SRE 寄りプロフィール
- engineering manager 寄りプロフィール
- カジュアル面談向け要約

Profile は、同じ facts から複数作られてよい。むしろ複数存在することが正常である。

### Must Not

- Fact の代わりになる正本として扱わない
- 確認されていない実績を紛れ込ませない
- export 固有のレイアウト責務まで抱え込まない

### Example

```md
id: backend-platform-profile
target_role: Backend / Platform Engineer
source_facts:
  - payment-platform-refresh
summary:
  複数チームをまたぐ基盤改善と開発プロセス整備を強みとするバックエンドエンジニア。
emphasis:
  - technical decision making
  - CI/CD improvement
  - cross-team coordination
```

## Export

### Responsibility

Export は、提出・共有・貼り付けのために整形された最終出力である。読み手や提出先の制約に合わせた形式上の責務を持つ。

たとえば:

- 職務経歴書 markdown
- PDF 用テキスト
- スカウト返信文
- フォーム貼り付け用プレーンテキスト

### Must Not

- 正本データの置き場にならない
- profile 設計の議論を内部に抱え込まない
- 生成後に唯一の真実として扱わない

### Example

```md
format: resume-markdown
profile: backend-platform-profile
audience: hiring-manager
output:
  株式会社Xにて決済基盤刷新に従事。技術選定、CI/CD改善、複数チーム調整を担当。
```

## Project

### Responsibility

Project は、fact を束ねる文脈単位である。会社、案件、プロダクト、期間、チーム構成などをまとめ、career facts を整理しやすくする。

### Must Not

- profile の代わりに自己紹介文を持たない
- tag の羅列だけで実体を表そうとしない

### Example

```md
id: payment-platform-refresh
company: Example Corp
period: 2022-04 to 2023-01
team_size: 8
domain: payment
```

## Tag

### Responsibility

Tag は、facts や projects を横断して整理するための分類ラベルである。検索、抽出、profile 生成時の観点統一に使う。

### Must Not

- 曖昧な流行語を無制限に増やさない
- role と同じ責務を持たせない

### Example

```md
- ci-cd
- platform-improvement
- stakeholder-alignment
```

## Role

### Responsibility

Role は、その時期や project において担っていた責務の型を表す。職種名そのものだけでなく、実際に果たした役割を説明するために使う。ただし Role 自体が primary truth ではなく、facts の action / decision / impact / context から整理・導出される補助概念に近い。

### Must Not

- 肩書きだけで実態を置き換えない
- project context を無視して固定ラベル化しすぎない
- title と role を confirmed fact の代替にしない

### Example

```md
id: tech-lead
responsibilities:
  - technical decision making
  - delivery coordination
  - review and mentoring
```

## Fact / Profile / Export の分離理由

この 3 つを分離する理由は、キャリア情報の再利用性を守るためである。

- Fact に戻れれば、見せ方を何度でも変えられる
- Profile を分ければ、derived role ごとの差分を自然に管理できる
- Export を分ければ、提出形式の都合で正本が汚れない

逆にこの分離がないと、毎回文章を修正する運用になり、どれが最新でどれが正本なのかが曖昧になる。

## Human-readable Design

この data model は AI のためだけの内部表現ではない。ユーザー、人間のレビュー担当者、将来の AI coding agent が見ても意図を理解できるように、人間可読であることを優先する。

## 未確定事項

- Fact schema をどこまで厳密に型定義するかは未確定
- Project と employer history の分離粒度は未確定
- Tag を自由入力主体にするか taxonomy 管理主体にするかは未確定
