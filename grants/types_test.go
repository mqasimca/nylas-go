package grants

import (
	"testing"
)

func TestListOptions_Values(t *testing.T) {
	t.Run("nil options", func(t *testing.T) {
		var o *ListOptions
		if o.Values() != nil {
			t.Error("expected nil for nil options")
		}
	})

	t.Run("empty options", func(t *testing.T) {
		o := &ListOptions{}
		v := o.Values()
		if len(v) != 0 {
			t.Errorf("expected empty map, got %d entries", len(v))
		}
	})

	t.Run("with all options", func(t *testing.T) {
		limit := 100
		offset := 20
		sortBy := "created_at"
		orderBy := "desc"
		since := int64(1700000000)
		before := int64(1710000000)
		email := "user@example.com"
		grantStatus := "valid"
		ip := "192.168.1.1"
		provider := "google"

		o := &ListOptions{
			Limit:       &limit,
			Offset:      &offset,
			SortBy:      &sortBy,
			OrderBy:     &orderBy,
			Since:       &since,
			Before:      &before,
			Email:       &email,
			GrantStatus: &grantStatus,
			IP:          &ip,
			Provider:    &provider,
		}
		v := o.Values()

		if v["limit"] != 100 {
			t.Errorf("expected limit=100, got %v", v["limit"])
		}
		if v["offset"] != 20 {
			t.Errorf("expected offset=20, got %v", v["offset"])
		}
		if v["sort_by"] != sortBy {
			t.Errorf("expected sort_by=%s, got %v", sortBy, v["sort_by"])
		}
		if v["order_by"] != orderBy {
			t.Errorf("expected order_by=%s, got %v", orderBy, v["order_by"])
		}
		if v["since"] != since {
			t.Errorf("expected since=%d, got %v", since, v["since"])
		}
		if v["before"] != before {
			t.Errorf("expected before=%d, got %v", before, v["before"])
		}
		if v["email"] != email {
			t.Errorf("expected email=%s, got %v", email, v["email"])
		}
		if v["grant_status"] != grantStatus {
			t.Errorf("expected grant_status=%s, got %v", grantStatus, v["grant_status"])
		}
		if v["ip"] != ip {
			t.Errorf("expected ip=%s, got %v", ip, v["ip"])
		}
		if v["provider"] != provider {
			t.Errorf("expected provider=%s, got %v", provider, v["provider"])
		}
	})

	t.Run("partial options", func(t *testing.T) {
		limit := 10
		provider := "microsoft"
		o := &ListOptions{
			Limit:    &limit,
			Provider: &provider,
		}
		v := o.Values()

		if len(v) != 2 {
			t.Errorf("expected 2 entries, got %d", len(v))
		}
		if v["limit"] != 10 {
			t.Errorf("expected limit=10, got %v", v["limit"])
		}
		if v["provider"] != provider {
			t.Errorf("expected provider=%s, got %v", provider, v["provider"])
		}
	})
}
