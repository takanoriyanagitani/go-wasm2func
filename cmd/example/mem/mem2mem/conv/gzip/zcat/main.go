package main

import (
	"context"
	"encoding/base64"
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
var wasmBytes []byte = must(fs.ReadFile(dfs, "rs_zcat.wasm"))

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

var inputOffset wa.Function = instance.ExportedFunction("offset_i")
var outputOffset wa.Function = instance.ExportedFunction("offset_o")

var inputResize wa.Function = instance.ExportedFunction("input_resize")
var outputReset wa.Function = instance.ExportedFunction("output_reset")

var gzipDecode wa.Function = instance.ExportedFunction("gzip_decode")

const testEncodedString64 string = "H4sIAAAAAAAEA8tIzckHAFlRj4UEAAAA"

var testEncodedBytes []byte = must(
	base64.StdEncoding.DecodeString(testEncodedString64),
)

func main() {
	defer func() {
		mustNil(instance.Close(context.Background()))
		mustNil(compiled.Close(context.Background()))
		mustNil(rtime.Close(context.Background()))
	}()

	var ctx context.Context = context.Background()

	var inputCapRaw uint64 = must(inputResize.Call(ctx, wa.EncodeI32(
		int32(len(testEncodedBytes)),
	)))[0]
	var inputCap int32 = wa.DecodeI32(inputCapRaw)
	fmt.Printf("input cap: %v\n", inputCap)

	var outputCapRaw uint64 = must(outputReset.Call(ctx, wa.EncodeI32(
		int32(len(testEncodedBytes)),
	)))[0]
	var outputCap int32 = wa.DecodeI32(outputCapRaw)
	fmt.Printf("output cap: %v\n", outputCap)

	var inputOffsetRaw uint64 = must(inputOffset.Call(ctx))[0]
	var ioff int32 = wa.DecodeI32(inputOffsetRaw)
	fmt.Printf("input offset: %v\n", ioff)

	var outputOffsetRaw uint64 = must(outputOffset.Call(ctx))[0]
	var optr int32 = wa.DecodeI32(outputOffsetRaw)
	fmt.Printf("output offset: %v\n", optr)

	var decCntRaw uint64 = must(gzipDecode.Call(ctx))[0]
	var decCnt int32 = wa.DecodeI32(decCntRaw)
	fmt.Printf("decCnt: %v\n", decCnt)

	if !mem.Write(uint32(ioff), testEncodedBytes) {
		panic("unable to write")
	}

	decCntRaw = must(gzipDecode.Call(ctx))[0]
	decCnt = wa.DecodeI32(decCntRaw)
	fmt.Printf("decCnt: %v\n", decCnt)

	decBytes, _ := mem.Read(uint32(optr), uint32(decCnt))
	fmt.Printf("decoded bytes: %v\n", decBytes)
	fmt.Printf("decoded string: %s\n", string(decBytes))

}
