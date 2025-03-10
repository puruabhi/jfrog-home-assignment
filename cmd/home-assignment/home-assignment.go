package main

import (
	"fmt"

	"github.com/puruabhi/jfrog/home-assignment/internal/process"
)

func main() {
	if finished, err := process.Setup(); err != nil {
		panic(err)
	} else {
		<-finished
		fmt.Println("Successfully finished process, exiting...")
	}
}
