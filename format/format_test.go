package format

import (
	"io"
	"testing"

	"github.com/SignorMercurio/attrezzi/cmd"
	"github.com/SignorMercurio/attrezzi/test"
)

var (
	in  = "./testdata/in.txt"
	out = "./testdata/out.txt"
	src = "Hello 世界 123"
)

func exec(args ...string) {
	rootCmd := cmd.NewRootCmd()
	fmtCmd := NewFmtCmd()
	fmtCmd.AddCommand(
		NewB64Cmd(),
		NewHexCmd(),
		NewBinCmd(),
		NewDecCmd(),
		NewUrlCmd(),
		NewHtmCmd(),
		NewUniCmd(),
		NewB32Cmd(),
		NewB58Cmd(),
		NewBsxCmd(),
		NewB85Cmd(),
	)
	rootCmd.AddCommand(fmtCmd)

	commonArgs := []string{"fmt", "-o", out, "-i"}
	rootCmd.SetArgs(append(commonArgs, args...))
	rootCmd.Execute()
}

func TestFormat(t *testing.T) {
	rootCmd := cmd.NewRootCmd()
	fmtCmd := NewFmtCmd()
	fmtCmd.AddCommand(NewB64Cmd())
	rootCmd.AddCommand(fmtCmd)

	cmd.Log.SetOutput(io.Discard)
	input = nil
	tests := []test.Test{
		// open output fail
		{Cmd: []string{"fmt", "-o", "bla/blabla.txt", "-i", in, "b64", "-e"}, Dst: ""},
		// open input fail
		{Cmd: []string{"fmt", "-i", "bla/blabla.txt", "-o", out, "b64", "-e"}, Dst: ""},
		// read input fail
		{Cmd: []string{"fmt", "-i", "", "-o", out, "b64", "-e"}, Dst: ""},
	}

	for _, tst := range tests {
		rootCmd.SetArgs(tst.Cmd)
		rootCmd.Execute()
		test.CheckResult(out, tst.Dst, t)
	}
}

func TestB64(t *testing.T) {
	tests := []test.Test{
		// std
		{Cmd: []string{in, "b64", "-e"}, Dst: "SGVsbG8g5LiW55WMIDEyMw=="},
		{Cmd: []string{out, "b64", "-d"}, Dst: src},
		// url
		{Cmd: []string{in, "b64", "-e", "-a", "url"}, Dst: "SGVsbG8g5LiW55WMIDEyMw=="},
		{Cmd: []string{out, "b64", "-d", "-a", "url"}, Dst: src},
		// custom alphabet
		{Cmd: []string{in, "b64", "-e", "-a", "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz+/"}, Dst: "I6LiR6yWvBYMvvMC834oCm=="},
		{Cmd: []string{out, "b64", "-d", "-a", "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz+/"}, Dst: src},
		// no padding
		{Cmd: []string{in, "b64", "-e", "-p", ""}, Dst: "SGVsbG8g5LiW55WMIDEyMw"},
		{Cmd: []string{out, "b64", "-d", "-p", ""}, Dst: src},
		// custom padding
		{Cmd: []string{in, "b64", "-e", "-p", "?"}, Dst: "SGVsbG8g5LiW55WMIDEyMw??"},
		{Cmd: []string{out, "b64", "-d", "-p", "?"}, Dst: src},
		// decode fail
		{Cmd: []string{in, "b64", "-d"}, Dst: ""},
		// no action
		{Cmd: []string{in, "b64"}, Dst: ""},
	}
	for _, tst := range tests {
		exec(tst.Cmd...)
		test.CheckResult(out, tst.Dst, t)
	}
}

func TestB32(t *testing.T) {
	tests := []test.Test{
		// std
		{Cmd: []string{in, "b32", "-e"}, Dst: "JBSWY3DPEDSLRFXHSWGCAMJSGM======"},
		{Cmd: []string{out, "b32", "-d"}, Dst: src},
		// hex
		{Cmd: []string{in, "b32", "-e", "-a", "hex"}, Dst: "91IMOR3F43IBH5N7IM620C9I6C======"},
		{Cmd: []string{out, "b32", "-d", "-a", "hex"}, Dst: src},
		// custom alphabet
		{Cmd: []string{in, "b32", "-e", "-a", "abcdefghijklmnopqrstuvwxyz234567"}, Dst: "jbswy3dpedslrfxhswgcamjsgm======"},
		{Cmd: []string{out, "b32", "-d", "-a", "abcdefghijklmnopqrstuvwxyz234567"}, Dst: src},
		// no padding
		{Cmd: []string{in, "b32", "-e", "-p", ""}, Dst: "JBSWY3DPEDSLRFXHSWGCAMJSGM"},
		{Cmd: []string{out, "b32", "-d", "-p", ""}, Dst: src},
		// custom padding
		{Cmd: []string{in, "b32", "-e", "-p", "!"}, Dst: "JBSWY3DPEDSLRFXHSWGCAMJSGM!!!!!!"},
		{Cmd: []string{out, "b32", "-d", "-p", "!"}, Dst: src},
		// decode fail
		{Cmd: []string{in, "b32", "-d"}, Dst: ""},
		// no action
		{Cmd: []string{in, "b32"}, Dst: ""},
	}
	for _, tst := range tests {
		exec(tst.Cmd...)
		test.CheckResult(out, tst.Dst, t)
	}
}

