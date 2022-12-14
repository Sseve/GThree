// 守护进程
package utils

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

func Daemon() {
	if syscall.Getppid() == 1 {
		if err := os.Chdir("./"); err != nil {
			panic(err)
		}
		syscall.Umask(0)
		return
	}
	fd, err := os.OpenFile("/dev/null", os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Println(err)
		return
	}
	defer fd.Close()
	cmd := exec.Command(os.Args[0], os.Args[1:]...)
	cmd.Stdout = fd
	cmd.Stderr = fd
	cmd.Stdin = nil
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	if err := cmd.Start(); err != nil {
		log.Println(err)
		return
	}
	os.Exit(0)
}
