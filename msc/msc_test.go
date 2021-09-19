package msc

import (
	"io"
	"testing"

	"github.com/SignorMercurio/attrezzi/cmd"
	"github.com/SignorMercurio/attrezzi/test"
	uuid "github.com/satori/go.uuid"
)

var (
	in  = "./testdata/in.txt"
	out = "./testdata/out.txt"
	src = "Hello 世界 123"
)

func exec(args ...string) {
	rootCmd := cmd.NewRootCmd()
	mscCmd := NewMscCmd()
	mscCmd.AddCommand(
		NewUidCmd(),
	)
	rootCmd.AddCommand(mscCmd)

	commonArgs := []string{"msc", "-o"}
	rootCmd.SetArgs(append(commonArgs, args...))
	rootCmd.Execute()
}

func TestMsc(t *testing.T) {
	cmd.Log.SetOutput(io.Discard)
	tests := []test.Test{
		// open output fail
		{Cmd: []string{"bla/blabla.txt", "uid"}, Dst: "*"}, // TODO
		// open input fail
		// {Cmd: []string{"bla/blabla.txt", "-o", out, "rot", "-e"}, Dst: ""},
	}

	for _, tst := range tests {
		exec(tst.Cmd...)
		test.CheckResult(out, tst.Dst, t)
	}
}

func TestUid(t *testing.T) {
	exec(out, "uid")

	uid := test.ReadOutput(out, t)
	_, err := uuid.FromString(string(uid))
	if err != nil {
		t.Error("Failed to Parse generated UUID")
	}
}