func TestB58(t *testing.T) {
	tests := []test.Test{
		// btc
		{Cmd: []string{in, "b58", "-e"}, Dst: "9wWTEnNTcvgeNTGbfmax8z"},
		{Cmd: []string{out, "b58", "-d"}, Dst: src},
		// flickr
		{Cmd: []string{in, "b58", "-e", "-a", "flickr"}, Dst: "9WvseMnsBVFDnsgAELzX8Z"},
		{Cmd: []string{out, "b58", "-d", "-a", "flickr"}, Dst: src},
		// custom alphabet
		{Cmd: []string{in, "b58", "-e", "-a", "rpshnaf39wBUDNEGHJKLM4PQRST7VWXYZ2bcdeCg65jkm8oFqi1tuvAxyz"}, Dst: "9AWTN84Tcvge4TGbCm2x3z"},
		{Cmd: []string{out, "b58", "-d", "-a", "rpshnaf39wBUDNEGHJKLM4PQRST7VWXYZ2bcdeCg65jkm8oFqi1tuvAxyz"}, Dst: src},
		// decode fail
		{Cmd: []string{in, "b58", "-d"}, Dst: ""},
		// no action
		{Cmd: []string{in, "b58"}, Dst: ""},
	}

	for _, tst := range tests {
		exec(tst.Cmd...)
		test.CheckResult(out, tst.Dst, t)
	}
}

func TestB85(t *testing.T) {
	var tests = []test.Test{
		// std
		{Cmd: []string{in, "b85", "-e"}, Dst: "87cURD]n,NQKONl+>GW-"},
		{Cmd: []string{out, "b85", "-d"}, Dst: src},
		// decode fail
		{Cmd: []string{in, "b85", "-d"}, Dst: ""},
		// no action
		{Cmd: []string{in, "b85"}, Dst: ""},
	}

	for _, tst := range tests {
		exec(tst.Cmd...)
		test.CheckResult(out, tst.Dst, t)
	}
}

func TestBsx62(t *testing.T) {
	var tests = []test.Test{
		// base 62
		{Cmd: []string{in, "bsx", "-e"}, Dst: "2CbnUNVhpxZqW7mkcOp2Ml"},
		{Cmd: []string{out, "bsx", "-d"}, Dst: src},
		// base 16
		{Cmd: []string{in, "bsx", "-e", "-b", "16"}, Dst: "48656C6C6F20E4B896E7958C20313233"},
		{Cmd: []string{out, "bsx", "-d", "-b", "16"}, Dst: src},
		// base 16 with alphabet only
		{Cmd: []string{in, "bsx", "-e", "-a", "0123456789abcdef"}, Dst: "48656c6c6f20e4b896e7958c20313233"},
		{Cmd: []string{out, "bsx", "-d", "-a", "0123456789abcdef"}, Dst: src},
		// invalid alphabet
		{Cmd: []string{in, "bsx", "-e", "-a", "00"}, Dst: ""},
		// decode fail
		{Cmd: []string{in, "bsx", "-d"}, Dst: ""},
		// no action
		{Cmd: []string{in, "bsx"}, Dst: ""},
	}

	for _, tst := range tests {
		exec(tst.Cmd...)
		test.CheckResult(out, tst.Dst, t)
	}
}

