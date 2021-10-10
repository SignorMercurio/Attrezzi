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
	"github.com/SignorMercurio/attrezzi/cmd"
	"github.com/spf13/cobra"
)

var (
	netCmd = NewNetCmd()
)

// NewNetCmd represents the net command
func NewNetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "net",
		Short: "net helps to deal with network-related operations",
	}

	return cmd
}

func Fail2Accept(err error) {
	cmd.Log.Panicf("Failed to accept the client: %s", err)
}

func Accepted(addr string) {
	cmd.Log.Infof("Client accepted from %s", addr)
}

func Connecting(addr string) {
	cmd.Log.Infof("Connecting to %s...", addr)
}

func Fail2Connect(addr string) {
	cmd.Log.Panicf("Failed to connect to %s. Is the service running?", addr)
}

func Connected(addr string) {
	cmd.Log.Infof("Connected to %s", addr)
}

func Fail2Resolve(addr string) {
	cmd.Log.Panicf("Failed to resolve %s. Please check the address format.", addr)
}

func Fail2Listen(addr string) {
	cmd.Log.Panicf("Failed to listen on %s. Is the port in use?", addr)
}

func Listening(addr string) {
	cmd.Log.Infof("Listening on %s...", addr)
}

func Fail2Read(addr string) {
	cmd.Log.Panicf("Failed to read from %s", addr)
}

func Fail2Write(addr string) {
	cmd.Log.Panicf("Failed to write to %s", addr)
}

func NoAddressSpecified() {
	cmd.Log.Panic("Please specify both addresses.")
}

func WrongAddressFormat() {
	cmd.Log.Panic("Please check the address format.")
}

func ScanCanceled() {
	cmd.Log.Info("Canceling scanning...")
}

func ScanStart() {
	cmd.Log.Info("Starting scanning...")
}

func ScanFinished(duration string) {
	cmd.Log.Infof("Scan finished in %s.", duration)
}

func init() {
	cmd.RootCmd.AddCommand(netCmd)
}
