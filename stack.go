package expansion_chan

import "sync"

var _ BufferData[int] = &Stack[int]{}

type Stack[T any] struct {
	data []T
	top  int
	lock *sync.Mutex
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{lock: &sync.Mutex{}}
}

func (s *Stack[T]) Push(v T) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.data = append(s.data, v)
	s.top++
}

func (s *Stack[T]) Pop() (res T, exist bool) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.top == 0 {
		return res, false
	}
	s.top--
	v := s.data[s.top]
	s.data = s.data[:s.top]
	return v, true
}

func (s *Stack[T]) GetAll() (res []T) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.top == 0 {
		return res
	}
	// 缩容
	res = make([]T, s.top)
	copy(res, s.data)
	// 清空
	s.data = make([]T, 0)
	return
}

func (s *Stack[T]) Len() int {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.top
}
