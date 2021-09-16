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
	"encoding/base64"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	encode   bool
	decode   bool
	alphabet string
)

// b64Cmd represents the b64 command
var b64Cmd = &cobra.Command{
	Use: "b64",
	Long: `base64 encode / decode
Example:
	echo -n "hello" | attrezzi fmt -o out.txt b64 -e
	attrezzi fmt -i in.txt b64 -d`,
	RunE: func(cmd *cobra.Command, args []string) error {
		enc := getEncoding()

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

// getEncoding gets the base64 encoding from user input
func getEncoding() *base64.Encoding {
	var enc *base64.Encoding
	switch alphabet {
	case "std":
		enc = base64.StdEncoding
	case "url":
		enc = base64.URLEncoding
	default:
		enc = base64.NewEncoding(alphabet)
	}

	return enc
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
	fmtCmd.AddCommand(b64Cmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// b64Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	b64Cmd.Flags().BoolVarP(&encode, "encode", "e", false, "Encode to base64")
	b64Cmd.Flags().BoolVarP(&decode, "decode", "d", false, "Decode from base64")
	b64Cmd.Flags().StringVarP(&alphabet, "alphabet", "a", "std", `Alphabet for base64, or "url" for URLEncoding`)
}
