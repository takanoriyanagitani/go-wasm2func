package unary

import (
	"context"
	"fmt"

	wa "github.com/tetratelabs/wazero/api"
)

type Fn32i struct {
	wa.Function
}

func (f Fn32i) Call(ctx context.Context) (int32, error) {
	ret, e := f.Function.Call(ctx)
	if nil != e {
		return 0, e
	}

	if 1 != len(ret) {
		return 0, fmt.Errorf("unexpected return value count: %v", len(ret))
	}

	return wa.DecodeI32(ret[0]), nil
}
