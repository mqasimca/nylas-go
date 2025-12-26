package contacts

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
		limit := 10
		email := "test@example.com"
		phone := "+1234567890"
		source := "google"
		group := "group-1"
		recurse := true

		o := &ListOptions{
			Limit:       &limit,
			PageToken:   "token123",
			Email:       &email,
			PhoneNumber: &phone,
			Source:      &source,
			Group:       &group,
			Recurse:     &recurse,
		}
		v := o.Values()

		if v["limit"] != 10 {
			t.Errorf("expected limit=10, got %v", v["limit"])
		}
		if v["page_token"] != "token123" {
			t.Errorf("expected page_token=token123, got %v", v["page_token"])
		}
		if v["email"] != email {
			t.Errorf("expected email=%s, got %v", email, v["email"])
		}
		if v["phone_number"] != phone {
			t.Errorf("expected phone_number=%s, got %v", phone, v["phone_number"])
		}
		if v["source"] != source {
			t.Errorf("expected source=%s, got %v", source, v["source"])
		}
		if v["group"] != group {
			t.Errorf("expected group=%s, got %v", group, v["group"])
		}
		if v["recurse"] != true {
			t.Errorf("expected recurse=true, got %v", v["recurse"])
		}
	})

	t.Run("partial options", func(t *testing.T) {
		limit := 5
		o := &ListOptions{
			Limit:     &limit,
			PageToken: "abc",
		}
		v := o.Values()

		if len(v) != 2 {
			t.Errorf("expected 2 entries, got %d", len(v))
		}
		if v["limit"] != 5 {
			t.Errorf("expected limit=5, got %v", v["limit"])
		}
		if v["page_token"] != "abc" {
			t.Errorf("expected page_token=abc, got %v", v["page_token"])
		}
	})
}
