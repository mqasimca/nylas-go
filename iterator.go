package nylas

import (
	"context"
	"errors"
)

// ErrDone is returned by Iterator.Next when iteration is complete.
var ErrDone = errors.New("no more items")

type Iterator[T any] struct {
	fetch     func(ctx context.Context, pageToken string) ([]T, string, error)
	ctx       context.Context
	buffer    []T
	pageToken string
	index     int
	done      bool
	err       error
}

func NewIterator[T any](ctx context.Context, fetch func(context.Context, string) ([]T, string, error)) *Iterator[T] {
	return &Iterator[T]{
		ctx:   ctx,
		fetch: fetch,
	}
}

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

func (it *Iterator[T]) Reset() {
	it.buffer = nil
	it.pageToken = ""
	it.index = 0
	it.done = false
	it.err = nil
}
