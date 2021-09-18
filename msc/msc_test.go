package msc

import (
	"io/ioutil"
	"testing"

	"github.com/SignorMercurio/attrezzi/cmd"
	uuid "github.com/satori/go.uuid"
)

var (
	in  = "./test/in.txt"
	out = "./test/out.txt"
	src = "Hello 世界 123"
)

func exec(args ...string) {
	rootCmd := cmd.NewRootCmd()
	mscCmd := NewMscCmd()
	mscCmd.AddCommand(
		NewUidCmd(),
	)
	rootCmd.AddCommand(mscCmd)

	commonArgs := []string{"msc", "-o", out, "-i"}
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

func TestUid(t *testing.T) {
	in := "/dev/null"

	exec(in, "uid")

	uid := readOutput(t)
	_, err := uuid.FromString(string(uid))
	if err != nil {
		t.Error("Failed to Parse generated UUID")
	}
}
