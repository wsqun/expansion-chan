package expansion_chan

import "log"

type Logger interface {
	Info(string)
	Error(string)
}

type Lg struct {
}

func (l *Lg) Info(s string) {
	log.Println(s)
}

func (l *Lg) Error(s string) {
	log.Println(s)
}
