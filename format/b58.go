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
	"github.com/mr-tron/base58"

	"github.com/spf13/cobra"
)

var (
	b58Alphabet string
)

// NewB58Cmd represents the b58 command
func NewB58Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "b58",
		Short: "Base58 encode / decode",
		Long: `Base58 encode / decode
Example:
	echo -n "hello" | att fmt -o out.txt b58 -e
	att fmt -i in.txt b58 -d`,
		RunE: func(cmd *cobra.Command, args []string) error {
			enc := getB58Alphabet()

			if encode {
				encoded := base58.EncodeAlphabet(inputBytes, enc)
				Echo(encoded)
			} else if decode {
				decoded, err := base58.DecodeAlphabet(string(inputBytes), enc)
				if err != nil {
					return err
				}
				Echo(string(decoded))
			} else {
				NoActionSpecified()
			}
			return nil
		},
	}
	cmd.Flags().BoolVarP(&encode, "encode", "e", false, "Encode to base58")
	cmd.Flags().BoolVarP(&decode, "decode", "d", false, "Decode from base58")
	cmd.Flags().StringVarP(&b58Alphabet, "alphabet", "a", "btc", `58-byte Alphabet for base58, or "flickr" for Flickr alphabet`)

	return cmd
}

// getB58Alphabet gets the base58 alphabet from user input
func getB58Alphabet() *base58.Alphabet {
	var enc *base58.Alphabet
	switch b58Alphabet {
	case "btc":
		enc = base58.BTCAlphabet
	case "flickr":
		enc = base58.FlickrAlphabet
	default:
		enc = base58.NewAlphabet(b58Alphabet)
	}

	return enc
}

func init() {
	fmtCmd.AddCommand(NewB58Cmd())
}
