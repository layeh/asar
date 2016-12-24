package asar // import "layeh.com/asar"

import (
	"io/ioutil"
	"strings"
	"testing"
)

func TestEncodeInvalidName(t *testing.T) {
	root := New(".", nil, 0, 0, FlagDir)
	root.Children = append(
		root.Children,
		New(".", strings.NewReader("test"), 4, 0, FlagNone),
	)
	if _, err := root.EncodeTo(ioutil.Discard); err == nil {
		t.Fatal("we should have had an error")
	}
}

func TestEncodeUnpacked(t *testing.T) {
	root := New(".", nil, 0, 0, FlagDir)
	root.Children = append(
		root.Children,
		New("sample", nil, 0, 0, FlagUnpacked),
	)
	if _, err := root.EncodeTo(ioutil.Discard); err != nil {
		t.Fatalf("err: %s", err)
	}
}
