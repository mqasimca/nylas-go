package nylas

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApplicationsService_GetDetails(t *testing.T) {
	tests := []struct {
		name       string
		response   string
		statusCode int
		wantErr    bool
		wantAppID  string
	}{
		{
			name: "success",
			response: `{
				"request_id": "req-123",
				"data": {
					"application_id": "app-123",
					"organization_id": "org-456",
					"region": "us",
					"environment": "production"
				}
			}`,
			statusCode: http.StatusOK,
			wantErr:    false,
			wantAppID:  "app-123",
		},
		{
			name:       "unauthorized",
			response:   `{"error": {"type": "unauthorized", "message": "Invalid API key"}}`,
			statusCode: http.StatusUnauthorized,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("Method = %s, want GET", r.Method)
				}
				if r.URL.Path != "/v3/applications" {
					t.Errorf("Path = %s, want /v3/applications", r.URL.Path)
				}
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.response))
			}))
			defer srv.Close()

			client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
			details, err := client.Applications.GetDetails(context.Background())

			if (err != nil) != tt.wantErr {
				t.Errorf("GetDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && details.ApplicationID != tt.wantAppID {
				t.Errorf("ApplicationID = %s, want %s", details.ApplicationID, tt.wantAppID)
			}
		})
	}
}
