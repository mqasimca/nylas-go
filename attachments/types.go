package attachments

import "io"

// Attachment represents an email attachment in the Nylas API.
type Attachment struct {
	// ID is the unique identifier for this attachment.
	ID string `json:"id"`
	// GrantID is the ID of the grant (connected account) this attachment belongs to.
	GrantID string `json:"grant_id"`
	// Filename is the original filename of the attachment.
	Filename string `json:"filename,omitempty"`
	// ContentType is the MIME type (e.g., "application/pdf", "image/png").
	ContentType string `json:"content_type,omitempty"`
	// Size is the file size in bytes.
	Size int `json:"size,omitempty"`
	// ContentID is the CID for inline attachments (used in HTML email bodies).
	ContentID string `json:"content_id,omitempty"`
	// ContentDisposition is "attachment" or "inline".
	ContentDisposition string `json:"content_disposition,omitempty"`
	// IsInline indicates whether this attachment is embedded in the message body.
	IsInline bool `json:"is_inline,omitempty"`
}

// DownloadResponse represents the response from downloading an attachment.
type DownloadResponse struct {
	// Content is the readable stream of attachment data. Caller must close.
	Content io.ReadCloser
	// ContentType is the MIME type of the attachment.
	ContentType string
	// Filename is the original filename.
	Filename string
	// Size is the file size in bytes.
	Size int64
}
