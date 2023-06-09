package repl

import (
	"bufio"
	"fmt"
	"io"
	"khanhanh_lang/evaluator"
	"khanhanh_lang/lexer"
	"khanhanh_lang/object"
	"khanhanh_lang/parser"
	"os"

	"github.com/diontr00/logStack"
	"github.com/fatih/color"
)

var PROMPT string

var log *logStack.Logger

func init() {
	cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
	green := color.New(color.FgHiGreen, color.Bold).SprintFunc()
	PROMPT = fmt.Sprintf("%s_%s ~>> ", cyan("khanhanh"), green("lang"))
	log = logStack.DefaultLogger()
}

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	tracker := object.NewTracker()
	logFile, err := os.OpenFile("./log.test", os.O_CREATE|os.O_APPEND, 0644)
	defer func() { logFile.Close() }()
	if err != nil {
		log.DPanic(err.Error())
	}

	log := logStack.NewLogger(logFile, logStack.DPanicLevel)
	logStack.ResetDefault(log)
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
		p := parser.New(l)
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printError(out, p.Errors())
			continue
		}
		evaluated := evaluator.Eval(program, tracker)

		if evaluated != nil {
			_, err := io.WriteString(out, evaluated.Inspect())
			if err != nil {
				log.Warn(err.Error())
			}
			_, err = io.WriteString(out, "\n")
			if err != nil {
				log.Warn(err.Error())
			}

		}

	}
}

func printError(out io.Writer, errors []string) {
	for _, msg := range errors {
		msg = color.RedString(msg)
		_, err := io.WriteString(out, "\t"+msg+"\n")
		if err != nil {
			log.Warn(err.Error())
		}
	}
}
