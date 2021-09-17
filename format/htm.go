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
	"html"

	"github.com/spf13/cobra"
)

// NewHtmCmd represents the htm command
func NewHtmCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "htm",
		Short: "HTML Entity encode / decode",
		Long: `HTML Entity encode / decode
Example:
	echo -n "hello" | att fmt -o out.txt htm -e
	att fmt -i in.txt htm -d
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if encode {
				encoded := html.EscapeString(string(inputBytes))
				Echo(encoded)
			} else if decode {
				decoded := html.UnescapeString(string(inputBytes))
				Echo(decoded)
			} else {
				NoActionSpecified()
			}
			return nil
		},
	}

	cmd.Flags().BoolVarP(&encode, "encode", "e", false, "HTML Entity encode")
	cmd.Flags().BoolVarP(&decode, "decode", "d", false, "HTML Entity decode")

	return cmd
}

func init() {
	fmtCmd.AddCommand(NewHtmCmd())
}
