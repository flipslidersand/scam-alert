package service

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"

	vision "cloud.google.com/go/vision/v2/apiv1"
	visionpb "cloud.google.com/go/vision/v2/apiv1/visionpb"
)

// VisionClient: Google Cloud Vision API クライアント
type VisionClient struct {
	client *vision.ImageAnnotatorClient
}

// NewVisionClient: Vision APIクライアントの初期化
func NewVisionClient(ctx context.Context) (*VisionClient, error) {
	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create vision client: %w", err)
	}
	return &VisionClient{client: client}, nil
}

// ExtractText: 画像からテキストを抽出
// 入力: ファイルストリーム（メモリ上のみ）
// 出力: 抽出テキスト + 信頼度
func (vc *VisionClient) ExtractText(ctx context.Context, imageData io.Reader) (string, float32, error) {
	if vc.client == nil {
		return "", 0, fmt.Errorf("vision client not initialized")
	}

	// ファイルをバイトに読み込む（メモリのみ）
	imageBytes, err := io.ReadAll(imageData)
	if err != nil {
		return "", 0, fmt.Errorf("failed to read image: %w", err)
	}

	// 画像を検証（decode）
	_, _, err = image.DecodeConfig(bytes.NewReader(imageBytes))
	if err != nil {
		return "", 0, fmt.Errorf("invalid image format: %w", err)
	}

	// Vision API に送信（DOCUMENT_TEXT_DETECTION）
	req := &visionpb.BatchAnnotateImagesRequest{
		Requests: []*visionpb.AnnotateImageRequest{
			{
				Image: &visionpb.Image{Content: imageBytes},
				Features: []*visionpb.Feature{
					{Type: visionpb.Feature_DOCUMENT_TEXT_DETECTION},
				},
			},
		},
	}

	// 画像はここで破棄（重要: メモリに保持しない）
	imageBytes = nil

	resp, err := vc.client.BatchAnnotateImages(ctx, req)
	if err != nil {
		return "", 0, fmt.Errorf("vision API error: %w", err)
	}
	if len(resp.GetResponses()) == 0 {
		return "", 0, fmt.Errorf("no response from vision API")
	}

	annotation := resp.GetResponses()[0]
	if annotation.GetError() != nil {
		return "", 0, fmt.Errorf("vision API error: %s", annotation.GetError().GetMessage())
	}

	fullText := annotation.GetFullTextAnnotation()
	if fullText == nil || fullText.GetText() == "" {
		return "", 0, fmt.Errorf("no text detected in image")
	}

	// TextAnnotation自体に信頼度は無いため、先頭ページの信頼度を代表値として使う
	var confidence float32
	if pages := fullText.GetPages(); len(pages) > 0 {
		confidence = pages[0].GetConfidence()
	}

	return fullText.GetText(), confidence, nil
}

// Close: Vision APIクライアントの後処理
func (vc *VisionClient) Close() error {
	if vc.client != nil {
		return vc.client.Close()
	}
	return nil
}
