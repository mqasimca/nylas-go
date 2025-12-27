package smartcompose

import (
	"encoding/json"
	"testing"
)

func TestComposeRequest_JSON(t *testing.T) {
	req := ComposeRequest{
		Prompt: "Write a professional email greeting",
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	var decoded ComposeRequest
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if decoded.Prompt != req.Prompt {
		t.Errorf("Prompt = %s, want %s", decoded.Prompt, req.Prompt)
	}
}

func TestComposeResponse_JSON(t *testing.T) {
	jsonData := `{
		"suggestion": "Hello, I hope this email finds you well. I wanted to reach out regarding..."
	}`

	var resp ComposeResponse
	if err := json.Unmarshal([]byte(jsonData), &resp); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	expected := "Hello, I hope this email finds you well. I wanted to reach out regarding..."
	if resp.Suggestion != expected {
		t.Errorf("Suggestion = %s, want %s", resp.Suggestion, expected)
	}
}

func TestComposeRequest_EmptyPrompt(t *testing.T) {
	req := ComposeRequest{
		Prompt: "",
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	var decoded ComposeRequest
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if decoded.Prompt != "" {
		t.Errorf("Prompt = %s, want empty string", decoded.Prompt)
	}
}

func TestComposeResponse_EmptySuggestion(t *testing.T) {
	jsonData := `{"suggestion": ""}`

	var resp ComposeResponse
	if err := json.Unmarshal([]byte(jsonData), &resp); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if resp.Suggestion != "" {
		t.Errorf("Suggestion = %s, want empty string", resp.Suggestion)
	}
}
