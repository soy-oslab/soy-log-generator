package main

import (
	"fmt"
	"github.com/soyoslab/soy_log_generator/compressor"
)

func main() {
	message := compressor.Hello("Damn")
	fmt.Println(message)
}
