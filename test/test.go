package test

import (
	"io/ioutil"
	"testing"
)

type Test struct {
	Cmd []string
	Dst string
}

func CheckResult(filename string, expected string, t *testing.T) {
	res, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}

	result := string(res)
	if result != expected {
		t.Errorf(`expected "%s", got "%s"`, expected, result)
	}
}
