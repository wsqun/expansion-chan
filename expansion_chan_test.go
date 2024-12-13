package expansion_chan

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestExpansionChan_Push(t *testing.T) {
	s := NewExpansionChan[string](context.Background(), Option[string]{
		Size:   3,
		Driver: 0,
		SetKit: nil,
	})

	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		s.Push("test:" + strconv.Itoa(i))
	}

	for {
		if v, ok := s.Pop(); ok {
			fmt.Println(v)
			time.Sleep(100 * time.Millisecond)
		} else {
			break
		}
	}

	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		s.Push("test:" + strconv.Itoa(i))
	}

	for {
		if v, ok := s.Pop(); ok {
			fmt.Println(v)
			time.Sleep(100 * time.Millisecond)
		} else {
			break
		}
	}

}
