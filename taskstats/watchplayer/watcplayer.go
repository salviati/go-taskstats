// This program watches the file a (video)player is playing.
// Should better be a made into a setuid program and acquire
// root only when necessary.
package main

import (
	"github.com/salviati/go-taskstats/taskstats"
	"log"
)

var (
	exes    = []string{"mplayer", "xine", "dragon"}
	exts    = []string{".avi", ".mpg", ".mpeg", ".divx", ".xvid", "wmv", "ogv", "mp4"}
	manager Manager
)

func init() {
	manager.tasks = make(map[int]*Task)
}

func main() {
	tw, err := taskstats.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	exec := make(chan int)
	exit := make(chan int)

	go func() {
		for {
			ev0, err := tw.Event()
			if err != nil {
				log.Fatal(err)
			}

			switch ev := ev0.(type) {
			case *taskstats.ExecEvent:
				exec <- ev.Pid
			case *taskstats.ExitEvent:
				exit <- ev.Pid
			}
		}
	}()

	for {
		select {
		case pid := <-exec:
			manager.AddTask(pid)
		case pid := <-exit:
			manager.RmTask(pid)
		}
	}
}
