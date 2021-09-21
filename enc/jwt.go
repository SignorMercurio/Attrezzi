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
	"encoding/json"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	method string
)

// NewJwtCmd represents the jwt command
func NewJwtCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "jwt",
		Short: "JWT-related operation",
		Long: `JWT-related operation
Example:
	att enc -i in.txt -o out.txt jwt -e
	att enc -i out.txt jwt -d`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if enc {
				token, err := sign()
				if err != nil {
					return err
				}
				Echo(token)
			} else if dec {
				token, err := verify(string(inputBytes))
				if err != nil {
					return err
				}
				Echo(token)
			} else {
				NoActionSpecified()
			}
			return nil
		},
	}
	cmd.Flags().BoolVarP(&enc, "sign", "s", false, "JWT sign")
	cmd.Flags().BoolVarP(&dec, "verify", "v", false, "JWT verify")
	cmd.Flags().StringVarP(&key, "key", "k", "", "File storing the JWT secret key")
	cmd.Flags().StringVarP(&method, "method", "m", "hs256", "JWT signing method: hs256 / hs384 / hs512 / rs256 / rs384 / rs512 / es256 / es384 / es512 / ps256 / ps384 / ps512")

	return cmd
}

// getSigningMethod gets the JWT signing method from user input
func getSigningMethod() jwt.SigningMethod {
	switch method {
	case "hs384":
		return jwt.SigningMethodHS384
	case "hs512":
		return jwt.SigningMethodHS512
	case "rs256":
		return jwt.SigningMethodRS256
	case "rs384":
		return jwt.SigningMethodRS384
	case "rs512":
		return jwt.SigningMethodRS512
	case "es256":
		return jwt.SigningMethodES256
	case "es384":
		return jwt.SigningMethodES384
	case "es512":
		return jwt.SigningMethodES512
	case "ps256":
		return jwt.SigningMethodPS256
	case "ps384":
		return jwt.SigningMethodPS384
	case "ps512":
		return jwt.SigningMethodPS512
	default:
		return jwt.SigningMethodHS256
	}
}

// readKey reads the JWT secret key from a file, parsing it from PEM if necessary
func readKey() (interface{}, error) {
	keyByte, err := os.ReadFile(key)
	if err != nil {
		return nil, errors.Wrap(err, "read key file")
	}

	switch method {
	case "rs256", "rs384", "rs512", "ps256", "ps384", "ps512":
		if enc {
			rsaPrivKey, err := jwt.ParseRSAPrivateKeyFromPEM(keyByte)
			if err != nil {
				return nil, errors.Wrap(err, "parse RSA private key")
			}
			return rsaPrivKey, nil
		} else {
			rsaPubKey, err := jwt.ParseRSAPublicKeyFromPEM(keyByte)
			if err != nil {
				return nil, errors.Wrap(err, "parse RSA public key")
			}
			return rsaPubKey, nil
		}
	case "es256", "es384", "es512":
		if enc {
			ecPrivKey, err := jwt.ParseECPrivateKeyFromPEM(keyByte)
			if err != nil {
				return nil, errors.Wrap(err, "parse EC private key")
			}
			return ecPrivKey, nil
		} else {
			ecPubKey, err := jwt.ParseECPublicKeyFromPEM(keyByte)
			if err != nil {
				return nil, errors.Wrap(err, "parse EC public key")
			}
			return ecPubKey, nil
		}
	default:
		return keyByte, nil
	}
}

// sign performs the JWT sign operation
func sign() (string, error) {
	var v = &jwt.MapClaims{}
	err := json.Unmarshal(inputBytes, v)
	if err != nil {
		return "", errors.Wrap(err, "unmarshal json")
	}

	token := jwt.NewWithClaims(getSigningMethod(), v)

	key, err := readKey()
	if err != nil {
		return "", err
	}

	return token.SignedString(key)
}

// verify performs the JWT verify operation
func verify(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if strings.ToLower(t.Method.Alg()) != method {
			return nil, errors.New("Invalid signing alg")
		}
		return readKey()
	})

	if token != nil {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			jsonContent, _ := json.MarshalIndent(claims, "", "  ")
			return string(jsonContent), nil
		}
	}

	return "", errors.Wrap(err, "validate the token")
}

func init() {
	encCmd.AddCommand(NewJwtCmd())
}
