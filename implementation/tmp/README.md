# tmp

`tmp/` は temporary artifacts の置き場である。ここに truth を保持してはいけない。

## Responsibility

- 一時生成物の隔離
- ローカル検証中の作業ファイル置き場
- 破棄可能な artifact の退避

## Must Not

- git-tracked runtime data を置かない
- canonical workspace data を置かない
- patch history や approval record の正本を置かない
- app が `tmp/` を前提にしか動かない設計にしない

## Why This Exists

- 一時ファイルと product asset を混在させないため
- accidental truth ownership を防ぐため

## Owned Concepts

- disposable artifact
- scratch output

## Dependencies

- なし。必要時に各種 script / local run が利用してよい

## Must Not Depend On

- `apps/*` の本番経路
- `packages/*` の canonical logic

## Future Considerations

- `.gitignore` と運用ルールで掃除しやすくする
- CI で残骸依存が起きないようにする

## Non-goals

- cache の正本化
- runtime state store
