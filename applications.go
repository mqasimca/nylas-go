package nylas

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mqasimca/nylas-go/applications"
)

// GetDetails returns the application configuration details.
func (s *ApplicationsService) GetDetails(ctx context.Context) (*applications.ApplicationDetails, error) {
	path := "/v3/applications"

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("applications.GetDetails: %w", err)
	}

	var details applications.ApplicationDetails
	_, err = s.client.Do(req, &details)
	if err != nil {
		return nil, fmt.Errorf("applications.GetDetails: %w", err)
	}

	return &details, nil
}
