package arr

func Fold[T, S any](s []T, init S, reducer func(state S, next T) S) S {
	var state = init
	for _, t := range s {
		state = reducer(state, t)
	}
	return state
}
