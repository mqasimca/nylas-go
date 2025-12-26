package nylas

import (
	"context"
	"errors"
)

// ErrDone is returned by Iterator.Next when iteration is complete.
var ErrDone = errors.New("no more items")

// Iterator provides paginated iteration over API resources.
// Use Next() to get items one at a time, or Collect() to get all items at once.
//
// Example:
//
//	iter := client.Messages.ListAll(ctx, grantID, nil)
//	for {
//	    msg, err := iter.Next()
//	    if errors.Is(err, nylas.ErrDone) {
//	        break
//	    }
//	    if err != nil {
//	        return err
//	    }
//	    process(msg)
//	}
type Iterator[T any] struct {
	fetch     func(ctx context.Context, pageToken string) ([]T, string, error)
	ctx       context.Context
	buffer    []T
	pageToken string
	index     int
	done      bool
	err       error
}

// NewIterator creates a new Iterator with the given fetch function.
func NewIterator[T any](ctx context.Context, fetch func(context.Context, string) ([]T, string, error)) *Iterator[T] {
	return &Iterator[T]{
		ctx:   ctx,
		fetch: fetch,
	}
}

// Next returns the next item in the iteration. Returns ErrDone when there are no more items.
func (it *Iterator[T]) Next() (*T, error) {
	if it.err != nil {
		return nil, it.err
	}

	if it.index < len(it.buffer) {
		item := &it.buffer[it.index]
		it.index++
		return item, nil
	}

	if it.done {
		return nil, ErrDone
	}

	items, nextToken, err := it.fetch(it.ctx, it.pageToken)
	if err != nil {
		it.err = err
		return nil, err
	}

	if len(items) == 0 {
		it.done = true
		return nil, ErrDone
	}

	it.buffer = items
	it.pageToken = nextToken
	it.index = 1
	it.done = nextToken == ""

	return &it.buffer[0], nil
}

// Collect returns all remaining items as a slice. Useful when you need all items at once.
func (it *Iterator[T]) Collect() ([]*T, error) {
	var all []*T
	for {
		item, err := it.Next()
		if errors.Is(err, ErrDone) {
			return all, nil
		}
		if err != nil {
			return all, err
		}
		all = append(all, item)
	}
}

// Reset clears the iterator state, allowing iteration to start over from the beginning.
func (it *Iterator[T]) Reset() {
	it.buffer = nil
	it.pageToken = ""
	it.index = 0
	it.done = false
	it.err = nil
}
