package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/layeh/asar"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [options] [command]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Commands:\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  l|list <archive>\n")
		fmt.Fprintf(os.Stderr, "    list contents of asar archive\n")
		fmt.Fprintf(os.Stderr, "\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() == 2 && (flag.Arg(0) == "l" || flag.Arg(0) == "list") {
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

		type stack struct {
			Entry      *asar.Entry
			ChildIndex int
		}
		s := []*stack{
			&stack{root, 0},
		}

		outer:
		for len(s) > 0 {
			top := s[len(s) - 1]
			for top.ChildIndex < len(top.Entry.Children)  {
				child := top.Entry.Children[top.ChildIndex]
				top.ChildIndex++

				fmt.Println(child.Path())
				next := &stack{child, 0}
				s = append(s, next)
				continue outer
			}
			s = s[:len(s) - 1]
		}

	} else {
		flag.Usage()
		os.Exit(3)
	}
}
