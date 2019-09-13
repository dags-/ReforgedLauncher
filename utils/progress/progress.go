package progress

import (
	"log"
	"time"
)

type Listener interface {
	Stat(string, float64)
	GlobalStatus(string)
	TaskStatus(string)
	TaskProgress(float64)
	GlobalProgress(float64)
	Wait()
}

type Progress struct {
	Status  func(string)
	Message func(string)
	Task    func(float64)
	Global  func(float64)
}

func Log(listener Listener) Listener {
	switch listener.(type) {
	case *logWrapper:
		return listener
	default:
		return &logWrapper{Listener: listener}
	}
}

type logWrapper struct {
	Listener
}

func (p *Progress) Stat(s string, f float64) {
	p.GlobalStatus(s)
	p.GlobalProgress(f)
}

func (p *Progress) GlobalStatus(s string) {
	fn := p.Status
	if fn != nil {
		fn(s)
	}
}

func (p *Progress) TaskStatus(s string) {
	fn := p.Message
	if fn != nil {
		fn(s)
	}
}

func (p *Progress) GlobalProgress(f float64) {
	fn := p.Global
	if fn != nil {
		fn(f)
	}
}

func (p *Progress) TaskProgress(f float64) {
	fn := p.Task
	if fn != nil {
		fn(f)
	}
}

func (p *Progress) Wait() {
	time.Sleep(time.Millisecond * 500)
}

func (l *logWrapper) Stat(s string, f float64) {
	l.Listener.Stat(s, f)
	log.Println(s)
}

func (l *logWrapper) GlobalStatus(s string) {
	l.Listener.GlobalStatus(s)
	log.Println(s)
}

func (l *logWrapper) TaskStatus(s string) {
	l.Listener.TaskStatus(s)
	log.Println(s)
}

func (l *logWrapper) Wait() {
	l.Listener.Wait()
}

func Listen(ch <-chan float64, h Listener) {
	for v := range ch {
		h.TaskProgress(v)
	}
}
