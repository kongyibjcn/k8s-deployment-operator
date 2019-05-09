FROM golang:1.11 as builder

RUN go get github.com/golang/dep/cmd/dep

WORKDIR /go/src/github.com/bluedocker

COPY . .

RUN go build -o /go/bin/bluedocker

ENTRYPOINT ["/go/bin/bluedocker"]

#FROM alpine:latest
#RUN apk --no-cache add ca-certificates
#WORKDIR /root
#COPY --from=builder /go/bin/bluedocker .
#ENTRYPOINT ["/root/bluedocker"]

#FROM scratch
# Copy our static executable.
#COPY --from=builder /go/bin/bluedocker /go/bin/bluedocker
# Run the hello binary.
#ENTRYPOINT ["/go/bin/bluedocker"]