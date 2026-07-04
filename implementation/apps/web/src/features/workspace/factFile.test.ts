import { describe, it, expect } from 'vitest'
import * as yaml from 'js-yaml'
import { upsertFactYaml, markFactStatusYaml } from './factFile'

describe('upsertFactYaml', () => {
  it('空ファイルへ新規 fact を配列形で追記する', () => {
    const out = upsertFactYaml('', { summary: 'CI 改善', type: 'experience' }, 'fact-001')
    const parsed = yaml.load(out)
    expect(Array.isArray(parsed)).toBe(true)
    expect(parsed).toEqual([{ fact_id: 'fact-001', summary: 'CI 改善', type: 'experience' }])
  })

  it('既存配列に別 fact_id の fact を追記する（既存は保持）', () => {
    const existing = yaml.dump([{ fact_id: 'fact-existing', summary: '既存' }])
    const out = upsertFactYaml(existing, { summary: '新規' }, 'fact-new')
    const parsed = yaml.load(out) as Array<Record<string, unknown>>
    expect(parsed).toHaveLength(2)
    expect(parsed.map((f) => f.fact_id)).toEqual(['fact-existing', 'fact-new'])
  })

  it('同じ fact_id は置換する（重複しない）', () => {
    const existing = yaml.dump([{ fact_id: 'fact-001', summary: '旧' }])
    const out = upsertFactYaml(existing, { summary: '新', status: 'proposed' }, 'fact-001')
    const parsed = yaml.load(out) as Array<Record<string, unknown>>
    expect(parsed).toHaveLength(1)
    expect(parsed[0]).toEqual({ fact_id: 'fact-001', summary: '新', status: 'proposed' })
  })

  it('引数 factId が fact 内の fact_id より優先される', () => {
    const out = upsertFactYaml('', { fact_id: 'fact-inner', summary: 's' }, 'fact-arg')
    const parsed = yaml.load(out) as Array<Record<string, unknown>>
    expect(parsed[0].fact_id).toBe('fact-arg')
  })

  it('配列・nested・array フィールドを保持する', () => {
    const fact = {
      summary: 's',
      tags: ['backend', 'ci'],
      team_size: 5,
      impact: { summary: 'CI 40% 短縮' },
    }
    const out = upsertFactYaml('', fact, 'fact-rich')
    const parsed = yaml.load(out) as Array<Record<string, unknown>>
    expect(parsed[0].tags).toEqual(['backend', 'ci'])
    expect(parsed[0].team_size).toBe(5)
    expect(parsed[0].impact).toEqual({ summary: 'CI 40% 短縮' })
  })

  it('配列でない / 壊れた既存本文は空配列から開始する（正本形へ寄せる）', () => {
    // 旧 mock が書いた平坦オブジェクト形
    const legacy = yaml.dump({ source: 'conversation', raw_text: '...', status: 'proposed' })
    const out = upsertFactYaml(legacy, { summary: 's' }, 'fact-001')
    const parsed = yaml.load(out)
    expect(parsed).toEqual([{ fact_id: 'fact-001', summary: 's' }])
  })
})

describe('markFactStatusYaml', () => {
  const existing = yaml.dump([
    { fact_id: 'fact-a', type: 'experience', status: 'proposed', summary: 'A', tags: ['backend'] },
    { fact_id: 'fact-b', type: 'skill', status: 'proposed', summary: 'B' },
  ])

  it('対象 fact の status を更新し、他フィールドは保持する', () => {
    const out = markFactStatusYaml(existing, 'fact-a', 'confirmed')
    const parsed = yaml.load(out) as Array<Record<string, unknown>>
    const a = parsed.find((f) => f.fact_id === 'fact-a')!
    expect(a.status).toBe('confirmed')
    expect(a.summary).toBe('A')
    expect(a.tags).toEqual(['backend'])
    // 他の fact は不変
    expect(parsed.find((f) => f.fact_id === 'fact-b')!.status).toBe('proposed')
  })

  it('updated_at を RFC3339（ミリ秒なし・末尾 Z）で付与する', () => {
    const out = markFactStatusYaml(existing, 'fact-a', 'confirmed')
    const parsed = yaml.load(out) as Array<Record<string, unknown>>
    const a = parsed.find((f) => f.fact_id === 'fact-a')!
    expect(typeof a.updated_at).toBe('string')
    expect(a.updated_at).toMatch(/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$/)
  })

  it('対象 fact が無ければ throw する', () => {
    expect(() => markFactStatusYaml(existing, 'fact-missing', 'confirmed')).toThrow()
  })
})
