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
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	rsaMode string
)

// NewRsaCmd represents the rsa command
func NewRsaCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rsa",
		Short: "RSA encryption / decryption",
		Long: `RSA encryption / decryption
Example:
	att enc -i in.txt -o out.txt rsa -p priv.key -e`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if enc {
				var enced []byte

				pub, err := importPubKey(pubKeyPath)
				if err != nil {
					return err
				}

				switch rsaMode {
				case "pkcs1v15":
					enced, _ = rsa.EncryptPKCS1v15(rand.Reader, pub, inputBytes)
				default:
					enced, _ = rsa.EncryptOAEP(sha256.New(), rand.Reader, pub, inputBytes, nil)
				}
				Echo(string(enced))
			} else if dec {
				var deced []byte

				priv, err := importPrivKey(privKeyPath)
				if err != nil {
					return err
				}

				switch rsaMode {
				case "pkcs1v15":
					deced, err = rsa.DecryptPKCS1v15(rand.Reader, priv, inputBytes)
				default:
					deced, err = rsa.DecryptOAEP(sha256.New(), rand.Reader, priv, inputBytes, nil)
				}
				Echo(string(deced))
			} else {
				NoActionSpecified()
			}

			return nil
		},
	}
	cmd.Flags().StringVar(&privKeyPath, "priv", "./priv.key", "Path of existing private key")
	cmd.Flags().StringVar(&pubKeyPath, "pub", "./pubkey.pub", "Path of existing public key")
	cmd.Flags().StringVarP(&rsaMode, "mode", "m", "oaep", "RSA encryption mode: oaep / pkcs1v15")
	cmd.Flags().BoolVarP(&enc, "encrypt", "e", false, "RSA encryption")
	cmd.Flags().BoolVarP(&dec, "decrypt", "d", false, "RSA decryption")

	return cmd
}

func importPrivKey(filename string) (*rsa.PrivateKey, error) {
	privKeyOut, err := os.Open(filename)
	defer privKeyOut.Close()
	if err != nil {
		return nil, errors.Wrap(err, "open privkey input file")
	}

	keyBytes, _ := ioutil.ReadAll(privKeyOut)

	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, errors.New("decode pubkey PEM")
	}

	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, errors.Wrap(err, "parse privkey")
	}

	return privKey, nil
}

func importPubKey(filename string) (*rsa.PublicKey, error) {
	pubKeyOut, err := os.Open(filename)
	defer pubKeyOut.Close()
	if err != nil {
		return nil, errors.Wrap(err, "open pubkey input file")
	}

	keyBytes, _ := ioutil.ReadAll(pubKeyOut)

	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, errors.New("decode pubkey PEM")
	}

	pubKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, errors.Wrap(err, "parse pubkey")
	}

	return pubKey, nil
}

func init() {
	encCmd.AddCommand(NewRsaCmd())
}
