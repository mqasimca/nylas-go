package nylas

import (
	"context"
	"errors"
	"testing"
)

func TestIterator_Next(t *testing.T) {
	pages := [][]string{
		{"a", "b", "c"},
		{"d", "e"},
	}
	pageIndex := 0

	iter := NewIterator(context.Background(), func(ctx context.Context, token string) ([]string, string, error) {
		if pageIndex >= len(pages) {
			return nil, "", nil
		}
		items := pages[pageIndex]
		pageIndex++
		nextToken := ""
		if pageIndex < len(pages) {
			nextToken = "next"
		}
		return items, nextToken, nil
	})

	var results []string
	for {
		item, err := iter.Next()
		if errors.Is(err, ErrDone) {
			break
		}
		if err != nil {
			t.Fatalf("Next() error = %v", err)
		}
		results = append(results, *item)
	}

	if len(results) != 5 {
		t.Errorf("Next() got %d items, want 5", len(results))
	}
}

func TestIterator_Collect(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}

	iter := NewIterator(context.Background(), func(ctx context.Context, token string) ([]int, string, error) {
		if token == "done" {
			return nil, "", nil
		}
		return items, "done", nil
	})

	results, err := iter.Collect()
	if err != nil {
		t.Fatalf("Collect() error = %v", err)
	}

	if len(results) != 5 {
		t.Errorf("Collect() got %d items, want 5", len(results))
	}

	for i, r := range results {
		if *r != items[i] {
			t.Errorf("Collect()[%d] = %d, want %d", i, *r, items[i])
		}
	}
}

func TestIterator_Error(t *testing.T) {
	expectedErr := errors.New("fetch error")

	iter := NewIterator(context.Background(), func(ctx context.Context, token string) ([]string, string, error) {
		return nil, "", expectedErr
	})

	_, err := iter.Next()
	if err != expectedErr {
		t.Errorf("Next() error = %v, want %v", err, expectedErr)
	}

	_, err = iter.Next()
	if err != expectedErr {
		t.Errorf("Next() after error = %v, want %v", err, expectedErr)
	}
}

func TestIterator_Reset(t *testing.T) {
	callCount := 0
	iter := NewIterator(context.Background(), func(ctx context.Context, token string) ([]string, string, error) {
		callCount++
		return []string{"a"}, "", nil
	})

	_, _ = iter.Next()
	_, _ = iter.Next()

	iter.Reset()

	_, err := iter.Next()
	if errors.Is(err, ErrDone) {
		t.Error("Next() after Reset() returned ErrDone")
	}

	if callCount != 2 {
		t.Errorf("fetch called %d times, want 2", callCount)
	}
}

func TestIterator_EmptyPage(t *testing.T) {
	iter := NewIterator(context.Background(), func(ctx context.Context, token string) ([]string, string, error) {
		return []string{}, "", nil
	})

	_, err := iter.Next()
	if !errors.Is(err, ErrDone) {
		t.Errorf("Next() on empty = %v, want ErrDone", err)
	}
}
