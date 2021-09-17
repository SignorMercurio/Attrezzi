package format

import (
	"io/ioutil"
	"testing"

	"github.com/SignorMercurio/attrezzi/cmd"
)

var (
	in  = "./test/in.txt"
	out = "./test/out.txt"
	src = "hello world"
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

func TestB64(t *testing.T) {
	dst := "aGVsbG8gd29ybGQ="
	exec(in, "b64", "-e")
	checkResult(dst, t)

	exec(out, "b64", "-d")
	checkResult(src, t)
}

func TestHex(t *testing.T) {
	dst := "68656c6c6f20776f726c64"
	exec(in, "hex", "-e")
	checkResult(dst, t)

	exec(out, "hex", "-d")
	checkResult(src, t)
}

func TestHexWith0x(t *testing.T) {
	dst := "0x680x650x6c0x6c0x6f0x200x770x6f0x720x6c0x64"
	exec(in, "hex", "-e", "--delim", "0x", "-p")
	checkResult(dst, t)

	exec(out, "hex", "-d", "--delim", "0x", "-p")
	checkResult(src, t)
}

func TestBinWithSpace(t *testing.T) {
	dst := "01101000 01100101 01101100 01101100 01101111 00100000 01110111 01101111 01110010 01101100 01100100"
	exec(in, "bin", "-e", "--delim", " ")
	checkResult(dst, t)

	exec(out, "bin", "-d", "--delim", " ")
	checkResult(src, t)
}

func TestDecWithLF(t *testing.T) {
	dst := "104\n101\n108\n108\n111\n32\n119\n111\n114\n108\n100\n"
	exec(in, "dec", "-e", "--delim", "\\n")
	checkResult(dst, t)

	exec(out, "dec", "-d", "--delim", "\\n")
	checkResult(src, t)
}

func TestURL(t *testing.T) {
	in := "./test/in_url.txt"
	src := "https://www.example.com/a/b/?c=d&e=f#g中文"
	dst := `https://www.example.com/a/b/?c=d&e=f#g%E4%B8%AD%E6%96%87`
	exec(in, "url", "-e")
	checkResult(dst, t)

	exec(out, "url", "-d")
	checkResult(src, t)
}

func TestURLAll(t *testing.T) {
	in := "./test/in_url.txt"
	src := "https://www.example.com/a/b/?c=d&e=f#g中文"
	dst := `https%3A%2F%2Fwww.example.com%2Fa%2Fb%2F%3Fc%3Dd%26e%3Df%23g%E4%B8%AD%E6%96%87`
	exec(in, "url", "-ea")
	checkResult(dst, t)

	exec(out, "url", "-d")
	checkResult(src, t)
}

func TestHTMLNamedAll(t *testing.T) {
	in := "./test/in_html.txt"
	src := "<script>alert('xss');</script>"
	dst := `&lt;&#115;&#99;&#114;&#105;&#112;&#116;&gt;&#97;&#108;&#101;&#114;&#116;&lpar;&apos;&#120;&#115;&#115;&apos;&rpar;&semi;&lt;&sol;&#115;&#99;&#114;&#105;&#112;&#116;&gt;`
	exec(in, "htm", "-ea")
	checkResult(dst, t)

	exec(out, "htm", "-d")
	checkResult(src, t)
}

func TestHTMLDec(t *testing.T) {
	in := "./test/in_html.txt"
	src := "<script>alert('xss');</script>"
	dst := `&#60;script&#62;alert&#40;&#39;xss&#39;&#41;&#59;&#60;&#47;script&#62;`
	exec(in, "htm", "-e", "-t", "dec")
	checkResult(dst, t)

	exec(out, "htm", "-d", "-t", "dec")
	checkResult(src, t)
}

func TestHTMLHex(t *testing.T) {
	in := "./test/in_html.txt"
	src := "<script>alert('xss');</script>"
	dst := `&#x3c;script&#x3e;alert&#x28;&#x27;xss&#x27;&#x29;&#x3b;&#x3c;&#x2f;script&#x3e;`
	exec(in, "htm", "-e", "-t", "hex")
	checkResult(dst, t)

	exec(out, "htm", "-d", "-t", "hex")
	checkResult(src, t)
}
