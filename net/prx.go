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
	"net"
	"net/http"

	"github.com/elazarl/goproxy"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/txthinking/socks5"
)

var (
	listen        string
	mode          string
	socksUsername string
	socksPassword string
	tcpTimeout    int
	udpTimeout    int
)

// NewPrxCmd represents the prx command
func NewPrxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "prx",
		Short: "Proxy",
		Long: `Proxy
Example:
	att net prx -t :8080
	att net prx -t :1080 -m socks5 -u merc -p 123`,
		RunE: func(cmd *cobra.Command, args []string) error {
			switch mode {
			case "socks5":
				host, _, err := net.SplitHostPort(listen)
				if err != nil {
					return errors.Wrap(err, "parse listen address")
				}

				if host == "" {
					host = "0.0.0.0"
				}

				proxy, err := socks5.NewClassicServer(listen, host, socksUsername, socksPassword, tcpTimeout, udpTimeout)
				if err != nil {
					return errors.Wrap(err, "create SOCKS5 proxy server")
				}
				Listening(listen)
				if err = proxy.ListenAndServe(nil); err != nil {
					return errors.Wrap(err, "start SOCKS5 proxy server")
				}
			default:
				proxy := goproxy.NewProxyHttpServer()
				proxy.Verbose = true
				Listening(listen)
				if err := http.ListenAndServe(listen, proxy); err != nil {
					return errors.Wrap(err, "start HTTP(s) proxy server")
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&listen, "listen", "l", ":8080", "Proxy listen address")
	cmd.Flags().StringVarP(&mode, "mode", "m", "http", "Proxy mode: http(s) / socks5")
	cmd.Flags().StringVarP(&socksUsername, "username", "u", "", "SOCKS5 username")
	cmd.Flags().StringVarP(&socksPassword, "password", "p", "", "SOCKS5 password")
	cmd.Flags().IntVarP(&tcpTimeout, "tcp-timeout", "t", 0, "SOCKS5 TCP timeout")
	cmd.Flags().IntVar(&udpTimeout, "udp-timeout", 60, "SOCK5 UDP timeout")

	return cmd
}

func init() {
	netCmd.AddCommand(NewPrxCmd())
}
