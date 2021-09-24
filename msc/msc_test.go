package msc

import (
	"io"
	"testing"

	"github.com/SignorMercurio/attrezzi/cmd"
	"github.com/SignorMercurio/attrezzi/test"
	uuid "github.com/satori/go.uuid"
)

var (
	base = "./testdata/"
	out  = base + "out.txt"
	bla  = "blabla/bla.txt"
)

func exec(args ...string) {
	rootCmd := cmd.NewRootCmd()
	mscCmd := NewMscCmd()
	mscCmd.AddCommand(
		NewUidCmd(),
		NewJpgCmd(),
	)
	rootCmd.AddCommand(mscCmd)

	commonArgs := []string{"msc", "-o"}
	rootCmd.SetArgs(append(commonArgs, args...))
	rootCmd.Execute()
}

func TestMsc(t *testing.T) {
	cmd.Log.SetOutput(io.Discard)
	input = nil
	tests := []test.Test{
		// open output fail
		{Cmd: []string{bla, "uid"}, Dst: ""},
		// open input fail
		{Cmd: []string{out, "-i", bla, "jpg"}, Dst: ""},
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

func TestJpg(t *testing.T) {
	tests := []test.Test{
		{Cmd: []string{out, "-i", base + "sample.jpg", "jpg"}, Dst: `Location: (46.241305, 24.849876)`},
		{Cmd: []string{out, "-i", base + "nogeo.jpg", "jpg"}, Dst: `XResolution: "26856/187"`},
		{Cmd: []string{out, "-i", base + "err.jpg", "jpg"}, Dst: `GPSLatitude: ["51/1","31/1","53584899/1000000"]`},
		{Cmd: []string{out, "-i", base + "in.txt", "jpg"}, Dst: ""},
	}

	for _, tst := range tests {
		exec(tst.Cmd...)
		test.CheckContains(out, tst.Dst, t)
	}
}
