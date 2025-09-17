package or

import "reflect"

// Or возвращает канал, который закрывается при закрытии любого из входных.
// nil-каналы игнорируются.
// 0 каналов -> возвращается уже закрытый канал.
// 1 канал   -> возвращается он же.
func Or(channels ...<-chan interface{}) <-chan interface{} {
	nonNil := make([]<-chan interface{}, 0, len(channels))
	for _, ch := range channels {
		if ch != nil {
			nonNil = append(nonNil, ch)
		}
	}

	switch len(nonNil) {
	case 0:
		c := make(chan interface{})
		close(c)
		return c
	case 1:
		return nonNil[0]
	}

	out := make(chan interface{})
	go func() {
		defer close(out)
		cases := make([]reflect.SelectCase, len(nonNil))
		for i, ch := range nonNil {
			cases[i] = reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(ch),
			}
		}
		_, _, _ = reflect.Select(cases)
	}()
	return out
}
