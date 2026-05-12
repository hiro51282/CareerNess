# Workspace Locking

CareerNess は single-user 前提だが、それでも patch apply 中の不整合は防ぐ必要がある。locking は多人数協調のためではなく、「承認した差分と実際に当たる差分を一致させる」ためにある。

## Goal

- apply 中の partial write を避ける
- stale patch の silent apply を避ける
- history append までを一貫処理に近づける

## Principles

- coarse-grained でよい
- Git 前提にしない
- complex distributed lock は不要
- single workspace owner の手元運用を壊さない

## Suggested Model

MVP では workspace 単位、または file group 単位の軽量 lock で十分である。

- acquire lock
- re-check revision
- apply files
- append history
- release lock

## Lock Scope

優先度は次の順でよい。

1. workspace-wide apply lock
2. fact/profile/export category lock
3. individual file lock

MVP は `workspace-wide apply lock` でよい可能性が高い。シンプルさを優先する。

## Conflict Cases

- patch review 中にユーザーが手で fact を編集した
- profile regenerate の前提 fact が変わった
- 別セッションで history が進んだ

これらは lock 以前に stale 判定で検出し、必要なら再承認に戻す。

## Failure Handling

- lock 取得失敗時は再試行可能
- apply 途中失敗時は部分適用を history に明示
- rollback patch が作れるなら保存する

## Prohibited Patterns

- lock なしで multi-file patch を apply する
- history append 前に成功扱いで応答する
- stale patch を lock 取得後に勝手に書き換える

## Open Questions

- lock 実装を file lock にするか process mutex にするか
