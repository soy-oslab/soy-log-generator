package main

import (
	"log"

	"sync"

	"github.com/soyoslab/soy_log_generator/pkg/transport"
)

var wg sync.WaitGroup

func run() {
	t, err := transport.InitTransport("config.json", nil)
	if err != nil {
		log.Fatalln(err)
	}
	defer t.Close()
	err = t.Run()
	if err != nil {
		log.Println(err)
	}
	wg.Done()
}

func main() {
	wg = sync.WaitGroup{}
	for {
		wg.Add(1)
		go run()
		wg.Wait()
	}
}
