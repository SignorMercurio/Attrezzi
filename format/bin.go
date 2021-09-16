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

const (
	byteLen = 8
)

// NewBinCmd represents the bin command
func NewBinCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bin",
		Short: "convert string to / from binary",
		Long: `convert string to / from binary
Example:
	echo -n "hello" | att fmt -o out.txt bin -e --delim=" "
	att fmt -i in.txt bin -d`,
		RunE: func(cmd *cobra.Command, args []string) error {
			delimiter := getDelimiter()
			if delimiter == nil {
				EmptyDelimiter()
			}

			if encode {
				encoded := encodeToBin(inputBytes)
				Echo(insertInto(encoded, byteLen, delimiter))
			} else if decode {
				arr := getDecodeArr(delimiter)
				err := bin2hex(arr)
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
	cmd.Flags().BoolVarP(&encode, "encode", "e", false, "Encode to binary")
	cmd.Flags().BoolVarP(&decode, "decode", "d", false, "Decode from binary")
	cmd.Flags().StringVar(&delim, "delim", " ", "Delimiter")
	cmd.Flags().BoolVarP(&delim_prefix, "prefix", "p", false, "Whether the delimiter is a prefix")

	return cmd
}

// encodeToBin converts a []byte to a binary string
func encodeToBin(src []byte) string {
	buf := bytes.NewBuffer([]byte{})

	for _, v := range src {
		buf.WriteString(fmt.Sprintf("%08b", v))
	}

	return buf.String()
}

// bin2hex converts a slice of binary string to a slice of hex string
func bin2hex(arr []string) error {
	for i, v := range arr {
		s, err := strconv.ParseInt(v, 2, 64)
		if err != nil {
			return errors.Wrap(err, "convert binary to hex")
		}
		arr[i] = strconv.FormatInt(s, 16)
	}
	return nil
}

func init() {
	fmtCmd.AddCommand(NewBinCmd())
}
