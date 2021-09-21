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
	"crypto/ecdsa"
	"crypto/elliptic"
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
	alg         string
	privKeyPath string
	pubKeyPath  string
)

// NewAkgCmd represents the akg command
func NewAkgCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "akg",
		Short: "Asymmetric encryption key generation",
		Long: `Asymmetric encryption key generation
Example:
	att enc -o pubkey.pub akg -p priv.key -b 4096`,
		RunE: func(cmd *cobra.Command, args []string) error {
			priv, pub := genKeyPair()
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
	cmd.Flags().IntVarP(&bits, "bits", "b", 2048, "Key bits, or Curve name in ECDSA: 224 / 256 / 384 / 521")
	cmd.Flags().StringVarP(&alg, "algorithm", "a", "rsa", "Encryption algorithm to use: rsa / ecdsa")
	cmd.Flags().StringVar(&privKeyPath, "priv", "./priv.pem", "Path to store private key")
	cmd.Flags().StringVar(&pubKeyPath, "pub", "./pub.pem", "Path to store public key")

	return cmd
}

func getCurve() elliptic.Curve {
	switch bits {
	case 224:
		return elliptic.P224()
	case 384:
		return elliptic.P384()
	case 521:
		return elliptic.P521()
	default:
		return elliptic.P256()
	}
}

// genKeyPair returns a Keypair
func genKeyPair() (interface{}, interface{}) {
	switch alg {
	case "ecdsa":
		privateKey, _ := ecdsa.GenerateKey(getCurve(), rand.Reader)
		return privateKey, &privateKey.PublicKey
	default:
		privateKey, _ := rsa.GenerateKey(rand.Reader, bits)
		return privateKey, &privateKey.PublicKey
	}
}

// exportPrivKey writes the private key to a file
func exportPrivKey(key interface{}) error {
	keyBytes, _ := x509.MarshalPKCS8PrivateKey(key)
	keyPem := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
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

// exportPubKey writes the public key to a file
func exportPubKey(key interface{}) error {
	keyBytes, _ := x509.MarshalPKIXPublicKey(key)
	keyPem := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
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
	encCmd.AddCommand(NewAkgCmd())
}
