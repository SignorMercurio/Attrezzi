package format

import (
	"io/ioutil"
	"testing"

	"github.com/SignorMercurio/attrezzi/cmd"
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
	)
	rootCmd.AddCommand(fmtCmd)

	rootCmd.SetArgs(args)
	rootCmd.Execute()
}

func checkResult(out string, expected string, t *testing.T) {
	res, err := ioutil.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}

	result := string(res)
	if result != expected {
		t.Errorf(`expected "%s", got "%s"`, expected, result)
	}
}

func TestB64Encode(t *testing.T) {
	in := "./test/in.txt"
	out := "./test/out.txt"
	expected := "aGVsbG8gd29ybGQ="
	exec("fmt", "-i", in, "-o", out, "b64", "-e")
	checkResult(out, expected, t)
}

func TestB64Decode(t *testing.T) {
	in := "./test/in.txt"
	out := "./test/out.txt"
	expected := "hello world"
	exec("fmt", "-i", in, "-o", out, "b64", "-e")
	exec("fmt", "-i", out, "-o", out, "b64", "-d")
	checkResult(out, expected, t)
}

func TestHexEncode(t *testing.T) {
	in := "./test/in.txt"
	out := "./test/out.txt"
	expected := "68656c6c6f20776f726c64"
	exec("fmt", "-i", in, "-o", out, "hex", "-e")
	checkResult(out, expected, t)
}

func TestHexDecode(t *testing.T) {
	in := "./test/in.txt"
	out := "./test/out.txt"
	expected := "hello world"
	exec("fmt", "-i", in, "-o", out, "hex", "-e")
	exec("fmt", "-i", out, "-o", out, "hex", "-d")

	checkResult(out, expected, t)
}

func TestHexEncodeWith0x(t *testing.T) {
	in := "./test/in.txt"
	out := "./test/out.txt"
	expected := "0x680x650x6c0x6c0x6f0x200x770x6f0x720x6c0x64"
	exec("fmt", "-i", in, "-o", out, "hex", "-e", "--delim", "0x", "-p")
	checkResult(out, expected, t)
}

func TestHexDecodeWith0x(t *testing.T) {
	in := "./test/in.txt"
	out := "./test/out.txt"
	expected := "hello world"
	exec("fmt", "-i", in, "-o", out, "hex", "-e", "--delim", "0x", "-p")
	exec("fmt", "-i", out, "-o", out, "hex", "-d", "--delim", "0x", "-p")

	checkResult(out, expected, t)
}

func TestBinEncodeWithSpace(t *testing.T) {
	in := "./test/in.txt"
	out := "./test/out.txt"
	expected := "01101000 01100101 01101100 01101100 01101111 00100000 01110111 01101111 01110010 01101100 01100100"
	exec("fmt", "-i", in, "-o", out, "bin", "-e", "--delim", " ")
	checkResult(out, expected, t)
}

func TestBinDecodeWithSpace(t *testing.T) {
	in := "./test/in.txt"
	out := "./test/out.txt"
	expected := "hello world"
	exec("fmt", "-i", in, "-o", out, "bin", "-e", "--delim", " ")
	exec("fmt", "-i", out, "-o", out, "bin", "-d", "--delim", " ")

	checkResult(out, expected, t)
}

func TestDecEncodeWithLF(t *testing.T) {
	in := "./test/in.txt"
	out := "./test/out.txt"
	expected := "01101000\n01100101\n01101100\n01101100\n01101111\n00100000\n01110111\n01101111\n01110010\n01101100\n01100100"
	exec("fmt", "-i", in, "-o", out, "bin", "-e", "--delim", "\\n")
	checkResult(out, expected, t)
}

func TestDecDecodeWithLF(t *testing.T) {
	in := "./test/in.txt"
	out := "./test/out.txt"
	expected := "hello world"
	exec("fmt", "-i", in, "-o", out, "bin", "-e", "--delim", "\\n")
	exec("fmt", "-i", out, "-o", out, "bin", "-d", "--delim", "\\n")

	checkResult(out, expected, t)
}

func TestURLEncode(t *testing.T) {
	in := "./test/in_url.txt"
	out := "./test/out.txt"
	expected := `https://www.example.com/a/b/?c=d&e=f#g%E4%B8%AD%E6%96%87`
	exec("fmt", "-i", in, "-o", out, "url", "-e")
	checkResult(out, expected, t)
}

func TestURLDecode(t *testing.T) {
	in := "./test/in_url.txt"
	out := "./test/out.txt"
	expected := "https://www.example.com/a/b/?c=d&e=f#g中文"
	exec("fmt", "-i", in, "-o", out, "url", "-e")
	exec("fmt", "-i", out, "-o", out, "url", "-d")

	checkResult(out, expected, t)
}

func TestURLEncodeAll(t *testing.T) {
	in := "./test/in_url.txt"
	out := "./test/out.txt"
	expected := `https%3A%2F%2Fwww.example.com%2Fa%2Fb%2F%3Fc%3Dd%26e%3Df%23g%E4%B8%AD%E6%96%87`
	exec("fmt", "-i", in, "-o", out, "url", "-ea")
	checkResult(out, expected, t)
}

func TestURLDecodeAll(t *testing.T) {
	in := "./test/in_url.txt"
	out := "./test/out.txt"
	expected := "https://www.example.com/a/b/?c=d&e=f#g中文"
	exec("fmt", "-i", in, "-o", out, "url", "-ea")
	exec("fmt", "-i", out, "-o", out, "url", "-d")

	checkResult(out, expected, t)
}
