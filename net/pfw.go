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
	"io"
	"net"
	"strings"
	"sync"

	"github.com/spf13/cobra"
)

var (
	srcAddr  string
	dstAddr  string
	timeout  int
	maxBytes = 32 * 1024
)

// NewPfwCmd represents the Pfw command
func NewPfwCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pfw",
		Short: "local / remote port forwarding",
		Long: `local / remote port forwarding
Note:
	The directions of application data flow and request flow are opposite.
	"src" means the **application data source**, so "src" also represents the **request destination**.
	Likewise, "dst" means the **application data destination**, which is actually the **request source**.
	Cascade port forwarding is also supported.
Example:
	Suppose you have a server at 100.100.100.100.
	To access a remote MySQL database on your server on your local port 33060:
		att net pfw -s c:100.100.100.100:3306 -d l::33060
	To present your local webpage on port 8080 on your server port 9090, using server port 7070:
		($ local)  att net pfw -s c::8080 -d c:100.100.100.100:7070
		($ server) att net pfw -s l::7070 -d l::9090`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if srcAddr == "" || dstAddr == "" {
				NoAddressSpecified()
			}

			srcType, src := parseAddr(srcAddr)
			dstType, dst := parseAddr(dstAddr)

			if srcType == "c" && dstType == "l" {
				listen2dial(dst, src)
			} else if srcType == "l" && dstType == "c" {
				listen2dial(src, dst)
			} else if srcType == "c" && dstType == "c" {
				dial2dial(dst, src) // actually handling request flow, so it's reversed
			} else if srcType == "l" && dstType == "l" {
				listen2listen(src, dst)
			} else {
				WrongAddressFormat()
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&srcAddr, "src-addr", "s", "", "Data source address, with the type (l / c) and a colon at the begining")
	cmd.Flags().StringVarP(&dstAddr, "dst-addr", "d", "", "Data destination address, with the type (l / c) and a colon at the begining")
	cmd.Flags().IntVarP(&timeout, "timeout", "t", 10, "Network timeout in seconds")

	return cmd
}

// parseAddr returns the type and the actual address
func parseAddr(addr string) (string, string) {
	parsed := strings.SplitN(addr, ":", 2)
	return parsed[0], parsed[1]
}

// connect prints messages about connection status and calls DialTCP()
func connect(addr string) (net.Conn, error) {
	Connecting(addr)
	dialConn, err := DialTCP(addr, timeout)
	if err != nil {
		return nil, err
	}
	Connected(addr)
	return dialConn, nil
}

// listen2dial links [listen] with [dial]
func listen2dial(listen, dial string) {
	client := make(chan net.Conn)
	go ListenTCP(listen, client)

	for {
		listenConn := <-client
		dialConn, err := connect(dial)
		if err != nil {
			listenConn.Close()
			Fail2Connect(dial)
		}
		go link(listenConn, dialConn)
	}
}

// dial2dial links [dial1] with [dial2]
func dial2dial(dial1, dial2 string) {
	for {
		dial1Conn, err := connect(dial1)
		if err != nil {
			Fail2Connect(dial1)
		}

		// Read request from dst
		buf := make([]byte, maxBytes)
		n, err := dial1Conn.Read(buf)
		if err != nil {
			Fail2Read(dial1)
		}

		dial2Conn, err := connect(dial2)
		if err != nil {
			Fail2Connect(dial2)
		}

		_, err = dial2Conn.Write(buf[:n])
		if err != nil {
			Fail2Write(dial2)
		}

		go link(dial1Conn, dial2Conn)
	}
}

// listen2listen links [listen1] with [listen2]
func listen2listen(listen1, listen2 string) {
	client1 := make(chan net.Conn)
	client2 := make(chan net.Conn)
	go ListenTCP(listen1, client1)
	go ListenTCP(listen2, client2)

	for {
		go link(<-client1, <-client2)
	}
}

// link forwards data between a & b and wait
func link(a, b net.Conn) {
	var wg sync.WaitGroup
	defer a.Close()
	defer b.Close()

	wg.Add(2)
	go func() {
		io.Copy(b, a)
		wg.Done()
	}()
	go func() {
		io.Copy(a, b)
		wg.Done()
	}()
	wg.Wait()
}

func init() {
	netCmd.AddCommand(NewPfwCmd())
}
