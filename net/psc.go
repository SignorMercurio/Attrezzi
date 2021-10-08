/*
Copyright © 2021 SignorMercurio

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package net

import (
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/SignorMercurio/attrezzi/scan"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	routines int
	scanType string
	portsStr string
)

// NewPscCmd represents the psc command
func NewPscCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "psc",
		Short: "Port scanning",
		Long: `Port scanning
Example:
	att net psc -t 192.168.1.1/24 -p 22,80,443,8000-8888
	att net psc -t example.com -r 100 --timeout 5 -s connect`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ports, err := parsePorts(portsStr)
			if err != nil {
				return err
			}

			ctx, cancel := context.WithCancel(context.Background())
			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt)
			go func() {
				<-c
				ScanCanceled()
				cancel()
			}()

			start := time.Now()
			ScanStart()
			targets, err := parseTargets(target)
			if err != nil {
				return err
			}
			scanner, err := createScanner(targets, scanType, time.Second*time.Duration(timeout), routines)
			if err != nil {
				return err
			}
			scanner.Start()

			results, err := scanner.Scan(ctx, ports)
			if err != nil {
				return err
			}
			for _, result := range results {
				if result.Latency > 0 {
					fmt.Println(result.String())
				}
			}
			ScanFinished(time.Since(start).String())
			return nil
		},
	}

	cmd.Flags().StringVarP(&target, "target", "t", "", "Target domain / IP / CIDR")
	cmd.Flags().StringVarP(&scanType, "scan-type", "s", "syn", "Scan type: connect / syn")
	cmd.Flags().IntVar(&timeout, "timeout", 2, "Scan timeout in seconds")
	cmd.Flags().IntVarP(&routines, "routines", "r", 1000, "Goroutines to use in scanning")
	cmd.Flags().StringVarP(&portsStr, "ports", "p", "22,80,443,8000-8888", "Ports to scan")

	return cmd
}

// create Scanner creates a connect or syn scanner
func createScanner(targets []net.IP, scanType string, timeout time.Duration, routines int) (scan.Scanner, error) {
	switch scanType {
	case "connect":
		return scan.NewConnectScanner(targets, timeout, routines), nil
	default:
		if os.Geteuid() > 0 {
			return nil, errors.New("You'll need root privilege to start an SYN scan.")
		}
		return scan.NewSYNScanner(targets, timeout, routines), nil
	}
}

// parsePorts parses the user input and returns a slice of ports
func parsePorts(portsStr string) ([]int, error) {
	var ports []int
	splitted := strings.Split(portsStr, ",")
	for _, s := range splitted {
		s = strings.TrimSpace(s)
		if strings.Contains(s, "-") {
			ranges := strings.Split(s, "-")
			if len(ranges) != 2 {
				return nil, errors.New("parse the ports")
			}

			from, err := strconv.Atoi(ranges[0])
			if err != nil {
				return nil, errors.Wrap(err, "parse the ports")
			}
			to, err := strconv.Atoi(ranges[1])
			if err != nil {
				return nil, errors.Wrap(err, "parse the ports")
			}
			if from > to {
				return nil, errors.New("parse the ports")
			}
			for port := from; port <= to; port++ {
				ports = append(ports, port)
			}
		} else {
			port, err := strconv.Atoi(s)
			if err != nil {
				return nil, errors.Wrap(err, "parse the ports")
			}
			ports = append(ports, port)
		}
	}
	return ports, nil
}

// parseTargets parses the user input and returns a slice of targets
func parseTargets(targetsStr string) ([]net.IP, error) {
	ip, inet, err := net.ParseCIDR(targetsStr)
	if err == nil { // CIDR
		cidr := &CIDR{ip: ip, ipnet: inet}
		network := binary.BigEndian.Uint32(cidr.ipnet.IP)
		broadcast := binary.BigEndian.Uint32(cidr.BroadcastIP())
		targets := []net.IP{}

		for addr := network; addr <= broadcast; addr++ {
			ip := make(net.IP, 4)
			binary.BigEndian.PutUint32(ip, addr)
			targets = append(targets, ip)
		}
		return targets, nil
	}
	ip = net.ParseIP(targetsStr)
	if ip != nil { // Single IP
		return []net.IP{ip}, nil
	}
	ips, err := net.LookupIP(targetsStr)
	if err == nil { // Domain
		if len(ips) == 0 {
			return nil, errors.New("lookup IP")
		}
		return ips, nil
	}
	return nil, errors.New("parse targets")
}

func init() {
	netCmd.AddCommand(NewPscCmd())
}
