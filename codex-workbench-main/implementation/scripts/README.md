# scripts

`scripts/` は developer utility scripts を置く。本番 runtime logic の居場所ではない。

## Responsibility

- 開発補助
- lint / format / setup / migration 補助
- ローカル検証や CI 補助

## Must Not

- production runtime の唯一実装を置かない
- core domain logic を script の中に閉じ込めない
- 手元 convenience のために boundary を破る自動操作を常態化しない

## Why This Exists

- 反復作業を減らし、開発速度を上げるため
- ただし product code と運用補助を分離するため

## Owned Concepts

- developer workflow automation
- repository maintenance helper

## Dependencies

- `apps/*`
- `packages/*`
- `infra/*`

## Must Not Depend On

- `tmp/` の存在を前提にした恒久運用
- 人手承認が必要な apply を無断で飛ばすロジック

## Future Considerations

- script は増えてよい
- 頻用 script が product behavior に近づいたら package / app へ昇格させる

## Non-goals

- 本番アプリケーション層
- hidden admin backdoor
