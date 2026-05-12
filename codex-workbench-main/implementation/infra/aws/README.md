# infra/aws

`infra/aws/` は AWS architecture memo と deploy consideration を置く future-oriented なディレクトリである。

## Responsibility

- AWS 利用方針メモ
- service selection の判断材料
- security / deploy consideration の整理

## Must Not

- 実装本体の責務をここで決め打ちしない
- AWS 依存を product core に漏らさない
- Local-first を崩して cloud persistence を正本化しない

## Why This Exists

- MVP 段階でも将来の運用選択肢を整理しておくため
- ただし infra memo を core architecture と混同しないため

## Owned Concepts

- AWS deploy option
- managed service tradeoff
- security / cost / operations memo

## Dependencies

- `docs/` の deploy / security / architecture 文脈
- `infra/terraform` と整合する方針

## Must Not Depend On

- `packages/*` の内部実装
- canonical workspace data

## Future Considerations

- actual IaC が増えたら `terraform/` へ責務を移す
- ここは memo 中心に留める

## Non-goals

- AWS 実装の唯一ソース
- product runtime logic の置き場
