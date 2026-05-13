# .github

`implementation/.github/` は GitHub 関連設定を置く。repository automation の層であり、product truth や権限委譲の中心ではない。

## Responsibility

- GitHub configuration
- CI/CD support
- repository automation 設定

## Must Not

- dangerous auto-apply を無条件で実行しない
- review を飛ばして workspace mutation を確定しない
- product secret handling の唯一責務を持たない

## Why This Exists

- repository 運用と product code を分離するため
- AI-generated unsafe automation を避けやすくするため

## Owned Concepts

- repository workflow policy
- automation trigger
- issue / PR support 設定

## Dependencies

- `.github/workflows`
- `scripts/`
- build / test 対象の `apps/*`, `packages/*`

## Must Not Depend On

- tmp artifact の継続利用
- hidden deploy state

## Future Considerations

- automation は増えてよい
- ただし least privilege と human review を前提にする

## Non-goals

- production control plane
- repository 外の真実の保管
