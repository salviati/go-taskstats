package main
import (
	"github.com/salviati/go-taskstats/taskstats"
	"log"
)

func main() {
	mon, err := taskstats.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	
	for {
		ev0, err := mon.Event()
		if err != nil {
			log.Fatal(err)
		}
		switch ev := ev0.(type) {
			case *taskstats.ExecEvent:
				log.Println("exec:", ev.Pid)
			case *taskstats.ExitEvent:
				log.Println("exit:", ev.Pid)
		}
	}
}
