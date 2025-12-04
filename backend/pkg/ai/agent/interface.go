package chat

import (
	"context"
)

// Agent defines the interface for chat agents
type Agent[T any, V any] interface {
	Stream(ctx context.Context, req T) (<-chan V, error)
}
