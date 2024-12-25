package main

import (
	"fmt"
	"github.com/kahvecikaan/monkey-lang/repl"
	"os"
	"os/user"
)

func main() {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s, this is the Monkey programming language!\n",
		usr.Username)
	fmt.Printf("Feel free to type any commands\n")
	repl.Start(os.Stdin, os.Stdout)
}
