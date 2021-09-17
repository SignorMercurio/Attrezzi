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
	"encoding/ascii85"

	"github.com/spf13/cobra"
)

// NewB85Cmd represents the b85 command
func NewB85Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "b85",
		Short: "Base85 encode / decode",
		Long: `Base85 encode / decode
Example:
	echo -n "hello" | att fmt -o out.txt b85 -e
	att fmt -i in.txt b85 -d`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if encode {
				dstLen := ascii85.MaxEncodedLen(len(inputBytes))
				encoded := make([]byte, dstLen)

				n := ascii85.Encode(encoded, inputBytes)
				Echo(string(encoded[:n]))
			} else if decode {
				decoded := make([]byte, 4*len(inputBytes))
				ndst, _, err := ascii85.Decode(decoded, inputBytes, true)
				if err != nil {
					return err
				}
				Echo(string(decoded[:ndst]))
			} else {
				NoActionSpecified()
			}
			return nil
		},
	}
	cmd.Flags().BoolVarP(&encode, "encode", "e", false, "Encode to base85")
	cmd.Flags().BoolVarP(&decode, "decode", "d", false, "Decode from base85")

	return cmd
}

func init() {
	fmtCmd.AddCommand(NewB85Cmd())
}
