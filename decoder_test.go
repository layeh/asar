package asar

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestExtractThis(t *testing.T) {
	f, err := os.Open("testdata/extractthis.asar")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	root, err := Decode(f)
	if err != nil {
		t.Fatal(err)
	}

	if root.Flags&FlagDir == 0 {
		t.Fatal("expecting root directory to have FlagDir")
	}

	{
		f1 := root.Find("dir1", "file1.txt")
		if f1 == nil {
			t.Fatal("could not find dir1/file1.txt")
		}
		if f1.Path() != "/dir1/file1.txt" {
			t.Fatal("unexpected path")
		}
		body, err := ioutil.ReadAll(f1.Open())
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(body, []byte(`file one.`)) {
			t.Fatalf("dir1/file1.txt body is incorrect (got %s)", body)
		}
	}

	{
		f2 := root.Find("dir2").Find("file3.txt")
		if f2 == nil {
			t.Fatal("could not find dir2/file3.txt")
		}
		s := f2.String()
		if s != `123` {
			t.Fatalf("dir2/file3.txt body is incorrect (got %s)", s)
		}
	}

	{
		empty := root.Find("emptyfile.txt")
		if empty == nil {
			t.Fatal("could not find emptyfile.txt")
		}
		if len(empty.Bytes()) != 0 {
			t.Fatal("expecting emptyfile.txt to be empty")
		}
	}
}
