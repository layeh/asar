package asar

import (
	"io/ioutil"
	"strings"
	"testing"
)

func TestEncodeInvalid(t *testing.T) {
	root := New(".", nil, 0, 0, FlagDir)
	root.Children = append(
		root.Children,
		New(".", strings.NewReader("test"), 4, 0, FlagNone),
	)
	if _, err := root.EncodeTo(ioutil.Discard); err == nil {
		t.Fatal("we should have had an error")
	}
}
