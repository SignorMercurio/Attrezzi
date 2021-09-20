package enc

import (
	"io"
	"strconv"
	"strings"
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
	encCmd := NewEncCmd()
	encCmd.AddCommand(
		NewRotCmd(),
		NewMorCmd(),
		NewXorCmd(),
		NewRndCmd(),
		NewAesCmd(),
		NewRkgCmd(),
		NewRsaCmd(),
		NewHshCmd(),
	)
	rootCmd.AddCommand(encCmd)

	commonArgs := []string{"enc", "-o", out, "-i"}
	rootCmd.SetArgs(append(commonArgs, args...))
	rootCmd.Execute()
}

func TestEnc(t *testing.T) {
	rootCmd := cmd.NewRootCmd()
	encCmd := NewEncCmd()
	encCmd.AddCommand(NewRotCmd())
	rootCmd.AddCommand(encCmd)

	cmd.Log.SetOutput(io.Discard)
	input = nil
	tests := []test.Test{
		// open output fail
		{Cmd: []string{"enc", "-o", "bla/blabla.txt", "-i", in, "rot", "-e"}, Dst: ""},
		// open input fail
		{Cmd: []string{"enc", "-i", "bla/blabla.txt", "-o", out, "rot", "-e"}, Dst: ""},
		// read input fail
		{Cmd: []string{"enc", "-i", "", "-o", out, "rot", "-e"}, Dst: ""},
	}

	for _, tst := range tests {
		rootCmd.SetArgs(tst.Cmd)
		rootCmd.Execute()
		test.CheckResult(out, tst.Dst, t)
	}
}

func TestRot(t *testing.T) {
	tests := []test.Test{
		// default
		{Cmd: []string{in, "rot", "-e"}, Dst: "Uryyb 世界 123"},
		{Cmd: []string{out, "rot", "-d"}, Dst: src},
		// caesar
		{Cmd: []string{in, "rot", "-e", "-n", "3"}, Dst: "Khoor 世界 123"},
		{Cmd: []string{out, "rot", "-d", "-n", "3"}, Dst: src},
		// no action
		{Cmd: []string{in, "rot"}, Dst: ""},
	}

	for _, tst := range tests {
		exec(tst.Cmd...)
		test.CheckResult(out, tst.Dst, t)
	}
}

func TestMor(t *testing.T) {
	in := "./testdata/in_mor.txt"
	src := "HELLO WORLD 123"

	tests := []test.Test{
		// custom
		{Cmd: []string{in, "mor", "-e", "--dash", "_", "--dot", "·", "-l", "/", "-w", `\n`}, Dst: "····/·/·_··/·_··/___\n·__/___/·_·/·_··/_··\n·____/··___/···__"},
		{Cmd: []string{out, "mor", "-d", "--dash", "_", "--dot", "·", "-l", "/", "-w", `\n`}, Dst: src},
		// with \r\n
		{Cmd: []string{in, "mor", "-e", "-l", "/", "-w", `\r\n`}, Dst: "...././.-../.-../---\r\n.--/---/.-./.-../-..\r\n.----/..---/...--"},
		{Cmd: []string{out, "mor", "-d", "-l", "/", "-w", `\r\n`}, Dst: src},
		// no action
		{Cmd: []string{in, "mor"}, Dst: ""},
	}

	for _, tst := range tests {
		exec(tst.Cmd...)
		test.CheckResult(out, tst.Dst, t)
	}
}

