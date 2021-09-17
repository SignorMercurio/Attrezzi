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
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// NewUniCmd represents the uni command
func NewUniCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uni",
		Short: "Unicode conversion",
		Long: `Unicode conversion
Example:
	echo -n "hello" | att fmt -o out.txt uni -e
	att fmt -i in.txt uni -d
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if encode {
				quoted := strconv.QuoteToASCII(string(inputBytes))
				encoded := quoted[1 : len(quoted)-1] // strip ""
				Echo(encoded)
			} else if decode {
				decoded, err := fromUnicode(string(inputBytes))
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

	cmd.Flags().BoolVarP(&encode, "encode", "e", false, "convert to unicode")
	cmd.Flags().BoolVarP(&decode, "decode", "d", false, "convert from unicode")

	return cmd
}

func fromUnicode(from string) (string, error) {
	str, err := strconv.Unquote(strings.Replace(
		strconv.Quote(from),
		`\\u`,
		`\u`,
		-1,
	))
	if err != nil {
		return "", errors.Wrap(err, "decode unicode")
	}
	return str, nil
}

func init() {
	fmtCmd.AddCommand(NewUniCmd())
}
