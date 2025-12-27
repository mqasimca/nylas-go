package attachments

import (
	"bytes"
	"io"
	"testing"
)

func TestAttachment_Structure(t *testing.T) {
	attachment := Attachment{
		ID:                 "attach-123",
		GrantID:            "grant-456",
		Filename:           "document.pdf",
		ContentType:        "application/pdf",
		Size:               12345,
		ContentID:          "cid:image001",
		ContentDisposition: "inline",
		IsInline:           true,
	}

	if attachment.ID != "attach-123" {
		t.Errorf("expected ID=attach-123, got %s", attachment.ID)
	}
	if attachment.GrantID != "grant-456" {
		t.Errorf("expected GrantID=grant-456, got %s", attachment.GrantID)
	}
	if attachment.Filename != "document.pdf" {
		t.Errorf("expected Filename=document.pdf, got %s", attachment.Filename)
	}
	if attachment.ContentType != "application/pdf" {
		t.Errorf("expected ContentType=application/pdf, got %s", attachment.ContentType)
	}
	if attachment.Size != 12345 {
		t.Errorf("expected Size=12345, got %d", attachment.Size)
	}
	if attachment.ContentID != "cid:image001" {
		t.Errorf("expected ContentID=cid:image001, got %s", attachment.ContentID)
	}
	if attachment.ContentDisposition != "inline" {
		t.Errorf("expected ContentDisposition=inline, got %s", attachment.ContentDisposition)
	}
	if !attachment.IsInline {
		t.Error("expected IsInline=true")
	}
}

func TestAttachment_Defaults(t *testing.T) {
	// Test zero values
	attachment := Attachment{}

	if attachment.ID != "" {
		t.Errorf("expected empty ID, got %s", attachment.ID)
	}
	if attachment.Size != 0 {
		t.Errorf("expected Size=0, got %d", attachment.Size)
	}
	if attachment.IsInline {
		t.Error("expected IsInline=false by default")
	}
}

func TestDownloadResponse_Structure(t *testing.T) {
	content := []byte("test file content")
	reader := io.NopCloser(bytes.NewReader(content))

	resp := DownloadResponse{
		Content:     reader,
		ContentType: "text/plain",
		Filename:    "test.txt",
		Size:        int64(len(content)),
	}

	if resp.ContentType != "text/plain" {
		t.Errorf("expected ContentType=text/plain, got %s", resp.ContentType)
	}
	if resp.Filename != "test.txt" {
		t.Errorf("expected Filename=test.txt, got %s", resp.Filename)
	}
	if resp.Size != int64(len(content)) {
		t.Errorf("expected Size=%d, got %d", len(content), resp.Size)
	}

	// Read content
	data, err := io.ReadAll(resp.Content)
	if err != nil {
		t.Fatalf("failed to read content: %v", err)
	}
	if string(data) != "test file content" {
		t.Errorf("expected content 'test file content', got '%s'", string(data))
	}

	// Close should work without error
	if err := resp.Content.Close(); err != nil {
		t.Errorf("unexpected error closing content: %v", err)
	}
}

func TestAttachment_InlineVsRegular(t *testing.T) {
	tests := []struct {
		name               string
		contentDisposition string
		isInline           bool
	}{
		{
			name:               "inline attachment",
			contentDisposition: "inline",
			isInline:           true,
		},
		{
			name:               "regular attachment",
			contentDisposition: "attachment",
			isInline:           false,
		},
		{
			name:               "no disposition",
			contentDisposition: "",
			isInline:           false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Attachment{
				ContentDisposition: tt.contentDisposition,
				IsInline:           tt.isInline,
			}

			if a.ContentDisposition != tt.contentDisposition {
				t.Errorf("expected ContentDisposition=%s, got %s", tt.contentDisposition, a.ContentDisposition)
			}
			if a.IsInline != tt.isInline {
				t.Errorf("expected IsInline=%v, got %v", tt.isInline, a.IsInline)
			}
		})
	}
}

func TestAttachment_CommonMimeTypes(t *testing.T) {
	mimeTypes := []struct {
		filename    string
		contentType string
	}{
		{"document.pdf", "application/pdf"},
		{"image.png", "image/png"},
		{"image.jpg", "image/jpeg"},
		{"spreadsheet.xlsx", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"},
		{"document.docx", "application/vnd.openxmlformats-officedocument.wordprocessingml.document"},
		{"archive.zip", "application/zip"},
		{"text.txt", "text/plain"},
	}

	for _, mt := range mimeTypes {
		t.Run(mt.filename, func(t *testing.T) {
			a := Attachment{
				Filename:    mt.filename,
				ContentType: mt.contentType,
			}

			if a.Filename != mt.filename {
				t.Errorf("expected Filename=%s, got %s", mt.filename, a.Filename)
			}
			if a.ContentType != mt.contentType {
				t.Errorf("expected ContentType=%s, got %s", mt.contentType, a.ContentType)
			}
		})
	}
}
