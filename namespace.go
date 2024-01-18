package main

import (
	"os"
	"os/exec"
	"syscall"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// mount -t overlay overlay -o lowerdir=/root/ns-mount/lower,upperdir=/root/ns-mount/upper,workdir=/root/ns-mount/worker/ /root/ns-mount/merged/
// chroot /root/ns-mount/merged/ /bin/sh
// mkdir /proc && mount -t proc proc /proc

func parseArgs(args *namespaceType) uintptr {
	var flags uintptr = 0x0

	if args.uts { flags |= syscall.CLONE_NEWUTS }
	if args.ipc { flags |= syscall.CLONE_NEWIPC }
	if args.network { flags |= syscall.CLONE_NEWNET }
	if args.pid { flags |= syscall.CLONE_NEWPID }
	if args.user { flags |= syscall.CLONE_NEWUSER }
	if args.mount { flags |= syscall.CLONE_NEWNS }
	if args.cgroup { flags |= syscall.CLONE_NEWCGROUP }

	if args.all { flags = syscall.CLONE_NEWUTS |
		syscall.CLONE_NEWIPC |
		syscall.CLONE_NEWNET |
		syscall.CLONE_NEWPID |
		syscall.CLONE_NEWUSER |
		syscall.CLONE_NEWNS |
		syscall.CLONE_NEWCGROUP }

	return flags
}

func createNamespace(args *namespaceType) error {
	flags := parseArgs(args)
	return newBash(flags)
}

func newBash(flags uintptr) error {
	cmd := exec.Command("bash")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: flags,
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Create bash error: %v\n", err)
		log.Errorf("Create bash error: %v", err)
		return err
	}
	
	return nil
}
