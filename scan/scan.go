package scan

import (
	"context"
	"fmt"
	"net"
	"time"
)

const (
	Unknown int = iota
	Open
	Closed
	Filtered
)

type Result struct {
	Host     net.IP
	Open     []int
	Filtered []int
	Closed   []int
	Latency  time.Duration
}

// NewResult returns a new Result based on [host]
func NewResult(host net.IP) Result {
	return Result{
		Host: host,
	}
}

// interface Stringer
func (r *Result) String() string {
	info := fmt.Sprintf("Scan results for %s\n", r.Host.String())

	if r.Latency <= 0 {
		info += "\tHost is down\n"
		return info
	}

	info += fmt.Sprintf("\tHost is up, with a latency of %s. Open ports:\n", r.Latency.String())
	for _, port := range r.Open {
		info += fmt.Sprintf(
			"\t%10s\t%s\n",
			fmt.Sprintf("%d/tcp", port),
			getService(port),
		)
	}

	return info
}

type Scanner interface {
	Start()
	Scan(ctx context.Context, ports []int) ([]Result, error)
}

type portScan struct {
	ip       net.IP
	port     int
	open     chan int
	filtered chan int
	closed   chan int
	done     chan struct{}
	ctx      context.Context
}

type hostScan struct {
	ip     net.IP
	ports  []int
	result chan *Result
	done   chan struct{}
	ctx    context.Context
}
