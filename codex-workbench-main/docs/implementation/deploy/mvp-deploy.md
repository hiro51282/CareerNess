# MVP Deploy

MVP では deploy complexity より、Local-first と workspace boundary を壊さないことを優先する。

## Goal

- Browser と AppServer を小さく動かす
- local workspace attach を前提にする
- cloud-side retention を最小限にする

## Suggested Shape

```text
Browser UI
  ↓
Small AppServer
  ↓
Local workspace integration
```

## MVP Assumptions

- single-user
- one attached workspace at a time
- no cloud sync
- no shared workspace

## Deployment Priorities

- セッションが成立すること
- patch proposal / approval / apply が通ること
- workspace 外アクセスを防げること
- failure 時に destructive にならないこと

## Not Required For MVP

- multi-region
- enterprise auth matrix
- background export farm
- centralized career data lake

## Operational Notes

- サーバーログは最小限にする
- prompt / workspace 内容の retention は明示的に絞る
- apply 失敗時の再試行は idempotent に寄せる
