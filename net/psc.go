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
	"github.com/spf13/cobra"
)

// NewPscCmd represents the psc command
func NewPscCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "psc",
		Short: "Port scanning",
		Long: `Port scanning
Example:
	att net psc -t google.com
	att net psc -r -t 142.250.187.206`,
		RunE: func(cmd *cobra.Command, args []string) error {

			return nil
		},
	}

	cmd.Flags().StringVarP(&target, "target", "t", "", "Target domain / IP")

	return cmd
}

func init() {
	netCmd.AddCommand(NewPscCmd())
}
