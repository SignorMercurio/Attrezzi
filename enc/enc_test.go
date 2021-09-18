package enc

import (
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
	)
	rootCmd.AddCommand(encCmd)

	commonArgs := []string{"enc", "-o", out, "-i"}
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
