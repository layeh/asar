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

		err := root.Walk(func(path string, info os.FileInfo, _ error) error {
			entry := info.Sys().(*asar.Entry)

			realPath := filepath.Join(target, path)
			if entry.Flags&asar.FlagDir != 0 {
				return os.MkdirAll(realPath, 0755)
			}
			if entry.Flags&asar.FlagUnpacked != 0 {
				return nil
			}

			perm := os.FileMode(0644)
			if entry.Flags&asar.FlagExecutable != 0 {
				perm |= 0111
			}

			f, err := os.OpenFile(realPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, perm)
			if err != nil {
				return err
			}

			_, err = entry.WriteTo(f)
			if err != nil {
				f.Close()
				return err
			}

			if err := f.Close(); err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "asar: %s\n", err)
			os.Exit(1)
		}
	}
}
