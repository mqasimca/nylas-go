package folders

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
		limit := 50
		parentID := "parent-123"

		o := &ListOptions{
			Limit:     &limit,
			PageToken: "token456",
			ParentID:  &parentID,
		}
		v := o.Values()

		if v["limit"] != 50 {
			t.Errorf("expected limit=50, got %v", v["limit"])
		}
		if v["page_token"] != "token456" {
			t.Errorf("expected page_token=token456, got %v", v["page_token"])
		}
		if v["parent_id"] != parentID {
			t.Errorf("expected parent_id=%s, got %v", parentID, v["parent_id"])
		}
	})

	t.Run("partial options", func(t *testing.T) {
		limit := 25
		o := &ListOptions{
			Limit: &limit,
		}
		v := o.Values()

		if len(v) != 1 {
			t.Errorf("expected 1 entry, got %d", len(v))
		}
		if v["limit"] != 25 {
			t.Errorf("expected limit=25, got %v", v["limit"])
		}
	})
}
