FROM golang:1.11

RUN go get github.com/golang/dep/cmd/dep

WORKDIR /go/src/github.com/greendocker

COPY . .

RUN go install

ENTRYPOINT ["/go/bin/greendocker"]

#FROM alpine:latest
#RUN apk --no-cache add ca-certificates
#WORKDIR /root/
#COPY --from=builder /go/bin/greendocker .
#ENTRYPOINT ["/root/greendocker"]
