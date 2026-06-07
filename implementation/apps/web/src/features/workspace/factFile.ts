import * as yaml from 'js-yaml'

// CareerVault の facts ファイルは正本として「fact オブジェクトの YAML 配列」で保持する。
// 参照: docs/workspace/careervault-yaml-structure.md, examples/CareerVault/facts/*.yaml,
// および Go 側 applier（[]*YAMLFact）。ブラウザ側 apply もこの形に揃える。

/**
 * upsertFactYaml は既存の facts ファイル本文へ 1 件の fact を upsert し、
 * 正本形（fact オブジェクトの配列）の YAML 文字列として返す純関数。
 *
 * - fact_id をキーに、一致する要素があれば置換、なければ末尾へ追記する。
 * - 引数 factId を必ず fact_id として採用する（op の entity_id を正とするため）。
 * - 既存本文が配列でない / パースできない場合は空配列から開始し、正本形へ寄せる。
 */
export function upsertFactYaml(
  existingText: string,
  fact: Record<string, unknown>,
  factId: string,
): string {
  const facts = parseFactList(existingText)

  // fact 側に別の fact_id があっても引数を優先し、かつ先頭キーとして並べる。
  const { fact_id: _ignored, ...rest } = fact
  const entry: Record<string, unknown> = { fact_id: factId, ...rest }

  const idx = facts.findIndex((f) => f['fact_id'] === factId)
  if (idx >= 0) {
    facts[idx] = entry
  } else {
    facts.push(entry)
  }

  return yaml.dump(facts, { lineWidth: 120 })
}

// 既存テキストを fact 配列として読む。配列でない / 壊れている場合は空配列から開始する。
function parseFactList(text: string): Record<string, unknown>[] {
  if (!text || !text.trim()) return []
  let loaded: unknown
  try {
    loaded = yaml.load(text)
  } catch {
    return []
  }
  if (!Array.isArray(loaded)) return []
  return loaded.filter(
    (f): f is Record<string, unknown> => f !== null && typeof f === 'object',
  )
}
