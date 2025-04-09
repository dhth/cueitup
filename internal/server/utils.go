package server

import (
	"errors"
	"fmt"
	"net"
	"os/exec"
	"runtime"
)

const (
	startPort = 8500
	endPort   = 9500
)

var (
	errNoPortOpen                    = errors.New("no open port found")
	errUnsupportedPlatformForURLOpen = errors.New("opening URL is not supported on this platform")
	errOpenURLCmdFailed              = errors.New("command for opening URL failed")
)

func findOpenPort(startPort, endPort int) (int, bool) {
	for port := startPort; port <= endPort; port++ {
		address := fmt.Sprintf("127.0.0.1:%d", port)
		listener, err := net.Listen("tcp", address)
		if err == nil {
			defer listener.Close()
			return port, true
		}
	}
	return 0, false
}

func openURL(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	default:
		return errUnsupportedPlatformForURLOpen
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w; command output: %s", errOpenURLCmdFailed, out)
	}

	return nil
}
