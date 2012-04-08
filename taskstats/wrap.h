#ifndef WRAP_H
#define WRAP_H

int event_data_size();
int kernel_pid_t_size();
int nl_init();
int nl_subscribe(int fd, int dosubscribe);
int nl_receive_event(int fd, int *what, int *event_data);

#endif
