package ui

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/mitchellh/go-ps"

	"github.com/Conquest-Reforged/ReforgedLauncher/utils/errs"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/files"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/platform"
)

var (
	dev = flag.Bool("dev", false, "")
)

type EventType int

type Event struct {
	Type    EventType
	Message string
}

type Window struct {
	l      *sync.RWMutex
	cmd    *exec.Cmd
	events chan *Event
}

type Settings struct {
	Title      string
	Url        string
	Icon       string
	Width      int
	Height     int
	Resizable  bool
	Borderless bool
}

const (
	INVOKE EventType = 0
	INJECT EventType = 1
	EXIT   EventType = 2
)

func (w *Window) Js(script string, args ...interface{}) {
	w.l.Lock()
	defer w.l.Unlock()
	if w.events == nil {
		return
	}

	script = fmt.Sprintf(script, args...)
	w.events <- &Event{
		Type:    INVOKE,
		Message: script,
	}
}

func (w *Window) CSS(css string) {
	w.l.Lock()
	defer w.l.Unlock()
	if w.events == nil {
		return
	}

	w.events <- &Event{
		Type:    INJECT,
		Message: css,
	}
}

func (w *Window) Exit() {
	w.l.Lock()
	defer w.l.Unlock()

	if w.events == nil {
		return
	}

	w.events <- &Event{
		Type:    EXIT,
		Message: "",
	}
}

func startWindow(settings *Settings) (*Window, error) {
	cmd, e := buildCommand(settings)
	if e != nil {
		return nil, e
	}

	in, e := cmd.StdinPipe()
	if e != nil {
		return nil, e
	}

	e = cmd.Start()
	if e != nil {
		return nil, e
	}

	w := &Window{
		l:      &sync.RWMutex{},
		cmd:    cmd,
		events: make(chan *Event),
	}

	go handleEvents(w.events, in)
	go handleClose(w)

	return w, nil
}

func handleEvents(events chan *Event, wr io.WriteCloser) {
	defer files.Close(wr)
	for event := range events {
		e := json.NewEncoder(wr).Encode(event)
		errs.Log("Encode event", e)
		_, e = fmt.Fprintln(wr)
		errs.Log("Send event", e)
	}
	log.Println("Window events thread shut down")
}

func handleClose(w *Window) {
	defer dispose(w)
	e := w.cmd.Wait()
	errs.Log("Wait", e)
	log.Println("Window listener thread shut down")
}

func terminate(w *Window) {
	w.l.Lock()
	if w.events != nil {
		// release lock as required in w.Exit
		w.l.Unlock()
		log.Println("Closing window")
		w.Exit()

		// obtain lock again to allow w.Exit to complete
		w.l.Lock()

		// close events channel, allows handleEvents to die
		close(w.events)
	}
	defer w.l.Unlock()

	// kill process
	if w.cmd != nil {
		log.Println("Stopping window process")
		_ = w.cmd.Process.Kill()
	}
}

func dispose(w *Window) {
	w.l.Lock()
	defer w.l.Unlock()
	if w.events != nil {
		close(w.events)
	}
	w.cmd = nil
	w.events = nil
}

func buildCommand(settings *Settings) (*exec.Cmd, error) {
	data, e := json.Marshal(settings)
	if e != nil {
		return nil, e
	}

	if *dev {
		return exec.Command("go", "run", "main.go", "-w", string(data)), nil
	}

	p, e := ps.FindProcess(os.Getpid())
	if e != nil {
		return nil, e
	}

	wd, e := os.Getwd()
	if e != nil {
		return nil, e
	}

	exe := filepath.Join(wd, p.Executable())
	return platform.RunExecutable(exe, "-w", string(data)), nil
}
