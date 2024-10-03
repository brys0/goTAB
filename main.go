package main

import (
	"github.com/brys0/goTAB/internal/memory"
)

func main() {
	_, err := memory.Get_memory_info()

	if err != nil {
		panic(err)
	}

}
