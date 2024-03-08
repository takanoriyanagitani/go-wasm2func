package simple

import (
	"context"
)

type Input interface {
	Set(ctx context.Context, data []byte) error
}
