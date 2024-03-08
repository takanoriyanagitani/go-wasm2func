package simple

import (
	"context"
)

type Caller interface {
	Call(ctx context.Context) error
}
