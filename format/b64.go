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
	"encoding/base64"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	encode   bool
	decode   bool
	alphabet string
	padding  string
)

// NewB64Cmd represents the b64 command
func NewB64Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "b64",
		Short: "Base64 encode / decode",
		Long: `Base64 encode / decode
Example:
	echo -n "hello" | att fmt -o out.txt b64 -e
	att fmt -i in.txt b64 -d
	echo -n "Attrezzi" | att fmt b64 -e | att fmt b64 -d`,
		RunE: func(cmd *cobra.Command, args []string) error {
			enc := getB64Encoding()

			if encode {
				encoded := enc.EncodeToString(inputBytes)
				Echo(encoded)
			} else if decode {
				decoded, err := decodeBase64(enc)
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
	cmd.Flags().BoolVarP(&encode, "encode", "e", false, "Encode to base64")
	cmd.Flags().BoolVarP(&decode, "decode", "d", false, "Decode from base64")
	cmd.Flags().StringVarP(&alphabet, "alphabet", "a", "std", `64-byte Alphabet for base64, or "url" for URLEncoding (See RFC4648)`)
	cmd.Flags().StringVarP(&padding, "padding", "p", "=", `Padding for base64, or "" for no padding`)

	return cmd
}

// getB64Encoding gets the base64 encoding from user input
func getB64Encoding() *base64.Encoding {
	var enc *base64.Encoding
	switch alphabet {
	case "std":
		enc = base64.StdEncoding
	case "url":
		enc = base64.URLEncoding
	default:
		enc = base64.NewEncoding(alphabet)
	}

	if padding == "" {
		return enc.WithPadding(base64.NoPadding)
	}

	return enc.WithPadding([]rune(padding)[0])
}

// decodeBase64 converts the inputBytes to a decoded string
func decodeBase64(enc *base64.Encoding) (string, error) {
	decoded, err := enc.DecodeString(string(inputBytes))
	if err != nil {
		return "", errors.Wrap(err, "decode base64")
	}
	return string(decoded), nil
}

func init() {
	fmtCmd.AddCommand(NewB64Cmd())
}
