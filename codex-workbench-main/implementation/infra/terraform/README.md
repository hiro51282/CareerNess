# infra/terraform

`infra/terraform/` は将来の AWS IaC を中心とした infrastructure as code を置く場所である。MVP 時点では minimal を前提とする。

## Responsibility

- 将来の cloud resource 定義
- environment provisioning の再現可能性
- deploy infrastructure のコード化

## Must Not

- MVP で過剰な enterprise infrastructure を先回りしない
- app logic や domain rule を Terraform に埋め込まない
- career truth の保存先を cloud に固定しない

## Why This Exists

- 運用拡大時に手動構築依存を減らすため
- ただし product architecture より先に infra を肥大化させないため

## Owned Concepts

- resource topology
- environment provisioning
- cloud-side least-privilege の表現

## Dependencies

- `infra/aws` の方針メモ
- deploy 対象となる `apps/*`

## Must Not Depend On

- workspace の canonical state
- examples や tmp の内容

## Future Considerations

- staging / production 分離は将来必要になる
- ただし MVP では最小構成を維持する

## Non-goals

- cloud-first への転換
- product domain architecture の代替
