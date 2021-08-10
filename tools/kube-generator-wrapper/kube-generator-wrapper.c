#include <errno.h>
#include <signal.h>
#include <stdio.h>
#include <stdlib.h>
#include <sys/types.h>
#include <sys/wait.h>
#include <unistd.h>

pid_t pid;

void handle_sigint(int sig) {
  int tpid, status;
  kill(pid, SIGKILL);
  while (1) {
    tpid = waitpid(-1, &status, WNOWAIT);
    if (tpid == -1) {
      break;
    }
    kill(tpid, SIGKILL);
  }
  if (sig == SIGKILL || sig == SIGTERM) {
    exit(0);
  }
}

int main(int argc, char *argv[]) {
  int status;
  if (argc == 1) {
    printf("number of the arguments must be over and equal 2\n");
    return -1;
  }
  signal(SIGINT, handle_sigint);
  signal(SIGKILL, handle_sigint);
  signal(SIGTERM, handle_sigint);
  while (1) {
    pid = fork();
    if (pid == 0) { // child
      execvp(argv[1], argv + 1);
      return 0;
    }
    waitpid(pid, &status, 0);
  }
  return 0;
}
