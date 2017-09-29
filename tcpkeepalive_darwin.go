package tcpkeepalive

import (
	"os"
	"runtime"
	"syscall"
	"time"
)

const sysTCP_KEEPINTVL = 0x101

func setIdle(fd uintptr, d time.Duration) error {
	// not possible with darwin
	return nil
}

func setCount(fd uintptr, n int) error {
	// not possible with darwin
	return nil
}

func setInterval(fd uintptr, d time.Duration) error {
	// # from https://golang.org/src/net/tcpsockopt_darwin.go
	// d += (time.Second - time.Nanosecond)
	// secs := int(d.Seconds())
	// return os.NewSyscallError("setsockopt", syscall.SetsockoptInt(int(fd), syscall.IPPROTO_TCP, sysTCP_KEEPINTVL, secs))
	// The kernel expects seconds so round to next highest second.
	d += (time.Second - time.Nanosecond)
	secs := int(d.Seconds())
	switch err := syscall.SetsockoptInt(int(fd), syscall.IPPROTO_TCP, sysTCP_KEEPINTVL, secs); err {
	case nil, syscall.ENOPROTOOPT: // OS X 10.7 and earlier don't support this option
	default:
		return os.NewSyscallError("setsockopt", err)
	}
	err := syscall.SetsockoptInt(int(fd), syscall.IPPROTO_TCP, syscall.TCP_KEEPALIVE, secs)
	runtime.KeepAlive(fd)
	return os.NewSyscallError("setsockopt", err)
}
