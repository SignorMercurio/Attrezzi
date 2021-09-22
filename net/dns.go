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
	"fmt"
	"net"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	target  string
	reverse bool
)

// NewDnsCmd represents the dns command
func NewDnsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dns",
		Short: "DNS lookup",
		Long: `DNS lookup
Example:
	att net dns -t google.com
	att net dns -r -t 142.250.187.206`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			os.Setenv("GODEBUG", "netdns=go")
			if reverse {
				if err = getAddr(); err != nil {
					return err
				}
				return nil
			}

			if err = getNS(); err != nil {
				return err
			}
			if err = getIP(); err != nil {
				return err
			}
			if err = getCNAME(); err != nil {
				return err
			}
			if err = getMX(); err != nil {
				return err
			}
			if err = getTXT(); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&target, "target", "t", "", "Target domain / IP")
	cmd.Flags().BoolVarP(&reverse, "reverse", "r", false, "Reverse lookup (target must be an IP)")

	return cmd
}

func getAddr() error {
	names, err := net.LookupAddr(target)
	if err != nil {
		return errors.Wrap(err, "reverse lookup domain names")
	}
	fmt.Println("Domain names:")
	for _, v := range names {
		fmt.Println(v)
	}

	return nil
}

func getNS() error {
	ns, err := net.LookupNS(target)
	if err != nil {
		return errors.Wrap(err, "lookup nameservers")
	}
	fmt.Println("Nameservers:")
	for _, v := range ns {
		fmt.Println(v.Host)
	}
	fmt.Println()

	return nil
}

func getIP() error {
	ip, err := net.LookupIP(target)
	if err != nil {
		return errors.Wrap(err, "lookup IPs")
	}
	fmt.Println("IPs:")
	for _, v := range ip {
		fmt.Println(v.String())
	}
	fmt.Println()

	return nil
}

func getCNAME() error {
	cname, err := net.LookupCNAME(target)
	if err != nil {
		return errors.Wrap(err, "lookup CNAME")
	}
	fmt.Println("CNAME:")
	fmt.Println(cname)
	fmt.Println()

	return nil
}

func getMX() error {
	mx, err := net.LookupMX(target)
	if err != nil {
		return errors.Wrap(err, "lookup MX records")
	}
	fmt.Println("MX records:")
	for _, v := range mx {
		fmt.Println(v.Host, v.Pref)
	}
	fmt.Println()

	return nil
}

func getTXT() error {
	txt, err := net.LookupTXT(target)
	if err != nil {
		return errors.Wrap(err, "lookup TXT records")
	}
	fmt.Println("TXT records:")
	for _, v := range txt {
		fmt.Println(v)
	}

	return nil
}

func init() {
	netCmd.AddCommand(NewDnsCmd())
}
