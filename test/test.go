package test

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"
)

type Test struct {
	Cmd []string
	Dst string
}

func ReadOutput(out string, t *testing.T) []byte {
	res, err := ioutil.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	return res
}

func CheckResult(filename string, expected string, t *testing.T) {
	if expected == "*" {
		return
	}
	res := ReadOutput(filename, t)

	result := string(res)
	if result != expected {
		t.Errorf(`expected "%s", got "%s"`, expected, result)
	}
}

func CheckContains(filename string, expected string, t *testing.T) {
	res := ReadOutput(filename, t)

	result := string(res)
	if !strings.Contains(result, expected) {
		t.Errorf(`expected to contain "%s", got "%s"`, expected, result)
	}
}

func CheckNotEmptyAndHasLen(filename string, expectedMinLen uint, expectedMaxLen uint, t *testing.T) {
	empty := make([]byte, expectedMaxLen)
	b := ReadOutput(filename, t)
	lenB := len(b)

	if bytes.Equal(empty, b) || lenB < int(expectedMinLen) || lenB > int(expectedMaxLen) {
		t.Errorf("b is empty OR of invalid length; expected length is %d ~ %d, got %d", expectedMinLen, expectedMaxLen, lenB)
	}
}
