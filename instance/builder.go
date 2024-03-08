package instance

import (
	"context"
)

type Builder[C, I any] interface {
	NewInstance(ctx context.Context, compiled C) (instance I, e error)
}
