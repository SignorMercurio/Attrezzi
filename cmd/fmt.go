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
package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	inputFile  string
	inputBytes []byte
	outputFile string
	output     io.Writer = os.Stdout
)

// fmtCmd represents the fmt command
var fmtCmd = &cobra.Command{
	Use:   "fmt",
	Short: "fmt helps to deal with data format.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		inputBytes, err = getInput()
		if err != nil {
			return err
		}
		if getOutput() != nil {
			return err
		}

		return nil
	},
}

func getInput() ([]byte, error) {
	var err error
	var inputBytes []byte
	input := os.Stdin

	if inputFile != "" {
		input, err = os.Open(inputFile)
		defer input.Close()
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

func getOutput() error {
	var err error

	if outputFile != "" {
		_, err = os.Stat(outputFile)
		if err != nil && os.IsNotExist(err) {
			output, err = os.Create(outputFile)
		} else {
			output, err = os.OpenFile(outputFile, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		}

		if err != nil {
			return errors.Wrap(err, "open output file")
		}
	}

	return nil
}

func Echo(content interface{}) {
	fmt.Fprintln(output, content)
}

func init() {
	rootCmd.AddCommand(fmtCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	fmtCmd.PersistentFlags().StringVarP(&inputFile, "input", "i", "", "Read input from file")
	fmtCmd.PersistentFlags().StringVarP(&outputFile, "output", "o", "", "Write output to file")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
}
