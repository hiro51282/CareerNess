# 職務経歴書

**氏名**: Example User  
**現在の職務**: シニアバックエンドエンジニア  
**経験年数**: 10 年以上

---

## 職務経歴

### Payment Platform 刷新 (2022年4月〜2023年1月) | 株式会社 Example

**職務**: バックエンド技術リード

決済サービスの大規模刷新に従事。アーキテクチャ設計から実装、3 チーム間の調整まで一貫してリード。

**主な責務・成果**:

- **Go への移行推進**: Python service から Go への技術選定・実装をリード。Type safety とパフォーマンスを優先した意思決定により、本番での安定性を確保。

- **パフォーマンス最適化**: Profiling データを基に hot path を最適化。Database connection pooling、gRPC 導入、キャッシング層実装により、p99 latency を **200ms → 140ms（30% 削減）** に改善。

- **CI/CD 最適化**: CI パイプラインを再設計。15 個の独立 job を 5 並列 stage に統合し、実行時間を **45 分 → 12 分（73% 削減）** に短縮。Developer feedback loop が高速化し、release frequency 向上。

- **Cross-team 調整**: 決済、platform、infra チーム（計 8 人）の技術的議論をリード。統一 API spec を策定し、service 間統合を加速。

**技術スタック**: Go, gRPC, PostgreSQL, Kubernetes, Kafka

**チーム規模**: 8 人（backend 5、frontend 2、infra 1）

---

### 内部ツール CLI Framework (2021年6月〜2022年3月) | 株式会社 Example

**職務**: Senior Engineer

5 つのチーム（計 20 人）で使用される内部 CLI tools framework を設計・開発。

**主な責務**:

- Framework アーキテクチャ設計
- リモートワーク・多拠点環境での環境構築自動化
- Deployment pipeline の一元化

**技術スタック**: Go, Python, Docker, AWS

---

### インフラ最適化・Kubernetes 移行 (2020年10月〜2021年5月) | 株式会社 Example

**職務**: Infrastructure Engineer (Lead)

Legacy on-prem インフラから Kubernetes（AWS EKS）への段階的移行をリード。

**主な実績**:

- Zero-downtime deployment の実現
- Auto-scaling 基盤の構築
- Service mesh / load balancing 設計

---

## スキル・専門知識

### 言語・フレームワーク
- **Go** (Advanced / 10年以上): microservices、performance-critical systems
- **Python** (Intermediate): Scripting、data processing
- **TypeScript / JavaScript**: Web systems（基本的な知識）

### インフラ・プラットフォーム
- **Kubernetes / Docker**: Container orchestration、deployment automation
- **AWS**: ECS、EKS、RDS、CloudWatch
- **CI/CD**: GitHub Actions、GitLab CI パイプライン設計・最適化

### 専門分野
- **Microservices Architecture**: Service boundary 定義、inter-service communication、consistency patterns
- **Performance Optimization**: CPU/Memory/Network profiling、bottleneck 分析、tuning
- **Distributed Systems**: Consistency、availability、concurrency patterns
- **Technical Leadership**: Architecture decision、cross-team coordination、knowledge transfer

### ソフトスキル
- Cross-team 技術的リーダーシップ
- Junior engineer メンタリング
- 社内文書化・knowledge sharing
- Data-driven decision making

---

## 強み・こだわり

- **パフォーマンス重視**: Profiling data に基づいた科学的な最適化
- **Data-driven decision**: 意思決定は metrics / evidence に基づく
- **チーム育成**: Junior engineer の成長支援、knowledge distribution
- **Quality 意識**: Testing strategy（unit、integration、benchmark）の設計

---

## 学位・資格

特になし（実務経験で スキル構築）

---

**最終更新**: 2026-05-24

> このレジュメは CareerVault から自動生成されました。
> 生成元：profile/backend-engineer.yaml
> 出典 fact: fact-proj-payment-platform、fact-ach-latency-improvement、fact-ach-ci-optimization
