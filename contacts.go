package nylas

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mqasimca/nylas-go/contacts"
)

// List returns contacts for a grant.
func (s *ContactsService) List(ctx context.Context, grantID string, opts *contacts.ListOptions) (*ListResponse[contacts.Contact], error) {
	path := fmt.Sprintf("/v3/grants/%s/contacts", grantID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("contacts.List: %w", err)
	}

	if opts != nil {
		q := req.URL.Query()
		setQueryParams(q, opts.Values())
		req.URL.RawQuery = q.Encode()
	}

	var data []contacts.Contact
	nextCursor, requestID, err := s.client.DoList(req, &data)
	if err != nil {
		return nil, fmt.Errorf("contacts.List: %w", err)
	}

	return &ListResponse[contacts.Contact]{
		Data:       data,
		NextCursor: nextCursor,
		RequestID:  requestID,
	}, nil
}

// Get returns a single contact.
func (s *ContactsService) Get(ctx context.Context, grantID, contactID string) (*contacts.Contact, error) {
	path := fmt.Sprintf("/v3/grants/%s/contacts/%s", grantID, contactID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("contacts.Get(%s): %w", contactID, err)
	}

	var contact contacts.Contact
	_, err = s.client.Do(req, &contact)
	if err != nil {
		return nil, fmt.Errorf("contacts.Get(%s): %w", contactID, err)
	}

	return &contact, nil
}

// Create creates a new contact.
func (s *ContactsService) Create(ctx context.Context, grantID string, create *contacts.CreateRequest) (*contacts.Contact, error) {
	path := fmt.Sprintf("/v3/grants/%s/contacts", grantID)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, create)
	if err != nil {
		return nil, fmt.Errorf("contacts.Create: %w", err)
	}

	var contact contacts.Contact
	_, err = s.client.Do(req, &contact)
	if err != nil {
		return nil, fmt.Errorf("contacts.Create: %w", err)
	}

	return &contact, nil
}

// Update updates a contact.
func (s *ContactsService) Update(ctx context.Context, grantID, contactID string, update *contacts.UpdateRequest) (*contacts.Contact, error) {
	path := fmt.Sprintf("/v3/grants/%s/contacts/%s", grantID, contactID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, update)
	if err != nil {
		return nil, fmt.Errorf("contacts.Update(%s): %w", contactID, err)
	}

	var contact contacts.Contact
	_, err = s.client.Do(req, &contact)
	if err != nil {
		return nil, fmt.Errorf("contacts.Update(%s): %w", contactID, err)
	}

	return &contact, nil
}

// Delete deletes a contact.
func (s *ContactsService) Delete(ctx context.Context, grantID, contactID string) error {
	path := fmt.Sprintf("/v3/grants/%s/contacts/%s", grantID, contactID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return fmt.Errorf("contacts.Delete(%s): %w", contactID, err)
	}

	_, err = s.client.Do(req, nil)
	if err != nil {
		return fmt.Errorf("contacts.Delete(%s): %w", contactID, err)
	}

	return nil
}

// ListAll returns an iterator for all contacts.
func (s *ContactsService) ListAll(ctx context.Context, grantID string, opts *contacts.ListOptions) *Iterator[contacts.Contact] {
	return NewIterator(ctx, func(ctx context.Context, pageToken string) ([]contacts.Contact, string, error) {
		o := opts
		if o == nil {
			o = &contacts.ListOptions{}
		}
		o.PageToken = pageToken

		resp, err := s.List(ctx, grantID, o)
		if err != nil {
			return nil, "", err
		}
		return resp.Data, resp.NextCursor, nil
	})
}

// ListGroups returns contact groups for a grant.
func (s *ContactsService) ListGroups(ctx context.Context, grantID string) ([]contacts.Group, error) {
	path := fmt.Sprintf("/v3/grants/%s/contacts/groups", grantID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("contacts.ListGroups: %w", err)
	}

	var groups []contacts.Group
	_, err = s.client.Do(req, &groups)
	if err != nil {
		return nil, fmt.Errorf("contacts.ListGroups: %w", err)
	}

	return groups, nil
}
