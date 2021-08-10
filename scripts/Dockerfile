FROM golang:1.16

COPY generator generator

RUN apt-get -qy \
	-o "Dpkg::Options::=--force-confdef" \
	-o "Dpkg::Options::=--force-confold" upgrade && \
	apt-get -qy autoclean

RUN apt-get update && apt-get install -y python3

WORKDIR /go/generator

RUN make

RUN make kube-generator-wrapper.c

CMD python3 ./scripts/config-file-generator.py; ./build/kube/wrapper ./build/generator/main
