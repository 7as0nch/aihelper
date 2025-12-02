package chat

import (
	"context"
)

// Agent defines the interface for chat agents
type Agent[K comparable, V any] interface {
	Stream(ctx context.Context, req K) (<-chan V, error)
}
