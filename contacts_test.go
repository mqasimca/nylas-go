package nylas

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mqasimca/nylas-go/contacts"
)

func TestContactsService_List(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		response   string
		wantErr    bool
	}{
		{
			name:       "success",
			statusCode: http.StatusOK,
			response: `{
				"request_id": "req-123",
				"data": [
					{
						"id": "contact-1",
						"grant_id": "grant-123",
						"given_name": "John",
						"surname": "Doe",
						"emails": [{"email": "john@example.com", "type": "work"}]
					}
				],
				"next_cursor": "next-page"
			}`,
			wantErr: false,
		},
		{
			name:       "unauthorized",
			statusCode: http.StatusUnauthorized,
			response:   `{"error": "unauthorized"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected GET, got %s", r.Method)
				}
				if r.URL.Path != "/v3/grants/grant-123/contacts" {
					t.Errorf("unexpected path: %s", r.URL.Path)
				}
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.response))
			}))
			defer srv.Close()

			client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
			resp, err := client.Contacts.List(context.Background(), "grant-123", nil)

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && len(resp.Data) == 0 {
				t.Error("expected contacts in response")
			}
		})
	}
}

func TestContactsService_Get(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		response   string
		wantErr    bool
	}{
		{
			name:       "success",
			statusCode: http.StatusOK,
			response: `{
				"request_id": "req-123",
				"data": {
					"id": "contact-1",
					"grant_id": "grant-123",
					"given_name": "John",
					"surname": "Doe"
				}
			}`,
			wantErr: false,
		},
		{
			name:       "not found",
			statusCode: http.StatusNotFound,
			response:   `{"error": "not found"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected GET, got %s", r.Method)
				}
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.response))
			}))
			defer srv.Close()

			client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
			contact, err := client.Contacts.Get(context.Background(), "grant-123", "contact-1")

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && contact.ID != "contact-1" {
				t.Errorf("expected contact-1, got %s", contact.ID)
			}
		})
	}
}

func TestContactsService_Create(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		var req contacts.CreateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("failed to decode request: %v", err)
		}
		if req.GivenName != "Jane" {
			t.Errorf("expected GivenName=Jane, got %s", req.GivenName)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"request_id": "req-123",
			"data": {
				"id": "contact-new",
				"grant_id": "grant-123",
				"given_name": "Jane",
				"surname": "Smith"
			}
		}`))
	}))
	defer srv.Close()

	client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
	contact, err := client.Contacts.Create(context.Background(), "grant-123", &contacts.CreateRequest{
		GivenName: "Jane",
		Surname:   "Smith",
	})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if contact.ID != "contact-new" {
		t.Errorf("expected contact-new, got %s", contact.ID)
	}
}

func TestContactsService_Update(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"request_id": "req-123",
			"data": {
				"id": "contact-1",
				"grant_id": "grant-123",
				"given_name": "John",
				"surname": "Updated"
			}
		}`))
	}))
	defer srv.Close()

	client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
	surname := "Updated"
	contact, err := client.Contacts.Update(context.Background(), "grant-123", "contact-1", &contacts.UpdateRequest{
		Surname: &surname,
	})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if contact.Surname != "Updated" {
		t.Errorf("expected Updated, got %s", contact.Surname)
	}
}

func TestContactsService_Delete(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"request_id": "req-123"}`))
	}))
	defer srv.Close()

	client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
	err := client.Contacts.Delete(context.Background(), "grant-123", "contact-1")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestContactsService_ListGroups(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/v3/grants/grant-123/contacts/groups" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"request_id": "req-123",
			"data": [
				{"id": "group-1", "name": "Work"},
				{"id": "group-2", "name": "Friends"}
			]
		}`))
	}))
	defer srv.Close()

	client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
	groups, err := client.Contacts.ListGroups(context.Background(), "grant-123")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(groups) != 2 {
		t.Errorf("expected 2 groups, got %d", len(groups))
	}
}

func TestContactsService_ListAll(t *testing.T) {
	callCount := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"request_id": "req-1",
				"data": [{"id": "contact-1", "grant_id": "grant-123"}],
				"next_cursor": "page2"
			}`))
		} else {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"request_id": "req-2",
				"data": [{"id": "contact-2", "grant_id": "grant-123"}],
				"next_cursor": ""
			}`))
		}
	}))
	defer srv.Close()

	client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
	iter := client.Contacts.ListAll(context.Background(), "grant-123", nil)

	all, err := iter.Collect()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(all) != 2 {
		t.Errorf("expected 2 contacts, got %d", len(all))
	}
}
