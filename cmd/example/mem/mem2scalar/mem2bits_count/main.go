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
var wasmBytes []byte = must(fs.ReadFile(dfs, "rs_mem2bcnt.wasm"))

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

var inputOffset wa.Function = instance.ExportedFunction("input_offset")

var countOnes wa.Function = instance.ExportedFunction("count_ones")

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

	var ones uint64 = must(countOnes.Call(ctx, wa.EncodeI32(4)))[0]
	fmt.Printf("ones: %v\n", ones)

	if ! mem.Write(uint32(ioff), []byte{
		0x37, // ones: 2+3 = 5
		0x76, // ones: 3+2 = 5
	}) {
		panic("unable to write")
	}

	ones = must(countOnes.Call(ctx, wa.EncodeI32(4)))[0]
	fmt.Printf("ones: %v\n", ones)

}
