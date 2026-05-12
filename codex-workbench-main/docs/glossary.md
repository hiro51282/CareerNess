# Glossary

## Fact

ユーザーのキャリアに関する確認可能な事実。CareerNess の正本に最も近い情報単位であり、profile や export の元になる。

## Profile

facts を特定の目的や role に合わせて再構成した見せ方。正本ではなく、複数存在してよい派生表現である。

## Export

提出・共有・貼り付けのために整形された出力。職務経歴書やスカウト返信文などが含まれる。正本ではない。

## Workspace

AI が触れられる範囲として attach されるローカル作業領域。CareerNess では CareerVault を指すことが多い。

## CareerVault

CareerNess における user-owned local workspace。facts、profiles、exports、projects などのキャリア資産を保持する。

## Patch

AI が提案する変更差分。CareerNess では直接上書きよりも patch proposal を基本とする。

## Capability

AI や各コンポーネントに許可された操作範囲。何を読めて、何を書けるかの境界を指す。

## Workspace Boundary

AI が attach された workspace の中だけを対象に動作するという制約。CareerNess の中核設計の一つ。

## Local-first

キャリア情報の正本をクラウドではなくユーザー所有のローカル workspace に置く方針。完全オフライン主義ではなく、ownership 優先の設計を指す。

## User-owned Data

ユーザーのキャリア情報の正本を、運営側ではなくユーザーが保持・管理するという考え方。

## Structured Career Data

自由文章ではなく、project、role、period、achievement などの意味単位に分けて扱うキャリア情報。

## Role

プロジェクトや時期において担っていた責務の型。肩書きだけでなく、実際の役割を説明するための概念。

## Tag

facts や projects を横断整理するための分類ラベル。

## Project

career fact を束ねる文脈単位。会社、案件、期間、チーム、ドメインなどを持つ。

## Browser

ユーザーが CareerNess を操作する UI 層。会話、確認、承認を担うが、キャリア正本は持たない。

## AppServer

Browser と Workspace の間で認証、セッション、AI orchestration を行う層。正本データストアではない。

## Cloud

認証や AI 利用のための補助的なサーバー側要素。CareerVault の代替ではない。
