package main

import (
	"fmt"
	"github.com/soyoslab/soy_log_generator/compressor"
)

func main() {
	buffer := compressor.Compress([]byte("Damn"))
	message := string(compressor.Decompress(buffer))
	fmt.Println(message)
}
