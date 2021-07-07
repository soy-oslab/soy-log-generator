module github.com/soyoslab/soy_log_generator

replace github.com/soyoslab/soy_log_generator/pkg/compressor => ./pkg/compressor

replace github.com/soyoslab/soy_log_generator/pkg/buffering => ./pkg/buffering

go 1.16

require github.com/klauspost/compress v1.13.1
