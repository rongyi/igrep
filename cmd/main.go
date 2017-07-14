package main

import (
	"igrep"
	"fmt"
	"os"
)

func main() {
	e, err := igrep.NewEngine()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(e.Run())
}
