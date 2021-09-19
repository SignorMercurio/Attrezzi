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
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	mode string
)

// NewAesCmd represents the aes command
func NewAesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "aes",
		Short: "AES encryption / decryption",
		Long: `AES encryption / decryption
Example:
	echo -n "hello" | att enc -o out.txt aes -e
	att enc -i in.txt aes -d`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var plainText, cipherText []byte
			var err error
			k, err := hex.DecodeString(key)
			if err != nil {
				return errors.Wrap(err, "parse AES key")
			}

			if enc {
				switch mode {
				case "cbc":
					cipherText, err = aesEncryptCBC(inputBytes, k)
				case "cfb":
					cipherText, err = aesEncryptCFB(inputBytes, k)
				case "ofb":
					cipherText, err = aesEncryptOFB(inputBytes, k)
				case "ctr":
					cipherText, err = aesEncryptCTR(inputBytes, k)
				default:
					cipherText, err = aesEncryptGCM(inputBytes, k)
				}
				if err != nil {
					return err
				}
				Echo(string(cipherText))
			} else if dec {
				switch mode {
				case "cbc":
					plainText, err = aesDecryptCBC(inputBytes, k)
				case "cfb":
					plainText, err = aesDecryptCFB(inputBytes, k)
				case "ofb":
					plainText, err = aesDecryptOFB(inputBytes, k)
				case "ctr":
					plainText, err = aesDecryptCTR(inputBytes, k)
				default:
					plainText, err = aesDecryptGCM(inputBytes, k)
				}
				if err != nil {
					return err
				}
				Echo(string(plainText))
			} else {
				NoActionSpecified()
			}
			return nil
		},
	}
	cmd.Flags().BoolVarP(&enc, "encrypt", "e", false, "AES encryption")
	cmd.Flags().BoolVarP(&dec, "decrypt", "d", false, "AES decryption")
	cmd.Flags().StringVarP(&key, "key", "k", "", "Encryption key in hex format, either 16 / 24 / 32 bytes to select AES-128 / AES-192 (GCM mode not supported) / AES-256")
	cmd.Flags().StringVarP(&mode, "mode", "m", "gcm", "Block mode to use: cbc / cfb / ofb / ctr / gcm")

	return cmd
}

// PKCS5Padding pads [plain] to [blockSize]
func PKCS5Padding(plain []byte, blockSize int) []byte {
	padding := blockSize - len(plain)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(plain, padtext...)
}

// PKCS5Unpadding unpads [padded] to the original plaintext
func PKCS5Unpadding(padded []byte) []byte {
	length := len(padded)
	toUnpad := int(padded[length-1])
	return padded[:(length - toUnpad)]
}

func newAES(key []byte) (cipher.Block, int, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, 0, errors.Wrap(err, "parse AES key")
	}
	return block, block.BlockSize(), nil
}

func genIV(blockSize int, plainLen int) ([]byte, []byte, error) {
	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	cipherText := make([]byte, blockSize+plainLen)
	iv := cipherText[:blockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, nil, errors.Wrap(err, "generate random IV")
	}

	return cipherText, iv, nil
}

func genNonce(nonceSize int) ([]byte, error) {
	nonce := make([]byte, nonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, errors.Wrap(err, "generate random nonce")
	}
	return nonce, nil
}

func validateCiphertext(cipherText []byte, size int) error {
	if len(cipherText) < size {
		return errors.New("validate cipherText")
	}
	return nil
}

// aesEncryptCBC encrypts [plainText] with [key] using CBC mode
func aesEncryptCBC(plainText, key []byte) ([]byte, error) {
	block, blockSize, err := newAES(key)
	if err != nil {
		return nil, err
	}
	plainText = PKCS5Padding(plainText, blockSize)

	cipherText, iv, err := genIV(blockSize, len(plainText))
	if err != nil {
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText[blockSize:], plainText)

	return cipherText, nil
}