func TestHex(t *testing.T) {
	var tests = []test.Test{
		// no delim
		{Cmd: []string{in, "hex", "-e"}, Dst: "48656c6c6f20e4b896e7958c20313233"},
		{Cmd: []string{out, "hex", "-d"}, Dst: src},
		// 0x
		{Cmd: []string{in, "hex", "-e", "--delim", "0x", "-p"}, Dst: "0x480x650x6c0x6c0x6f0x200xe40xb80x960xe70x950x8c0x200x310x320x33"},
		{Cmd: []string{out, "hex", "-d", "--delim", "0x", "-p"}, Dst: src},
		// \r\n
		{Cmd: []string{in, "hex", "-e", "--delim", `\r\n`}, Dst: "48\r\n65\r\n6c\r\n6c\r\n6f\r\n20\r\ne4\r\nb8\r\n96\r\ne7\r\n95\r\n8c\r\n20\r\n31\r\n32\r\n33"},
		{Cmd: []string{out, "hex", "-d", "--delim", `\r\n`}, Dst: src},
		// decode fail
		{Cmd: []string{in, "hex", "-d"}, Dst: ""},
		// no action
		{Cmd: []string{in, "hex"}, Dst: ""},
	}

	for _, tst := range tests {
		exec(tst.Cmd...)
		test.CheckResult(out, tst.Dst, t)
	}
}

func TestBin(t *testing.T) {
	var tests = []test.Test{
		// empty delim
		{Cmd: []string{in, "bin", "-e"}, Dst: "01001000 01100101 01101100 01101100 01101111 00100000 11100100 10111000 10010110 11100111 10010101 10001100 00100000 00110001 00110010 00110011"},
		{Cmd: []string{out, "bin", "-d"}, Dst: src},
		// decode fail
		{Cmd: []string{in, "bin", "-d"}, Dst: ""},
		// no action
		{Cmd: []string{in, "bin"}, Dst: ""},
	}

	for _, tst := range tests {
		exec(tst.Cmd...)
		test.CheckResult(out, tst.Dst, t)
	}
}

func TestDec(t *testing.T) {
	var tests = []test.Test{
		// \n delim as a prefix
		{Cmd: []string{in, "dec", "-e", "--delim", `\n`, "-p"}, Dst: "\n72\n101\n108\n108\n111\n32\n228\n184\n150\n231\n149\n140\n32\n49\n50\n51"},
		{Cmd: []string{out, "dec", "-d", "--delim", `\n`, "-p"}, Dst: src},
		// decode fail
		{Cmd: []string{in, "dec", "-d"}, Dst: ""},
		// no action
		{Cmd: []string{in, "dec"}, Dst: ""},
	}

	for _, tst := range tests {
		exec(tst.Cmd...)
		test.CheckResult(out, tst.Dst, t)
	}
}

func TestURL(t *testing.T) {
	in := "./testdata/in_url.txt"
	in_fail := "./testdata/in_url_fail.txt"
	src := "https://www.example.com/a/b/?c=d&e=f#g中文"

	var tests = []test.Test{
		// normal
		{Cmd: []string{in, "url", "-e"}, Dst: `https://www.example.com/a/b/?c=d&e=f#g%E4%B8%AD%E6%96%87`},
		{Cmd: []string{out, "url", "-d"}, Dst: src},
		// all
		{Cmd: []string{in, "url", "-ea"}, Dst: `https%3A%2F%2Fwww.example.com%2Fa%2Fb%2F%3Fc%3Dd%26e%3Df%23g%E4%B8%AD%E6%96%87`},
		{Cmd: []string{out, "url", "-d"}, Dst: src},
		// parse URL fail
		{Cmd: []string{in_fail, "url", "-e"}, Dst: ""},
		// decode fail
		{Cmd: []string{in_fail, "url", "-d"}, Dst: ""},
		// no action
		{Cmd: []string{in, "url"}, Dst: ""},
	}

	for _, tst := range tests {
		exec(tst.Cmd...)
		test.CheckResult(out, tst.Dst, t)
	}
}

func TestHTML(t *testing.T) {
	in := "./testdata/in_html.txt"
	src := "<script>alert('xss');</script>"

	var tests = []test.Test{
		// normal
		{Cmd: []string{in, "htm", "-e"}, Dst: `&lt;script&gt;alert(&#39;xss&#39;);&lt;/script&gt;`},
		{Cmd: []string{out, "htm", "-d"}, Dst: src},
		// no action
		{Cmd: []string{in, "htm"}, Dst: ""},
	}

	for _, tst := range tests {
		exec(tst.Cmd...)
		test.CheckResult(out, tst.Dst, t)
	}
}

func TestUni(t *testing.T) {
	var tests = []test.Test{
		// normal
		{Cmd: []string{in, "uni", "-e"}, Dst: `Hello \u4e16\u754c 123`},
		{Cmd: []string{out, "uni", "-d"}, Dst: src},
		// no action
		{Cmd: []string{in, "uni"}, Dst: ""},
	}

	for _, tst := range tests {
		exec(tst.Cmd...)
		test.CheckResult(out, tst.Dst, t)
	}
}
