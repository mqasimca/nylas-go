package nylas

import (
	"encoding/json"
	"testing"
)

func TestResponse_JSON(t *testing.T) {
	jsonData := `{
		"data": {"id": "msg-123", "subject": "Test"},
		"request_id": "req-456"
	}`

	type TestData struct {
		ID      string `json:"id"`
		Subject string `json:"subject"`
	}

	var resp Response[TestData]
	if err := json.Unmarshal([]byte(jsonData), &resp); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if resp.RequestID != "req-456" {
		t.Errorf("RequestID = %s, want req-456", resp.RequestID)
	}
	if resp.Data.ID != "msg-123" {
		t.Errorf("Data.ID = %s, want msg-123", resp.Data.ID)
	}
	if resp.Data.Subject != "Test" {
		t.Errorf("Data.Subject = %s, want Test", resp.Data.Subject)
	}
}

func TestListResponse_JSON(t *testing.T) {
	jsonData := `{
		"data": [
			{"id": "msg-1"},
			{"id": "msg-2"}
		],
		"request_id": "req-789",
		"next_cursor": "cursor-abc"
	}`

	type TestItem struct {
		ID string `json:"id"`
	}

	var resp ListResponse[TestItem]
	if err := json.Unmarshal([]byte(jsonData), &resp); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if resp.RequestID != "req-789" {
		t.Errorf("RequestID = %s, want req-789", resp.RequestID)
	}
	if resp.NextCursor != "cursor-abc" {
		t.Errorf("NextCursor = %s, want cursor-abc", resp.NextCursor)
	}
	if len(resp.Data) != 2 {
		t.Errorf("len(Data) = %d, want 2", len(resp.Data))
	}
	if resp.Data[0].ID != "msg-1" {
		t.Errorf("Data[0].ID = %s, want msg-1", resp.Data[0].ID)
	}
}

func TestListResponse_EmptyData(t *testing.T) {
	jsonData := `{
		"data": [],
		"request_id": "req-empty"
	}`

	type TestItem struct {
		ID string `json:"id"`
	}

	var resp ListResponse[TestItem]
	if err := json.Unmarshal([]byte(jsonData), &resp); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if len(resp.Data) != 0 {
		t.Errorf("len(Data) = %d, want 0", len(resp.Data))
	}
	if resp.NextCursor != "" {
		t.Errorf("NextCursor = %s, want empty", resp.NextCursor)
	}
}

func TestListResponse_NoCursor(t *testing.T) {
	jsonData := `{
		"data": [{"id": "msg-1"}],
		"request_id": "req-nocursor"
	}`

	type TestItem struct {
		ID string `json:"id"`
	}

	var resp ListResponse[TestItem]
	if err := json.Unmarshal([]byte(jsonData), &resp); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if resp.NextCursor != "" {
		t.Errorf("NextCursor = %s, want empty string", resp.NextCursor)
	}
}
