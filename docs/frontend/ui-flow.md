# UI Flow

CareerNess の UI は、派手な AI chat 体験よりも、workspace boundary と承認フローを自然に理解できることを優先する。

基本フローは次の通り。

```text
login
  ↓
workspace attach
  ↓
chat
  ↓
fact extraction
  ↓
profile generation
  ↓
export
```

## 1. Login

ユーザーはまず認証を行う。ここで重要なのは、認証が済んだだけでは AI がキャリア情報へアクセスできない点である。login はあくまで AppServer を利用するための入口であり、workspace access とは別である。

## 2. Workspace Attach

ユーザーは CareerVault を明示的に attach する。UI 上では、今どの workspace が attach されているかを常に見える状態にするべきである。

この段階で伝わるべきこと:

- AI が扱うのは attach された workspace だけである
- attach していないローカル領域は対象ではない
- workspace を切り替えれば、扱う正本も切り替わる

## 3. Chat

chat は自由会話ではなく、career facts を掘り起こすための作業面として位置づける。UI は「会話して終わり」ではなく、「どの fact が取れたか」「何が未確定か」が見える形が望ましい。

## 4. Fact Extraction

AI は会話や添付資料から facts 候補を抽出する。ここでは次が重要である。

- 確認済みと未確定を分ける
- 変更案を review 可能にする
- facts への保存を承認フロー付きにする

## 5. Profile Generation

facts がある程度揃ったら、用途別に profile を生成する。UI では、どの facts を参照してその profile を作ったかが見えるべきである。生成された profile は「正本」ではなく「見せ方」として扱う。

## 6. Export

最後に提出先や用途に応じた export を生成する。ここでは format に合わせた整形を行うが、正本の修正と混同しない UI が必要である。

## UI の重要な設計原則

- attach 状態が常にわかる
- facts / profiles / exports の区別が視覚的にもわかる
- AI が断定したのか、ユーザー確認済みなのかがわかる
- overwrite より patch review を優先する

## 避けたい UI

- chat 履歴だけが残り、facts がどこにも定着しない UI
- workspace boundary が見えない UI
- export を作った瞬間に正本が更新されたように見える UI
- AI が自由にローカル全体を見ているように誤解させる UI
