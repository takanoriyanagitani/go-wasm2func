package compiler

import (
	"context"
	"io"
	"io/fs"
)

type Compiler[C any] interface {
	Compile(ctx context.Context, wasmBytes []byte) (compiled C, e error)
	io.Closer
}

func CompileFromPath[C any](
	ctx context.Context,
	compiler Compiler[C],
	wasmPath string,
	f fs.FS,
) (compiled C, e error) {
	wasmBytes, e := fs.ReadFile(f, wasmPath)
	if nil != e {
		return
	}
	return compiler.Compile(ctx, wasmBytes)
}
