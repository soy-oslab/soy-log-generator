FROM golang:1.16

COPY generator generator

RUN apt-get -qy \
	-o "Dpkg::Options::=--force-confdef" \
	-o "Dpkg::Options::=--force-confold" upgrade && \
	apt-get -qy autoclean

WORKDIR /go/generator

RUN make
