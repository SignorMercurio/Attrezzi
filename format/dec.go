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
package format

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// NewDecCmd represents the dec command
func NewDecCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dec",
		Short: "convert string to / from decimal",
		Long: `convert string to / from decimal
Example:
	echo -n "hello" | att fmt -o out.txt dec -e --delim=" "
	att fmt -i in.txt dec -d`,
		RunE: func(cmd *cobra.Command, args []string) error {
			delimiter := getDelimiter()
			if delimiter == nil {
				EmptyDelimiter()
			}

			if encode {
				encoded := encodeToDec(inputBytes, delimiter)
				Echo(encoded)
			} else if decode {
				arr := getDecodeArr(delimiter)
				err := dec2hex(arr)
				if err != nil {
					return err
				}
				decoded, err := decodeHex(arr)
				if err != nil {
					return err
				}
				Echo(decoded)
			} else {
				NoActionSpecified()
			}
			return nil
		},
	}
	cmd.Flags().BoolVarP(&encode, "encode", "e", false, "Encode to decimal")
	cmd.Flags().BoolVarP(&decode, "decode", "d", false, "Decode from decimal")
	cmd.Flags().StringVar(&delim, "delim", "", "Delimiter")
	cmd.Flags().BoolVarP(&delim_prefix, "prefix", "p", false, "Whether the delimiter is a prefix")

	return cmd
}

// encodeTodec converts a []byte to a decimal string
func encodeToDec(src []byte, delimiter []byte) string {
	buf := bytes.NewBuffer([]byte{})

	if delim_prefix {
		buf.Write(delimiter)
	}

	for _, v := range src {
		buf.WriteString(fmt.Sprintf("%d", v))
		buf.Write(delimiter)
	}

	return buf.String()
}

// dec2hex converts a slice of decimal string to a slice of hex string
func dec2hex(arr []string) error {
	for i, v := range arr {
		s, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return errors.Wrap(err, "convert decimal to hex")
		}
		arr[i] = strconv.FormatInt(s, 16)
	}
	return nil
}

func init() {
	fmtCmd.AddCommand(NewDecCmd())
}