// aesDecryptCBC decrypts [cipherText] with [key] using CBC mode
func aesDecryptCBC(cipherText, key []byte) ([]byte, error) {
	block, blockSize, err := newAES(key)
	if err != nil {
		return nil, err
	}

	if err := validateCiphertext(cipherText, blockSize); err != nil {
		return nil, err
	}
	iv, cipherText := cipherText[:blockSize], cipherText[blockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipherText, cipherText)
	plainText := PKCS5Unpadding(cipherText)

	return plainText, nil
}

// aesEncryptCFB encrypts [plainText] with [key] using CFB mode
func aesEncryptCFB(plainText, key []byte) ([]byte, error) {
	block, blockSize, err := newAES(key)
	if err != nil {
		return nil, err
	}

	cipherText, iv, err := genIV(blockSize, len(plainText))
	if err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[blockSize:], plainText)

	return cipherText, nil
}

// aesDecryptCFB decrypts [cipherText] with [key] using CFB mode
func aesDecryptCFB(cipherText, key []byte) ([]byte, error) {
	block, blockSize, err := newAES(key)
	if err != nil {
		return nil, err
	}

	if err := validateCiphertext(cipherText, blockSize); err != nil {
		return nil, err
	}
	iv, cipherText := cipherText[:blockSize], cipherText[blockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return cipherText, nil
}

// aesEncryptOFB encrypts [plainText] with [key] using OFB mode
func aesEncryptOFB(plainText, key []byte) ([]byte, error) {
	block, blockSize, err := newAES(key)
	if err != nil {
		return nil, err
	}

	cipherText, iv, err := genIV(blockSize, len(plainText))
	if err != nil {
		return nil, err
	}

	stream := cipher.NewOFB(block, iv)
	stream.XORKeyStream(cipherText[blockSize:], plainText)

	return cipherText, nil
}

// aesDecryptOFB decrypts [cipherText] with [key] using OFB mode
func aesDecryptOFB(cipherText, key []byte) ([]byte, error) {
	block, blockSize, err := newAES(key)
	if err != nil {
		return nil, err
	}

	if err := validateCiphertext(cipherText, blockSize); err != nil {
		return nil, err
	}
	iv, cipherText := cipherText[:blockSize], cipherText[blockSize:]

	stream := cipher.NewOFB(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return cipherText, nil
}

// aesEncryptCTR encrypts [plainText] with [key] using CTR mode
func aesEncryptCTR(plainText, key []byte) ([]byte, error) {
	block, blockSize, err := newAES(key)
	if err != nil {
		return nil, err
	}

	cipherText, iv, err := genIV(blockSize, len(plainText))
	if err != nil {
		return nil, err
	}

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(cipherText[blockSize:], plainText)

	return cipherText, nil
}

// aesDecryptCTR decrypts [cipherText] with [key] using CTR mode
func aesDecryptCTR(cipherText, key []byte) ([]byte, error) {
	block, blockSize, err := newAES(key)
	if err != nil {
		return nil, err
	}

	if err := validateCiphertext(cipherText, blockSize); err != nil {
		return nil, err
	}
	iv, cipherText := cipherText[:blockSize], cipherText[blockSize:]

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return cipherText, nil
}

// aesEncryptGCM encrypts [plainText] with [key] using GCM mode
func aesEncryptGCM(plainText, key []byte) ([]byte, error) {
	block, _, err := newAES(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.Wrap(err, "generate new GCM")
	}

	nonce, err := genNonce(gcm.NonceSize())
	if err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plainText, nil), nil
}

// aesDecryptGCM decrypts [cipherText] with [key] using GCM mode
func aesDecryptGCM(cipherText, key []byte) ([]byte, error) {
	block, _, err := newAES(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.Wrap(err, "generate new GCM")
	}

	nonceSize := gcm.NonceSize()
	if err := validateCiphertext(cipherText, nonceSize); err != nil {
		return nil, err
	}
	nonce, cipherText := cipherText[:nonceSize], cipherText[nonceSize:]

	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return nil, errors.Wrap(err, "decrypt AES-GCM cipherText")
	}
	return plainText, nil
}

func init() {
	encCmd.AddCommand(NewAesCmd())
}
