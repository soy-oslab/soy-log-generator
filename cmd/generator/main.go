package main

import (
	"fmt"
	"github.com/soyoslab/soy_log_generator/compressor"
)

func main() {
	{
		fmt.Println("Zstd CASE")
		c := &compressor.ZstdComp{}
		fmt.Println([]byte("Damn"))
		buffer := c.Compress([]byte("Damn"))
		fmt.Println(buffer)
		fmt.Println(c.Decompress(buffer))
	}
	{
		fmt.Println("Gzip CASE")
		c := &compressor.GzipComp{}
		fmt.Println([]byte("Damn"))
		buffer := c.Compress([]byte("Damn"))
		fmt.Println(buffer)
		fmt.Println(c.Decompress(buffer))
	}
}
