package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"
)

type Task struct {
	pid   int
	files []string
}

func (t *Task) RefreshFiles() error {
	oldfiles := t.files
	t.files = make([]string, 0)

	log.Println("RefreshFiles for pid", t.pid)

	dirpath := fmt.Sprint("/proc/", t.pid, "/fd/")
	file, err := os.Open(dirpath)
	if err != nil {
		return err
	}
	defer file.Close()

	names, err := file.Readdirnames(0)
	if err != nil {
		return err
	}

	for _, name := range names {
		target, err := os.Readlink(dirpath + name)
		if err != nil {
			continue
		} // FIXME(utkan)
		ext := strings.ToLower(filepath.Ext(target))
		if !in(ext, exts) {
			continue
		}

		t.files = append(t.files, target)
	}

	if !seq(t.files, oldfiles) {
		t.FilesChanged()
	}

	return nil
}

func (t *Task) FilesChanged() {
	log.Println("new files:", t.files)
}

func (t *Task) Tracer() error {
	runtime.LockOSThread()

	if err := syscall.PtraceAttach(t.pid); err != nil {
		return err
	}
	defer func() {
		// 		syscall.PtraceCont(t.pid, 0)
		syscall.PtraceDetach(t.pid)
	}()

	regsout := new(syscall.PtraceRegs)
	status := new(syscall.WaitStatus)
	var timer *time.Timer
	refresh := func() { t.RefreshFiles(); timer.Stop(); timer = nil }
	for {
		if _, err := syscall.Wait4(t.pid, status, 0, nil); err != nil {
			log.Println("wait failed", err)
			return err
		}
		if status.Exited() {
			log.Println("exited")
			return nil
		}

		if err := syscall.PtraceGetRegs(t.pid, regsout); err != nil {
			log.Println("getregs failed", err)
			return err
		}

		if regsout.Orig_rax == syscall.SYS_OPEN {
			if timer != nil {
				if timer.Stop() == false {
					log.Println("cannot stop the timer")
				}
			}
			timer = time.AfterFunc(1e9, refresh) // Wait until open()s "settle".
		}

		if err := PtraceSyscall(t.pid); err != nil {
			log.Println("PtraceSyscall failed", err)
			return err
		}
	}
	panic("can't reach")
}

func NewTask(pid int) *Task {
	t := new(Task)
	t.pid = pid
	go func() {
		time.Sleep(1e8)
		t.Tracer()
	}()
	t.RefreshFiles()
	return t
}
