#include <sys/socket.h>
#include <linux/netlink.h>
#include <linux/connector.h>
#include <linux/cn_proc.h>
#include <signal.h>
#include <errno.h>
#include <stdbool.h>
#include <unistd.h>
#include <string.h>
#include <stdlib.h>
#include <stdio.h>

struct __attribute__ ((aligned(NLMSG_ALIGNTO))) op_msg {
	struct nlmsghdr nl_hdr;
	struct __attribute__ ((__packed__)) {
		struct cn_msg cn_msg;
		enum proc_cn_mcast_op op;
	};
};

 struct __attribute__ ((aligned(NLMSG_ALIGNTO))) ev_msg {
	struct nlmsghdr nl_hdr;
	struct __attribute__ ((__packed__)) {
		struct cn_msg cn_msg;
		struct proc_event proc_ev;
	};
};

int event_data_size() {
	struct proc_event proc_ev;
	return sizeof(proc_ev.event_data);
}

int kernel_pid_t_size() {
	return sizeof(__kernel_pid_t);
}

int nl_init()
{
	int fd = socket(PF_NETLINK, SOCK_DGRAM, NETLINK_CONNECTOR);
	if (fd == -1) { return -errno; }

	struct sockaddr_nl sa = { .nl_family = AF_NETLINK, .nl_groups = CN_IDX_PROC, .nl_pid = getpid() };
    int err = bind(fd, (struct sockaddr*)&sa, sizeof(struct sockaddr_nl));
	if (err != 0) { close(fd); return -errno; }

    return fd;
}


int nl_subscribe(int fd, int dosubscribe)
{
	struct op_msg m = {
		.nl_hdr = { .nlmsg_len = sizeof(struct op_msg), .nlmsg_type = NLMSG_DONE, .nlmsg_flags = 0, .nlmsg_seq = 0, .nlmsg_pid = getpid() },
		.cn_msg = { .id = {.idx = CN_IDX_PROC, .val = CN_VAL_PROC}, .seq = 0, .ack = 0, .len = sizeof(enum proc_cn_mcast_op) },
		.op = dosubscribe ? PROC_CN_MCAST_LISTEN : PROC_CN_MCAST_IGNORE
	};
    int err = send(fd, &m, sizeof(struct op_msg), 0);
	if (err != 0) { return -errno; }

    return 0;
}

int nl_receive_event(int fd, int *what, int *event_data)
{
	struct ev_msg m;
	ssize_t n = recv(fd, &m, sizeof(struct ev_msg), 0);

	if (n == -1) { return -errno; }
	if (n == 0) { return 0; } // orderly shutdown

	*what = m.proc_ev.what;
	memcpy(event_data, &m.proc_ev.event_data, sizeof(m.proc_ev.event_data));
	return 1;
} 
