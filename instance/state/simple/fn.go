package simple

import (
	"context"
)

type Fn struct {
	Input
	Caller
	Output
}

func (f Fn) Call(ctx context.Context, input []byte) (output []byte, e error) {
	e = f.Input.Set(ctx, input)
	if nil != e {
		return nil, e
	}
	e = f.Caller.Call(ctx)
	if nil != e {
		return nil, e
	}
	return f.Output.Get(ctx)
}
