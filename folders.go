package nylas

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mqasimca/nylas-go/folders"
)

// List returns folders for a grant.
func (s *FoldersService) List(ctx context.Context, grantID string, opts *folders.ListOptions) (*ListResponse[folders.Folder], error) {
	path := fmt.Sprintf("/v3/grants/%s/folders", grantID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("folders.List: %w", err)
	}

	if opts != nil {
		q := req.URL.Query()
		setQueryParams(q, opts.Values())
		req.URL.RawQuery = q.Encode()
	}

	var data []folders.Folder
	nextCursor, requestID, err := s.client.DoList(req, &data)
	if err != nil {
		return nil, fmt.Errorf("folders.List: %w", err)
	}

	return &ListResponse[folders.Folder]{
		Data:       data,
		NextCursor: nextCursor,
		RequestID:  requestID,
	}, nil
}

// Get returns a single folder.
func (s *FoldersService) Get(ctx context.Context, grantID, folderID string) (*folders.Folder, error) {
	path := fmt.Sprintf("/v3/grants/%s/folders/%s", grantID, folderID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("folders.Get(%s): %w", folderID, err)
	}

	var folder folders.Folder
	_, err = s.client.Do(req, &folder)
	if err != nil {
		return nil, fmt.Errorf("folders.Get(%s): %w", folderID, err)
	}

	return &folder, nil
}

// Create creates a new folder.
func (s *FoldersService) Create(ctx context.Context, grantID string, create *folders.CreateRequest) (*folders.Folder, error) {
	path := fmt.Sprintf("/v3/grants/%s/folders", grantID)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, create)
	if err != nil {
		return nil, fmt.Errorf("folders.Create: %w", err)
	}

	var folder folders.Folder
	_, err = s.client.Do(req, &folder)
	if err != nil {
		return nil, fmt.Errorf("folders.Create: %w", err)
	}

	return &folder, nil
}

// Update updates a folder.
func (s *FoldersService) Update(ctx context.Context, grantID, folderID string, update *folders.UpdateRequest) (*folders.Folder, error) {
	path := fmt.Sprintf("/v3/grants/%s/folders/%s", grantID, folderID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, update)
	if err != nil {
		return nil, fmt.Errorf("folders.Update(%s): %w", folderID, err)
	}

	var folder folders.Folder
	_, err = s.client.Do(req, &folder)
	if err != nil {
		return nil, fmt.Errorf("folders.Update(%s): %w", folderID, err)
	}

	return &folder, nil
}

// Delete deletes a folder.
func (s *FoldersService) Delete(ctx context.Context, grantID, folderID string) error {
	path := fmt.Sprintf("/v3/grants/%s/folders/%s", grantID, folderID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return fmt.Errorf("folders.Delete(%s): %w", folderID, err)
	}

	_, err = s.client.Do(req, nil)
	if err != nil {
		return fmt.Errorf("folders.Delete(%s): %w", folderID, err)
	}

	return nil
}

// ListAll returns an iterator for all folders.
func (s *FoldersService) ListAll(ctx context.Context, grantID string, opts *folders.ListOptions) *Iterator[folders.Folder] {
	return NewIterator(ctx, func(ctx context.Context, pageToken string) ([]folders.Folder, string, error) {
		o := opts
		if o == nil {
			o = &folders.ListOptions{}
		}
		o.PageToken = pageToken

		resp, err := s.List(ctx, grantID, o)
		if err != nil {
			return nil, "", err
		}
		return resp.Data, resp.NextCursor, nil
	})
}
