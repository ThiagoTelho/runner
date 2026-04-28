//go:build windows

package process

import (
	"os"
	"os/exec"
	"syscall"
)

func detachProcess(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}
}

func isProcessRunning(p *os.Process) bool {
	// No Windows, tentamos obter o código de saída — processo vivo retorna erro "ainda em execução".
	handle, err := syscall.OpenProcess(syscall.PROCESS_QUERY_INFORMATION, false, uint32(p.Pid))
	if err != nil {
		return false
	}
	defer syscall.CloseHandle(handle)
	var code uint32
	err = syscall.GetExitCodeProcess(handle, &code)
	return err == nil && code == 259 // STILL_ACTIVE
}

func terminateProcess(p *os.Process) error {
	return p.Kill()
}
