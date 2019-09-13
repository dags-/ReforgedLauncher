package progress

import (
	"log"
	"time"
)

type logger struct {
}

func Debug() Listener {
	return &logger{}
}

func (l *logger) Stat(s string, f float64) {
	l.GlobalStatus(s)
	l.GlobalProgress(f)
}

func (l *logger) GlobalStatus(s string) {
	log.Println(s)
}

func (l *logger) GlobalProgress(f float64) {
	log.Printf("\r%.2f", f*100)
}

func (l *logger) TaskStatus(s string) {
	log.Println(s)
}

func (l *logger) TaskProgress(f float64) {
	log.Printf("\r%.2f", f*100)
}

func (l *logger) Wait() {
	time.Sleep(time.Millisecond * 500)
}
