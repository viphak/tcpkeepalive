package tcpkeepalive

import (
	"os"
	"syscall"
	"time"
)

const _TCP_KEEPALIVE = syscall.TCP_KEEPALIVE
const _TCP_KEEPINTVL = 0x101 /* interval between keepalives */
const _TCP_KEEPCNT = 0x102   /* number of keepalives before close */

func setIdle(fd uintptr, d time.Duration) error {
	secs := durToSecs(d)
	err := syscall.SetsockoptInt(int(fd), syscall.IPPROTO_TCP, _TCP_KEEPALIVE, secs)
	return os.NewSyscallError("setsockopt", err)
}

func setCount(fd uintptr, n int) error {
	err := syscall.SetsockoptInt(int(fd), syscall.IPPROTO_TCP, _TCP_KEEPCNT, n)
	return os.NewSyscallError("setsockopt", err)
}

func setInterval(fd uintptr, d time.Duration) error {
	// # from https://golang.org/src/net/tcpsockopt_darwin.go
	secs := durToSecs(d)
	err := syscall.SetsockoptInt(int(fd), syscall.IPPROTO_TCP, _TCP_KEEPINTVL, secs)

	// OS X 10.7 and earlier don't support this option
	if err == nil || err == syscall.ENOPROTOOPT {
		return nil
	}

	return os.NewSyscallError("setsockopt", err)
}

func durToSecs(d time.Duration) int {
	d += (time.Second - time.Nanosecond)
	secs := int(d.Seconds())
	return secs
}