func TestXor(t *testing.T) {
	in := "./testdata/in_xor.txt"
	dst := "成了"

	tests := []test.Test{
		// hex ^ hex
		{Cmd: []string{in, "xor", "-k", "deadbeefcafe"}, Dst: dst},
		// hex ^ bin
		{Cmd: []string{in, "xor", "-k", "110111101010110110111110111011111100101011111110", "--key-fmt", "bin"}, Dst: dst},
		// invalid bin
		{Cmd: []string{in, "xor", "-k", "f110111101010110110111110111011111100101011111110", "--key-fmt", "bin"}, Dst: ""},
		// hex ^ dec
		{Cmd: []string{in, "xor", "-k", "244837814094590", "--key-fmt", "dec"}, Dst: dst},
		// invalid dec
		{Cmd: []string{in, "xor", "-k", "f244837814094590", "--key-fmt", "dec"}, Dst: ""},
		// hex ^ b64
		{Cmd: []string{in, "xor", "-k", "3q2+78r+", "--key-fmt", "b64"}, Dst: dst},
		// invalid base 64
		{Cmd: []string{in, "xor", "-k", "3q2+78r+===", "--key-fmt", "b64"}, Dst: ""},
		// invalid input
		{Cmd: []string{"./testdata/in_xor_utf8.txt", "xor", "-k", "3q2+78r+", "--key-fmt", "b64", "--input-fmt", "hex"}, Dst: ""},
		// utf8 ^ utf8
		{Cmd: []string{"./testdata/in_xor_utf8.txt", "xor", "-k", `!"#$%&'`, "--key-fmt", "utf8", "--input-fmt", "utf8"}, Dst: "@@@@@@@"},
		// no key
		{Cmd: []string{in, "xor"}, Dst: ""},
	}

	for _, tst := range tests {
		exec(tst.Cmd...)
		test.CheckResult(out, tst.Dst, t)
	}
}

func TestRnd(t *testing.T) {
	in := "/dev/null" // not really needed

	tests := []test.Test{
		// hex
		{Cmd: []string{in, "rnd", "-l", "16"}, Dst: "32,32"},
		// bin
		{Cmd: []string{in, "rnd", "-l", "8", "-f", "bin"}, Dst: "64,64"},
		// dec
		{Cmd: []string{in, "rnd", "-l", "1", "-f", "dec"}, Dst: "1,3"},
	}

	for _, tst := range tests {
		exec(tst.Cmd...)
		dst := strings.Split(tst.Dst, ",")
		minLenExpected, _ := strconv.ParseUint(dst[0], 10, 64)
		maxLenExpected, _ := strconv.ParseUint(dst[1], 10, 64)
		test.CheckNotEmptyAndHasLen(out, uint(minLenExpected), uint(maxLenExpected), t)
	}
}

