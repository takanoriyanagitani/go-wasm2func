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
	wfim "github.com/takanoriyanagitani/go-wasm2func/instance/memory"

	wzc "github.com/takanoriyanagitani/go-wasm2func/runtime/rt-wazero/compiler"
	wzf "github.com/takanoriyanagitani/go-wasm2func/runtime/rt-wazero/function"
	wzi "github.com/takanoriyanagitani/go-wasm2func/runtime/rt-wazero/instance"

	// revive:disable:line-length-limit
	wzfu "github.com/takanoriyanagitani/go-wasm2func/runtime/rt-wazero/function/unary"
	wzf0 "github.com/takanoriyanagitani/go-wasm2func/runtime/rt-wazero/function/unit"
	wzim "github.com/takanoriyanagitani/go-wasm2func/runtime/rt-wazero/instance/memory"
	// revive:enable:line-length-limit
)

var rtime wz.Runtime = wz.NewRuntime(context.Background())
var wtime wzc.Runtime = wzc.Runtime{
	Runtime: rtime,
}
var wcomp wfc.Compiler[wz.CompiledModule] = wtime.AsIf()

var fsys fs.FS = os.DirFS("./")

const pwasDefault string = "rs_hex2bytes.wasm"

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

var mem wa.Memory = mdl.Memory()
var rwm wfim.ReadWriteMemory = wzim.RwMem{
	Memory: mem,
}.AsIf()

type ConvFn struct {
	inputResize  func(int32) (int32, error)
	outputResize func(int32) (int32, error)

	inputPtr  func() (int32, error)
	outputPtr func() (int32, error)

	emsgPtr func() (int32, error)
	emsgSize func() (int32, error)
    errorMessageSize func()(int32, error)

	hex2bytes func() (int32, error)
}

var fbld wf2.Builder[wa.Module, ConvFn] = wzf.Builder[ConvFn](
	func(_ context.Context, m wa.Module) (ConvFn, error) {
		var mf ConvFn

		var isz wa.Function = m.ExportedFunction("input_resize")
		var osz wa.Function = m.ExportedFunction("output_resize")

		var iptr wa.Function = m.ExportedFunction("input_ptr")
		var optr wa.Function = m.ExportedFunction("output_ptr")

		var eptr wa.Function = m.ExportedFunction("emsg_ptr")
		var esize wa.Function = m.ExportedFunction("emsg_sz")
		var eSize wa.Function = m.ExportedFunction("emsg_size")

		var h2b wa.Function = m.ExportedFunction("hex2bytes")

		var found bool = ua.All(
			[]wa.Function{
				isz, osz,
				iptr, optr, eptr,
				h2b,
			},
			func(f wa.Function) bool { return nil != f },
		)
		if !found {
			return mf, fmt.Errorf("some functions not found")
		}

		mf.inputResize = func(sz int32) (capacity int32, e error) {
			return wzfu.Fn32i{Function: isz}.Call(context.Background(), sz)
		}

		mf.outputResize = func(sz int32) (capacity int32, e error) {
			return wzfu.Fn32i{Function: osz}.Call(context.Background(), sz)
		}

		mf.inputPtr = func() (offset int32, e error) {
			return wzf0.Fn32i{Function: iptr}.Call(context.Background())
		}
		mf.outputPtr = func() (offset int32, e error) {
			return wzf0.Fn32i{Function: optr}.Call(context.Background())
		}

		mf.emsgPtr = func() (offset int32, e error) {
			return wzf0.Fn32i{Function: eptr}.Call(context.Background())
		}
		mf.emsgSize = func() (offset int32, e error) {
			return wzf0.Fn32i{Function: esize}.Call(context.Background())
		}
		mf.errorMessageSize = func() (offset int32, e error) {
			return wzf0.Fn32i{Function: eSize}.Call(context.Background())
		}

		mf.hex2bytes = func() (offset int32, e error) {
			return wzf0.Fn32i{Function: h2b}.Call(context.Background())
		}

		return mf, nil
	},
).AsIf()

var convFn ConvFn = must(wf2.FromCompiled(
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

    _, _ = mem.Grow(1)

	var icap int32 = must(convFn.inputResize(32))
	fmt.Printf("input cap: %v\n", icap)
	var ocap int32 = must(convFn.outputResize(16))
	fmt.Printf("output cap: %v\n", ocap)

	var iptr int32 = must(convFn.inputPtr())
	fmt.Printf("input ptr: %v\n", iptr)
	fmt.Printf("mem sz: %v\n", mem.Size())

    mustNil(rwm.Write(
        uint32(iptr),
        []byte("cafef00ddeadbeafface864299792458"),
    ))

    var eptr int32 = must(convFn.emsgPtr())
	fmt.Printf("emsg ptr: %v\n", eptr)
    var esz int32 = must(convFn.emsgSize())
	fmt.Printf("emsg size: %v\n", esz)
    var errorMessageSize int32 = must(convFn.errorMessageSize())
	fmt.Printf("err msg size: %v\n", errorMessageSize)

	var hsz int32 = must(convFn.hex2bytes())
	fmt.Printf("hex size: %v\n", hsz)
}
