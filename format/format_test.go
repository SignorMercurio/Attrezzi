package format

import (
	"io/ioutil"
	"testing"

	"github.com/SignorMercurio/attrezzi/cmd"
)

var (
	in  = "./test/in.txt"
	out = "./test/out.txt"
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

	commonArgs := []string{"fmt", "-o", out, "-i"}
	rootCmd.SetArgs(append(commonArgs, args...))
	rootCmd.Execute()
}

func checkResult(expected string, t *testing.T) {
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
	expected := "aGVsbG8gd29ybGQ="
	exec(in, "b64", "-e")
	checkResult(expected, t)
}

func TestB64Decode(t *testing.T) {
	expected := "hello world"
	exec(in, "b64", "-e")
	exec(out, "b64", "-d")
	checkResult(expected, t)
}

func TestHexEncode(t *testing.T) {
	expected := "68656c6c6f20776f726c64"
	exec(in, "hex", "-e")
	checkResult(expected, t)
}

func TestHexDecode(t *testing.T) {
	expected := "hello world"
	exec(in, "hex", "-e")
	exec(out, "hex", "-d")
	checkResult(expected, t)
}

func TestHexEncodeWith0x(t *testing.T) {
	expected := "0x680x650x6c0x6c0x6f0x200x770x6f0x720x6c0x64"
	exec(in, "hex", "-e", "--delim", "0x", "-p")
	checkResult(expected, t)
}

func TestHexDecodeWith0x(t *testing.T) {
	expected := "hello world"
	exec(in, "hex", "-e", "--delim", "0x", "-p")
	exec(out, "hex", "-d", "--delim", "0x", "-p")
	checkResult(expected, t)
}

func TestBinEncodeWithSpace(t *testing.T) {
	expected := "01101000 01100101 01101100 01101100 01101111 00100000 01110111 01101111 01110010 01101100 01100100"
	exec(in, "bin", "-e", "--delim", " ")
	checkResult(expected, t)
}

func TestBinDecodeWithSpace(t *testing.T) {
	expected := "hello world"
	exec(in, "bin", "-e", "--delim", " ")
	exec(out, "bin", "-d", "--delim", " ")
	checkResult(expected, t)
}

func TestDecEncodeWithLF(t *testing.T) {
	expected := "104\n101\n108\n108\n111\n32\n119\n111\n114\n108\n100\n"
	exec(in, "dec", "-e", "--delim", "\\n")
	checkResult(expected, t)
}

func TestDecDecodeWithLF(t *testing.T) {
	expected := "hello world"
	exec(in, "dec", "-e", "--delim", "\\n")
	exec(out, "dec", "-d", "--delim", "\\n")
	checkResult(expected, t)
}

func TestURLEncode(t *testing.T) {
	in := "./test/in_url.txt"
	expected := `https://www.example.com/a/b/?c=d&e=f#g%E4%B8%AD%E6%96%87`
	exec(in, "url", "-e")
	checkResult(expected, t)
}

func TestURLDecode(t *testing.T) {
	in := "./test/in_url.txt"
	expected := "https://www.example.com/a/b/?c=d&e=f#g中文"
	exec(in, "url", "-e")
	exec(out, "url", "-d")
	checkResult(expected, t)
}

func TestURLEncodeAll(t *testing.T) {
	in := "./test/in_url.txt"
	expected := `https%3A%2F%2Fwww.example.com%2Fa%2Fb%2F%3Fc%3Dd%26e%3Df%23g%E4%B8%AD%E6%96%87`
	exec(in, "url", "-ea")
	checkResult(expected, t)
}

func TestURLDecodeAll(t *testing.T) {
	in := "./test/in_url.txt"
	expected := "https://www.example.com/a/b/?c=d&e=f#g中文"
	exec(in, "url", "-ea")
	exec(out, "url", "-d")
	checkResult(expected, t)
}
