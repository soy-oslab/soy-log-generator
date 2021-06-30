package compressor

// Compressor is generic interface for compressor algorithms
type Compressor interface {
	Compress(data []byte) []byte
	Decompress(data []byte) []byte
}
