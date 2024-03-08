package main

import (
	"context"
	"fmt"
	"io/fs"
	"os"

	wz "github.com/tetratelabs/wazero"
	wa "github.com/tetratelabs/wazero/api"

	ua "github.com/takanoriyanagitani/go-wasm2func/util/arr"

	wfc "github.com/takanoriyanagitani/go-wasm2func/compiler"
	wf2 "github.com/takanoriyanagitani/go-wasm2func/function"
	wfi "github.com/takanoriyanagitani/go-wasm2func/instance"

	wzc "github.com/takanoriyanagitani/go-wasm2func/runtime/rt-wazero/compiler"
	wzf "github.com/takanoriyanagitani/go-wasm2func/runtime/rt-wazero/function"
	wzi "github.com/takanoriyanagitani/go-wasm2func/runtime/rt-wazero/instance"

	// revive:disable:line-length-limit
	wzfb "github.com/takanoriyanagitani/go-wasm2func/runtime/rt-wazero/function/binary"
	// revive:enable:line-length-limit
)

var rtime wz.Runtime = wz.NewRuntime(context.Background())
var wtime wzc.Runtime = wzc.Runtime{
	Runtime: rtime,
}
var wcomp wfc.Compiler[wz.CompiledModule] = wtime.AsIf()

var fsys fs.FS = os.DirFS("./wasm.d")

const pwasDefault string = "rs_mul.wasm"

func getenvOrAlt(key string, alt string) string {
	val, ok := os.LookupEnv(key)
	switch ok {
	case true:
		return val
	default:
		return alt
	}
}

var pwas string = getenvOrAlt("ENV_WASM_PATH", pwasDefault)

func must[T any](t T, e error) T {
	if nil == e {
		return t
	}
	panic(e)
}

func mustNil(e error) {
	if nil != e {
		panic(e)
	}
}

var compiled wz.CompiledModule = must(wfc.CompileFromPath(
	context.Background(),
	wcomp,
	pwas,
	fsys,
))

var builder wzi.Builder = wzi.Builder{
	Runtime:      rtime,
	ModuleConfig: wz.NewModuleConfig(),
}
var ibld wfi.Builder[wz.CompiledModule, wa.Module] = builder.AsIf()

var mdl wa.Module = must(builder.NewInstance(
	context.Background(),
	compiled,
))

type MulFn struct {
	mul32i func(int32, int32) (int32, error)
	mul64i func(int64, int64) (int64, error)

	mul32f func(float32, float32) (float32, error)
	mul64f func(float64, float64) (float64, error)
}

var fbld wf2.Builder[wa.Module, MulFn] = wzf.Builder[MulFn](
	func(_ context.Context, m wa.Module) (MulFn, error) {
		var mf MulFn

		var m5i wa.Function = m.ExportedFunction("mul32i")
		var m6i wa.Function = m.ExportedFunction("mul64i")
		var m5f wa.Function = m.ExportedFunction("mul32f")
		var m6f wa.Function = m.ExportedFunction("mul64f")

		var found bool = ua.All(
			[]wa.Function{
				m5i,
				m6i,
				m5f,
				m6f,
			},
			func(f wa.Function) bool { return nil != f },
		)
		if !found {
			return mf, fmt.Errorf("some functions not found")
		}

		mf.mul32f = func(x, y float32) (float32, error) {
			return wzfb.Fn32f{Function: m5f}.Call(
				context.Background(),
				x,
				y,
			)
		}

		return mf, nil
	},
).AsIf()

var mulFn MulFn = must(wf2.FromCompiled(
	context.Background(),
	compiled,
	ibld,
	fbld,
))

func main() {
	defer func() {
		mustNil(mdl.Close(context.Background()))
		mustNil(compiled.Close(context.Background()))
		mustNil(wtime.Close())
	}()

	var c float32 = must(mulFn.mul32f(
		42.0,
		0.25,
	))
	fmt.Printf("answer: %v\n", c)
}
