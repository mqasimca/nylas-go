package webhooks

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
		limit := 25

		o := &ListOptions{
			Limit:     &limit,
			PageToken: "webhook-page-token",
		}
		v := o.Values()

		if v["limit"] != 25 {
			t.Errorf("expected limit=25, got %v", v["limit"])
		}
		if v["page_token"] != "webhook-page-token" {
			t.Errorf("expected page_token=webhook-page-token, got %v", v["page_token"])
		}
	})

	t.Run("partial options - limit only", func(t *testing.T) {
		limit := 50
		o := &ListOptions{
			Limit: &limit,
		}
		v := o.Values()

		if len(v) != 1 {
			t.Errorf("expected 1 entry, got %d", len(v))
		}
		if v["limit"] != 50 {
			t.Errorf("expected limit=50, got %v", v["limit"])
		}
	})

	t.Run("partial options - page_token only", func(t *testing.T) {
		o := &ListOptions{
			PageToken: "next-page",
		}
		v := o.Values()

		if len(v) != 1 {
			t.Errorf("expected 1 entry, got %d", len(v))
		}
		if v["page_token"] != "next-page" {
			t.Errorf("expected page_token=next-page, got %v", v["page_token"])
		}
	})
}
