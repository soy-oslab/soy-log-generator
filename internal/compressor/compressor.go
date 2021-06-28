package compressor

type Compressor interface {
	Compress(data []byte) []byte
	Decompress(data []byte) []byte
}
