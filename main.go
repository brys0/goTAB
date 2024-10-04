package main

import (
	"fmt"
	"github.com/brys0/goTAB/internal/app"
	"time"
)

func main() {
	defer timer("Application")()
	application := app.CreateNewTAB(true)

	//verbose := true

	application.GetHardware()

	application.PromptGPU()
}

func timer(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", name, time.Since(start))
	}
}
