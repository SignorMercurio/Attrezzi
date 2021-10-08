package scan

import (
	"net"
	"time"
)

func NewSYNScanner(targets []net.IP, timeout time.Duration, routines int) *ConnectScanner {
	return &ConnectScanner{}
}