func TestAES(t *testing.T) {
	in_fail := "./testdata/in_xor_utf8.txt"

	key16 := "f5f73713bc57d1cec7deb623b292bbc6"
	key24 := "dd4ebf0e6ced5fb8a356d1acf843d672656d2261590195d2"
	key32 := "0d94f846deac35f48e8055413c556263e647f36feb939f0c49562dcb6a718d9c"

	tests := []test.Test{
		// aes-128-cbc
		{Cmd: []string{in, "aes", "-e", "-m", "cbc", "-k", key16}, Dst: "*"},
		{Cmd: []string{out, "aes", "-d", "-m", "cbc", "-k", key16}, Dst: src},
		// aes-128-cbc invalid key
		{Cmd: []string{in, "aes", "-e", "-m", "cbc", "-k", "123"}, Dst: ""},
		{Cmd: []string{in, "aes", "-e", "-m", "cbc", "-k", "1234"}, Dst: ""},
		{Cmd: []string{in, "aes", "-d", "-m", "cbc", "-k", "1234"}, Dst: ""},
		// aes-128-cbc invalid ciphertext
		{Cmd: []string{in_fail, "aes", "-d", "-m", "cbc", "-k", key16}, Dst: ""},
		// aes-192-cfb
		{Cmd: []string{in, "aes", "-e", "-m", "cfb", "-k", key24}, Dst: "*"},
		{Cmd: []string{out, "aes", "-d", "-m", "cfb", "-k", key24}, Dst: src},
		// aes-192-cfb invalid key
		{Cmd: []string{in, "aes", "-e", "-m", "cfb", "-k", "1234"}, Dst: ""},
		{Cmd: []string{in, "aes", "-d", "-m", "cfb", "-k", "1234"}, Dst: ""},
		// aes-192-cfb invalid ciphertext
		{Cmd: []string{in_fail, "aes", "-d", "-m", "cfb", "-k", key24}, Dst: ""},
		// aes-128-ofb
		{Cmd: []string{in, "aes", "-e", "-m", "ofb", "-k", key16}, Dst: "*"},
		{Cmd: []string{out, "aes", "-d", "-m", "ofb", "-k", key16}, Dst: src},
		// aes-128-ofb invalid key
		{Cmd: []string{in, "aes", "-e", "-m", "ofb", "-k", "1234"}, Dst: ""},
		{Cmd: []string{in, "aes", "-d", "-m", "ofb", "-k", "1234"}, Dst: ""},
		// aes-128-ofb invalid ciphertext
		{Cmd: []string{in_fail, "aes", "-d", "-m", "ofb", "-k", key16}, Dst: ""},
		// aes-192-ctr
		{Cmd: []string{in, "aes", "-e", "-m", "ctr", "-k", key24}, Dst: "*"},
		{Cmd: []string{out, "aes", "-d", "-m", "ctr", "-k", key24}, Dst: src},
		// aes-192-ctr invalid key
		{Cmd: []string{in, "aes", "-e", "-m", "ctr", "-k", "1234"}, Dst: ""},
		{Cmd: []string{in, "aes", "-d", "-m", "ctr", "-k", "1234"}, Dst: ""},
		// aes-192-ctr invalid ciphertext
		{Cmd: []string{in_fail, "aes", "-d", "-m", "ctr", "-k", key24}, Dst: ""},
		// aes-256-gcm
		{Cmd: []string{in, "aes", "-e", "-k", key32}, Dst: "*"},
		{Cmd: []string{out, "aes", "-d", "-k", key32}, Dst: src},
		// aes-256-gcm invalid key
		{Cmd: []string{in, "aes", "-e", "-k", "1234"}, Dst: ""},
		{Cmd: []string{in, "aes", "-d", "-k", "1234"}, Dst: ""},
		// aes-256-gcm invalid ciphertext
		{Cmd: []string{in_fail, "aes", "-d", "-k", key32}, Dst: ""},
		// aes-256-gcm wrong key
		{Cmd: []string{in, "aes", "-d", "-k", key24}, Dst: ""},
		//no action
		{Cmd: []string{in, "aes"}, Dst: ""},
	}

	for _, tst := range tests {
		exec(tst.Cmd...)
		test.CheckResult(out, tst.Dst, t)
	}
}

