package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"handler/function"
)

func main() {
	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("Unable to read standard input: %s", err.Error())
	}

	var wg sync.WaitGroup
	fmt.Println(function.Handle(input, &wg))
	wg.Wait()
}
