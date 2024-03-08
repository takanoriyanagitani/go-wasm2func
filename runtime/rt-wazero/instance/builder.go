package instance

import (
	"context"

	inst "github.com/takanoriyanagitani/go-wasm2func/instance"

	wz "github.com/tetratelabs/wazero"
	wa "github.com/tetratelabs/wazero/api"
)

type Builder struct {
	wz.Runtime
	wz.ModuleConfig
}

func (b Builder) NewInstance(
	ctx context.Context,
	compiled wz.CompiledModule,
) (wa.Module, error) {
	return b.Runtime.InstantiateModule(ctx, compiled, b.ModuleConfig)
}

func (b Builder) AsIf() inst.Builder[wz.CompiledModule, wa.Module] { return b }
