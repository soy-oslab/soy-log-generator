package main

import (
	"fmt"

	"github.com/soyoslab/soy_log_generator/pkg/compressor"
)

func main() {
	{
		fmt.Println("Zstd CASE")
		c := &compressor.ZstdComp{}
		fmt.Println([]byte("Damn"))
		buffer, _ := c.Compress([]byte("Damn"))
		fmt.Println(buffer)
		bytes, _ := c.Decompress(buffer)
		fmt.Println(bytes)
	}
	{
		fmt.Println("Gzip CASE")
		c := &compressor.GzipComp{}
		fmt.Println([]byte("Damn"))
		buffer, _ := c.Compress([]byte("Damn"))
		fmt.Println(buffer)
		bytes, _ := c.Decompress(buffer)
		fmt.Println(bytes)
	}
}
