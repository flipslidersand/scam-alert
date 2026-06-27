package model

import "time"

// AnalyzeRequest: フロントエンドからの画像アップロード（マルチパート）
type AnalyzeRequest struct {
	// ファイルはハンドラで直接処理
}

// AnalyzeResponse: フロントエンドへの返却
type AnalyzeResponse struct {
	PatternID      string `json:"patternId"`
	Count          int    `json:"count"`
	NormalizedText string `json:"normalizedText"`
	Error          string `json:"error,omitempty"`
}

// ReportRequest: ユーザーが「記録する」を押した時のリクエスト
type ReportRequest struct {
	PatternID string `json:"patternId"`
}

// Pattern: DBの patterns テーブル行
type Pattern struct {
	ID              string
	TextHash        string
	NormalizedText  string
	Count           int
	FirstSeenAt     time.Time
	LastSeenAt      time.Time
}

// Report: DBの reports テーブル行
type Report struct {
	ID        string
	PatternID string
	SessionID string
	ReportedAt time.Time
}

// VisionAPIResponse: Google Cloud Vision APIのレスポンス
type VisionAPIResponse struct {
	Text       string // 抽出テキスト
	Confidence float32
}
