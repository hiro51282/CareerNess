# Prompting Strategy

CareerNess における prompting の目的は、派手な文章を出すことではない。構造化された career facts を安定して引き出し、誤った断定を避けながら profile と export を組み立てることにある。

## 基本方針

- structure extraction を優先する
- clarification first で進める
- avoid assumption を徹底する
- fact-first generation を守る
- 出力よりも境界と根拠を重視する

## Structure Extraction

AI への指示は、自由作文よりも「何を抽出するか」を明示した方がよい。

例:

- プロジェクト名
- 期間
- role
- 責務
- 技術
- 成果
- 未確定事項

この形にすると、会話から facts を取り出す際に、どこが埋まっていてどこが空欄かを判断しやすい。

## Clarification First

情報が足りないときは、補完して書くのではなく質問する。

優先される挙動:

- 期間が曖昧なら確認する
- role と肩書きがずれていそうなら確認する
- 実装、調整、意思決定のどれが中心だったかを確認する
- 数値成果が不明なら不明のまま扱う

CareerNess では、質問が一手増えることよりも、誤った fact が固定されることの方が問題である。

## Avoid Assumption

prompt では、AI に「自然に補完してよい」と読める余地を減らす必要がある。

避けたい指示:

- いい感じにまとめて
- 足りない部分は補って
- 一般的な表現に整えて

望ましい指示:

- 事実と推測を分けて出力する
- 不明点は未確定として残す
- 根拠のない数字や肩書きは追加しない

## Fact-First Generation

profile や export を作る場合でも、まず参照元 facts を明示してから生成するのが望ましい。

基本順序:

1. source facts を列挙する
2. 不足や矛盾を確認する
3. profile を作る
4. export へ整形する

この順序を守ると、「なぜこの文章になったか」を後から説明しやすい。

## 期待する AI の出力形式

### Fact extraction 時

- 抽出した facts
- 未確定事項
- 追加で聞くべき質問

### Profile generation 時

- 参照した facts
- 想定 audience / role
- 強調した観点
- 生成したプロフィール本文

### Export generation 時

- 参照 profile
- format
- 制約
- 出力本文

## Prompt に含めるべき制約

- workspace-scoped であること
- attach されていない情報を使わないこと
- hallucinated fact を作らないこと
- patch proposal ベースで更新すること
- 人間がレビューできる形で出力すること

## なぜこの戦略か

CareerNess では、prompt は品質調整テクニックではなく、設計思想を AI に伝えるインターフェースでもある。したがって、文章品質だけでなく、責務分離と能力境界が prompt に反映されている必要がある。

## 未確定事項

- schema を固定した extraction prompt をどこまで厳密にするかは未確定
- 対話型とバッチ型の prompting をどう分けるかは未確定
