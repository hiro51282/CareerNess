import { parseFactList } from '../workspace/factFile'

// FactRecord は表示用に主要フィールドを型付けした fact。
// YAML はユーザー所有の正本であり任意フィールドを持ち得るため、未知キーも許容する。
export interface FactRecord {
  fact_id?: string
  type?: string
  status?: string
  summary?: string
  description?: string
  confidence?: string
  company?: string
  period?: string
  tags?: string[]
  created_at?: string
  [key: string]: unknown
}

export interface FactEntry {
  file: string
  fact: FactRecord
}

// collectFacts はワークスペースの全ファイルから facts/*.yaml を抽出し、
// 各ファイルの fact を平坦化して返す（読み取り専用）。
// パースできない / 配列でないファイルはスキップする（Fact ビューアを落とさない）。
export function collectFacts(files: Map<string, string>): FactEntry[] {
  const entries: FactEntry[] = []
  for (const [path, text] of files) {
    if (!isFactFile(path)) continue
    for (const fact of parseFactList(text)) {
      entries.push({ file: path, fact: fact as FactRecord })
    }
  }
  return entries
}

// isFactFile は facts/ 配下の YAML ファイルパスかを判定する。
function isFactFile(path: string): boolean {
  return /^facts\/.+\.(ya?ml)$/i.test(path)
}
