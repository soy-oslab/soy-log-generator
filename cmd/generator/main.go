package main

import (
	"log"

	"github.com/soyoslab/soy_log_generator/pkg/scheduler"
)

func main() {
	ops := scheduler.SubmitOperations{}
	ops.Hot = func(messages []scheduler.Message) error {
		for _, message := range messages {
			log.Println("hot", message.Info, string(message.Data))
		}
		return nil
	}
	ops.Cold = func(messages []scheduler.Message) error {
		for _, message := range messages {
			log.Println("cold", message.Info, string(message.Data))
		}
		return nil
	}
	s, err := scheduler.InitScheduler("config.json", ops, nil)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(s)
	defer s.Close()
	err = s.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
