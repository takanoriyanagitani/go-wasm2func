package function

import (
	"context"

	fn "github.com/takanoriyanagitani/go-wasm2func/function"

	wa "github.com/tetratelabs/wazero/api"
)

type Builder[F any] func(context.Context, wa.Module) (fn F, e error)

func (b Builder[F]) Build(ctx context.Context, i wa.Module) (F, error) {
	return b(ctx, i)
}

func (b Builder[F]) AsIf() fn.Builder[wa.Module, F] { return b }
