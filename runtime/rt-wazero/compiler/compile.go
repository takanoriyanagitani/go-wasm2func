package compiler

import (
	"context"

	comp "github.com/takanoriyanagitani/go-wasm2func/compiler"

	wz "github.com/tetratelabs/wazero"
)

type Runtime struct {
	wz.Runtime
}

func (r Runtime) Compile(
	ctx context.Context,
	wasmBytes []byte,
) (wz.CompiledModule, error) {
	return r.Runtime.CompileModule(ctx, wasmBytes)
}

func (r Runtime) AsIf() comp.Compiler[wz.CompiledModule] { return r }

func (r Runtime) Close() error { return r.Runtime.Close(context.Background()) }
