module github.com/soyoslab/soy_log_generator

replace github.com/soyoslab/soy_log_generator/pkg/compressor => ./pkg/compressor

replace github.com/soyoslab/soy_log_generator/pkg/buffering => ./pkg/buffering

replace github.com/soyoslab/soy_log_generator/pkg/watcher => ./pkg/watcher

go 1.16

require (
	github.com/fsnotify/fsnotify v1.4.9
	github.com/klauspost/compress v1.13.1
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
)
