FROM golang:1.11

RUN go get github.com/golang/dep/cmd/dep

WORKDIR /go/src/github.com/greendocker

COPY . .

RUN go install

ENTRYPOINT [ "/go/bin/greendocker" ]