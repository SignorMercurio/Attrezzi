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
	"github.com/spf13/cobra"
)

var (
	encode    bool
	decode    bool
	rotNumber uint8
)

// NewRotCmd represents the rot command
func NewRotCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rot",
		Short: "ROT13-like encryption / decryption",
		Long: `ROT13-like encryption / decryption
Example:
	echo -n "hello" | att enc -o out.txt rot -e
	att enc -i in.txt rot -n 13 -d
	echo -n "Attrezzi" | att enc rot -e | att enc rot -d`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if encode {
				encoded := rotEncode(inputBytes)
				Echo(encoded)
			} else if decode {
				decoded := rotDecode(inputBytes)
				Echo(decoded)
			} else {
				NoActionSpecified()
			}
			return nil
		},
	}
	cmd.Flags().BoolVarP(&encode, "encode", "e", false, "Encode to base64")
	cmd.Flags().BoolVarP(&decode, "decode", "d", false, "Decode from base64")
	cmd.Flags().Uint8VarP(&rotNumber, "number", "n", 13, "Number to shift")

	return cmd
}

func rot(b byte, shift uint8) byte {
	var a, z byte
	switch {
	case 'a' <= b && b <= 'z':
		a, z = 'a', 'z'
	case 'A' <= b && b <= 'Z':
		a, z = 'A', 'Z'
	default:
		return b
	}

	return a + (b-a+shift)%(z-a+1)
}

func rotEncode(src []byte) string {
	for i, v := range src {
		src[i] = rot(v, rotNumber)
	}
	return string(src)
}

func rotDecode(src []byte) string {
	for i, v := range src {
		src[i] = rot(v, 26-rotNumber)
	}
	return string(src)
}

func init() {
	encCmd.AddCommand(NewRotCmd())
}
