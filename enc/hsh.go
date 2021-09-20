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
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"

	"github.com/spf13/cobra"
)

// NewHshCmd represents the hsh command
func NewHshCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hsh",
		Short: "Hash function calculation",
		Long: `Hash function calculation
Example:
	echo -n "hello" | att enc -o out.txt hsh --hash sha512`,
		RunE: func(cmd *cobra.Command, args []string) error {
			hash := getHash()
			hash.Write(inputBytes)
			Echo(fmt.Sprintf("%x", hash.Sum(nil)))
			return nil
		},
	}
	cmd.Flags().StringVar(&hashFunc, "hash", "sha256", "Hash function")

	return cmd
}

// getHash chooses the hash function according to the user input
func getHash() hash.Hash {
	switch hashFunc {
	case "md5":
		return md5.New()
	case "sha1":
		return sha1.New()
	case "sha384":
		return sha512.New384()
	case "sha512":
		return sha512.New()
	default:
		return sha256.New()
	}
}

func init() {
	encCmd.AddCommand(NewHshCmd())
}
