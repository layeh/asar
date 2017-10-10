package main // import "layeh.com/asar/cmd/asar"

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"layeh.com/asar"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [options] [command]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Commands:\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  l|list <archive>\n")
		fmt.Fprintf(os.Stderr, "    list contents of asar archive\n")
		fmt.Fprintf(os.Stderr, "  x|extract <archive> <dir>\n")
		fmt.Fprintf(os.Stderr, "    extract contents of asar archive to directory\n")
		fmt.Fprintf(os.Stderr, "\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() < 2 {
		flag.Usage()
		os.Exit(3)
	}

	file, err := os.Open(flag.Arg(1))
	if err != nil {
		fmt.Fprintf(os.Stderr, "asar: %s\n", err)
		os.Exit(1)
	}
	defer file.Close()

	root, err := asar.Decode(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "asar: %s\n", err)
		os.Exit(1)
	}

	switch command := flag.Arg(0); command {
	case "l", "list":
		root.Walk(func(path string, _ os.FileInfo, _ error) error {
			fmt.Println("/" + path)
			return nil
		})

	case "x", "extract":
		if flag.NArg() < 3 {
			flag.Usage()
			os.Exit(1)
		}

		target := flag.Arg(2)

		err := root.WalkEntry(func(entry *asar.Entry, _ error) error {
			if entry.Flags&asar.FlagDir == asar.FlagDir {
				return os.Mkdir(entry.Path(), 0755)
			}
			if entry.Flags&asar.FlagUnpacked != 0 {
				return nil
			}

			realPath := filepath.Join(target, entry.Path())

			f, err := os.Create(realPath)
			if err != nil {
				return err
			}

			_, err = entry.WriteTo(f)
			return err
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "asar: %s\n", err)
			os.Exit(1)
		}
	}
}
