package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"

	"github.com/flipslidersand/scam-alert/backend/internal/model"
)

// DBService: PostgreSQL との連携
type DBService struct {
	db *sql.DB
}

// NewDBService: データベース接続を初期化
func NewDBService(databaseURL string) (*DBService, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// コネクションプールのテスト
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DBService{db: db}, nil
}

// GetPatternByHash: テキストハッシュでパターンを検索
func (ds *DBService) GetPatternByHash(ctx context.Context, textHash string) (*model.Pattern, error) {
	pattern := &model.Pattern{}
	err := ds.db.QueryRowContext(ctx,
		`SELECT id, text_hash, normalized_text, count, first_seen_at, last_seen_at
		 FROM patterns WHERE text_hash = $1`,
		textHash).Scan(
		&pattern.ID,
		&pattern.TextHash,
		&pattern.NormalizedText,
		&pattern.Count,
		&pattern.FirstSeenAt,
		&pattern.LastSeenAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // Not found
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query pattern: %w", err)
	}

	return pattern, nil
}

// CreateOrUpdatePattern: パターンを作成または更新
func (ds *DBService) CreateOrUpdatePattern(ctx context.Context, textHash, normalizedText string) (*model.Pattern, error) {
	patternID := uuid.New().String()
	now := time.Now()

	pattern := &model.Pattern{}

	// INSERT or UPDATE（UPSERT）
	err := ds.db.QueryRowContext(ctx,
		`INSERT INTO patterns (id, text_hash, normalized_text, count, first_seen_at, last_seen_at)
		 VALUES ($1, $2, $3, 1, $4, $5)
		 ON CONFLICT (text_hash) DO UPDATE
		 SET count = count + 1, last_seen_at = EXCLUDED.last_seen_at
		 RETURNING id, text_hash, normalized_text, count, first_seen_at, last_seen_at`,
		patternID, textHash, normalizedText, now, now).Scan(
		&pattern.ID,
		&pattern.TextHash,
		&pattern.NormalizedText,
		&pattern.Count,
		&pattern.FirstSeenAt,
		&pattern.LastSeenAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create/update pattern: %w", err)
	}

	return pattern, nil
}

// LogReport: ユーザーの報告をログに記録
func (ds *DBService) LogReport(ctx context.Context, patternID, sessionID string) error {
	reportID := uuid.New().String()
	_, err := ds.db.ExecContext(ctx,
		`INSERT INTO reports (id, pattern_id, session_id, reported_at)
		 VALUES ($1, $2, $3, NOW())`,
		reportID, patternID, sessionID,
	)
	if err != nil {
		return fmt.Errorf("failed to log report: %w", err)
	}
	return nil
}

// Close: データベース接続を閉じる
func (ds *DBService) Close() error {
	if ds.db != nil {
		return ds.db.Close()
	}
	return nil
}
