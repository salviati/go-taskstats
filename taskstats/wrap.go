package taskstats

// #include "wrap.h"
import "C"
import (
	"errors"
	"syscall"
	"unsafe"
)

const (
	PROC_EVENT_NONE   = 0x00000000
	PROC_EVENT_FORK   = 0x00000001
	PROC_EVENT_EXEC   = 0x00000002
	PROC_EVENT_UID    = 0x00000004
	PROC_EVENT_GID    = 0x00000040
	PROC_EVENT_SID    = 0x00000080
	PROC_EVENT_PTRACE = 0x00000100
	PROC_EVENT_COMM   = 0x00000200
	PROC_EVENT_EXIT   = 0x80000000
)

var ErrOrderly = errors.New("GetExecPid: orderly shutdown")

var event_data_size int
var kernel_pid_t_size int

func init() {
	event_data_size = int(C.event_data_size())
	kernel_pid_t_size = int(C.kernel_pid_t_size())
}

func nl_init() (int, error) {
	fd := int(C.nl_init())
	if fd < 0 {
		return 0, syscall.Errno(-fd)
	}
	return fd, nil
}

func nl_subscribe(fd int, on bool) error {
	con := C.int(0)
	if on {
		con = 1
	}

	err := int(C.nl_subscribe(C.int(fd), con))
	if err < 0 {
		return syscall.Errno(fd)
	}
	return nil
}

func nl_receive_event(fd int) (ev *RawEvent, err error) {
	ev = new(RawEvent)
	ev.Data = make([]byte, event_data_size)
	cwhatp := (*C.int)(unsafe.Pointer(&ev.What))
	cevent_datap := (*C.int)(unsafe.Pointer(&ev.Data[0]))
	cfd := C.int(fd)

	e := int(C.nl_receive_event(cfd, cwhatp, cevent_datap))
	if e == 0 {
		err = ErrOrderly
		return nil, err
	}
	if e < 0 {
		err = syscall.Errno(-fd)
		return nil, err
	}
	return
}
