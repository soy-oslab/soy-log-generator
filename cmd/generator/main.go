package main

import (
	"log"

	"flag"
	"sync"
	"time"

	"github.com/soyoslab/soy_log_generator/pkg/classifier"
	"github.com/soyoslab/soy_log_generator/pkg/transport"
)

var wg sync.WaitGroup
var c *classifier.Classifier

func run() {
	t, err := transport.InitTransport("config.json", func(str string, isHot bool) bool {
		if isHot {
			c.Learn(str, classifier.Hot)
			return true
		}
		c.Learn(str, classifier.Cold)
		result, _ := c.Classify(str)
		if result[classifier.Hot] > 1e-10 {
			return true
		}
		log.Println(result, str, false)
		return false
	})
	if err != nil {
		log.Fatalln(err)
	}
	defer t.Close()
	err = t.Run()
	if err != nil {
		log.Println(":test", err)
	}
	wg.Done()
}

func backup() {
	for {
		err := c.Backup()
		if err != nil {
			break
		}
		time.Sleep(time.Duration(1) * time.Second)
	}
}

func main() {
	var err error

	saveFile := flag.String("file", "model.sav", "Bayesian model's save path")
	flag.Parse()
	wg = sync.WaitGroup{}
	c, err = classifier.InitClassfier(*saveFile)
	if err != nil {
		log.Panicf("initialize the classifier failed %v (filepath: %s)", err, *saveFile)
	}

	go backup()
	for {
		wg.Add(1)
		go run()
		wg.Wait()
		time.Sleep(time.Duration(1) * time.Second)
	}
}
