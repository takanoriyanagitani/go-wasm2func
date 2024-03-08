package simple

import (
	"context"
)

type Output interface {
	Get(ctx context.Context) (result []byte, e error)
}
