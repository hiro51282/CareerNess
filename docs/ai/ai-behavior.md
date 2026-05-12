# AI Behavior

CareerNess における AI は、便利な自動生成器ではなく、明示的な能力境界の中で動く支援者である。この文書は、AI が何をしてよく、何をしてはいけないかを定義する。

## 基本原則

- workspace-scoped である
- unrestricted shell access を持たない
- patch proposal oriented である
- hallucinated fact を作らない
- user approval を前提にする
- AI-assisted, not AI-owned を守る

## Workspace Scoped

AI が参照・更新対象にできるのは、ユーザーが attach した workspace に限定される。これは UX 上の都合ではなく、CareerNess の設計原則である。

AI がしてよいこと:

- attach 済み workspace 内の facts / profiles / exports を読む
- attach 済み workspace 内に対する patch proposal を作る
- workspace 内の不足情報を見つけて質問する

AI がしてはいけないこと:

- attach されていないフォルダを読むこと
- 端末全体を検索対象にすること
- 無関係なファイルを参照して profile を補完すること

## Unrestricted Shell Access を持たない

CareerNess は、AI に unrestricted shell access を与える前提で設計しない。

これは単に危険だからではなく、責務がぼやけるからでもある。AI が自由な shell を持つと、workspace boundary が説明上だけのものになりやすい。

## Patch Proposal Oriented

AI の更新は、原則として patch proposal として表現されるべきである。

この方針の目的:

- 何が変わるかをユーザーが読める
- facts の追加や修正に承認フローを挟める
- AI の変更責務を限定できる

AI が直接上書きを前提にすると、生成の速さは上がっても ownership は下がる。

## Hallucinated Fact 禁止

AI は、会話や既存 facts から導けない内容を、もっともらしい career fact として追加してはならない。

禁止例:

- 期間が不明なのに断定する
- 実績の数値を推測で補う
- role の比重を根拠なく決める
- 「たぶんこうだった」を confirmed fact にする

許容されること:

- 不明点を未確定として残す
- clarification question を返す
- 仮説を仮説として分離して提示する

## User Approval Flow

CareerNess における AI の変更は、ユーザー承認を経る前提で設計する。

最低限必要な段階:

1. AI が変更案を作る
2. 何を追加・修正するかを示す
3. ユーザーが承認または修正する
4. 承認後に workspace へ反映する

facts への変更は特に慎重に扱う。profile や export であっても、正本との関係が見えないまま確定保存しない。

## AI の主な責務

- structure extraction
- clarification question
- fact normalization
- profile generation
- export drafting
- inconsistency detection

## AI の禁止事項

- ユーザー未確認の経歴を fact として保存すること
- workspace 外のデータを暗黙に利用すること
- 承認を飛ばして既存 facts を上書きすること
- profile を canonical source として扱うこと
- export だけを見て facts を更新すること

## Human-readable Design との関係

AI の出力は、AI だけが理解できる内部状態に閉じるべきではない。ユーザーや別の agent が読んでも、何を根拠にどう変えたのかがわかる状態を優先する。

## 未確定事項

- patch proposal の具体フォーマットは未確定
- AI が自動実行できる変更の範囲は MVP では狭く保つ
