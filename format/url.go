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
	"net/url"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	all bool
)

// NewUrlCmd represents the url command
func NewUrlCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "url",
		Short: "URL encode / decode",
		Long: `URL encode / decode
Example:
	echo -n "hello" | att fmt -o out.txt url -ea
	att fmt -i in.txt url -d
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if encode {
				var encoded string
				if all {
					encoded = url.QueryEscape(string(inputBytes))
				} else {
					resURL, err := url.Parse(string(inputBytes))
					if err != nil {
						return errors.Wrap(err, "parse URL")
					}
					encoded = resURL.String()
				}
				Echo(encoded)
			} else if decode {
				decoded, err := url.QueryUnescape(string(inputBytes))
				if err != nil {
					return errors.Wrap(err, "decode URL")
				}
				Echo(decoded)
			} else {
				NoActionSpecified()
			}
			return nil
		},
	}

	cmd.Flags().BoolVarP(&encode, "encode", "e", false, "URL encode")
	cmd.Flags().BoolVarP(&decode, "decode", "d", false, "URL decode")
	cmd.Flags().BoolVarP(&all, "all", "a", false, "URL encode all special characters")

	return cmd
}

func init() {
	fmtCmd.AddCommand(NewUrlCmd())
}
