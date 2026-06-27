package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/flipslidersand/scam-alert/backend/internal/handler"
	"github.com/flipslidersand/scam-alert/backend/internal/service"
)

func main() {
	ctx := context.Background()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	log.Printf("Starting Scam Alert API on port %s...\n", port)

	// 1. Vision API クライアント初期化
	visionClient, err := service.NewVisionClient(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize Vision API: %v", err)
	}
	defer visionClient.Close()

	// 2. PostgreSQL 接続
	dbService, err := service.NewDBService(databaseURL)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer dbService.Close()

	// 3. ビジネスロジックサービス初期化
	patternService := service.NewPatternService()

	// 4. HTTP ハンドラー登録
	mux := http.NewServeMux()

	analyzeHandler := handler.NewAnalyzeHandler(visionClient, patternService, dbService)
	reportHandler := handler.NewReportHandler(dbService)

	mux.Handle("/api/analyze", analyzeHandler)
	mux.Handle("/api/report", reportHandler)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	})

	// 5. サーバー起動
	server := &http.Server{
		Addr:         net.JoinHostPort("", port),
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	log.Printf("Server listening on %s", server.Addr)

	// Graceful shutdown 待機
	<-sigChan
	log.Println("Shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Shutdown error: %v", err)
	}

	log.Println("Server stopped")
}
