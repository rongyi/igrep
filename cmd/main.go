package main

import (
	"github.com/rongyi/igrep"
	"fmt"
	"os"
)

func main() {
	e, err := igrep.NewEngine(os.Stdin)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(e.RunWithOutput())
}
