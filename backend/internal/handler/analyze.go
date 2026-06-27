package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/flipslidersand/scam-alert/backend/internal/model"
	"github.com/flipslidersand/scam-alert/backend/internal/service"
)

// AnalyzeHandler: スクショアップロード → OCR → パターン照合
type AnalyzeHandler struct {
	visionService  *service.VisionClient
	patternService *service.PatternService
	dbService      *service.DBService
}

// NewAnalyzeHandler: ハンドラー初期化
func NewAnalyzeHandler(
	visionService *service.VisionClient,
	patternService *service.PatternService,
	dbService *service.DBService,
) *AnalyzeHandler {
	return &AnalyzeHandler{
		visionService:  visionService,
		patternService: patternService,
		dbService:      dbService,
	}
}

// ServeHTTP: POST /api/analyze
func (h *AnalyzeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// マルチパート画像を読み込む
	err := r.ParseMultipartForm(32 << 20) // 32MB max
	if err != nil {
		respondJSON(w, http.StatusBadRequest, model.AnalyzeResponse{Error: "Failed to parse form"})
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		respondJSON(w, http.StatusBadRequest, model.AnalyzeResponse{Error: "No file provided"})
		return
	}
	defer file.Close()

	// ファイル形式チェック
	if !strings.HasPrefix(fileHeader.Header.Get("Content-Type"), "image/") {
		respondJSON(w, http.StatusBadRequest, model.AnalyzeResponse{Error: "File must be an image"})
		return
	}

	// Vision API でテキスト抽出
	extractedText, confidence, err := h.visionService.ExtractText(r.Context(), file)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, model.AnalyzeResponse{
			Error: fmt.Sprintf("Vision API error: %v", err),
		})
		return
	}

	if confidence < 0.5 {
		respondJSON(w, http.StatusBadRequest, model.AnalyzeResponse{
			Error: "Confidence too low",
		})
		return
	}

	// テキスト正規化
	normalizedText := h.patternService.NormalizeText(extractedText)
	if len(normalizedText) == 0 {
		respondJSON(w, http.StatusBadRequest, model.AnalyzeResponse{Error: "No text extracted"})
		return
	}

	// ハッシュ計算
	textHash := h.patternService.ComputeHash(normalizedText)

	// DB でパターン照合・作成
	pattern, err := h.dbService.CreateOrUpdatePattern(r.Context(), textHash, normalizedText)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, model.AnalyzeResponse{
			Error: fmt.Sprintf("Database error: %v", err),
		})
		return
	}

	// 成功レスポンス
	respondJSON(w, http.StatusOK, model.AnalyzeResponse{
		PatternID:      pattern.ID,
		Count:          pattern.Count,
		NormalizedText: pattern.NormalizedText,
	})
}

// ReportHandler: POST /api/report — ユーザーが「記録する」を押した時
type ReportHandler struct {
	dbService *service.DBService
}

// NewReportHandler: ハンドラー初期化
func NewReportHandler(dbService *service.DBService) *ReportHandler {
	return &ReportHandler{dbService: dbService}
}

// ServeHTTP: POST /api/report
func (h *ReportHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var req model.ReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		return
	}

	// セッション ID（簡易版: Cookie から取得。本来は RFC に従うべき）
	sessionID := "session_default" // TODO: Cookie から取得

	// DB に記録
	err := h.dbService.LogReport(r.Context(), req.PatternID, sessionID)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to log report: %v", err),
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"status": "recorded"})
}

// Helper: JSON レスポンス
func respondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
