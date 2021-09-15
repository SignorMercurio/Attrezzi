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
	"bytes"
	"encoding/hex"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	delim        string
	delim_prefix bool
)

// hexCmd represents the hex command
var hexCmd = &cobra.Command{
	Use:   "hex",
	Short: "string to hex",
	RunE: func(cmd *cobra.Command, args []string) error {
		var delimiter []byte
		switch delim {
		case `\n`:
			delimiter = []byte("\n")
		case `\r\n`:
			delimiter = []byte("\r\n")
		default:
			delimiter = []byte(delim)
		}

		if encode {
			encoded := hex.EncodeToString(inputBytes)
			Echo(insertInto(encoded, 2, delimiter))
		} else if decode {
			arr := strings.Split(string(inputBytes), string(delimiter))
			if arr[0] == "" {
				arr = arr[1:]
			}
			toDecode := strings.Join(arr, "")
			decoded, err := hex.DecodeString(toDecode)
			if err != nil {
				return errors.Wrap(err, "decode hex")
			}
			Echo(string(decoded))
		} else {
			log.Error("Please specify -e or -d")
		}
		return nil
	},
}

func insertInto(s string, interval int, delimiter []byte) string {
	var buffer bytes.Buffer
	before := interval - 1
	last := len(s) - 1

	if delim_prefix {
		buffer.Write(delimiter)
	}

	for i, char := range s {
		buffer.WriteRune(char)
		if i%interval == before && i != last {
			buffer.Write(delimiter)
		}
	}

	return buffer.String()
}

func init() {
	fmtCmd.AddCommand(hexCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// hexCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	hexCmd.Flags().BoolVarP(&encode, "encode", "e", false, "Encode to hex")
	hexCmd.Flags().BoolVarP(&decode, "decode", "d", false, "Decode from hex")
	hexCmd.Flags().StringVar(&delim, "delim", "", "Delimiter")
	hexCmd.Flags().BoolVarP(&delim_prefix, "prefix", "p", false, "Whether the delimiter is a prefix")
}
