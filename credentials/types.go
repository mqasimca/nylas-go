// Package credentials provides types for the Nylas Credentials API.
package credentials

import "github.com/mqasimca/nylas-go/connectors"

// CredentialType represents the type of credential.
type CredentialType string

const (
	CredentialTypeAdminConsent   CredentialType = "adminconsent"
	CredentialTypeServiceAccount CredentialType = "serviceaccount"
	CredentialTypeConnector      CredentialType = "connector"
)

// Credential represents a connector credential.
type Credential struct {
	ID             string         `json:"id"`
	Name           string         `json:"name"`
	CredentialType CredentialType `json:"credential_type,omitempty"`
	HashedData     string         `json:"hashed_data,omitempty"`
	CreatedAt      int64          `json:"created_at,omitempty"`
	UpdatedAt      int64          `json:"updated_at,omitempty"`
}

// CreateRequest contains the data to create a credential.
type CreateRequest struct {
	Name           string                 `json:"name"`
	CredentialType CredentialType         `json:"credential_type"`
	CredentialData map[string]interface{} `json:"credential_data"`
}

// MicrosoftCredentialData contains Microsoft admin consent credential data.
type MicrosoftCredentialData struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

// GoogleCredentialData contains Google service account credential data.
type GoogleCredentialData struct {
	PrivateKeyID string `json:"private_key_id"`
	PrivateKey   string `json:"private_key"`
	ClientEmail  string `json:"client_email"`
}

// CreateMicrosoftRequest creates a Microsoft admin consent credential request.
func CreateMicrosoftRequest(name, clientID, clientSecret string) *CreateRequest {
	return &CreateRequest{
		Name:           name,
		CredentialType: CredentialTypeAdminConsent,
		CredentialData: map[string]interface{}{
			"client_id":     clientID,
			"client_secret": clientSecret,
		},
	}
}

// CreateGoogleRequest creates a Google service account credential request.
func CreateGoogleRequest(name, privateKeyID, privateKey, clientEmail string) *CreateRequest {
	return &CreateRequest{
		Name:           name,
		CredentialType: CredentialTypeServiceAccount,
		CredentialData: map[string]interface{}{
			"private_key_id": privateKeyID,
			"private_key":    privateKey,
			"client_email":   clientEmail,
		},
	}
}

// UpdateRequest contains the data to update a credential.
type UpdateRequest struct {
	Name           *string                `json:"name,omitempty"`
	CredentialData map[string]interface{} `json:"credential_data,omitempty"`
}

// ListOptions contains options for listing credentials.
type ListOptions struct {
	Limit     *int                `json:"limit,omitempty"`
	PageToken string              `json:"page_token,omitempty"`
	Provider  connectors.Provider `json:"provider,omitempty"`
}

// Values returns the options as a map for query parameters.
func (o *ListOptions) Values() map[string]any {
	if o == nil {
		return nil
	}
	m := make(map[string]any)
	if o.Limit != nil {
		m["limit"] = *o.Limit
	}
	if o.PageToken != "" {
		m["page_token"] = o.PageToken
	}
	return m
}
