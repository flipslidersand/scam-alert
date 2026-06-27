-- Patterns table: スクショから抽出したテキストパターンの統計
CREATE TABLE IF NOT EXISTS patterns (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  text_hash VARCHAR(64) NOT NULL UNIQUE,  -- SHA256ハッシュ
  normalized_text TEXT NOT NULL,          -- 正規化済みテキスト
  count INTEGER DEFAULT 1,                -- 記録件数
  first_seen_at TIMESTAMP DEFAULT NOW(),
  last_seen_at TIMESTAMP DEFAULT NOW(),

  INDEX idx_text_hash (text_hash),
  INDEX idx_last_seen (last_seen_at)
);

-- Reports table: ユーザーが「記録する」ボタンを押したログ
CREATE TABLE IF NOT EXISTS reports (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  pattern_id UUID REFERENCES patterns(id) ON DELETE CASCADE,
  session_id VARCHAR(100),                -- セッション識別（個人情報なし）
  reported_at TIMESTAMP DEFAULT NOW(),

  INDEX idx_pattern_id (pattern_id),
  INDEX idx_reported_at (reported_at)
);

-- 将来: ユーザーのホワイトリスト（Phase 2)
-- CREATE TABLE IF NOT EXISTS trusted_patterns (
--   id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--   pattern_id UUID REFERENCES patterns(id),
--   user_id VARCHAR(100),
--   reason TEXT,
--   created_at TIMESTAMP DEFAULT NOW()
-- );
