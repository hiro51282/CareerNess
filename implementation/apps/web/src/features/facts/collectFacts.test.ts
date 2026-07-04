import { describe, it, expect } from 'vitest'
import { collectFacts } from './collectFacts'

const experiences = `
- fact_id: fact-proj-a
  type: experience
  status: proposed
  summary: A の経験
  company: ABC
- fact_id: fact-proj-b
  type: experience
  status: confirmed
  summary: B の経験
`

const skills = `
- fact_id: fact-skill-go
  type: skill
  status: confirmed
  summary: Go
  tags: [language]
`

describe('collectFacts', () => {
  it('facts/*.yaml から全 fact を平坦化して集約する', () => {
    const files = new Map<string, string>([
      ['facts/experiences.yaml', experiences],
      ['facts/skills.yaml', skills],
    ])
    const entries = collectFacts(files)
    expect(entries).toHaveLength(3)
    // ファイル帰属とフィールド保持
    const a = entries.find((e) => e.fact.fact_id === 'fact-proj-a')
    expect(a?.file).toBe('facts/experiences.yaml')
    expect(a?.fact.summary).toBe('A の経験')
    expect(a?.fact.company).toBe('ABC')
    const go = entries.find((e) => e.fact.fact_id === 'fact-skill-go')
    expect(go?.fact.tags).toEqual(['language'])
  })

  it('facts/ 以外・非 YAML は無視する', () => {
    const files = new Map<string, string>([
      ['facts/experiences.yaml', experiences],
      ['meta/workspace.yaml', '- fact_id: should-ignore'],
      ['README.md', '# not a fact'],
      ['facts/notes.md', '- fact_id: also-ignore'],
    ])
    const entries = collectFacts(files)
    expect(entries.every((e) => e.file === 'facts/experiences.yaml')).toBe(true)
    expect(entries).toHaveLength(2)
  })

  it('パース不能 / 配列でないファイルはスキップして落ちない', () => {
    const files = new Map<string, string>([
      ['facts/broken.yaml', ':::not: valid: yaml:::'],
      ['facts/notarray.yaml', 'foo: bar'],
      ['facts/skills.yaml', skills],
    ])
    const entries = collectFacts(files)
    expect(entries).toHaveLength(1)
    expect(entries[0].fact.fact_id).toBe('fact-skill-go')
  })

  it('空ファイル集合では空配列を返す', () => {
    expect(collectFacts(new Map())).toEqual([])
  })
})
