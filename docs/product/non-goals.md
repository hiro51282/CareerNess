# Non-Goals

この文書は、CareerNess が意図的に目指さないものを明示するためのものである。MVP 段階では特に、やらないことを先に固定した方が設計がぶれにくい。

## AI を unrestricted local agent にしない

CareerNess は、AI にユーザー端末全体への unrestricted access を与えるプロダクトではない。

やらないこと:

- ホームディレクトリ全体を自由に探索させること
- attach していないフォルダを暗黙に読ませること
- unrestricted shell access を前提にすること
- なんでも自動修正するローカル常駐 agent にすること

CareerNess の AI は、workspace-scoped であり、能力境界を持つ支援者として扱う。

## local LLM 必須プロダクトにしない

CareerNess は Local-first だが、local LLM mandatory ではない。

やらないこと:

- GPU を持つローカルマシンを前提にすること
- Ollama や各種 local model を必須要件にすること
- AI 機能の利用条件として高度なローカル推論環境を要求すること

ユーザー自身の OpenAI / ChatGPT アカウントを利用し、軽量な構成でも動くことを優先する。

## GPU hosting 基盤を運営しない

CareerNess 運営側が、大規模 GPU hosting を前提とした AI SaaS になることは目標ではない。

やらないこと:

- 独自モデル運用を前提にした GPU クラスタ設計
- 推論コスト最適化のための大規模 infra 先行投資
- モデル提供事業化

## enterprise HR platform を目指さない

CareerNess は、個人のキャリア情報管理と活用を支援するためのものであり、企業向け HR 基盤を作ることが主目的ではない。

やらないこと:

- 採用管理システム全体
- 人事評価ワークフロー全体
- 組織内人材データ統合基盤
- 大企業向け権限管理と監査を中心にした製品設計

将来的に法人利用がありえても、出発点は個人の workspace とキャリア資産である。

## huge cloud storage を持たない

CareerNess は、ユーザーのキャリア正本を大量にクラウド保管するプロダクトを目指さない。

やらないこと:

- ユーザーの職歴・草稿・export を運営側ストレージへ常時集約すること
- クラウド DB を canonical source にすること
- ファイル同期サービス自体を主製品にすること

## 文章生成だけのツールにしない

CareerNess は「それっぽい職務経歴書を一発生成する」ことだけを目的にしない。

やらないこと:

- facts を持たずに最終文章だけを保存すること
- 見た目のよいレジュメ文面を正本として扱うこと
- profile と export の責務を混ぜること

## 先回りした巨大設計をしない

CareerNess は MVP を前提にする。

やらないこと:

- Kubernetes 前提の設計
- distributed systems 前提の説明
- 将来の大規模 multi-tenant を仮定した複雑な境界設計
- 実需要が見えていない機能を先に docs へ固定すること
