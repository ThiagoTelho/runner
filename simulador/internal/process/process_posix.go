//go:build !windows

package process

import (
	"os"
	"os/exec"
	"syscall"
)

func detachProcess(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
}

func isProcessRunning(p *os.Process) bool {
	return p.Signal(syscall.Signal(0)) == nil
}

func terminateProcess(p *os.Process) error {
	if err := p.Signal(syscall.SIGTERM); err != nil {
		return p.Kill()
	}
	return nil
}
