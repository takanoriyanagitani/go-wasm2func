package compiler

import (
	"context"
	"io"
)

type Compiler[C any] interface {
	Compile(ctx context.Context, wasmBytes []byte) (compiled C, e error)
	io.Closer
}
