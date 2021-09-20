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
package enc

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/SignorMercurio/attrezzi/cmd"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	inputFile  string
	input      io.Reader = os.Stdin
	inputBytes []byte
	outputFile string
	output     io.Writer = os.Stdout
	encCmd               = NewEncCmd()
)

// NewEncCmd represents the enc command
func NewEncCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enc",
		Short: "enc helps to deal with cryptographic operations",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var err error
			if cmd.Use != "rnd" && cmd.Use != "rkg" {
				inputBytes, err = getInput()
				if err != nil {
					return err
				}
			}

			if getOutput() != nil {
				return err
			}

			return nil
		},
	}

	cmd.PersistentFlags().StringVarP(&inputFile, "input", "i", "", "Read input from file")
	cmd.PersistentFlags().StringVarP(&outputFile, "output", "o", "", "Write output to file")

	return cmd
}

// getInput gets the input []byte
func getInput() ([]byte, error) {
	var err error
	var inputBytes []byte

	if inputFile != "" {
		input, err = os.Open(inputFile)
		if err != nil {
			return nil, errors.Wrap(err, "open input file")
		}
	}

	inputBytes, err = ioutil.ReadAll(input)
	if err != nil {
		return nil, errors.Wrap(err, "read input file")
	}

	return inputBytes, nil
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

func NoActionSpecified() {
	cmd.Log.Error("No action specified. Please specify -e or -d.")
}

func NoKeySpecified() {
	cmd.Log.Error("No key specified. Please specify one.")
}

func init() {
	cmd.RootCmd.AddCommand(encCmd)
}
