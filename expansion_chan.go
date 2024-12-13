package expansion_chan

import (
	"context"
	"sync"
)

type ExpansionChan[T any] struct {
	ch              chan T
	buffer          BufferData[T]
	onceBufferStart *sync.Once
	lock            *sync.RWMutex
	ctx             context.Context

	logger Logger
}

type TypeBuffer int

const (
	TypeBufferStack TypeBuffer = iota
)

type BufferData[T any] interface {
	Push(v T)
	Pop() (res T, exist bool)
	GetAll() (res []T)
	Len() int
}

type Option[T any] struct {
	Size   int
	Driver TypeBuffer

	SetKit []SetOpt[T]
}

type SetOpt[T any] func(*ExpansionChan[T])

func SetLogger[T any](lg Logger) SetOpt[T] {
	return func(c *ExpansionChan[T]) {
		c.logger = lg
	}
}

func NewExpansionChan[T any](ctx context.Context, opt Option[T]) *ExpansionChan[T] {

	var driver BufferData[T]
	switch opt.Driver {
	case TypeBufferStack:
		driver = NewStack[T]()
	}

	ec := &ExpansionChan[T]{
		ch:              make(chan T, opt.Size),
		onceBufferStart: &sync.Once{},
		buffer:          driver,
		lock:            &sync.RWMutex{},
		ctx:             ctx,
	}

	for _, v := range opt.SetKit {
		v(ec)
	}

	if ec.logger == nil {
		ec.logger = &Lg{}
	}

	return ec
}

func (c *ExpansionChan[T]) Push(v T) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	select {
	case c.ch <- v:
	default:
		c.onceBufferStart.Do(func() {
			go c.startBuffer()
		})
		c.buffer.Push(v)
	}
}

func (c *ExpansionChan[T]) startBuffer() {
	c.logger.Info("buffer start")
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
		}
		v, exist := c.buffer.Pop()
		if !exist {
			break
		}
		c.ch <- v
	}
	c.lock.Lock()
	defer c.lock.Unlock()

	list := c.buffer.GetAll()
	for _, v := range list {
		c.ch <- v
	}
	c.onceBufferStart = &sync.Once{}
	c.logger.Info("buffer end")
}

func (c *ExpansionChan[T]) Pop() (res T, exist bool) {
	select {
	case res = <-c.ch:
		return res, true
	default:
		return res, false
	}
}
