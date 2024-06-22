FROM golang:1.22.4

ENV GO111MODULE=auto
ENV GOPROXY="https://goproxy.cn"

RUN go get k8s.io/code-generator; exit 0
WORKDIR /go/src/k8s.io/code-generator
RUN go get -d ./...

RUN mkdir -p /go/src///
VOLUME /go/src///

WORKDIR /go/src///
FROM golang:1.22.4

ENV GO111MODULE=auto
ENV GOPROXY="https://goproxy.cn"

RUN go get k8s.io/code-generator; exit 0
WORKDIR /go/src/k8s.io/code-generator
RUN go get -d ./...

RUN mkdir -p /go/src///
VOLUME /go/src///

WORKDIR /go/src///
FROM golang:1.22.4

ENV GO111MODULE=auto
ENV GOPROXY="https://goproxy.cn"

RUN go get k8s.io/code-generator; exit 0
WORKDIR /go/src/k8s.io/code-generator
RUN go get -d ./...

RUN mkdir -p /go/src///
VOLUME /go/src///

WORKDIR /go/src///
FROM golang:1.22.4

ENV GO111MODULE=auto
ENV GOPROXY="https://goproxy.cn"

RUN go get k8s.io/code-generator; exit 0
WORKDIR /go/src/k8s.io/code-generator
RUN go get -d ./...

RUN mkdir -p /go/src///
VOLUME /go/src///

WORKDIR /go/src///
FROM golang:1.22.4

ENV GO111MODULE=auto
ENV GOPROXY="https://goproxy.cn"

RUN go get k8s.io/code-generator; exit 0
WORKDIR /go/src/k8s.io/code-generator
RUN go get -d ./...

RUN mkdir -p /go/src///
VOLUME /go/src///

WORKDIR /go/src///
FROM golang:1.22.4

ENV GO111MODULE=auto
ENV GOPROXY="https://goproxy.cn"

RUN go get k8s.io/code-generator; exit 0
WORKDIR /go/src/k8s.io/code-generator
RUN go get -d ./...

RUN mkdir -p /go/src///
VOLUME /go/src///

WORKDIR /go/src///
FROM golang:1.22.4

ENV GO111MODULE=auto
ENV GOPROXY="https://goproxy.cn"

RUN go get k8s.io/code-generator; exit 0
WORKDIR /go/src/k8s.io/code-generator
RUN go get -d ./...

RUN mkdir -p /go/src///
VOLUME /go/src///

WORKDIR /go/src///
FROM golang:1.22.4

ENV GO111MODULE=auto
ENV GOPROXY="https://goproxy.cn"

RUN go get k8s.io/code-generator; exit 0
WORKDIR /go/src/k8s.io/code-generator
RUN go get -d ./...

RUN mkdir -p /go/src///
VOLUME /go/src///

WORKDIR /go/src///
