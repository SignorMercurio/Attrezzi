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
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	bits        int
	privKeyPath string
	pubKeyPath  string
)

// NewRkgCmd represents the rkg command
func NewRkgCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rkg",
		Short: "RSA key generation",
		Long: `RSA key generation
Example:
	att enc -o pubkey.pub rkg -p priv.key -b 4096`,
		RunE: func(cmd *cobra.Command, args []string) error {
			priv, pub := genRSAKeyPair()
			err := exportPrivKey(priv)
			if err != nil {
				return err
			}
			err = exportPubKey(pub)
			if err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().IntVarP(&bits, "bits", "b", 2048, "Key bits")
	cmd.Flags().StringVar(&privKeyPath, "priv", "./priv.key", "Path to store private key")
	cmd.Flags().StringVar(&pubKeyPath, "pub", "./priv.key", "Path to store private key")

	return cmd
}

func genRSAKeyPair() (*rsa.PrivateKey, *rsa.PublicKey) {
	privateKey, _ := rsa.GenerateKey(rand.Reader, bits)
	return privateKey, &privateKey.PublicKey
}

func exportPrivKey(key *rsa.PrivateKey) error {
	keyBytes := x509.MarshalPKCS1PrivateKey(key)
	keyPem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: keyBytes,
	})
	privKeyOut, err := os.OpenFile(privKeyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return errors.Wrap(err, "open privkey output file")
	}
	defer privKeyOut.Close()

	privKeyOut.WriteString(string(keyPem))
	return nil
}

func exportPubKey(key *rsa.PublicKey) error {
	keyBytes := x509.MarshalPKCS1PublicKey(key)
	keyPem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: keyBytes,
	})
	pubKeyOut, err := os.OpenFile(pubKeyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return errors.Wrap(err, "open pubkey output file")
	}
	defer pubKeyOut.Close()

	pubKeyOut.WriteString(string(keyPem))
	return nil
}

func init() {
	encCmd.AddCommand(NewRkgCmd())
}
