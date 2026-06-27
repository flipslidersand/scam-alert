package service

import (
	"crypto/sha256"
	"fmt"
	"regexp"
	"strings"
)

// PatternService: テキスト正規化とパターンマッチのビジネスロジック
type PatternService struct{}

// NewPatternService: パターンサービス初期化
func NewPatternService() *PatternService {
	return &PatternService{}
}

// NormalizeText: テキスト正規化
// - 全角スペース → 半角スペース
// - 改行 → スペース
// - 連続スペース → 単一スペース
// - 首尾の空白削除
func (ps *PatternService) NormalizeText(text string) string {
	// 全角スペース（U+3000）を半角スペース（U+0020）に
	text = strings.ReplaceAll(text, "　", " ")

	// 改行・タブをスペースに統一
	text = regexp.MustCompile(`[\r\n\t]+`).ReplaceAllString(text, " ")

	// 連続スペースを単一スペースに
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")

	// 首尾の空白を削除
	text = strings.TrimSpace(text)

	return text
}

// ComputeHash: テキストの SHA256 ハッシュ計算
func (ps *PatternService) ComputeHash(text string) string {
	hash := sha256.Sum256([]byte(text))
	return fmt.Sprintf("%x", hash)
}

// CreatePatternID: パターンIDを生成
// 通常はDB側で UUID を生成するが、テスト用に簡易版を提供
func (ps *PatternService) CreatePatternID() string {
	// 実運用では database/sql 層で gen_random_uuid() を使う
	// ここはスタブ
	return "pattern_id_stub"
}

// IsSimilar: テキスト相似判定（将来: Levenshtein距離で部分一致検出）
// Phase 1ではスキップし、完全一致のみを使用
func (ps *PatternService) IsSimilar(text1, text2 string, threshold float64) bool {
	// Phase 2 で実装: 編集距離 < threshold の場合
	// for now: 完全一致のみ
	return text1 == text2
}
