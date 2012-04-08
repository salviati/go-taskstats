// See http://www.kernel.org/doc/Documentation/accounting/taskstats.txt
// package taskstats

package taskstats

import (
// 	"exp/inotify"
// 	"log"
)

func NewWatcher() (*Watcher, error) {
	var err error
	w := new(Watcher)

	w.fd, err = nl_init()
	if err != nil {
		return nil, err
	}

	if err = w.Control(true); err != nil {
		return nil, err
	}

	return w, nil
}

func (w *Watcher) Control(on bool) error {
	if w.on == on {
		return nil
	}

	if err := nl_subscribe(w.fd, on); err != nil {
		return err
	}

	w.on = on
	return nil
}
