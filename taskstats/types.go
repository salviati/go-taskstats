/* /usr/include/linux/cn_proc.h:
 * 
 * From the user's point of view, the process
 * ID is the thread group ID and thread ID is the internal
 * kernel "pid". So, fields are assigned as follow:
 *
 *  In user space     -  In  kernel space
 *
 * parent process ID  =  parent->tgid
 * parent thread  ID  =  parent->pid
 * child  process ID  =  child->tgid
 * child  thread  ID  =  child->pid
 */
package taskstats

type Watcher struct {
	fd int
	on bool
}

type RawEvent struct {
	What int
	Data []byte
}

type ForkEvent struct {
	ParentGid  int
	ParentTGid int
	ChildPid   int
	ChildTPid  int
}

type ExecEvent struct {
	Pid  int
	TGid int
}

type IdEvent struct {
	Pid  int
	TGid int
	RUid uint32 // Also EUid
	RGid uint32 // Also EGid
}

type SidEvent struct {
	Pid  int
	TGid int
}

type PtraceEvent struct {
	Pid        int
	TGid       int
	TracerPid  int
	TracerTGid int
}

type CommEvent struct {
	Pid  int
	TGid int
	Comm [16]byte
}

type ExitEvent struct {
	Pid        int
	TGid       int
	ExitCode   uint32
	ExitSignal uint32
}
