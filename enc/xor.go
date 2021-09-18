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
	"encoding/base64"

	"github.com/SignorMercurio/attrezzi/format"
	"github.com/lukechampine/fastxor"

	"github.com/spf13/cobra"
)

var (
	key    string
	inFmt  string
	keyFmt string
)

// NewXorCmd represents the xor command
func NewXorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "xor",
		Short: "XOR operation",
		Long: `XOR operation
Example:
	echo -n "hello" | att enc -o out.txt xor -k deadbeef`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if key != "" {
				inputByte, err := getByte(string(inputBytes), inFmt)
				if err != nil {
					return err
				}
				keyByte, err := getByte(key, keyFmt)
				if err != nil {
					return err
				}

				res := make([]byte, len(inputByte))
				fastxor.Bytes(res, inputByte, keyByte)
				Echo(string(res))
			} else {
				NoKeySpecified()
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&key, "key", "k", "", "Key to XOR with")
	cmd.Flags().StringVar(&inFmt, "input-fmt", "hex", "Format of input: hex / dec / bin / b64 / utf8")
	cmd.Flags().StringVar(&keyFmt, "key-fmt", "hex", "Format of key: hex / dec / bin / b64 / utf8")

	return cmd
}

// getByte gets the []byte form of the input or the key
func getByte(target string, fmt string) ([]byte, error) {
	arr := []string{target}

	switch fmt {
	case "bin":
		err := format.Bin2hex(arr)
		if err != nil {
			return nil, err
		}
		fallthrough
	case "hex":
		return format.DecodeHex(arr)
	case "dec":
		err := format.Dec2hex(arr)
		if err != nil {
			return nil, err
		}
		return format.DecodeHex(arr)
	case "b64":
		return base64.StdEncoding.DecodeString(target)
	default: // utf8
		return []byte(target), nil
	}
}

func init() {
	encCmd.AddCommand(NewXorCmd())
}
