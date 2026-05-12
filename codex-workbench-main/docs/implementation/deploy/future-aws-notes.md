# Future AWS Notes

この文書は将来の AWS 配置案メモであり、現時点の設計原則を上書きするものではない。特に、AWS を使うことは CareerVault のクラウド正本化を意味しない。

## Invariant

- career data truth は local workspace に残る
- AppServer は AI orchestration layer のまま
- cloud は session / auth / delivery を助ける

## Reasonable Cloud Use

- web frontend hosting
- AppServer runtime
- auth integration
- operational logging

## Avoid

- facts / profiles / exports の恒久保管
- workspace mirror bucket
- AI convenience のための長期全文 retention

## Future Considerations

- temporary artifact storage の TTL 設計
- encryption / secret handling
- minimal audit trails
- regional deployment より先に retention 明確化

## Open Questions

- AWS を使うなら local helper とどう接続するか
- approval apply を完全ローカルに寄せるか
