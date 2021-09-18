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
	"strings"

	"github.com/spf13/cobra"
)

var (
	dash        string
	dot         string
	letterDelim string
	wordDelim   string
)

// NewMorCmd represents the mor command
func NewMorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mor",
		Short: "Morse code transformation",
		Long: `Morse code transformation
Example:
	echo -n "hello" | att enc -o out.txt mor -e
	att enc -i in.txt mor -d --dash "DASH" --dot "DOT" -l "/" -w "\n"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			letter := getDelimiter(letterDelim)
			word := getDelimiter(wordDelim)

			if enc {
				enced := morEncode(string(inputBytes), string(letter), string(word))
				Echo(enced)
			} else if dec {
				deced := morDecode(string(inputBytes), string(letter), string(word))
				Echo(deced)
			} else {
				NoActionSpecified()
			}
			return nil
		},
	}
	cmd.Flags().BoolVarP(&enc, "encode", "e", false, "Encode to morse code")
	cmd.Flags().BoolVarP(&dec, "decode", "d", false, "Decode from morse code")
	cmd.Flags().StringVar(&dash, "dash", "-", "Dash")
	cmd.Flags().StringVar(&dot, "dot", ".", "Dot")
	cmd.Flags().StringVarP(&letterDelim, "letter-delim", "l", " ", "Letter delimiter")
	cmd.Flags().StringVarP(&wordDelim, "word-delim", "w", "\n", "Word delimiter")

	return cmd
}

var alpha2mor = map[rune]string{
	'A':  ".-",
	'B':  "-...",
	'C':  "-.-.",
	'D':  "-..",
	'E':  ".",
	'F':  "..-.",
	'G':  "--.",
	'H':  "....",
	'I':  "..",
	'J':  ".---",
	'K':  "-.-",
	'L':  ".-..",
	'M':  "--",
	'N':  "-.",
	'O':  "---",
	'P':  ".--.",
	'Q':  "--.-",
	'R':  ".-.",
	'S':  "...",
	'T':  "-",
	'U':  "..-",
	'V':  "...-",
	'W':  ".--",
	'X':  "-..-",
	'Y':  "-.--",
	'Z':  "--..",
	'1':  ".----",
	'2':  "..---",
	'3':  "...--",
	'4':  "....-",
	'5':  ".....",
	'6':  "-....",
	'7':  "--...",
	'8':  "---..",
	'9':  "----.",
	'0':  "-----",
	'.':  ".-.-.-",  // period
	':':  "---...",  // colon
	',':  "--..--",  // comma
	';':  "-.-.-",   // semicolon
	'?':  "..--..",  // question
	'=':  "-...-",   // equals
	'\'': ".----.",  // apostrophe
	'/':  "-..-.",   // slash
	'!':  "-.-.--",  // exclamation
	'-':  "-....-",  // dash
	'_':  "..--.-",  // underline
	'"':  ".-..-.",  // quotation marks
	'(':  "-.--.",   // parenthesis (open)
	')':  "-.--.-",  // parenthesis (close)
	'$':  "...-..-", // dollar
	'&':  ".-...",   // ampersand
	'@':  ".--.-.",  // at
	'+':  ".-.-.",   // plus
}

var mor2alpha = map[string]rune{}

func morEncode(src string, letter string, word string) string {
	res := ""
	src = strings.ToUpper(src)
	for i, v := range src {
		if v == ' ' {
			res += word
		} else {
			morse := strings.ReplaceAll(strings.ReplaceAll(alpha2mor[v], "-", dash), ".", dot)
			res += morse
			if i != len(src)-1 && src[i+1] != ' ' {
				res += letter
			}
		}
	}

	return res
}

func morDecode(src string, letter string, word string) string {
	res := []rune{}
	src = strings.ReplaceAll(strings.ReplaceAll(src, dash, "-"), dot, ".")
	words := strings.Split(src, word)
	for i, wd := range words {
		letters := strings.Split(wd, letter)
		for _, lt := range letters {
			res = append(res, mor2alpha[lt])
		}
		if i != len(words)-1 {
			res = append(res, ' ')
		}
	}

	return string(res)
}

// getDelimiter gets the delimiter from user input, mainly dealing with LF & CRLF
func getDelimiter(delim string) []byte {
	var delimiter []byte

	switch delim {
	case `\n`:
		delimiter = []byte("\n")
	case `\r\n`:
		delimiter = []byte("\r\n")
	default:
		delimiter = []byte(delim)
	}

	return delimiter
}

func init() {
	encCmd.AddCommand(NewMorCmd())

	for k, v := range alpha2mor {
		mor2alpha[v] = k
	}
}
