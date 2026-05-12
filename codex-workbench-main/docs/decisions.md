# Decisions

この文書は、CareerNess の設計上すでに強く置いている判断を短くまとめるためのメモである。詳細は各文書を参照する。

## 固定した判断

- CareerNess は Local-first を前提とする
- ユーザーの career data は user-owned である
- AI は workspace-scoped で動作する
- facts / profiles / exports は分離する
- AppServer は orchestration を担うが正本は保持しない
- AI は unrestricted shell / file access を前提にしない
- patch proposal と user approval を重視する

## まだ未確定のもの

- profiles と narratives の命名統一
- fact schema の厳密度
- sync の扱い
- embeddings の位置づけ
