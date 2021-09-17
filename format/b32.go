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
	"encoding/base32"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// NewB32Cmd represents the b32 command
func NewB32Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "b32",
		Short: "Base32 encode / decode",
		Long: `Base32 encode / decode
Example:
	echo -n "hello" | att fmt -o out.txt b32 -e
	att fmt -i in.txt b32 -d`,
		RunE: func(cmd *cobra.Command, args []string) error {
			enc := getB32Encoding()

			if encode {
				encoded := enc.EncodeToString(inputBytes)
				Echo(encoded)
			} else if decode {
				decoded, err := decodeBase32(enc)
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
	cmd.Flags().BoolVarP(&encode, "encode", "e", false, "Encode to base32")
	cmd.Flags().BoolVarP(&decode, "decode", "d", false, "Decode from base32")
	cmd.Flags().StringVarP(&alphabet, "alphabet", "a", "std", `Alphabet for base32, or "hex" for hex encoding (See RFC 4648)`)
	cmd.Flags().StringVarP(&padding, "padding", "p", "=", `Padding for base32, or "" for no padding`)

	return cmd
}

// getB32Encoding gets the base32 encoding from user input
func getB32Encoding() *base32.Encoding {
	var enc *base32.Encoding
	switch alphabet {
	case "std":
		enc = base32.StdEncoding
	case "hex":
		enc = base32.HexEncoding
	default:
		enc = base32.NewEncoding(alphabet)
	}
	if padding == "" {
		return enc.WithPadding(base32.NoPadding)
	}

	return enc.WithPadding([]rune(padding)[0])
}

// decodeBase32 converts the inputBytes to a decoded string
func decodeBase32(enc *base32.Encoding) (string, error) {
	decoded, err := enc.DecodeString(string(inputBytes))
	if err != nil {
		return "", errors.Wrap(err, "decode base32")
	}
	return string(decoded), nil
}

func init() {
	fmtCmd.AddCommand(NewB32Cmd())
}
