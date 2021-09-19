package enc

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/SignorMercurio/attrezzi/cmd"
)

var (
	in  = "./test/in.txt"
	out = "./test/out.txt"
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
	)
	rootCmd.AddCommand(encCmd)

	commonArgs := []string{"enc", "-o", out, "-i"}
	rootCmd.SetArgs(append(commonArgs, args...))
	rootCmd.Execute()
}

func readOutput(t *testing.T) []byte {
	res, err := ioutil.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	return res
}

func checkResult(expected string, t *testing.T) {
	result := string(readOutput(t))
	if result != expected {
		t.Errorf(`expected "%s", got "%s"`, expected, result)
	}
}

func checkNotEmptyAndHasLen(expectedLen uint, t *testing.T) {
	empty := make([]byte, expectedLen)
	b := readOutput(t)
	lenB := len(b)

	if bytes.Equal(empty, b) || lenB != int(expectedLen) {
		t.Errorf("b is empty OR of invalid length; expected length is %d, got %d", expectedLen, lenB)
	}
}

func TestRot(t *testing.T) {
	dst := "Uryyb 世界 123"
	exec(in, "rot", "-e")
	checkResult(dst, t)

	exec(out, "rot", "-d")
	checkResult(src, t)
}

func TestRot3(t *testing.T) {
	dst := "Khoor 世界 123"
	exec(in, "rot", "-e", "-n", "3")
	checkResult(dst, t)

	exec(out, "rot", "-d", "-n", "3")
	checkResult(src, t)
}

func TestMor(t *testing.T) {
	in := "./test/in_mor.txt"
	src := "HELLO WORLD 123"
	dst := "····/·/·_··/·_··/___\n·__/___/·_·/·_··/_··\n·____/··___/···__"
	exec(in, "mor", "-e", "--dash", "_", "--dot", "·", "-l", "/", "-w", "\n")
	checkResult(dst, t)

	exec(out, "mor", "-d", "--dash", "_", "--dot", "·", "-l", "/", "-w", "\n")
	checkResult(src, t)
}

func TestXorHex(t *testing.T) {
	in := "./test/in_xor.txt"
	dst := "成了"

	exec(in, "xor", "-k", "deadbeefcafe")
	checkResult(dst, t)

	exec(in, "xor", "-k", "110111101010110110111110111011111100101011111110", "--key-fmt", "bin")
	checkResult(dst, t)

	exec(in, "xor", "-k", "244837814094590", "--key-fmt", "dec")
	checkResult(dst, t)

	exec(in, "xor", "-k", "3q2+78r+", "--key-fmt", "b64")
	checkResult(dst, t)
}

func TestXorUTF8(t *testing.T) {
	in := "./test/in_xor_utf8.txt"
	dst := "@@@@@@@"

	exec(in, "xor", "-k", `!"#$%&'`, "--key-fmt", "utf8", "--input-fmt", "utf8")
	checkResult(dst, t)
}

func TestRndHex(t *testing.T) {
	in := "/dev/null" // not really needed

	exec(in, "rnd", "-l", "16")
	checkNotEmptyAndHasLen(32, t)
}

func TestRndBin(t *testing.T) {
	in := "/dev/null"

	exec(in, "rnd", "-l", "8", "-f", "bin")
	checkNotEmptyAndHasLen(64, t)
}

func TestAes128CBC(t *testing.T) {
	key := "f5f73713bc57d1cec7deb623b292bbc6"
	exec(in, "aes", "-e", "-m", "cbc", "-k", key)
	exec(out, "aes", "-d", "-m", "cbc", "-k", key)
	checkResult(src, t)
}

func TestAes192CFB(t *testing.T) {
	key := "dd4ebf0e6ced5fb8a356d1acf843d672656d2261590195d2"
	exec(in, "aes", "-e", "-m", "cfb", "-k", key)
	exec(out, "aes", "-d", "-m", "cfb", "-k", key)
	checkResult(src, t)
}

func TestAes128OFB(t *testing.T) {
	key := "f5f73713bc57d1cec7deb623b292bbc6"
	exec(in, "aes", "-e", "-m", "ofb", "-k", key)
	exec(out, "aes", "-d", "-m", "ofb", "-k", key)
	checkResult(src, t)
}

func TestAes192CTR(t *testing.T) {
	key := "dd4ebf0e6ced5fb8a356d1acf843d672656d2261590195d2"
	exec(in, "aes", "-e", "-m", "ctr", "-k", key)
	exec(out, "aes", "-d", "-m", "ctr", "-k", key)
	checkResult(src, t)
}

func TestAes256GCM(t *testing.T) {
	key := "0d94f846deac35f48e8055413c556263e647f36feb939f0c49562dcb6a718d9c"
	exec(in, "aes", "-e", "-k", key)
	exec(out, "aes", "-d", "-k", key)
	checkResult(src, t)
}
