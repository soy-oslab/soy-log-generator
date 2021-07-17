module github.com/soyoslab/soy_log_generator

replace github.com/soyoslab/soy_log_generator/pkg/compressor => ./pkg/compressor

replace github.com/soyoslab/soy_log_generator/pkg/buffering => ./pkg/buffering

replace github.com/soyoslab/soy_log_generator/pkg/watcher => ./pkg/watcher

replace github.com/soyoslab/soy_log_generator/pkg/scheduler => ./pkg/scheduler

replace github.com/soyoslab/soy_log_generator/pkg/ring => ./pkg/ring

go 1.16

require (
	github.com/Workiva/go-datastructures v1.0.53
	github.com/apache/thrift v0.14.2 // indirect
	github.com/armon/go-metrics v0.3.9 // indirect
	github.com/cloudflare/ahocorasick v0.0.0-20210425175752-730270c3e184
	github.com/edwingeng/doublejump v0.0.0-20200330080233-e4ea8bd1cbed // indirect
	github.com/fatih/color v1.12.0 // indirect
	github.com/fsnotify/fsnotify v1.4.9
	github.com/go-ping/ping v0.0.0-20210506233800-ff8be3320020 // indirect
	github.com/go-redis/redis/v8 v8.11.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/grandcat/zeroconf v1.0.0 // indirect
	github.com/hashicorp/consul/api v1.9.1 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-hclog v0.16.2 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/klauspost/compress v1.13.1
	github.com/klauspost/cpuid/v2 v2.0.8 // indirect
	github.com/klauspost/reedsolomon v1.9.12 // indirect
	github.com/lucas-clemente/quic-go v0.21.2 // indirect
	github.com/mattn/go-isatty v0.0.13 // indirect
	github.com/mcuadros/go-defaults v1.2.0
	github.com/miekg/dns v1.1.43 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/rs/cors v1.8.0 // indirect
	github.com/smallnest/quick v0.0.0-20210406061658-4bf95e372fbd // indirect
	github.com/smallnest/rpcx v1.6.4
	github.com/soheilhy/cmux v0.1.5 // indirect
	github.com/soyoslab/soy_log_collector v0.0.0-20210719112509-b29462985a6e
	github.com/tjfoc/gmsm v1.4.1 // indirect
	github.com/valyala/fastrand v1.0.0 // indirect
	go.opencensus.io v0.23.0 // indirect
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97 // indirect
	golang.org/x/net v0.0.0-20210716203947-853a461950ff // indirect
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	golang.org/x/tools v0.1.5 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)
