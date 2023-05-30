package main

import (
	"khanhanh_lang/repl"
	"os"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}
