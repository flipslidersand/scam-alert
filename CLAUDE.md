# scam-alert CLAUDE.md

## Project Overview

スクショ → OCR → パターンマッチ → 件数表示 の Web アプリ。

**設計思想**: AIは判定しない。ユーザーが判断。画像は保存しない。

## Tech Stack

| Layer | Tech | Purpose |
|-------|------|---------|
| **Frontend** | React (Vite) | UI: Upload + Results |
| **Backend** | Go (Cloud Run) | Vision API 連携、パターンマッチ |
| **DB** | PostgreSQL | パターン統計 |
| **OCR** | Google Cloud Vision | テキスト抽出 |

## Key Implementation Details

### 1. Vision API Integration

- DOCUMENT_TEXT_DETECTION モード（通常の Document Text Detection）
- テキスト抽出 + ボックス座標（レイアウト復元用）
- 画像は即座に破棄（メモリ上のみ）

### 2. Text Normalization

```go
// 正規化ルール
- 全角→半角スペース
- 改行を単一スペースに
- 連続スペース → 単一スペース
- 首尾の空白を削除
- ハッシュ化（重複検出用）
```

### 3. Pattern Matching

1. **完全一致**: `text_hash` で直接マッチ
2. **部分一致**: 編集距離 < threshold の場合カウント（フェーズ2で実装検討）

### 4. Database Schema

```sql
-- パターンテーブル
CREATE TABLE patterns (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  text_hash VARCHAR(64) NOT NULL,
  normalized_text TEXT NOT NULL,
  count INTEGER DEFAULT 1,
  first_seen_at TIMESTAMP DEFAULT NOW(),
  last_seen_at TIMESTAMP DEFAULT NOW()
);

-- 記録ログ
CREATE TABLE reports (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  pattern_id UUID REFERENCES patterns(id),
  reported_at TIMESTAMP DEFAULT NOW()
);
```

## Phase Breakdown

| Phase | Scope | Effort | Status |
|-------|-------|--------|--------|
| **1** | MVP: Upload → Analyze → Record | 2–3w | ⏳ Planning |
| **2** | LLM Context (傾向表示) | 1w | 🔜 Future |
| **3** | Monetize (Ad targeting) | 2w | 🔜 Future |

## Development Notes

- Vision API quota: Check GCP billing
- Session ID (no user tracking)
- Image deletion: Important for privacy story

## Links

- **API Spec**: docs/api.md
- **GitHub**: github.com/flipslidersand/scam-alert (TBD)
