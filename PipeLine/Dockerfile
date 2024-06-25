FROM golang
RUN mkdir -p /go/src/pipeline
WORKDIR /go/src/pipeline
ADD *.go .
ADD go.mod .
RUN go install . 

FROM alpine
LABEL version="2.1"
LABEL maintainer="alexandr.snisarenko@mail.com>"
WORKDIR /root/
COPY --from=0 /go/bin/pipeline .
#ENTRYPOINT ./pipeline
#автоматически не стартует, запускать вручную, в интерактивном режиме: docker run -ti --rm pipeline