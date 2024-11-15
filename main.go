package main

import (
	"os"

	"github.com/carlmjohnson/exitcode"
	"github.com/spotlightpa/moreofa/commentthan"
)

func main() {
	exitcode.Exit(commentthan.CLI(os.Args[1:]))
}
