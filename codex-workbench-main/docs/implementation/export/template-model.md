# Template Model

template は export の器であり、career truth ではない。template 設計では、見せ方の制約を facts や profiles に侵食させないことが重要である。

## Responsibility

- format ごとの構造を与える
- audience ごとの emphasis を反映する
- media constraints を吸収する

## Must Not

- fact schema を定義する
- profile の canonical meaning を決める
- export から逆算して facts を歪める

## Template Inputs

- selected profile
- selected facts or fact ranges
- audience
- format
- optional length constraint

## Template Layers

### Structural Template

- sections
- ordering
- required fields

### Style Template

- tone
- density
- emphasis rules

### Media Constraint Template

- line limits
- field labels
- plain text only rules

## Example Directions

- technical depth first
- impact first
- platform / infra emphasis
- leadership / decision emphasis

ここで重要なのは「肩書」より「行動・判断・成果」が前に出ることだ。

## Human-Readable Templates

template 定義も人間が読める形を優先する。rule-heavy な DSL を急いで作るより、YAML + markdown guidance 程度で十分な可能性が高い。

## Open Questions

- template をファイル保存するか、UI preset に寄せるか
- hard validation と soft guidance の境界をどこに置くか
