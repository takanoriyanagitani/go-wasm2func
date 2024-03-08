package binary

import (
	"context"
	"fmt"

	wa "github.com/tetratelabs/wazero/api"
)

type Fn32f struct {
	wa.Function
}

func (f Fn32f) Call(ctx context.Context, x, y float32) (float32, error) {
	var ux uint64 = wa.EncodeF32(x)
	var uy uint64 = wa.EncodeF32(y)
	ret, e := f.Function.Call(ctx, ux, uy)
	if nil != e {
		return 0.0, e
	}

	if 1 != len(ret) {
		return 0.0, fmt.Errorf("unexpected return value count: %v", len(ret))
	}

	return wa.DecodeF32(ret[0]), nil
}
