package taskstats

// #include <linux/types.h>
import "C"
import (
	"unsafe"
	"reflect"
)

func makePidArray(data []byte, n int) []C.__kernel_pid_t {
	pidArray := (*[1 << 30]C.__kernel_pid_t)(unsafe.Pointer(&data[0]))[:n]
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&pidArray))
	hdr.Cap = n
	hdr.Len = n
	
	return pidArray
}

func makeUint32Array(data []byte, n int) []uint32 {
	pidArray := (*[1 << 30]uint32)(unsafe.Pointer(&data[0]))[:n]
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&pidArray))
	hdr.Cap = n
	hdr.Len = n
	
	return pidArray
}

func (w *Watcher) Event() (interface{}, error) {
	ev, err := nl_receive_event(w.fd)
	if err != nil {
		return nil, err
	}
	
	switch uint(ev.What) {
		case PROC_EVENT_NONE:
			return nil, nil
		case PROC_EVENT_FORK:
			data := makePidArray(ev.Data, 4)
			return &ForkEvent{ParentGid: int(data[0]), ParentTGid: int(data[1]), ChildPid: int(data[2]), ChildTPid: int(data[3])}, nil
		case PROC_EVENT_EXEC:
			data := makePidArray(ev.Data, 2)
			return &ExecEvent{Pid: int(data[0]), TGid: int(data[1])}, nil
		case PROC_EVENT_UID, PROC_EVENT_GID:
			data1 := makePidArray(ev.Data, 2)
			data2 := makeUint32Array(ev.Data[kernel_pid_t_size*2:], 2)
			return &IdEvent{Pid: int(data1[0]), TGid: int(data1[1]), RUid: uint32(data2[0]), RGid: uint32(data2[1])}, nil
		case PROC_EVENT_SID:
			data1 := makePidArray(ev.Data, 2)
			return &SidEvent{Pid: int(data1[0]), TGid: int(data1[1])}, nil
		case PROC_EVENT_PTRACE:
			data := makePidArray(ev.Data, 4)
			return &PtraceEvent{Pid: int(data[0]), TGid: int(data[1]), TracerPid: int(data[2]), TracerTGid: int(data[3])}, nil
		case PROC_EVENT_COMM:
			data1 := makePidArray(ev.Data, 2)
			data2 := ev.Data[kernel_pid_t_size*2:kernel_pid_t_size*2+16]
			commEvent := &CommEvent{Pid: int(data1[0]), TGid: int(data1[1])}
			copy(commEvent.Comm[0:16], data2)
			return commEvent, nil
		case PROC_EVENT_EXIT:
			data1 := makePidArray(ev.Data, 2)
			data2 := makeUint32Array(ev.Data[kernel_pid_t_size*2:], 2)
			return &ExitEvent{Pid: int(data1[0]), TGid: int(data1[1]), ExitCode: uint32(data2[0]), ExitSignal: uint32(data2[1])}, nil
	}


	return nil, nil
	panic("shouldn't happen")
}