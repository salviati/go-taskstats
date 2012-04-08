package main

// #include <sys/ptrace.h>
// #include <stdio.h>
// #include <errno.h>
// long PtraceSyscall(int pid) {
// 	return ptrace(PTRACE_SYSCALL, pid, NULL, NULL); 
// }
// int get_errno() { return errno; }
//  
import "C"

import "syscall"

func PtraceSyscall(pid int) error {
	ret := int(C.PtraceSyscall(C.int(pid)))
	if ret == -1 {
		return syscall.Errno(C.get_errno()) // BUG(utkan): shouldn't be doing this!
	}
	return nil
}
