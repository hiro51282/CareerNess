# packages/schema

`packages/schema/` は Fact / Profile / Export を中心とした structured career data contract を定義する。

## Responsibility

- fact schema
- profile schema
- export schema
- validation contract
- derived metadata と truth の区別を支える定義

## Must Not

- role/title を truth の中心にしない
- AI prompt logic を持たない
- workspace mutation orchestration を持たない
- presentation wording の正本を持たない

## Why This Exists

- facts / profiles / exports の責務を分離し、データの意味を固定するため
- app や package ごとに解釈がずれないようにするため

## Owned Concepts

- action / decision / impact / context 中心の fact model
- profile は derived view であるという前提
- export は presentation artifact であるという前提

## Dependencies

- 必要最小限の `packages/shared`

## Must Not Depend On

- `apps/*`
- `packages/workspace-core`
- AI provider 実装
- frontend / backend framework details

## Future Considerations

- schema versioning は将来必要になる
- taxonomy / role model の厳密度は拡張余地がある

## Non-goals

- 全文章生成ルールの保持
- product workflow の orchestration
