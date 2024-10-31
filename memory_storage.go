package main

type MemoryStorage[T any] struct {
	data T
}

// [T any] is a type parameter declaration where T is a type parameter, like a variable but for type
func NewMemoryStorage[T any](initialData T) *MemoryStorage[T] {
	//   ^function name    ^type param  ^arg     ^return type
	//                    & constraint   of type T  pointer to MemoryStorage
	//                                            that works with type T
	return &MemoryStorage[T]{
		data: initialData,
	}
}

func (m *MemoryStorage[T]) Save(data T) error {
	m.data = data
	return nil
}

func (m *MemoryStorage[T]) Load(target *T) error {
	*target = m.data
	return nil
}
