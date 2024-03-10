package main

import (
	"context"
	"fmt"
	"io/fs"
	"os"

	wz "github.com/tetratelabs/wazero"
	wa "github.com/tetratelabs/wazero/api"
)

func must[T any](t T, e error) T {
	if nil != e {
		panic(e)
	}
	return t
}

func mustNil(e error) {
	if nil != e {
		panic(e)
	}
}

var dfs fs.FS = os.DirFS("./")
var wasmBytes []byte = must(fs.ReadFile(dfs, "rs_hex2bytes_static.wasm"))

var rtime wz.Runtime = wz.NewRuntime(context.Background())

var compiled wz.CompiledModule = must(rtime.CompileModule(
	context.Background(),
	wasmBytes,
))

var mcfg wz.ModuleConfig = wz.NewModuleConfig()

var instance wa.Module = must(rtime.InstantiateModule(
	context.Background(),
	compiled,
	mcfg,
))

var mem wa.Memory = instance.Memory()

var inputOffset wa.Function = instance.ExportedFunction("i_ptr")
var outputOffset wa.Function = instance.ExportedFunction("o_ptr")

var hexStr2bytes wa.Function = instance.ExportedFunction("hex_string2bytes")

func main() {
	defer func() {
		mustNil(instance.Close(context.Background()))
		mustNil(compiled.Close(context.Background()))
		mustNil(rtime.Close(context.Background()))
	}()

	var ctx context.Context = context.Background()

	var inputOffsetRaw uint64 = must(inputOffset.Call(ctx))[0]
	var ioff int32 = wa.DecodeI32(inputOffsetRaw)
	fmt.Printf("input offset: %v\n", ioff)

	var outputOffsetRaw uint64 = must(outputOffset.Call(ctx))[0]
	var optr int32 = wa.DecodeI32(outputOffsetRaw)
	fmt.Printf("output offset: %v\n", optr)

	var hbcntRaw uint64 = must(hexStr2bytes.Call(ctx, wa.EncodeI32(4)))[0]
	var hbcnt int32 = wa.DecodeI32(hbcntRaw)
	fmt.Printf("hbcnt: %v\n", hbcnt)

	var hexBytes []byte = []byte("254c")

	if !mem.Write(uint32(ioff), hexBytes) {
		panic("unable to write")
	}

	hbcntRaw = must(
		hexStr2bytes.Call(
			ctx,
			wa.EncodeI32(int32(len(hexBytes))),
		),
	)[0]
	hbcnt = wa.DecodeI32(hbcntRaw)
	fmt.Printf("hbcnt: %v\n", hbcnt)
	if hbcnt < 0 {
		panic("invalid hbcnt")
	}

	out, _ := mem.Read(uint32(optr), uint32(hbcnt))
	fmt.Printf("bytes: %v\n", out)
}
