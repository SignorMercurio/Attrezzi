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
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/SignorMercurio/attrezzi/format"
	"github.com/spf13/cobra"
)

var (
	byteLength uint
	numFmt     string
)

// NewRndCmd represents the rnd command
func NewRndCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rnd",
		Short: "Random number generation",
		Long: `Random number generation
Example:
	att enc rnd -l 16 -f hex
	att enc -o out.txt rnd -l 8 -f bin`,
		Run: func(cmd *cobra.Command, args []string) {
			b := make([]byte, byteLength)
			rand.Read(b)

			rnd := formatRnd(b)
			Echo(rnd)
		},
	}
	cmd.Flags().UintVarP(&byteLength, "length", "l", 8, "Byte length of generated number")
	cmd.Flags().StringVarP(&numFmt, "format", "f", "hex", "Format of generated number: hex / bin / dec")

	return cmd
}

func formatRnd(b []byte) string {
	rnd := hex.EncodeToString(b)
	switch numFmt {
	case "bin":
		rnd = format.EncodeToBin(b)
	case "dec":
		num, _ := strconv.ParseUint(rnd, 16, 64)
		rnd = fmt.Sprintf("%d", num)
	}
	return rnd
}

func init() {
	encCmd.AddCommand(NewRndCmd())
}
