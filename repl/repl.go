package repl

import (
	"bufio"
	"fmt"
	"github.com/kahvecikaan/monkey-lang/evaluator"
	"github.com/kahvecikaan/monkey-lang/lexer"
	"github.com/kahvecikaan/monkey-lang/object"
	"github.com/kahvecikaan/monkey-lang/parser"
	"io"
)

const PROMPT = ">> "

// ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorOrange = "\033[38;5;208m"
)

const MONKEY_FACE = ColorOrange + `
            __,__
   .--.  .-"     "-.  .--.
  / .. \/  .-. .-.  \/ .. \
 | |  '|  /   Y   \  |'  | |
 | \   \  \ 0 | 0 /  /   / |
  \ '- ,\.-"""""""-./, -' /
   ''-' /_   ^ ^   _\ '-''
       |  \._   _./  |
       \   \ '~' /   /
        '._ '-=-' _.'
           '-----'

 ███╗   ███╗ ██████╗ ███╗   ██╗██╗  ██╗███████╗██╗   ██╗
 ████╗ ████║██╔═══██╗████╗  ██║██║ ██╔╝██╔════╝╚██╗ ██╔╝
 ██╔████╔██║██║   ██║██╔██╗ ██║█████╔╝ █████╗   ╚████╔╝ 
 ██║╚██╔╝██║██║   ██║██║╚██╗██║██╔═██╗ ██╔══╝    ╚██╔╝  
 ██║ ╚═╝ ██║╚██████╔╝██║ ╚████║██║  ██╗███████╗   ██║   
 ╚═╝     ╚═╝ ╚═════╝ ╚═╝  ╚═══╝╚═╝  ╚═╝╚══════╝   ╚═╝   
         SYNTAX ERROR - TIME TO DEBUG!
` + ColorReset

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, MONKEY_FACE)
	io.WriteString(out, "Whoops! We ran into some monkey business here!\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
