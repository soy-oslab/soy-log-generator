module github.com/soyoslab/soy_log_generator

replace github.com/soyoslab/soy_log_generator/pkg/compressor => ./pkg/compressor

replace github.com/soyoslab/soy_log_generator/pkg/buffering => ./pkg/buffering

replace github.com/soyoslab/soy_log_generator/pkg/watcher => ./pkg/watcher

replace github.com/soyoslab/soy_log_generator/pkg/scheduler => ./pkg/scheduler

replace github.com/soyoslab/soy_log_generator/pkg/ring => ./pkg/ring

go 1.16

require (
	github.com/Workiva/go-datastructures v1.0.53
	github.com/cloudflare/ahocorasick v0.0.0-20210425175752-730270c3e184
	github.com/fsnotify/fsnotify v1.4.9
	github.com/klauspost/compress v1.13.1
	github.com/mcuadros/go-defaults v1.2.0
	github.com/soyoslab/soy_log_collector v0.0.0-20210713045715-bfa6cd2fdf5a
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)