func TestRkgAndRsa(t *testing.T) {
	in_empty := "/dev/null"
	bla := "blabla/bla.txt"
	pubKeyOut := "./testdata/pubkey.pub"
	privKeyOut := "./testdata/priv.key"

	tests := []test.Test{
		// rkg
		{Cmd: []string{in_empty, "rkg", "--pub", pubKeyOut, "--priv", bla}, Dst: ""},
		{Cmd: []string{in_empty, "rkg", "--pub", bla, "--priv", privKeyOut}, Dst: ""},
		{Cmd: []string{in_empty, "rkg", "--pub", pubKeyOut, "--priv", privKeyOut}, Dst: ""},
		// invalid pubkey
		{Cmd: []string{in, "rsa", "--pub", bla, "--priv", privKeyOut, "-e"}, Dst: ""},
		// damaged pubkey
		{Cmd: []string{in, "rsa", "--pub", "./testdata/pubkey_b1.pub", "--priv", privKeyOut, "-e"}, Dst: ""},
		// modified pubkey
		{Cmd: []string{in, "rsa", "--pub", privKeyOut, "--priv", privKeyOut, "-e"}, Dst: ""},
		// invalid privkey
		{Cmd: []string{in, "rsa", "--pub", pubKeyOut, "--priv", bla, "-d"}, Dst: ""},
		// damaged privkey
		{Cmd: []string{in, "rsa", "--pub", pubKeyOut, "--priv", "./testdata/priv_b1.key", "-d"}, Dst: ""},
		// modified privkey
		{Cmd: []string{in, "rsa", "--pub", pubKeyOut, "--priv", pubKeyOut, "-d"}, Dst: ""},
		// rsa-oaep with sha256
		{Cmd: []string{in, "rsa", "--pub", pubKeyOut, "--priv", privKeyOut, "-e"}, Dst: "*"},
		{Cmd: []string{out, "rsa", "--pub", pubKeyOut, "--priv", privKeyOut, "-d"}, Dst: src},
		// with md5
		{Cmd: []string{in, "rsa", "--pub", pubKeyOut, "--priv", privKeyOut, "-e", "--hash", "md5"}, Dst: "*"},
		{Cmd: []string{out, "rsa", "--pub", pubKeyOut, "--priv", privKeyOut, "-d", "--hash", "md5"}, Dst: src},
		// with sha1
		{Cmd: []string{in, "rsa", "--pub", pubKeyOut, "--priv", privKeyOut, "-e", "--hash", "sha1"}, Dst: "*"},
		{Cmd: []string{out, "rsa", "--pub", pubKeyOut, "--priv", privKeyOut, "-d", "--hash", "sha1"}, Dst: src},
		// with sha384
		{Cmd: []string{in, "rsa", "--pub", pubKeyOut, "--priv", privKeyOut, "-e", "--hash", "sha384"}, Dst: "*"},
		{Cmd: []string{out, "rsa", "--pub", pubKeyOut, "--priv", privKeyOut, "-d", "--hash", "sha384"}, Dst: src},
		// with sha512
		{Cmd: []string{in, "rsa", "--pub", pubKeyOut, "--priv", privKeyOut, "-e", "--hash", "sha512"}, Dst: "*"},
		{Cmd: []string{out, "rsa", "--pub", pubKeyOut, "--priv", privKeyOut, "-d", "--hash", "sha512"}, Dst: src},
		// rsa-pkcs1v15
		{Cmd: []string{in, "rsa", "--pub", pubKeyOut, "--priv", privKeyOut, "-e", "-m", "pkcs1v15"}, Dst: "*"},
		{Cmd: []string{out, "rsa", "--pub", pubKeyOut, "--priv", privKeyOut, "-d", "-m", "pkcs1v15"}, Dst: src},
		// no action
		{Cmd: []string{in, "rsa"}, Dst: ""},
	}

	for _, tst := range tests {
		exec(tst.Cmd...)
		test.CheckResult(out, tst.Dst, t)
	}
}

func TestHsh(t *testing.T) {
	tests := []test.Test{
		// sha256
		{Cmd: []string{in, "hsh"}, Dst: "982b10efe4fece5c4d91b7e90bfc6c1b5c0ada421ad67689d6c19c2b2873b0a5"},
		// md5
		{Cmd: []string{in, "hsh", "--hash", "md5"}, Dst: "512ece16e11bcacb827a923093e5ea80"},
		// sha1
		{Cmd: []string{in, "hsh", "--hash", "sha1"}, Dst: "a9fa680cdc75a63f562d4f76860d51ef572e6eb4"},
		// sha384
		{Cmd: []string{in, "hsh", "--hash", "sha384"}, Dst: "2228c508f652b7f1e9b06b87d76b9a23c4e732f14b2c81e39fb35d080e5f981fa9e13fa6536ee680b179ab2b74785edc"},
		// sha512
		{Cmd: []string{in, "hsh", "--hash", "sha512"}, Dst: "da03b6f9510a7325fdd38677e1332e4179bc99ab4c828e44307434e29e8ac7fcf5a7f0077632797041e689b2f9cd9067d92c49b208255514c66b5bc86ce4e5ec"},
	}
	for _, tst := range tests {
		exec(tst.Cmd...)
		test.CheckResult(out, tst.Dst, t)
	}
}
