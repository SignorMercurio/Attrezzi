package scan

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

type ConnectScanner struct {
	scan     chan portScan
	targets  []net.IP
	timeout  time.Duration
	routines int
}

func NewConnectScanner(targets []net.IP, timeout time.Duration, routines int) *ConnectScanner {
	return &ConnectScanner{
		targets:  targets,
		timeout:  timeout,
		routines: routines,
		scan:     make(chan portScan, routines),
	}
}

// Start starts a bunch of routines waiting for port scanning jobs
func (cs *ConnectScanner) Start() {
	for i := 0; i < cs.routines; i++ {
		go func() {
			for {
				scan := <-cs.scan
				if scan.port == 0 {
					break
				}

				select {
				case <-scan.ctx.Done():
					close(scan.done)
					return
				default:
				}

				if portState, err := cs.portScan(scan.ip, scan.port); err == nil {
					switch portState {
					case Open:
						scan.open <- scan.port
					case Filtered:
						scan.filtered <- scan.port
					case Closed:
						scan.closed <- scan.port
					}
				}
				close(scan.done)
			}
		}()
	}
}

// isClosed returns if the port is closed
func isClosed(err error) bool {
	return strings.Contains(err.Error(), "refused")
}

// portScan use connect() to scan ports
func (cs *ConnectScanner) portScan(ip net.IP, port int) (int, error) {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip.String(), port), cs.timeout)
	if err != nil {
		if isClosed(err) {
			return Closed, nil
		}
		return Unknown, err
	}
	conn.Close()
	return Open, err
}

// hostScan scans a single host and pushs portScan jobs into channel
func (cs *ConnectScanner) hostScan(ctx context.Context, host net.IP, ports []int) Result {
	var (
		wg           = &sync.WaitGroup{}
		result       = NewResult(host)
		openChan     = make(chan int)
		filteredChan = make(chan int)
		closedChan   = make(chan int)
		done         = make(chan struct{})
	)
	start := time.Now()

	go func() {
		for {
			select {
			case open := <-openChan:
				if open == 0 {
					close(done)
					return
				}
				if result.Latency <= 0 {
					result.Latency = time.Since(start)
				}
				result.Open = append(result.Open, open)
			case filtered := <-filteredChan:
				if result.Latency <= 0 {
					result.Latency = time.Since(start)
				}
				result.Filtered = append(result.Filtered, filtered)
			case closed := <-closedChan:
				if result.Latency <= 0 {
					result.Latency = time.Since(start)
				}
				result.Closed = append(result.Closed, closed)
			}
		}
	}()

	for _, port := range ports {
		wg.Add(1)
		go func(port int, wg *sync.WaitGroup) {
			done := make(chan struct{})
			cs.scan <- portScan{
				open:     openChan,
				filtered: filteredChan,
				closed:   closedChan,
				ip:       host,
				port:     port,
				done:     done,
				ctx:      ctx,
			}
			<-done
			wg.Done()
		}(port, wg)
	}
	wg.Wait()
	close(openChan)
	<-done

	return result
}

// Scan starts scanning all the hosts
func (cs *ConnectScanner) Scan(ctx context.Context, ports []int) ([]Result, error) {
	var (
		wg         = &sync.WaitGroup{}
		resultChan = make(chan *Result)
		results    = []Result{}
		done       = make(chan struct{})
	)

	go func() {
		for {
			result := <-resultChan
			if result == nil {
				close(done)
				break
			}
			results = append(results, *result)
		}
	}()

	for _, ip := range cs.targets {
		wg.Add(1)
		go func(ip net.IP, ports []int, wg *sync.WaitGroup) {
			result := cs.hostScan(ctx, ip, ports)
			resultChan <- &result
			wg.Done()
		}(ip, ports, wg)
	}
	wg.Wait()
	close(resultChan)
	close(cs.scan)
	<-done

	return results, nil
}
