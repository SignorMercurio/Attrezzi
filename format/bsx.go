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
	"github.com/eknkc/basex"
	"github.com/pkg/errors"

	"github.com/spf13/cobra"
)

const (
	b62Alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

var (
	baseXAlphabet string
	base          uint8
)

// NewBsxCmd represents the bsx command
func NewBsxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bsx",
		Short: "BaseX encode / decode (default Base62)",
		Long: `BaseX encode / decode (default Base62)
Note: Do not use this module to deal with Base32 / Base58 / Base64 / Base85 !

Example:
	echo -n "hello" | att fmt -o out.txt bsx -e -a 0123456789abcdef
	att fmt -i in.txt bsx -d -a 0123456789ABCDEFGHJKMNPQRSTVWXYZ`,
		RunE: func(cmd *cobra.Command, args []string) error {
			enc, err := getBsxEncoding()
			if err != nil {
				return err
			}

			if encode {
				encoded := enc.Encode(inputBytes)
				Echo(encoded)
			} else if decode {
				decoded, err := enc.Decode(string(inputBytes))
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
	cmd.Flags().BoolVarP(&encode, "encode", "e", false, "Encode to baseX")
	cmd.Flags().BoolVarP(&decode, "decode", "d", false, "Decode from baseX")
	cmd.Flags().Uint8VarP(&base, "base", "b", 62, `Value of X`)
	cmd.Flags().StringVarP(&baseXAlphabet, "alphabet", "a", "", `X-byte Alphabet for baseX`)

	return cmd
}

// getBsxEncoding gets the baseX encoding from user input
func getBsxEncoding() (*basex.Encoding, error) {
	if baseXAlphabet == "" { // alphabet has higher priority than base
		baseXAlphabet = b62Alphabet[:base]
	}

	enc, err := basex.NewEncoding(baseXAlphabet)
	if err != nil {
		return nil, errors.Wrap(err, "parse baseX alphabet")
	}

	return enc, nil
}

func init() {
	fmtCmd.AddCommand(NewBsxCmd())
}
