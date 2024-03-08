package function

import (
	"context"

	util "github.com/takanoriyanagitani/go-wasm2func/util"

	inst "github.com/takanoriyanagitani/go-wasm2func/instance"
)

type Builder[I, F any] interface {
	Build(ctx context.Context, instance I) (fn F, e error)
}

type BuildFn[I, F any] func(context.Context, I) (F, error)

func (b BuildFn[I, F]) Build(
	ctx context.Context,
	instance I,
) (F, error) {
	return b(ctx, instance)
}

func (b BuildFn[I, F]) AsIf() Builder[I, F] { return b }

func FromCompiled[C, I, F any](
	ctx context.Context,
	compiled C,
	ib inst.Builder[C, I],
	fb Builder[I, F],
) (fn F, e error) {
	return util.ComposeCtx(
		ib.NewInstance,
		fb.Build,
	)(ctx, compiled)
}
