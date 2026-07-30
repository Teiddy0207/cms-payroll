package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

// MeetingAIClient giao tiếp với Python AI Service (Flask, port 5050).
type MeetingAIClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewMeetingAIClient khởi tạo client với timeout 5 phút
// (cuộc họp dài có thể mất thời gian xử lý Whisper + RL).
func NewMeetingAIClient(baseURL string) *MeetingAIClient {
	return &MeetingAIClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Minute,
		},
	}
}

// MeetingAIResponse là response JSON từ Python /analyze-meeting.
type MeetingAIResponse struct {
	Status          string   `json:"status"`
	Transcript      []string `json:"transcript"`
	Summary         string   `json:"summary"`
	KeyDecisions    string   `json:"key_decisions"`
	ActionItems     string   `json:"action_items"`
	Sentiment       string   `json:"sentiment"`
	EfficiencyScore string   `json:"efficiency_score"`
	Message         string   `json:"message"` // chỉ có khi status = "error"
}

// AnalyzeMeeting gửi file audio lên Python service và trả về kết quả AI.
//
// Params:
//   - audioData: nội dung file âm thanh (bytes)
//   - filename:  tên file gốc (ví dụ "meeting.wav"), dùng để xác định MIME type
//
// Returns:
//   - *MeetingAIResponse: kết quả tóm tắt từ AI
//   - error: lỗi kết nối hoặc lỗi từ phía Python service
func (c *MeetingAIClient) AnalyzeMeeting(audioData []byte, filename string) (*MeetingAIResponse, error) {
	// Xây dựng multipart form body
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile("audio", filename)
	if err != nil {
		return nil, fmt.Errorf("ai_client: tạo form file thất bại: %w", err)
	}

	if _, err = io.Copy(part, bytes.NewReader(audioData)); err != nil {
		return nil, fmt.Errorf("ai_client: ghi audio vào form thất bại: %w", err)
	}

	writer.Close()

	// Gửi request tới Python service
	url := c.baseURL + "/analyze-meeting"
	req, err := http.NewRequest(http.MethodPost, url, &body)
	if err != nil {
		return nil, fmt.Errorf("ai_client: tạo HTTP request thất bại: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ai_client: gọi Python service thất bại: %w", err)
	}
	defer resp.Body.Close()

	// Parse JSON response
	var result MeetingAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("ai_client: parse response thất bại: %w", err)
	}

	if result.Status != "success" {
		return nil, fmt.Errorf("ai_client: Python service trả lỗi: %s", result.Message)
	}

	return &result, nil
}
