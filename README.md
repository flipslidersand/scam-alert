# Scam Alert

![Language](https://img.shields.io/badge/language-Go%20%7C%20React-blue)
![License](https://img.shields.io/badge/license-MIT-green)
![Platform](https://img.shields.io/badge/platform-GCP-orange)

---

## English

A web app that lets you upload a screenshot and instantly see how many times the same text pattern has been recorded in the database — powered by OCR, not AI judgment.

- **No AI verdict** — the user decides whether it's a scam
- **No image storage** — the image is discarded after OCR; privacy is preserved
- **Network effect** — the database grows as more users contribute

### Tech Stack

| Layer | Technology |
|-------|-----------|
| Backend | Go (Cloud Run) |
| Frontend | React (Vite) |
| Database | PostgreSQL (Cloud SQL) |
| OCR | Google Cloud Vision API |
| Deployment | GCP |

### Roadmap

#### Phase 1 (MVP)
- Upload screenshot → OCR via Vision API
- Text normalization → pattern matching
- Show match count → record button

#### Phase 2
- LLM-based context analysis (trend display, no danger score)

#### Phase 3
- Monetization (job ad targeting)

### Project Structure

```
scam-alert/
├── backend/          # Go API
│   ├── cmd/api/
│   ├── internal/
│   └── go.mod
├── frontend/         # React (Vite)
│   ├── src/
│   └── package.json
├── db/
│   └── schema.sql
└── docs/
    └── api.md
```

### Getting Started

#### Backend

```bash
cd backend
go mod download
go run cmd/api/main.go
```

#### Frontend

```bash
cd frontend
npm install
npm run dev
```

#### Database (Local PostgreSQL)

```bash
psql -U postgres -f db/schema.sql
```

### API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/analyze` | Upload screenshot → OCR + pattern match |
| `POST` | `/api/report` | Record a pattern to the DB |
| `GET` | `/api/patterns` | List all patterns (for statistics) |

### License

MIT

---

## 日本語

スクリーンショットを投げるだけで、同じ文言パターンが何件記録されているかを返す Web アプリ。OCR ベースであり、AI が判定するわけではありません。

- **AIが判定しない** — ユーザーが判断する
- **画像は保存しない** — OCR 後は破棄。プライバシー担保
- **ネットワーク効果** — ユーザーが増えるほど DB が育つ

### 技術スタック

| レイヤー | 技術 |
|----------|------|
| バックエンド | Go (Cloud Run) |
| フロントエンド | React (Vite) |
| DB | PostgreSQL (Cloud SQL) |
| OCR | Google Cloud Vision API |
| デプロイ | GCP |

### フェーズ設計

#### Phase 1 (MVP)
- スクショアップロード → Vision API で OCR
- テキスト正規化 → パターン照合
- 件数表示 → 記録ボタン

#### Phase 2
- LLM で文脈判定（傾向表示、危険度スコアなし）

#### Phase 3
- マネタイズ（求人広告ターゲティング）

### プロジェクト構成

```
scam-alert/
├── backend/          # Go API
│   ├── cmd/api/
│   ├── internal/
│   └── go.mod
├── frontend/         # React (Vite)
│   ├── src/
│   └── package.json
├── db/
│   └── schema.sql
└── docs/
    └── api.md
```

### セットアップ

#### バックエンド

```bash
cd backend
go mod download
go run cmd/api/main.go
```

#### フロントエンド

```bash
cd frontend
npm install
npm run dev
```

#### データベース（ローカル PostgreSQL）

```bash
psql -U postgres -f db/schema.sql
```

### API エンドポイント

| メソッド | エンドポイント | 説明 |
|----------|----------------|------|
| `POST` | `/api/analyze` | スクショアップロード → OCR + 照合 |
| `POST` | `/api/report` | パターンを DB に記録 |
| `GET` | `/api/patterns` | パターン一覧（統計用） |

### ライセンス

MIT
