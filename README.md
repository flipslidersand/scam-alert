# Scam Alert — スクショで件数検索

スクリーンショットを投げるだけで、同じ文言パターンが何件記録されているかを返す Web アプリ。

- **AIが判定しない** — ユーザーが判断する
- **画像は保存しない** — OCR後は破棄。プライバシー担保
- **ネットワーク効果** — ユーザーが増えるほどDB が育つ

## 技術スタック

- **Backend**: Go (Cloud Run)
- **Frontend**: React
- **DB**: PostgreSQL (Cloud SQL)
- **OCR**: Google Cloud Vision API
- **Deployment**: GCP

## フェーズ設計

### Phase 1 (MVP)
- スクショアップロード → Vision API で OCR
- テキスト正規化 → パターン照合
- 件数表示 → 記録ボタン

### Phase 2
- LLM で文脈判定（傾向表示、危険度スコアなし）

### Phase 3
- マネタイズ（求人広告ターゲティング）

## Project Structure

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

## Getting Started

### Backend

```bash
cd backend
go mod download
go run cmd/api/main.go
```

### Frontend

```bash
cd frontend
npm install
npm run dev
```

### Database (Local PostgreSQL)

```bash
psql -U postgres -f db/schema.sql
```

## API Endpoints

- `POST /api/analyze` — スクショアップロード → OCR + 照合
- `POST /api/report` — パターンをDBに記録
- `GET /api/patterns` — パターン一覧（統計用）

## License

MIT
