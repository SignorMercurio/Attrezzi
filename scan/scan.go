package scan

import (
	"context"
	"fmt"
	"net"
	"time"
)

type Result struct {
	Host     net.IP
	Open     []int
	Filtered []int
	Closed   []int
	Latency  time.Duration
}

// NewResult returns a new Result based on [host]
func NewResult(host string) *Result {
	return &Result{
		Host: net.ParseIP(host),
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
	Start() error
	Scan(ctx context.Context, ports []int) ([]Result, error)
}
