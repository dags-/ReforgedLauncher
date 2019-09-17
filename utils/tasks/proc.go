package tasks

import (
	"log"
	"os"
	"os/exec"
	"sync"
	"time"
)

var (
	lock      = &sync.RWMutex{}
	tick      = time.NewTicker(time.Second * 15)
	processes []*os.Process
)

func init() {
	go func() {
		for range tick.C {
			clean()
		}
	}()
}

func Shutdown() {
	log.Println("Stopping active processes")
	tick.Stop()

	lock.Lock()
	defer lock.Unlock()
	for _, p := range processes {
		if p != nil {
			_ = p.Kill()
		}
	}

	log.Println("Stopped")
	os.Exit(0)
}

func Running(p *os.Process) bool {
	return p != nil && RunningId(p.Pid)
}

func RunningId(pid int) bool {
	if _, e := os.FindProcess(pid); e == nil {
		return true
	}
	return false
}

func Start(cmd *exec.Cmd) error {
	e := cmd.Start()
	if e != nil {
		return e
	}
	addProcess(cmd.Process)
	return nil
}

func Run(cmd *exec.Cmd) error {
	e := Start(cmd)
	if e != nil {
		return e
	}
	addProcess(cmd.Process)
	return cmd.Wait()
}

func addProcess(p *os.Process) {
	lock.Lock()
	defer lock.Unlock()
	processes = append(processes, p)
}

func clean() {
	lock.Lock()
	defer lock.Unlock()
	var running []*os.Process
	for _, p := range processes {
		if Running(p) {
			running = append(running, p)
		}
	}
	processes = running
}
