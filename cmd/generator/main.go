package main

import (
	"io"
	"log"
	"os"

	"flag"
	"sync"
	"time"

	"github.com/soyoslab/soy_log_generator/pkg/classifier"
	"github.com/soyoslab/soy_log_generator/pkg/transport"
)

var wg sync.WaitGroup
var c *classifier.Classifier
var backupInterval int
var mutex *sync.Mutex

func filter(str string, isHot bool) bool {
	if isHot {
		mutex.Lock()
		c.Learn(str, classifier.Hot)
		mutex.Unlock()
		return true
	}
	mutex.Lock()
	c.Learn(str, classifier.Cold)
	mutex.Unlock()
	result, _ := c.Classify(str)
	if result[classifier.Hot] > 1e-10 {
		return true
	}
	return false

}

func run(configFilePath string) {
	t, err := transport.InitTransport(configFilePath, filter)
	if err != nil {
		log.Fatalln(err)
	}
	defer t.Close()
	log.Println("transport running start")
	err = t.Run()
	if err != nil {
		log.Println("error detected: ", err)
	}
	wg.Done()
}

func backup() {
	for {
		mutex.Lock()
		err := c.Backup()
		mutex.Unlock()
		if err != nil {
			break
		}
		time.Sleep(time.Duration(backupInterval) * time.Second)
	}
}

func main() {
	var err error

	fp, err := os.OpenFile("generator.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	multiWriter := io.MultiWriter(fp, os.Stdout)
	log.SetOutput(multiWriter)

	mutex = new(sync.Mutex)

	configFilePath := flag.String("config", "config.json", "transport config path")
	modelFilePath := flag.String("model", "model.sav", "Bayesian model's save path")
	interval := flag.Int("interval", 1, "Bayesian model's save interval(sec)")
	flag.Parse()
	log.Println("model save path is", *modelFilePath)
	log.Println("config file path is", *configFilePath)
	wg = sync.WaitGroup{}
	c, err = classifier.InitClassfier(*modelFilePath)
	if err != nil {
		log.Panicf("initialize the classifier failed %v (filepath: %s)", err, *modelFilePath)
	}
	log.Println("successfully classifier generated")

	backupInterval = *interval
	go backup()
	log.Println("backup runs every", backupInterval, "seconds")
	for {
		defer func() {
			err := recover()
			log.Printf("recover detected: %v\n", err)
		}()
		wg.Add(1)
		go run(*configFilePath)
		wg.Wait()
		log.Println("retry the running sequence after 10 seconds")
		time.Sleep(time.Duration(10) * time.Second)
	}
}
