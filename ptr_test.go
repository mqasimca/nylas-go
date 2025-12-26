package nylas

import "testing"

func TestPtr(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		p := Ptr(42)
		if *p != 42 {
			t.Errorf("Ptr(42) = %v, want 42", *p)
		}
	})

	t.Run("string", func(t *testing.T) {
		p := Ptr("hello")
		if *p != "hello" {
			t.Errorf("Ptr(hello) = %v, want hello", *p)
		}
	})

	t.Run("bool", func(t *testing.T) {
		p := Ptr(true)
		if *p != true {
			t.Errorf("Ptr(true) = %v, want true", *p)
		}
	})

	t.Run("struct", func(t *testing.T) {
		type item struct{ ID string }
		p := Ptr(item{ID: "123"})
		if p.ID != "123" {
			t.Errorf("Ptr(item).ID = %v, want 123", p.ID)
		}
	})
}
