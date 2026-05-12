# Profile Generation Flow

profile は保存してよいが正本ではない。CareerNess における profile generation は、「事実から見せ方を作る」ための派生生成であり、facts の代替生成ではない。

## Principles

- profile は role-oriented view である
- editable である
- regenerate 可能である
- export より上流だが fact より下流である

## Source Order

```text
facts
  ↓
selection / filtering
  ↓
profile synthesis
  ↓
user review
  ↓
save as derived profile
```

## Inputs

- confirmed facts
- 必要に応じて proposed / inferred facts
- target orientation
- emphasis rules

orientation 例:

- platform
- DevOps
- backend
- tech lead

## Generation Rules

- 肩書だけから profile を作らない
- action / decision / impact / context から軸を作る
- どの facts を強調したか追える方がよい
- 未確認の内容は profile 上でもそれと分かるようにする

## Save Model

profile は workspace に保存してよい。ただし扱いは派生データである。

- overwrite 可能
- regenerate 可能
- manual edit 可能

manual edit したからといって source of truth にはならない。

## Suggested Contents

- profile id
- target orientation
- source fact references
- summary
- emphasis themes
- optional narrative blocks
- generated_at / updated_at

## Relationship To Export

export は profile を入力としてよいが、profile を通さず facts から直接 export を作る場合もありうる。重要なのは、どちらの場合も export が truth にならないことだ。

## Prohibited Patterns

- profile 文面を fact に逆流させる
- role label だけで profile の価値を決める
- 古い export 文面を profile 正本として保存する

## Open Questions

- profile を複数粒度で持つか
- profile source references を必須にするか
