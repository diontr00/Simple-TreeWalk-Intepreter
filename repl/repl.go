package repl

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"io"
	"khanhanh_lang/lexer"
	"khanhanh_lang/token"
)

var PROMPT string

func init() {
	cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
	green := color.New(color.FgHiGreen, color.Bold).SprintFunc()
	PROMPT = fmt.Sprintf("%s_%s ~ ", cyan("khanhanh"), green("lang"))
}

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	for {
		fmt.Printf("%s", PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		if line == "quit" {
			return
		}
		l := lexer.New(line)
		for o_token := l.NextToken(); o_token.Type != token.EOF; o_token = l.NextToken() {
			fmt.Printf("%+v\n", o_token.Type)
		}
	}

}
