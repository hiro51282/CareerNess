# infra

`infra/` は deployment / infrastructure 関連の定義とメモを置く。MVP では lightweight を前提とし、product truth を持たない。

## Responsibility

- deploy strategy の表現
- environment definition
- runtime infrastructure の構成メモ
- local / cloud 実行環境の補助

## Must Not

- core domain logic を持たない
- career data の hidden persistence を増やさない
- infra convenience のために architecture boundary を破らない

## Why This Exists

- product code と deploy concern を分離するため
- MVP の運用を軽くしつつ、将来の拡張余地を残すため

## Owned Concepts

- container / IaC / deployment note
- environment-specific setup

## Dependencies

- `apps/*` の runtime 前提
- repository の build / deploy artifact

## Must Not Depend On

- examples を本番構成とみなすこと
- tmp の一時成果物
- product truth を補う hidden state

## Future Considerations

- AWS 前提の構成は増えてよい
- ただし Local-first と lightweight MVP を壊さない

## Non-goals

- repository の中心責務になること
- canonical data storage layer
