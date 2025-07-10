package fan_in

import (
	"context"
	"reflect"
	"sync"
)

// поменять use на одну из 2х реализаций [MANY_GO_FABRIC_KEY, REFLECT_FABRIC_KEY]
const (
	USE                = MANY_GO_FABRIC_KEY
	MANY_GO_FABRIC_KEY = "many"
	REFLECT_FABRIC_KEY = "reflect"
)

type MergeFabric[T any] interface {
	Merge(context.Context, ...<-chan T) <-chan T
}

// MergeChannels - принимает несколько каналов на вход и объединяет их в один
// Fan-in и merge channels синонимы
func MergeChannels(channels ...<-chan int) <-chan int {
	ctx := context.TODO()
	f := fabricMerge[int](USE)

	return f.Merge(ctx, channels...)
}

// Сами решения ниже
// Фабрика просто для удобства

// Sloution 1
// на каждый входной канал создаётся горутина, которая его читает и пишет в единый out
// особеность решения - простота и чиатемость
// минус - неоправданое использование множества горутин, но всё во имя простоты и читаемости
func mergeManyGo[T any](ctx context.Context, ins ...<-chan T) <-chan T {
	out := make(chan T)

	wg := sync.WaitGroup{}
	wg.Add(len(ins))

	go func() {
		wg.Wait()
		close(out)
	}()

	for _, in := range ins {
		go func() {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case v, ok := <-in:
					if !ok {
						return
					}

					out <- v
				}
			}
		}()
	}

	return out
}

// Sloution 2
// отличие - запускается только одна горутина, которая читает разом все входные каналы и пишет в out
func mergeReflect[T any](ctx context.Context, ins ...<-chan T) <-chan T {
	out := make(chan T)

	// Что-то такое можно написать используя рефлексию
	//	closed := 0
	//	for closed < len(ins) {
	//		select {
	// 		case <-ctx.Done():
	// 			return
	//		case v, ok := <-ins[0]:
	//			if !ok {
	//				closed++
	//				ins[0] = nil
	//				continue
	//			}
	//
	//			out <- v
	//		case v, ok := <-ins[0]:
	//			if !ok {
	//				closed++
	//				ins[0] = nil
	//				continue
	//			}
	//
	//			out <- v
	//		}
	//		// .....
	//		// .....
	//	}

	go func() {
		defer close(out)

		cases := make([]reflect.SelectCase, 0, len(ins)+1)
		cases = append(cases, reflect.SelectCase{ // context append
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ctx.Done()),
			Send: reflect.ValueOf(nil),
		})
		for _, in := range ins {
			cases = append(cases, reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(in),
				Send: reflect.ValueOf(nil),
			})
		}

		closed := 0
		for closed < len(ins) {
			i, vRefAny, ok := reflect.Select(cases)
			if !ok {
				if i == 0 {
					return // ctx.Done()
				}

				closed++
				cases[i].Chan = reflect.ValueOf(nil)
				continue // return // если по условию с kaiten закрывать out когда один из входных закрылся
			}

			v, ok := vRefAny.Interface().(T)
			if !ok {
				continue
			}

			out <- v
		}
	}()

	return out
}
