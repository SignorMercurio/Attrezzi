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
package msc

import (
	"fmt"
	"io"
	"os"

	"github.com/SignorMercurio/attrezzi/cmd"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	inputFile  string
	input      *os.File
	outputFile string
	output     io.WriteCloser = os.Stdout
	mscCmd                    = NewMscCmd()
)

// NewMscCmd represents the msc command
func NewMscCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "msc",
		Short: "msc helps to deal with miscellaneous operations",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var err error
			if cmd.Use != "uid" {
				if err = getInput(); err != nil {
					return err
				}
			}

			if err = getOutput(); err != nil {
				return err
			}

			return nil
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			output.Close()
		},
	}

	cmd.PersistentFlags().StringVarP(&inputFile, "input", "i", "", "Read input from file")
	cmd.PersistentFlags().StringVarP(&outputFile, "output", "o", "", "Write output to file")

	return cmd
}

// getInput gets the inputFile
func getInput() error {
	var err error
	input, err = os.Open(inputFile)
	if err != nil {
		return errors.Wrap(err, "open input file")
	}
	return nil
}

// getOutput gets the output fd
func getOutput() error {
	var err error

	if outputFile != "" {
		output, err = os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return errors.Wrap(err, "open output file")
		}
	}

	return nil
}

func Echo(content interface{}) {
	fmt.Fprint(output, content)
}

func Warn(message string) {
	cmd.Log.Warn(message)
}

func init() {
	cmd.RootCmd.AddCommand(mscCmd)
}
