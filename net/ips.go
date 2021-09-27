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
	"encoding/hex"
	"fmt"
	"math/big"
	"net"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	cidr         string
	checkPrivate string
)

type CIDR struct {
	ip    net.IP
	ipnet *net.IPNet
}

func (c *CIDR) Mask() string {
	mask, _ := hex.DecodeString(c.ipnet.Mask.String())
	return net.IP(mask).String()
}

func (c *CIDR) BroadcastIP() string {
	mask, network := c.ipnet.Mask, c.ipnet.IP
	maskLen, networkLen := len(mask), len(network)
	b := network

	for i := 0; i < maskLen; i++ {
		idx := networkLen - i - 1
		maskIdx := maskLen - i - 1
		b[idx] = network[idx] | ^mask[maskIdx]
	}

	return b.String()
}

func (c *CIDR) Count() *big.Int {
	ones, bits := c.ipnet.Mask.Size()
	return big.NewInt(0).Lsh(big.NewInt(1), uint(bits-ones))
}

func ParseCIDR(s string) (*CIDR, error) {
	ip, inet, err := net.ParseCIDR(s)
	if err != nil {
		return nil, errors.Wrap(err, "parse CIDR")
	}

	return &CIDR{ip: ip, ipnet: inet}, nil
}

// NewIpsCmd represents the ips command
func NewIpsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ips",
		Short: "Show IP ranges",
		Long: `Show IP ranges
Example:
	att net ips --cidr 10.0.0.0/24`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if cidr != "" {
				return parseCIDR()
			}
			if checkPrivate != "" {
				ip := net.ParseIP(checkPrivate)
				fmt.Println(ip.IsPrivate())
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&cidr, "cidr", "", "CIDR")
	cmd.Flags().StringVar(&checkPrivate, "chk-priv", "", "Check whether the IP address is a private address")

	return cmd
}

func parseCIDR() error {
	cidr, err := ParseCIDR(cidr)
	if err != nil {
		return errors.Wrap(err, "parse CIDR")
	}

	network := cidr.ipnet.IP.String()
	s := fmt.Sprintf(`CIDR: %s
Network: %s
Mask: %s
Range: %s - %s
Numbers of IP: %d`,
		cidr.ipnet.String(),
		network,
		cidr.Mask(),
		network,
		cidr.BroadcastIP(),
		cidr.Count(),
	)
	fmt.Println(s)
	return nil
}

func init() {
	netCmd.AddCommand(NewIpsCmd())
}
