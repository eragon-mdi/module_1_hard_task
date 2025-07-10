package fan_in

import "context"

type mergeChans_ManyGo[T any] struct{}
type mergeChans_Reflect[T any] struct{}

func fabricMerge[T any](key string) MergeFabric[T] {
	var fabricMap = map[string]MergeFabric[T]{
		MANY_GO_FABRIC_KEY: mergeChans_ManyGo[T]{},
		REFLECT_FABRIC_KEY: mergeChans_Reflect[T]{},
	}

	res, ok := fabricMap[key]
	if !ok {
		panic("undefined mergeFabric type")
	}

	return res
}

func (mergeChans_ManyGo[T]) Merge(ctx context.Context, ins ...<-chan T) <-chan T {
	return mergeManyGo(ctx, ins...)
}

func (mergeChans_Reflect[T]) Merge(ctx context.Context, ins ...<-chan T) <-chan T {
	return mergeReflect(ctx, ins...)
}
