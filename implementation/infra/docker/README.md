# infra/docker

`infra/docker/` は local development と lightweight deployment のための Docker 関連定義を置く。

## Responsibility

- ローカル起動用コンテナ定義
- 最小限の deploy packaging
- 開発環境の再現性向上

## Must Not

- Docker image に canonical workspace data を埋め込まない
- production truth を container layer に持たせない
- environment secret の唯一の保管庫にしない

## Why This Exists

- 環境差分を下げ、MVP の起動・検証を簡単にするため

## Owned Concepts

- container build definition
- compose / local runtime orchestration

## Dependencies

- `apps/*`
- 必要なら `scripts/`

## Must Not Depend On

- user workspace の実データ
- `tmp/` の偶発的生成物

## Future Considerations

- devcontainer や preview environment 向け拡張余地がある
- ただし Docker を architecture の中心にしない

## Non-goals

- canonical deployment platform
- data persistence strategy の定義
